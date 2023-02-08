FROM node

RUN cd ./web && npm install && npm run build

RUN cp -r ./web/dist /

FROM golang:latest

COPY --from=0 /dist ./static

RUN go build -trimpath -ldflags="-s -w" -o /dns-board main.go

FROM alpine:latest

COPY --from=1 /dns-board /dns-board

EXPOSE 80

ENTRYPOINT ["/dns-board"]
