package awsconfig

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

func Connection() *s3.Client {
	client, err := S3Client()
	if err != nil {
		fmt.Println("Error create client S3: " + err.Error())
	}
	return client
}

// configS3 creates the S3 client
func S3Client() (*s3.Client, error) {
	AWS_S3_KEY := viper.GetString("KEY")
	AWS_S3_SECRET := viper.GetString("SECRET")
	AWS_S3_REGION := viper.GetString("REGION")
	creds := credentials.NewStaticCredentialsProvider(AWS_S3_KEY, AWS_S3_SECRET, "")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(creds), config.WithRegion(AWS_S3_REGION))
	if err != nil {
		fmt.Println("\nConection failed with credential")
		log.Fatal(err)
	}
	//fmt.Println("Conection S3 Sucefull ;)")
	awsS3Client := s3.NewFromConfig(cfg)
	return awsS3Client, nil
}

func NewS3Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return client, nil
}
