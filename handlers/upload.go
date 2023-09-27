package handlers

import (
	"davigo/s3uploader/awsUtils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func HandleUpload(c *fiber.Ctx) error {
	// Parse the multipart form:
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	// => *multipart.Form

	// Get all files from "documents" key:
	files := form.File["images"]
	// => []*multipart.FileHeader
	fmt.Println("Got files")
	// Loop through files:
	for _, file := range files {
		fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
		// Send file to s3
		go awsUtils.UploadToS3("voodireito", file)

	}
	return c.SendString("Uploaded")
}
