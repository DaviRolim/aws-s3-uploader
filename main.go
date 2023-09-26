package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func uploadImageToS3(bucketName string, fileHeader *multipart.FileHeader, cfg *aws.Config) error {
	// Open the image file
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new session with default session credentials
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(cfg))

	// S3 service client the Upload manager will use.
	svc := s3.New(sess)

	// Configure the S3 object input parameters
	input := &s3.PutObjectInput{
		Body:          file,
		Bucket:        aws.String(bucketName),
		Key:           aws.String(fileHeader.Filename),
		ContentType:   aws.String(fileHeader.Header["Content-Type"][0]),
		ContentLength: aws.Int64(fileHeader.Size),
	}

	fmt.Println("Created input")

	// Upload the image to S3
	_, err = svc.PutObject(input)
	if err != nil {
		return err
	}

	fmt.Println("Uploaded image to s3")

	return nil
}

func main() {
	// Load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	awsKey := os.Getenv("AWS_KEY")
	awsSecret := os.Getenv("AWS_SECRET")
	// Initialize AWS Credentials
	credentials := credentials.NewStaticCredentials(awsKey, awsSecret, "")

	cfg := aws.NewConfig().WithRegion("us-west-1").WithCredentials(credentials)

	// Fiber instance
	app := fiber.New()

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	// Routes
	app.Get("/", hello)

	app.Post("/", func(c *fiber.Ctx) error {
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

			// Save the files to disk:
			// err := c.SaveFile(file, fmt.Sprintf("./%s", file.Filename))
			// Check for errors
			// if err != nil {
			// 	return err
			// }
			// Send file to s3

			fmt.Println("Started go routine")
			go uploadImageToS3("voodireito", file, cfg)

		}
		return nil
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Start server
	log.Fatal(app.Listen(":" + port))
}

// Handler
func hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
