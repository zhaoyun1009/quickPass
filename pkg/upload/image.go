package upload

import (
	"QuickPass/pkg/logf"
	"fmt"
	"github.com/minio/minio-go/v6"
	"io"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"QuickPass/pkg/file"
	"QuickPass/pkg/setting"
	"QuickPass/pkg/util"
)

var MinioClient *minio.Client

func SetUp() {
	var err error
	minioSetting := setting.MinioSetting
	// 初使化minio client对象。
	MinioClient, err = minio.New(minioSetting.Endpoint, minioSetting.AccessKey, minioSetting.SecretKey, minioSetting.UseSSL)
	if err != nil {
		logf.Fatal(err)
	}

	bucketName := minioSetting.BucketName
	//创建存储桶并设置只读权限
	exists, err := MinioClient.BucketExists(bucketName)
	if err != nil {
		logf.Fatal("BucketExists", err)
	}
	if !exists {
		location := "us-east-1"
		err := MinioClient.MakeBucket(bucketName, location)
		if err != nil {
			logf.Fatal("MakeBucket", err)
		}
	}
	logf.Info("minio SetUp Successfully")
}

// GetImageName get image name
func GetImageName(name string) string {
	ext := path.Ext(name)
	id := util.GetUniqueName()

	datetime := time.Now().Format(util.TIME_TEMPLATE_5)
	return fmt.Sprintf("%s/%s%s", datetime, id, ext)
}

// CheckImageExt check image file ext
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

// CheckImageSize check image size
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		logf.Error(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}

func PutObject(objectName, contentType, contentDisposition string, reader io.Reader, objectSize int64) (string, string, int64, error) {
	var opts minio.PutObjectOptions
	opts.ContentType = contentType
	opts.ContentDisposition = contentDisposition
	object, err := MinioClient.PutObject(setting.MinioSetting.BucketName, objectName, reader, objectSize, opts)
	return setting.MinioSetting.PreUrl, setting.MinioSetting.BucketName, object, err
}
