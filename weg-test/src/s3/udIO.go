package s3

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var MyBucket string
var filepath string

//GetEnvWithKey : get env value
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

func ConnectAws() *session.Session {
	AccessKeyID = GetEnvWithKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	MyRegion = GetEnvWithKey("AWS_REGION")

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})

	if err != nil {
		panic(err)
	}

	fmt.Println("ConnectAws", AccessKeyID," - ", SecretAccessKey," - ", MyRegion," - ",sess)

	return sess
}

func Upload(c *gin.Context) {
	//body := c.Request.Body
	//
	//b, err :=ioutil.ReadAll(body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//buffer := []byte(b)
	//
	//fmt.Println("Uploaded", len(b))
	//
	//sess := c.MustGet("sess").(*session.Session)
	//uploader := s3manager.NewUploader(sess)
	//
	//filename := "test"
	//
	//fmt.Println("Upload", uploader," - ", filename)
	//
	////upload to the s3 bucket
	//up, err := uploader.Upload(&s3manager.UploadInput{
	//	Bucket: aws.String(GetEnvWithKey("BUCKET_NAME")),
	//	Key:    aws.String(filename),
	//	Body:   bytes.NewReader(buffer),
	//})
	//
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error":    "Failed to upload file",
	//		"uploader": up,
	//	})
	//	return
	//}
	//filepath = "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/" + filename
	//c.JSON(http.StatusOK, gin.H{
	//	"filepath":    filepath,
	//})

	startTime := time.Now()

	body := c.Request.Body

	b, err :=ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	var elapsedTime = time.Since(startTime)
	fmt.Printf("body read byte: %s\n", elapsedTime.Seconds())

	var filename = "c:/tmp/test"

	fmt.Println("fileName : ",filename)

	if err := EnsureDir("c:/tmp"); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}


	err = ioutil.WriteFile(filename, b, 0644)

	elapsedTime = time.Since(startTime)
	fmt.Printf("write file: %s\n", elapsedTime.Seconds())

	file, err := os.Open(filename)
	if err != nil {

	}
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	elapsedTime = time.Since(startTime)
	fmt.Printf("open file: %s\n", elapsedTime.Seconds())

	sess := c.MustGet("sess").(*session.Session)
	uploader := s3manager.NewUploader(sess)

	fmt.Println("Upload", uploader," - ", filename)

	//upload to the s3 bucket
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(GetEnvWithKey("BUCKET_NAME")),
		Key:    aws.String(filename),
		Body:   file,
	})

	elapsedTime = time.Since(startTime)
	fmt.Printf("s3 upload file: %s\n", elapsedTime.Seconds())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to upload file",
			"uploader": up,
		})
		return
	}
	filepath = "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/" + filename
	c.JSON(http.StatusOK, gin.H{
		"filepath":    filepath,
	})

	elapsedTime = time.Since(startTime)
	fmt.Printf("response: %s\n", elapsedTime.Seconds())

	os.Remove(filename)
}

func Download(c *gin.Context)  {
//	sess := c.MustGet("sess").(*session.Session)
//
//	item := "test"
//
//	svc := s3.New(sess)
//
//	rawObject, _ := svc.GetObject(
//		&s3.GetObjectInput{
//			Bucket: aws.String(GetEnvWithKey("BUCKET_NAME")),
//			Key:    aws.String(item),
//		})
//
//	buf := new(bytes.Buffer)
//	buf.ReadFrom(rawObject.Body)
//	myFileContentAsString := buf.String()
//
//	fmt.Println("Downloaded ", myFileContentAsString, " bytes")
//
//	c.Header("Content-Type", "application/octet-stream")
//	c.Header("Content-Transfer-Encoding", "binary")
//	c.Writer.Write(buf.Bytes())
////	c.FileAttachment(item,item)
//
////	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	sess := c.MustGet("sess").(*session.Session)
	downloader := s3manager.NewDownloader(sess)

	path := "c:/tmp/"

	item := "test"

	file, err := os.Create(path+item)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", item, err)
	}

	defer file.Close()

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(GetEnvWithKey("BUCKET_NAME")),
			Key:    aws.String(item),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", item, err)
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.FileAttachment(path+item,item)

	os.Remove(path+item)

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func EnsureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}