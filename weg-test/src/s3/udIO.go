package s3

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	"strconv"
	"sync/atomic"
	"time"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var MyBucket string
var filepath string

const (
	maxPartSize        = int64(2 * 1024 * 1024)
	maxRetries         = 3
)

var ops uint64

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

func UploadSingleBuffer(c *gin.Context) {
	startTime := time.Now()

	body := c.Request.Body

	b, err :=ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	body = nil
	c.Request.Body.Close()

	buffer := []byte(b)

	var elapsedTime = time.Since(startTime)

	s3sess := c.MustGet("s3sess").(*session.Session)
	uploader := s3manager.NewUploader(s3sess)

	atomic.AddUint64(&ops, 1)

	filename := "test"+strconv.FormatUint(ops,10)+""+strconv.FormatInt(time.Now().UnixNano(),10)


	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(GetEnvWithKey("BUCKET_NAME")),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(buffer),
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

	up = nil

	filepath = "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/" + filename
	c.JSON(http.StatusOK, gin.H{
		"filepath":    filepath,
	})
	elapsedTime = time.Since(startTime)
	fmt.Printf("response: %s", elapsedTime.Seconds())
}

func UploadSingle(c *gin.Context)  {
	startTime := time.Now()

	body := c.Request.Body

	b, err :=ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	var elapsedTime = time.Since(startTime)

	var filename = "c:/tmp/test"

	fmt.Println("fileName : ",filename)

	if err := EnsureDir("c:/tmp"); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}

	err = ioutil.WriteFile(filename, b, 0644)

	file, err := os.Open(filename)
	if err != nil {

	}
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	elapsedTime = time.Since(startTime)
	fmt.Printf("open file: %s\n", elapsedTime.Seconds())

	s3sess := c.MustGet("s3sess").(*session.Session)
	uploader := s3manager.NewUploader(s3sess)

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
	fmt.Printf("response: %s", elapsedTime.Seconds())

	os.Remove(filename)
}

func UploadMulti(c *gin.Context)  {
	// FILE TO S3
	startTime := time.Now()

	body := c.Request.Body

	b, err :=ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	var elapsedTime = time.Since(startTime)

	atomic.AddUint64(&ops, 1)

	var filePath = "c:/tmp"

	var fileName = "test"+strconv.FormatUint(ops,10)+""+strconv.FormatInt(time.Now().UnixNano(),10)

	fmt.Println("fileName : ",fileName)

	if err := EnsureDir(filePath); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}

	var fileFullName = filePath+"/"+fileName

	err = ioutil.WriteFile(fileFullName, b, 0644)
	file, err := os.Open(fileFullName)
	if err != nil {

	}
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	s3sess := c.MustGet("s3sess").(*session.Session)

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)
	fileType := http.DetectContentType(buffer)
	file.Read(buffer)

	path := file.Name()
	input := &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(GetEnvWithKey("BUCKET_NAME")),
		Key:         aws.String(path),
		ContentType: aws.String(fileType),
	}

	svc := s3.New(s3sess, s3sess.Config)

	resp, err := svc.CreateMultipartUpload(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Created multipart upload request")

	var curr, partLength int64
	var remaining = size
	var completedParts []*s3.CompletedPart
	partNumber := 1
	for curr = 0; remaining != 0; curr += partLength {
		if remaining < maxPartSize {
			partLength = remaining
		} else {
			partLength = maxPartSize
		}
		completedPart, err := uploadPart(svc, resp, buffer[curr:curr+partLength], partNumber)
		if err != nil {
			fmt.Println(err.Error())
			err := abortMultipartUpload(svc, resp)
			if err != nil {
				fmt.Println(err.Error())
			}
			return
		}
		remaining -= partLength
		partNumber++
		completedParts = append(completedParts, completedPart)
	}

	completeResponse, err := completeMultipartUpload(svc, resp, completedParts)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Successfully uploaded file: %s\n", completeResponse.String())

	elapsedTime = time.Since(startTime)
	fmt.Printf("response: %s", elapsedTime.Seconds())
	os.Remove(fileFullName)
}

func UploadGoRoutine(c *gin.Context) {
	startTime := time.Now()

	body := c.Request.Body

	b, err :=ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	body = nil
	c.Request.Body.Close()

	buffer := []byte(b)

	var elapsedTime = time.Since(startTime)

	s3sess := c.MustGet("s3sess").(*session.Session)
	uploader := s3manager.NewUploader(s3sess)

	atomic.AddUint64(&ops, 1)

	filename := "test"+strconv.FormatUint(ops,10)+""+strconv.FormatInt(time.Now().UnixNano(),10)

	uploadGoRoutineS3(buffer,filename,uploader)

	filepath = "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/" + filename
	c.JSON(http.StatusOK, gin.H{
		"filepath":    filepath,
	})

	elapsedTime = time.Since(startTime)
	fmt.Printf("response: %s", elapsedTime.Seconds())
}

func UploadEfs(c *gin.Context)  {
	startTime := time.Now()

	body := c.Request.Body

	b, err :=ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	var elapsedTime = time.Since(startTime)

	var filename = "c:/tmp/test"

	fmt.Println("fileName : ",filename)

	if err := EnsureDir("c:/tmp"); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}

	err = ioutil.WriteFile(filename, b, 0644)

	file, err := os.Open(filename)
	if err != nil {

	}
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	elapsedTime = time.Since(startTime)
	fmt.Printf("open file: %s\n", elapsedTime.Seconds())

	s3sess := c.MustGet("s3sess").(*session.Session)
	uploader := s3manager.NewUploader(s3sess)

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
	fmt.Printf("response: %s", elapsedTime.Seconds())
}

func Download(c *gin.Context)  {
	s3sess := c.MustGet("s3sess").(*session.Session)
	downloader := s3manager.NewDownloader(s3sess)

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

func uploadGoRoutineS3(b []byte,filename string, uploader *s3manager.Uploader) {
	go func() {
		_, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(GetEnvWithKey("BUCKET_NAME")),
			Key:    aws.String(filename),
			Body:   bytes.NewReader(b),
		})
		if err != nil {
			log.Fatal(err)
			b = nil
			return
		}
		b = nil
	}()
}

func completeMultipartUpload(svc *s3.S3, resp *s3.CreateMultipartUploadOutput, completedParts []*s3.CompletedPart) (*s3.CompleteMultipartUploadOutput, error) {
	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}
	return svc.CompleteMultipartUpload(completeInput)
}

func uploadPart(svc *s3.S3, resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNumber int) (*s3.CompletedPart, error) {
	tryNum := 1
	partInput := &s3.UploadPartInput{
		Body:          bytes.NewReader(fileBytes),
		Bucket:        resp.Bucket,
		Key:           resp.Key,
		PartNumber:    aws.Int64(int64(partNumber)),
		UploadId:      resp.UploadId,
		ContentLength: aws.Int64(int64(len(fileBytes))),
	}

	for tryNum <= maxRetries {
		uploadResult, err := svc.UploadPart(partInput)
		if err != nil {
			if tryNum == maxRetries {
				if aerr, ok := err.(awserr.Error); ok {
					return nil, aerr
				}
				return nil, err
			}
			fmt.Printf("Retrying to upload part #%v\n", partNumber)
			tryNum++
		} else {
			fmt.Printf("Uploaded part #%v\n", partNumber)
			return &s3.CompletedPart{
				ETag:       uploadResult.ETag,
				PartNumber: aws.Int64(int64(partNumber)),
			}, nil
		}
	}
	return nil, nil
}

func abortMultipartUpload(svc *s3.S3, resp *s3.CreateMultipartUploadOutput) error {
	fmt.Println("Aborting multipart upload for UploadId#" + *resp.UploadId)
	abortInput := &s3.AbortMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
	}
	_, err := svc.AbortMultipartUpload(abortInput)
	return err
}