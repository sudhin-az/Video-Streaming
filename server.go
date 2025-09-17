package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const uploadDir = "./uploads"

func main() {
	r := gin.Default()
	r.POST("/upload", uploadFile)
	r.GET("/stream/:filename", streamFile)
	r.DELETE("/delete/:filename", deleteFile)

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}
	fmt.Println("Server started running at port:8080")
	r.Run(":8080")
}

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}
	dst := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully", 
		"filename": file.Filename,
	})
}

func streamFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(uploadDir, filename)

	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	http.ServeFile(c.Writer, c.Request, filePath)
}

func deleteFile(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join(uploadDir, filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	err := os.Remove(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
