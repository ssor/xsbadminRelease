package api

import (
	"github.com/gin-gonic/gin"
	// "io"
	"net/http"
	log "xsbAdmin/libs/log"
	// "xsbAdmin/libs/tools"
	"xsbAdmin/models"
)

func AddCompany(c *gin.Context) {
	type companyInfo struct {
		Name string `form:"name" json:"name" binding:"required"`
	}

	var err error
	var ci companyInfo
	err = c.BindJSON(&ci)
	if err != nil {
		log.InfoF("AddCompany bindjson error: %s", err.Error())
		c.JSON(http.StatusOK, ResponseMap["需要提供单位名称"])
		return
	}
	com := models.NewCompany(ci.Name)

	err = company_manager.Add(com)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	c.JSON(http.StatusOK, ResponseMap["操作成功"])
}
