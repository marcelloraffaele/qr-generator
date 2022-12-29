package main

import (
	"fmt"
	"net/http"
    "github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"image"
	"image/png"
	"image/draw"
	"bytes"
)

type QRCodeForm struct {
    Txt string `form:"txt"`
	Size int `form:"size"`
}

func main() {
	
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Static("/ui", "./ui")
	

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/qrcode", func(c *gin.Context) {

		png, err := EncodeQRCode(c)
		if err != nil {
			fmt.Println("error:", err)
			c.AbortWithStatusJSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Data(200,"image/png", png)
		
	})

	router.POST("/qrcode-overlay", func(c *gin.Context) {

		png, err := EncodeQRCodeOverlay(c)
		if err != nil {
			fmt.Println("error:", err)
			c.AbortWithStatusJSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Data(200,"image/png", png)
		
	})

	fmt.Println("QR-generator is listening on port", 8080)
	router.Run("0.0.0.0:8080")
	
}

func EncodeQRCode( c *gin.Context ) ([]byte, error) {

	//form binding
	form, err := parseQRCodeRequest(c)	
	if err != nil {
		return nil, err
	}

	png, errEncode := qrcode.Encode(form.Txt, qrcode.Highest, form.Size)
	if errEncode != nil {
		return nil, errEncode
	}

	return png, nil
}

func EncodeQRCodeOverlay( c *gin.Context ) ([]byte, error) {
	
	//form binding
	form, err := parseQRCodeRequest(c)	
	if err != nil {
		return nil, err
	}

	// Single file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}

	openedFile, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	
	logo, err := png.Decode( openedFile )
	if err != nil {
		fmt.Println("logo decode")
		return nil, err
	}

	//qrcode
	qrcodeByte, err := qrcode.Encode(form.Txt, qrcode.Highest, form.Size)
	if err != nil {
		fmt.Println("qrencode")
		return nil, err
	}
	qrcode, err := png.Decode( bytes.NewReader(qrcodeByte) )
	if err != nil {
		fmt.Println("qr png decode")
		return nil, err
	}
	
	
	qrcodeBounds := qrcode.Bounds()
    new_image := image.NewRGBA(qrcodeBounds)    
	draw.Draw(new_image, qrcodeBounds, qrcode, image.ZP, draw.Src)

	logoBounds := logo.Bounds()
	offset := image.Pt(qrcodeBounds.Max.X/2 - logoBounds.Max.X/2, qrcodeBounds.Max.Y/2 - logoBounds.Max.Y/2)
	draw.Draw(new_image, logo.Bounds().Add(offset), logo, image.ZP, draw.Over)

	buf := new(bytes.Buffer)
	err1 := png.Encode(buf, new_image)
	if err1 != nil {
		return nil, err
	}
	file := buf.Bytes()

	return file, nil
}

func parseQRCodeRequest(c *gin.Context) (QRCodeForm, error) {
	var form QRCodeForm
	if err := c.ShouldBind(&form); err != nil {
		return form, err
	}	
	if len:= len(form.Txt); len < 3 {
		return form, fmt.Errorf("Txt length invalid %d should greater than 3", len)
	}
	if form.Size < 100 {
		return form, fmt.Errorf("Size value invalid %d should greater than 100", form.Size)
	}
	return form, nil
}