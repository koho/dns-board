package server

import (
	"github.com/gin-gonic/gin"
	"github.com/koho/dns-board/controllers"
	"github.com/koho/dns-board/middleware"
	"github.com/koho/dns-board/static"
	"io/fs"
	"log"
	"net/http"
)

func NewRouter() (r *gin.Engine) {
	if gin.Mode() == gin.ReleaseMode {
		r = gin.New()
		r.Use(gin.Recovery())
	} else {
		r = gin.Default()
	}
	setupStatic(r)
	r.POST("/auth", controllers.Login)
	api := r.Group("/api", middleware.AuthRequired())
	{
		api.GET("/meta", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"identity": identity,
				"version":  version,
				"startup":  startup,
				"mapUrl":   mapURL,
			})
		})
		api.GET("/domain", controllers.GetDomainTable)
		api.GET("/duration", controllers.GetRequestDuration)
		api.GET("/client", controllers.GetQueryIPStat)
		api.GET("/count", controllers.GetQueryCountStat)
		api.GET("/qtype", controllers.GetQueryTypeStat)
		api.GET("/search", controllers.SearchRecord)
		api.GET("/rcode", controllers.GetRCodeStat)
		api.GET("/size", controllers.GetResponseSize)
		api.GET("/cache", controllers.GetCacheHit)
		api.GET("/raw", controllers.GetRawData)
	}
	return
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
