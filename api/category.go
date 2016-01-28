package api

import (
	"github.com/gin-gonic/gin"
	// "io"
	"net/http"
	log "xsbAdmin/libs/log"
	// "xsbAdmin/libs/tools"
	"xsbAdmin/models"
)

//categories created by this company
func CompanyCategories(c *gin.Context) {
	companyID := c.Param("company")

	if len(companyID) <= 0 {
		c.JSON(http.StatusOK, ResponseMap["需要提供单位编号"])
		return
	}

	// company, err := company_manager.GetCompanyWithID(companyID)
	// if err != nil {
	// 	c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
	// 	return
	// }

	categories, err := category_manager.GetCompanyCategories(companyID)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "", categories))
}

func Categories(c *gin.Context) {
	companyID := c.Param("company")

	if len(companyID) <= 0 {
		c.JSON(http.StatusOK, ResponseMap["需要提供单位编号"])
		return
	}

	// company, err := company_manager.GetCompanyWithID(companyID)
	// if err != nil {
	// 	c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
	// 	return
	// }

	categories, err := category_manager.GetCategoriesWithCompany(companyID)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	for _, cat := range categories {
		if len(cat.Image) <= 0 {
			cat.Image = "default_cover.png"
		}
	}

	c.JSON(http.StatusOK, NewResponse(0, "", categories))
}

func AddCategory(c *gin.Context) {
	type categoryInfo struct {
		Company      string `form:"company" json:"company" binding:"required"`
		Name         string `form:"name" json:"name" binding:"required"`
		Image        string `form:"image" json:"image"`
		CategoryType int    `form:" categoryType" json:"categoryType"`
		Index        int    `form:" index" json:"index"`
	}
	var err error
	var ci categoryInfo

	err = c.BindJSON(&ci)
	if err != nil {
		log.InfoF("AddCategory bindjson error: %s", err.Error())
		c.JSON(http.StatusOK, ResponseMap["需要提供用户手机或者邮箱、名称和所在支部编号"])
		return
	}

	cm, err := company_manager.GetCompanyWithID(ci.Company)
	if err != nil {
		if err == models.Err_mongodb_error {
			c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
			return
		}

		if err == models.Err_not_found {
			c.JSON(http.StatusOK, ResponseMap["该单位不存在或者单位编号错误"])
			return
		}
	}

	cat := models.NewCategory(cm.ID, ci.Name, ci.Image, cm.Companies, ci.Index, ci.CategoryType)
	err = category_manager.Add(cat)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	c.JSON(http.StatusOK, ResponseMap["操作成功"])
}

func RemoveCategory(c *gin.Context) {
	id := c.Param("id")

	if len(id) <= 0 {
		c.JSON(http.StatusOK, ResponseMap["需要提供单位编号"])
		return
	}

	err := category_manager.Remove(id)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}
	c.JSON(http.StatusOK, ResponseMap["操作成功"])
}
