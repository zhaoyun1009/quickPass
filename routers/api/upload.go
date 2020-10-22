package api

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/e"
	"QuickPass/pkg/upload"
	"fmt"
	"github.com/gin-gonic/gin"
)

// @Summary 图片上传
// @Produce  json
// @Param image formData file true "Image File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/upload/image [post]
func UploadImage(c *gin.Context) {
	f, image, err := c.Request.FormFile("image")
	if err != nil {
		app.ErrorResp(c, e.ERROR, err.Error())
		return
	}

	if image == nil {
		app.ErrorResp(c, e.INVALID_PARAMS, "")
		return
	}

	imageName := upload.GetImageName(image.Filename)

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(f) {
		app.ErrorResp(c, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, "")
		return
	}

	src, _ := image.Open()

	contentType := image.Header.Get("Content-Type")
	contentDisposition := image.Header.Get("Content-Disposition")
	endPoint, bucketName, _, err := upload.PutObject(imageName, contentType, contentDisposition, src, image.Size)
	if err != nil {
		app.ErrorResp(c, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, err.Error())
		return
	}

	app.SuccessResp(c, RespUploadImage{
		ImageUrl:     fmt.Sprintf("http://%s/%s/%s", endPoint, bucketName, imageName),
		ImageSaveUrl: fmt.Sprintf("%s/%s", bucketName, imageName),
	})
}

type RespUploadImage struct {
	ImageUrl     string `json:"image_url"`
	ImageSaveUrl string `json:"image_save_url"`
}
