package service

import (
	"bytes"
	"context"
	Base64 "encoding/base64"
	"fmt"
	"os"
	"portfolio-website/internals/model"
	"portfolio-website/internals/reposiotry"
	"regexp"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
)

type AllowedExtension string

const (
	// MaxFileSize defines maximum allowed decoded file size (20 MB).
	MaxFileSize = 20 * 1024 * 1024

	// Images
	ExtPNG  AllowedExtension = "png"
	ExtJPEG AllowedExtension = "jpeg"
	ExtJPG  AllowedExtension = "jpg"
	ExtWEBP AllowedExtension = "webp"

	// Videos
	ExtMP4 AllowedExtension = "mp4"
	ExtMOV AllowedExtension = "mov"
	ExtAVI AllowedExtension = "avi"

	// GIFs
	ExtGIF AllowedExtension = "gif"
)

var validExtensions = map[AllowedExtension]bool{
	ExtPNG:  true,
	ExtJPEG: true,
	ExtJPG:  true,
	ExtWEBP: true,
	ExtMP4:  true,
	ExtMOV:  true,
	ExtAVI:  true,
	ExtGIF:  true,
}

type UploadService struct {
	UploadRepo *reposiotry.UploadRepo
}

func NewUploadService(uploadRepo *reposiotry.UploadRepo) *UploadService {
	return &UploadService{UploadRepo: uploadRepo}
}

func (s *UploadService) Upload_Service(ctx context.Context, base64 string) (string, error) {
	re1 := regexp.MustCompile(`^data:(?P<type>[\w/+-]+);base64,(?P<string>.+)$`)
	matches := re1.FindStringSubmatch(base64)
	if len(matches) < 2 {
		return "", model.ErrInvalidMimeStructure
	}

	var mimeType, base64Str string
	for i, name := range re1.SubexpNames() {
		switch name {
		case "type":
			mimeType = matches[i]
		case "string":
			base64Str = matches[i]
		}
	}

	re2 := regexp.MustCompile(`^[^/]+/([\w+-]+)$`)
	extension := re2.FindStringSubmatch(mimeType)
	if len(extension) < 2 {
		return "", model.ErrInvalidMimeDetails
	}

	extractedExt := AllowedExtension(extension[1])
	if !validExtensions[extractedExt] {
		return "", fmt.Errorf("%w: extension '%s' is not allowed", model.ErrExtensionNotAllowed, extractedExt)
	}

	data, err := Base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", model.ErrDecodingBase64
	}

	if len(data) > MaxFileSize {
		return "", model.ErrFileTooLarge
	}

	// 1. Get Cloudinary client
	cld, err := credentials()
	if err != nil {
		return "", fmt.Errorf("cloudinary credentials error: %w", err)
	}

	// 2. Upload to Cloudinary and get secure delivery URL
	imageUrl, err := uploadImage(cld, data, ctx)
	if err != nil {
		return "", fmt.Errorf("cloudinary upload error: %w", err)
	}

	// Optional: save record to database if you want to track uploads
	if s.UploadRepo != nil {
		_ = s.UploadRepo.SaveFiles(ctx, imageUrl, string(extractedExt), mimeType)
	}

	// 3. Return the real Cloudinary URL
	return imageUrl, nil
}

func credentials() (*cloudinary.Cloudinary, error) {
	_ = godotenv.Load()
	cldUrl := os.Getenv("CLOUDINARY_URL")
	var cld *cloudinary.Cloudinary
	var err error

	if cldUrl != "" {
		cld, err = cloudinary.NewFromURL(cldUrl)
	} else {
		cld, err = cloudinary.New()
	}

	if err != nil {
		return nil, err
	}

	cld.Config.URL.Secure = true
	return cld, nil
}

func uploadImage(cld *cloudinary.Cloudinary, file []byte, ctx context.Context) (string, error) {
	resp, err := cld.Upload.Upload(ctx, bytes.NewReader(file), uploader.UploadParams{
		UniqueFilename: api.Bool(true),
		Overwrite:      api.Bool(true),
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}
