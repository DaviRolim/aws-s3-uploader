package awsUtils

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/joho/godotenv"
)

var cfg *aws.Config

func Setup() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	awsKey := os.Getenv("AWS_KEY")
	awsSecret := os.Getenv("AWS_SECRET")
	// Initialize AWS Credentials
	credentials := credentials.NewStaticCredentials(awsKey, awsSecret, "")

	cfg = aws.NewConfig().WithRegion("us-west-1").WithCredentials(credentials)
}
