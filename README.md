# QR-Generator
## Introduction
QR-Generator is a small REST application that can be used to generate your QR-Code.
It's API supports two kind of QR Code:
- `/qrcode`: generates a simple black and white QR Code
- `/qrcode-overlay`: generates a QR Code with logo overlay.

It is written in Golang using [Gin Web Framework](https://pkg.go.dev/github.com/gin-gonic/gin) and [go-qrcode](https://pkg.go.dev/github.com/skip2/go-qrcode) library for this reason it's very fast and lightweight.

The project is opensource and if you want you can build it again or change it's logic.

## How to use it

### Build and run locally
You can build the project using the following command:
```shell
go build -o app .
```
And run it locally:
```shell
./app
```

### Docker
The image is already present on [Docker Hub](https://hub.docker.com/repository/docker/rmarcello/qr-generator).

Otherwise you can build your own:
```shell
docker build -t qr-generator:latest .
```

To run with docker you can use this:
```shell
docker run -p 8080:8080 -it --rm --name qr-generator rmarcello/qr-generator:latest
```

## How it works
The REST application listen on HTTP 8080 port. It provedes a simple REST API and a UI.

### UI
You can access the UI from a browser using the `/ui` url. From here you can choose a QRCode type and generate it.

![UI Example](test/ui.png)

### API

If you need a simple QRCode you can use the `/qrcode` API using the GET method.
The following parameters are required:
  - `txt`: represents the string that the user wants to encode;
  - `size`: represents the size of the encoded QR code express in pixels.
The API returns a PNG file that can be saved in a file.

Here's an example:

```shell
curl "http://localhost:8080/qrcode?txt=https://marcelloraffaele.github.io&size=512" --output qrcode-without-logo.png
```

If you need a QR Code with logo overlay, you can use the `/qrcode-overlay` API using the POST method in order to upload the file that you want to use as logo. 
The following parameters are required:
  - `txt`: represents the string that the user wants to encode;
  - `size`: represents the size of the encoded QR code express in pixels.
  - `file`: represents the logo file in png format that will be overlayed on the encoded QR code.
The image must be in PNG format, it will be resized and overlayed on the generated QR Code.
For the logo, choose 150x150px size. 

```shell
cd test
curl "http://localhost:8080/qrcode-overlay" \
    -X POST \
    -F txt=https://marcelloraffaele.github.io \
    -F size=512 \
    -F 'file=@"logo.png"' \
    --output qrcode-with-logo.png
```