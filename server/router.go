package server

import (
	"github.com/gin-gonic/gin"
	"github.com/koho/dnstap-web/controllers"
	"github.com/koho/dnstap-web/middleware"
	"github.com/koho/dnstap-web/static"
	"io/fs"
	"log"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	setupStatic(r)
	r.POST("/auth", controllers.Login)
	api := r.Group("/api", middleware.AuthRequired())
	{
		api.GET("/meta", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"identity": Identity,
				"version":  Version,
				"startup":  Startup,
				"mapUrl":   MapURL,
			})
		})
		api.GET("/domain", controllers.GetDomainTable)
		api.GET("/duration", controllers.GetRequestDuration)
		api.GET("/client", controllers.GetQueryIPStat)
		api.GET("/count", controllers.GetQueryCountStat)
		api.GET("/qtype", controllers.GetQueryTypeStat)
		api.GET("/search", controllers.SearchRecord)
	}
	return r
}

func setupStatic(r *gin.Engine) {
	// Web files
	dist, err := fs.Sub(static.FS, "dist")
	if err != nil {
		log.Fatal(err)
	}
	staticFS := http.FS(dist)
	assets, err := fs.Sub(dist, "assets")
	if err != nil {
		log.Fatal(err)
	}
	r.StaticFS("/assets", http.FS(assets))
	r.GET("/", func(c *gin.Context) {
		c.FileFromFS("/", staticFS)
	})
	r.GET("/login", func(c *gin.Context) {
		c.FileFromFS("login.html", staticFS)
	})
	entries, err := fs.ReadDir(dist, ".")
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() {
			r.GET("/"+name, func(c *gin.Context) {
				c.FileFromFS(name, staticFS)
			})
		}
	}
}
