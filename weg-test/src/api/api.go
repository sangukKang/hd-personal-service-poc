package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	db "weg-test/src/db"
)

type CloudFileInfoReq struct {
	name string `json:"name"`
	abcd int `json:"abcd"`
}

var fileName string = ""

func Router() {
//	db.Insert()
	router := gin.Default()
	router.GET("/cloud/fileInfo", getFileInfo)
	router.GET("/cloud/file/:id", getFileReq)
	router.GET("/kafka", testKafka)
	router.GET("/cloud/file/sending", getFileDownload)
	router.GET("/cloud/fileSync", getFileSync)
	router.POST("/cloud/file", insertFileUpdate)
	router.POST("/cloud/file/sending", fileUpload)
	router.DELETE("/cloud/file", deleteFile)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("localhost:8080")
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

	if err := ensureDir("c:/tmp"); err != nil {
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

func getFileSync(c *gin.Context) {
//	id := c.Param("id")

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "id not found"})
}

func deleteFile(c *gin.Context) {
	id := c.Param("id")
	db.Delete(id)
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "id not found"})
}

func ensureDir(dirName string) error {
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