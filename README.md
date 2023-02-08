# DNS Board

A mini dns dashboard for DNSTAP.

## Compile

1. Build web pages

    ```shell
    cd web
    npm run build
    cd ..
    ```

2. Copy dist directory

    ```shell
    cp -r ./web/dist ./static/
    ```

3. Build the binary

    ```shell
    go build main.go
    ```
