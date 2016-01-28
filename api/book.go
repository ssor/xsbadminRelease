package api

import (
	"github.com/gin-gonic/gin"
	// "io"
	"net/http"
	log "xsbAdmin/libs/log"
	// "xsbAdmin/libs/tools"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
	"xsbAdmin/models"
)

var defaultMaxMemory int64 = 5 << 20 // 32 MB

func AddBookToCategory(c *gin.Context) {
	type addBookInfo struct {
		Books    []string `form:"books" json:"books" binding:"required"`
		Category string   `form:"category" json:"category" binding:"required"`
	}

	var abi addBookInfo
	var err error

	err = c.BindJSON(&abi)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["需要提供栏目和图书编号"])
		return
	}

	err = book_manager.AddCategory(abi.Books, abi.Category)
	if err != nil {
		if err == models.Err_not_found {
			c.JSON(http.StatusOK, ResponseMap["图书不存在"])
			return
		}
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	c.JSON(http.StatusOK, ResponseMap["操作成功"])
}

func RemoveBookFromCategory(c *gin.Context) {
	type addBookInfo struct {
		Books    []string `form:"books" json:"books" binding:"required"`
		Category string   `form:"category" json:"category" binding:"required"`
	}

	var abi addBookInfo
	var err error

	err = c.BindJSON(&abi)
	if err != nil {
		if err == models.Err_not_found {
			c.JSON(http.StatusOK, ResponseMap["图书不存在"])
			return
		}
		c.JSON(http.StatusOK, ResponseMap["需要提供栏目和图书编号"])
		return
	}

	err = book_manager.RemoveCategory(abi.Books, abi.Category)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	c.JSON(http.StatusOK, ResponseMap["操作成功"])

}

func BooksOfCategory(c *gin.Context) {
	category := c.Param("category")
	in := c.Query("in")
	inCategory := true
	if in == "false" {
		inCategory = false
	}

	if len(category) <= 0 {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	list, err := book_manager.BooksOfCategory(category, inCategory)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "", list))
}

func Books(c *gin.Context) {

	list, err := book_manager.All()
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "", list))
}

func UploadSingleBook(c *gin.Context) {
	log.TraceF("UploadSingleBook ->")

	header1, err := getUploadFileHeader(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseMap["上传出错"])
		return
	}

	log.TraceF("UploadSingleBook -> fileName: %s", header1.Filename)

	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(header1.Filename)
	tofile := path.Join(Book_dir, fileName)

	err = receiveUploadFile(tofile, header1)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseMap["上传出错"])
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "OK", fileName))
}

func receiveUploadFile(newFileFullPath string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	defer file.Close()

	f, err := os.OpenFile(newFileFullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}
	return nil
}

func getUploadFileHeader(req *http.Request) (*multipart.FileHeader, error) {

	err := req.ParseMultipartForm(defaultMaxMemory)
	if err != nil {
		log.Sys(err.Error())
		return nil, err
	}

	if req.MultipartForm == nil {
		return nil, errors.New("no file uploaded")
	}

	mf := req.MultipartForm
	if len(mf.File) <= 0 {
		return nil, errors.New("no file uploaded")
	}

	var header1 *multipart.FileHeader

	for _, headers := range mf.File {
		header1 = headers[0]
		break
	}
	return header1, nil
}

//upload file and add record to db
func AddBook(c *gin.Context) {
	type bookInfo struct {
		Book string `form:"book" json:"book" binding:"required"`
		Name string `form:"name" json:"name" binding:"required"`
	}

	var err error
	var bi bookInfo

	err = c.BindJSON(&bi)
	if err != nil {
		log.InfoF("AddGroup bindjson error: %s", err.Error())
		c.JSON(http.StatusOK, ResponseMap["需要提供图书服务器编号和图书名称"])
		return
	}
	log.TraceF("AddBook -> Name: %s File: %s", bi.Name, bi.Book)

	book := models.NewBook(bi.Name, bi.Book)

	err = book_manager.Add(book)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	c.JSON(http.StatusOK, ResponseMap["操作成功"])
}

// remove file and record
func RemoveBook(c *gin.Context) {

	id := c.Param("book")
	if len(id) <= 0 {
		c.JSON(http.StatusOK, ResponseMap["需要提供图书编号"])
		return
	}

	err := book_manager.Remove(id)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	c.JSON(http.StatusOK, ResponseMap["操作成功"])
}
