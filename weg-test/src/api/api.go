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

	db "weg-test/src/db"
)

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
	router.Run("localhost:8080")
}

func getFileInfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, db.SelectFileInfo())
}

func getFileReq(c *gin.Context) {
	id := c.Param("id")
	c.IndentedJSON(http.StatusOK, db.SelectFileReq(id))
}

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