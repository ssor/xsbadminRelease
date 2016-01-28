package api

import (
	"github.com/gin-gonic/gin"
	// "io"
	"net/http"
	log "xsbAdmin/libs/log"
	// "xsbAdmin/libs/tools"
	"xsbAdmin/models"
)

//get group info with users
func GroupInfo(c *gin.Context) {
	type group struct {
		Info  *models.GroupInfo
		Users []*models.AccountInfo
	}
	var err error
	id := c.Param("id")

	if len(id) <= 0 {
		c.JSON(http.StatusOK, ResponseMap["支部编号错误"])
		return
	}

	g, err := group_manager.GetGroupInfo(id)
	if err != nil {
		if err == models.Err_not_found {
			c.JSON(http.StatusOK, ResponseMap["支部编号错误"])
			return
		}
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	users, err := user_manager.All(id)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	for _, u := range users {
		if len(u.Image) <= 0 {
			u.Image = default_avatar
		}
	}
	c.JSON(http.StatusOK, NewResponse(0, "", &group{g, users.ToAccount()}))
}

//all groups in a company
func Groups(c *gin.Context) {
	id := c.Param("id")

	if len(id) <= 0 {
		c.JSON(http.StatusOK, ResponseMap["需要提供单位编号"])
		return
	}

	list, err := group_manager.All(id)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "", list))
}

func AddGroup(c *gin.Context) {
	type groupInfo struct {
		Name    string `form:"name" json:"name" binding:"required"`
		Company string `form:"company" json:"company" binding:"required"`
	}

	var err error
	var gi groupInfo

	err = c.BindJSON(&gi)
	if err != nil {
		log.InfoF("AddGroup bindjson error: %s", err.Error())
		c.JSON(http.StatusOK, ResponseMap["需要提供支部名称和所在单位编号"])
		return
	}

	company, err := company_manager.GetCompanyWithID(gi.Company)
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

	group := models.NewGroupInfo(gi.Name, company.ID)

	err = company_manager.AddGroup(company, group)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	c.JSON(http.StatusOK, ResponseMap["操作成功"])
}

func RemoveGroup(c *gin.Context) {
	type groupInfo struct {
		Group   string `form:"group" json:"group" binding:"required"`
		Company string `form:"company" json:"company" binding:"required"`
	}

	var err error
	var gi groupInfo

	err = c.BindJSON(&gi)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["需要支部和单位编号"])
		return
	}

	company, err := company_manager.GetCompanyWithID(gi.Company)
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

	if company.HasGroup(gi.Group) == false {
		c.JSON(http.StatusOK, ResponseMap["支部不属于该单位"])
		return
	}

	err = company_manager.RemoveGroup(company, gi.Group)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	c.JSON(http.StatusOK, ResponseMap["操作成功"])
}
