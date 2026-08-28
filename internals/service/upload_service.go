package service

import (
	"context"
	Base64 "encoding/base64"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type AllowedExtension string

const (
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

func Upload_Service(ctx context.Context, base64 string) error {
	re1 := regexp.MustCompile(`^data:(?P<type>[\w/+-]+);base64,(?P<string>.+)$`)
	matches := re1.FindStringSubmatch(base64)

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

	if len(matches) < 2 {
		return fmt.Errorf("Invalid MIME type structure")
	}

	extractedExt := AllowedExtension(extension[1])
	if !validExtensions[extractedExt] {
		return fmt.Errorf("Error: extension '%s' is not allowed!\n", extractedExt)
	}

	data, err := Base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return fmt.Errorf("Error Decoding String")
	}

	if err := os.Mkdir("data", 0750); err != nil && !os.IsExist(err) {
		return fmt.Errorf("Error making directory")
	}
	if err := os.Mkdir("data/uploads", 0750); err != nil && !os.IsExist(err) {
		return fmt.Errorf("Error making directory")
	}

	id := uuid.NewString()
	folder := "data/uploads/"
	filename := id + "." + string(extractedExt)
	if err := os.WriteFile(folder+filename, data, 0666); err != nil {
		fmt.Println(err)
		return fmt.Errorf("Error writing file")
	}

	return nil
}

func FindFile_Service(ctx context.Context, filename string) error {
	filenameWithoutExt := strings.TrimSuffix(filename, ".")

	file, err := os.OpenFile(filenameWithoutExt, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("File Not Exist")
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}

	return nil
}
