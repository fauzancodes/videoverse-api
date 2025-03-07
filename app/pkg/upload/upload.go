package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Backblaze/blazer/b2"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/fauzancodes/videoverse-api/app/config"
)

func UploadFile(file any, folder string, filename string) (secureUrl, publicID, cloudName string, err error) {
	cloudName = config.LoadConfig().CloudinaryCloudName
	apiKey := config.LoadConfig().CloudinaryAPIKey
	apiSecret := config.LoadConfig().CLoudinaryAPISecret
	folder = config.LoadConfig().CloudinaryFolder + "/" + folder
	request, _ := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	response, err := request.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder:   folder,
		PublicID: filename,
	})
	if err != nil {
		err = fmt.Errorf("failed to upload file to cloudinary: %s", err.Error())
		return
	}

	secureUrl = response.SecureURL
	publicID = response.PublicID

	return
}

func InitBackbalze(ctx context.Context) (bucket *b2.Bucket, statusCode int, err error) {
	keyID := config.LoadConfig().BackblazeKeyID
	applicationKey := config.LoadConfig().BackblazeApplicationKey
	bucketName := config.LoadConfig().BackblazeBucketName

	b2, err := b2.NewClient(ctx, keyID, applicationKey)
	if err != nil {
		statusCode = http.StatusInternalServerError
		err = fmt.Errorf("failed to connect to Backblaze: %s", err.Error())
		return
	}

	bucket, err = b2.Bucket(ctx, bucketName)
	if err != nil {
		statusCode = http.StatusNotFound
		err = fmt.Errorf("backblaze bucket not found: %s", err.Error())
		return
	}

	return
}

func GetRemoteFile(filename string) (buf bytes.Buffer, statusCode int, err error) {
	ctx := context.Background()
	folder := config.LoadConfig().BackblazeFolder

	bucket, statusCode, err := InitBackbalze(ctx)
	if err != nil {
		statusCode = http.StatusNotFound
		err = fmt.Errorf("failed to initialize Backblaze: %s", err.Error())
		return
	}

	reader := bucket.Object(folder + filename).NewReader(ctx)

	_, err = io.Copy(&buf, reader)
	if err != nil {
		statusCode = http.StatusNotFound
		err = fmt.Errorf("failed to read file: %s", err.Error())
		return
	}

	return
}

func WriteRemoteFile(file bytes.Buffer, filename string) (statusCode int, err error) {
	ctx := context.Background()
	folder := config.LoadConfig().BackblazeFolder

	bucket, statusCode, err := InitBackbalze(ctx)
	if err != nil {
		statusCode = http.StatusNotFound
		err = fmt.Errorf("failed to initialize Backblaze: %s", err.Error())
		return
	}

	fmt.Println("filename:", filename)

	obj := bucket.Object(folder + filename)
	if _, err = obj.Attrs(ctx); err == nil {
		if err = obj.Delete(ctx); err != nil {
			err = fmt.Errorf("failed to delete existing file in backblaze: %s", err.Error())
			return
		}
	}

	w := obj.NewWriter(ctx)
	defer w.Close()

	if _, err = file.WriteTo(w); err != nil {
		w.Close()
		statusCode = http.StatusNotFound
		err = fmt.Errorf("failed to write to backblaze: %s", err.Error())
		return
	}

	return
}
