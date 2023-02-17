FROM node

COPY ./web /web

WORKDIR /web

RUN npm install && npm run build

FROM golang:alpine

COPY . /board

WORKDIR /board

COPY --from=0 /web/dist ./static/dist

RUN apk add --no-cache build-base

RUN CGO_ENABLED=1 go build -tags=go_json -trimpath -ldflags="-s -w" -o /dns-board main.go

FROM alpine:latest

COPY --from=1 /dns-board /

EXPOSE 80

ENTRYPOINT ["/dns-board"]
