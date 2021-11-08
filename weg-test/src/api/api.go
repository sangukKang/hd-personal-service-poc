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
	router.GET("/psnsvc", getPsnsvc)
	router.GET("/psnsvc/:id", getPsnsvc)
	router.GET("/psnsvc/download/:id", getFileDownload)
	router.GET("/psnsvc/sync/:id", getPsnsvcSyncByID)
	router.POST("/psnsvc", postPsnsvc)
	router.POST("/psnsvc/upload", postUpload)
	router.DELETE("/psnsvc/:id", deletePsnsvcByID)
	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getPsnsvc(c *gin.Context) {
	id := c.Param("id")

	c.IndentedJSON(http.StatusOK, db.Select(id))
}


// postAlbums adds an album from JSON received in the request body.
func postPsnsvc(c *gin.Context) {
	tid := c.Request.Header.Get("tid")

	fmt.Println("tid : ",tid)

	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var data map[string]interface{}


	json.Unmarshal([]byte(value), &data) // JSON을 Go언어 자료형으로 변환(여기서는 map으로 변환)

	data["Tid"] = tid

	db.Insert(data)
	c.IndentedJSON(http.StatusCreated, data)
}

// postAlbums adds an album from JSON received in the request body.
func postUpload(c *gin.Context) {
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

func getPsnsvcSyncByID(c *gin.Context) {
//	id := c.Param("id")

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "id not found"})
}

func deletePsnsvcByID(c *gin.Context) {
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