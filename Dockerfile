FROM node

COPY ./web /web

RUN cd /web && npm install && npm run build

FROM golang:latest

COPY . /board

COPY --from=0 /web/dist /board/static

RUN cd /board && go build -trimpath -ldflags="-s -w" -o /dns-board main.go

FROM alpine:latest

COPY --from=1 /dns-board /dns-board

EXPOSE 80

ENTRYPOINT ["/dns-board"]
