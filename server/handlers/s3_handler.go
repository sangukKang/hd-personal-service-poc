package handlers

import (
	"context"
	"fmt"
	"os"

	"api-server-poc/config"
	pb "api-server-poc/proto/generated"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	ErrorInvalidArgument = "invalid argument = %s"
)

type S3ManagerServiceServer struct {
	config config.Streamingnvironment
}

func NewS3ManagerServiceServer() *S3ManagerServiceServer {
	return &S3ManagerServiceServer{}
}

func (server *S3ManagerServiceServer) newAwsSession() (*session.Session, error) {
	awsCfg := aws.NewConfig()
	awsCfg.Region = &server.config.S3Region

	return session.NewSession(awsCfg)

}

func (server *S3ManagerServiceServer) Download(ctx context.Context, request *pb.DownloadRequest) (*pb.Result, error) {
	type tempConfig struct {
		filename string
		myBucket string
		myString string
	}

	tempCfg := tempConfig{
		filename: "test",
		myBucket: "test-bucket",
		myString: "test-string",
	}

	// The session the S3 Downloader will use
	sess := session.Must(server.newAwsSession())

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	f, err := os.Create(tempCfg.filename)
	if err != nil {
		// return fmt.Errorf("failed to create file %q, %v", tempCfg.filename, err)
	}

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(tempCfg.myBucket),
		Key:    aws.String(tempCfg.myString),
	})
	if err != nil {
		// return fmt.Errorf("failed to download file, %v", err)
	}
	fmt.Printf("file downloaded, %d bytes\n", n)
	return &pb.Result{}, nil
}

func (server *S3ManagerServiceServer) Upload(ctx context.Context, request *pb.UploadRequest) (*pb.Result, error) {
	type tempConfig struct {
		filename string
		myBucket string
		myString string
	}

	tempCfg := tempConfig{
		filename: "test",
		myBucket: "test-bucket",
		myString: "test-string",
	}

	// The session the S3 Uploader will use
	sess := session.Must(server.newAwsSession())

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(tempCfg.filename)
	if err != nil {
		// return fmt.Errorf("failed to open file %q, %v", tempCfg.filename, err)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(tempCfg.myBucket),
		Key:    aws.String(tempCfg.myString),
		Body:   f,
	})
	if err != nil {
		// return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	return &pb.Result{}, nil
}
