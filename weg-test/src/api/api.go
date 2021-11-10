package api

import (
	"encoding/json"
	"fmt"
	//	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/easonlin404/limit"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	db "weg-test/src/db"
	s3 "weg-test/src/s3"
)

type CloudFileInfoReq struct {
	name string `json:"name"`
	abcd int `json:"abcd"`
}

var fileName string = ""

func Router() {
	s3.LoadEnv()
	s3sess := s3.ConnectAws()

//	uploader := s3manager.NewUploader(sess)
	router := gin.Default()
	s := &http.Server{
		Addr:           "localhost:8080",
		Handler:        router,
//		ReadTimeout:    30 * time.Second,
//		WriteTimeout:   30 * time.Second,
		//  MaxHeaderBytes: 1 << 20,
	}
//	s.SetKeepAlivesEnabled(false)

	_ = &http.Transport{
		IdleConnTimeout:     10 * time.Second,
		MaxIdleConnsPerHost: 1000,
	}


	router.Use(limit.Limit(2000))
	router.GET("/cloud/fileInfo", getFileInfo)
	router.GET("/cloud/file/:id", getFileReq)
	router.GET("/kafka", testKafka)
	router.GET("/cloud/file/sending", getFileDownload)
	router.GET("/cloud/fileSync", getFileSync)
	router.POST("/cloud/file", insertFileUpdate)
	router.POST("/cloud/file/sending", fileUpload)
	router.DELETE("/cloud/file", deleteFile)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//S3 Test
	router.Use(func(c *gin.Context) {
		c.Set("s3sess", s3sess)
		c.Next()
	})

	router.POST("/s3/upload/direct/single/buffer", func(context *gin.Context) {
		s3.UploadSingleBuffer(context)
		context = nil
//		runtime.GC()
	})

	router.POST("/s3/upload/direct/single/file", func(context *gin.Context) {
		s3.UploadSingle(context)
		context = nil
		//		runtime.GC()
	})

	router.POST("/s3/upload/direct/multi", func(context *gin.Context) {
		s3.UploadMulti(context)
		context = nil
	})

	router.POST("/s3/upload/direct/goroutine", func(context *gin.Context) {
		s3.UploadGoRoutine(context)
		context = nil
	})

	router.POST("/s3/upload/direct/efs", func(context *gin.Context) {
		s3.UploadEfs(context)
		context = nil
	})
	
	router.GET("/s3/download", getS3FileDownload)

	s.ListenAndServe()
}

// CFS godoc
// @Summary 파일 리스트 조회
// @Description 파일 리스트 조회
// @name get-string-by-int
// @Accept  json
// @Produce  json
// @Param CloudFileInfoReq body CloudFileInfoReq false "cloudFileInfoReq"
// @Router /cloud/fileInfo [get]
// @Success 200 {object} string
func getFileInfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, db.SelectFileInfo())
}

// CFS godoc
// @Summary 파일 다운로드 요청 조회
// @Description 파일 다운로드 요청 조회
// @name get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path string false "id입니다."
// @Router /cloud/file/{id} [get]
// @Success 200 {object} string
func getFileReq(c *gin.Context) {
	id := c.Param("id")
	c.IndentedJSON(http.StatusOK, db.SelectFileReq(id))
}

// CFS godoc
// @Summary kafka test
// @Description kafka test
// @name get-string-by-int
// @Accept  json
// @Produce  json
// @Router /kafka [get]
// @Success 200 {object} string
func testKafka(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, db.TestKafka())
}

func insertFileUpdate(c *gin.Context) {
	tid := c.Request.Header.Get("tid")

	fmt.Println("tid : ",tid)

	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var data map[string]interface{}


	json.Unmarshal([]byte(value), &data)

	data["Tid"] = tid

	db.Insert(data)
	c.IndentedJSON(http.StatusCreated, data)
}

// CFS godoc
// @Summary fileUpload
// @Description 파일 업로드
// @name get-string-by-int
// @Accept  octet-stream
// @Produce  octet-stream
// @Format binary
// @Param fileSelect formData file false "uploadFile"
// @Router /cloud/file/sending [post]
// @Success 200 {object} string
func fileUpload(c *gin.Context) {
	body := c.Request.Body

	b, err :=ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	fileName = strconv.FormatInt(time.Now().UnixNano(),10)

	fmt.Println("fileName : ","c:/tmp/data"+fileName)

	if err := s3.EnsureDir("c:/tmp"); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}

	err = ioutil.WriteFile("c:/tmp/data"+fileName, b, 0644)

	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusCreated, "200")
}

// CFS godoc
// @Summary filedownload
// @Description 다운로드는 업로드 먼저 진행 하고 할것
// @name get-string-by-int
// @Accept  octet-stream
// @Produce  octet-stream
// @Router /cloud/file/sending [get]
// @Success 200 {object} string
func getFileDownload(c *gin.Context) {
//	id := c.Param("id")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.FileAttachment("c:/tmp/data"+fileName,fileName)


//	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "data not found"})
}


// CFS godoc
// @Summary s3 fileUpload
// @Description 파일 업로드
// @name get-string-by-int
// @Accept  octet-stream
// @Produce  octet-stream
// @Format binary
// @Param fileSelect formData file false "uploadFile"
// @Router /s3/upload [post]
// @Success 200 {object} string
func s3FileUpload(c *gin.Context) {
	//body := c.Request.Body
	//
	//b, err :=ioutil.ReadAll(body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fileName := "c:/tmp/data"+strconv.FormatInt(time.Now().UnixNano(),10)
	//
	//fmt.Println("fileName : ","c:/tmp/data"+fileName)
	//
	//if err := ensureDir("c:/tmp"); err != nil {
	//	fmt.Println("Directory creation failed with error: " + err.Error())
	//	os.Exit(1)
	//}
	//
	//
	//err = ioutil.WriteFile(fileName, b, 0644)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

//	s3.Upload(c)
}

// CFS godoc
// @Summary s3 filedownload
// @Description 다운로드는 업로드 먼저 진행 하고 할것
// @name get-string-by-int
// @Accept  octet-stream
// @Produce  octet-stream
// @Router /s3/download [get]
// @Success 200 {object} string
func getS3FileDownload(c *gin.Context) {
	s3.Download(c)
}

func getFileSync(c *gin.Context) {
//	id := c.Param("id")

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "id not found"})
}

func deleteFile(c *gin.Context) {
	id := c.Param("id")
	db.Delete(id)
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "id not found"})
}
