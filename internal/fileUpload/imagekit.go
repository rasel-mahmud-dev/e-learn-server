package fileUpload

import (
	"context"
	"e-learn/internal/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"log"
	"mime/multipart"
)

func UploadImage2(file multipart.File, fileName string) *uploader.UploadResult {

	cld, err := cloudinary.NewFromParams(
		config.Cfg.CLOUDINARY_CLOUD_NAME,
		config.Cfg.CLOUDINARY_API_KEY,
		config.Cfg.CLOUDINARY_API_SECRET,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
		return nil
	}

	uploadParams := uploader.UploadParams{
		PublicID: fileName,
	}

	// Perform the upload
	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
		return nil
	}

	return uploadResult
}
