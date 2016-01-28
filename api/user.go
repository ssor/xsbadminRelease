package api

import (
	"github.com/gin-gonic/gin"
	// "io"
	"net/http"
	log "xsbAdmin/libs/log"
	"xsbAdmin/libs/tools"
	"xsbAdmin/models"
)

var (
	verification_code_test = "123456789qwertyuiop"
)

type loginInfo struct {
	Account  string `form:"account" json:"account" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Device   string `form:"device" json:"device" binding:"required"`
}

// func newResponse(code int, message string, data interface{}) gin.H {
// 	return gin.H{"code": code, "message": message, "data": data}
// }

func AllUsers(c *gin.Context) {
	group := c.Param("group")
	ul, err := user_manager.All(group)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}
	for _, u := range ul {
		if len(u.Image) <= 0 {
			u.Image = default_avatar
		}
	}
	c.JSON(http.StatusOK, NewResponse(0, "", ul))
}

func Login(c *gin.Context) {
	code, h := validateLogin(c)
	c.JSON(code, h)
}

func validateLogin(c *gin.Context) (int, *Response) {
	var li loginInfo

	err := c.BindJSON(&li)
	if err != nil {
		return http.StatusOK, ResponseMap["需要提供用户账号、新密码和设备唯一码"]
	}

	userInfo, err := user_manager.GetAccount(li.Account)
	if err == models.Err_mongodb_error {
		return http.StatusOK, ResponseMap["网络错误，请重试"]
	}

	if err == models.Err_not_found {
		return http.StatusOK, ResponseMap["用户不存在"]
	}

	if userInfo.Password != li.Password {
		return http.StatusOK, ResponseMap["密码错误"]
	}

	if userInfo.Actived != models.Actived {
		return http.StatusOK, ResponseMap["用户尚未激活"]
	}

	// //IMEI check
	// imeis, err := imei_manager.GetAll(li.Account)
	// if err == models.Err_mongodb_error {
	// 	return http.StatusOK, ResponseMap["网络错误，请重试"]
	// }

	// if err == models.Err_not_found || (imeis != nil && len(imeis) <= 0) {
	// 	return http.StatusOK, ResponseMap["用户尚未激活"]
	// }

	// if imeis.Contains(li.IMEI) == false {
	// 	return http.StatusOK, ResponseMap["设备尚未绑定"]
	// }

	b, res := checkIfImeiBoundUser(li.Account, li.Device)
	if b == false {
		return http.StatusOK, res
	}
	return http.StatusOK, NewResponse(0, "登陆成功", userInfo.ToAccount())
}

func checkIfImeiBoundUser(userID, imei string) (bool, *Response) {
	//IMEI check
	imeis, err := imei_manager.GetAll(userID)
	if err == models.Err_mongodb_error {
		return false, ResponseMap["网络错误，请重试"]
	}

	if err == models.Err_not_found || (imeis != nil && len(imeis) <= 0) {
		return false, ResponseMap["用户尚未激活"]
	}

	if imeis.Contains(imei) == false {
		return false, ResponseMap["设备尚未绑定"]
	}
	return true, nil
}

func ResetPwdByVerification(c *gin.Context) {
	type resetInfo struct {
		// UID      string `form:"uid" json:"uid" binding:"required"`
		Account  string `form:"account" json:"account" binding:"required"`
		Code     string `form:"code" json:"code" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
		Device   string `form:"device" json:"device" binding:"required"`
	}

	var ri resetInfo

	err := c.BindJSON(&ri)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["需要用户账号、密码、短信验证码和手机deviceID"])
		return
	}

	if ri.Code != verification_code_test {
		c.JSON(http.StatusOK, ResponseMap["验证码错误"])
		return
	}

	userInfo, err := user_manager.GetAccount(ri.Account)
	if err == models.Err_mongodb_error {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	if err == models.Err_not_found {
		c.JSON(http.StatusOK, ResponseMap["用户不存在"])
		return
	}

	//only if imei is bound with user
	b, res := checkIfImeiBoundUser(ri.Account, ri.Device)
	if b == false {
		c.JSON(http.StatusOK, res)
		return
	}

	err = user_manager.SetNewPwdDirectly(ri.Account, ri.Password)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["重置密码出错，请重试"])
		return
	}
	userInfo, err = user_manager.GetAccount(ri.Account)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}
	c.JSON(http.StatusOK, NewResponse(0, "重置密码成功", userInfo.ToAccount()))
}

func AddUser(c *gin.Context) {
	var err error

	type userInfo struct {
		Phone string `form:"phone" json:"phone"`
		Email string `form:"email" json:"email"`
		Name  string `form:"name" json:"name" binding:"required"`
		Group string `form:"group" json:"group" binding:"required"`
		Image string `form:"image" json:"image"`
		Index int    `form:"index" json:"index"`
	}

	var ui userInfo

	// body := c.Request.Body
	// bytes := make([]byte, 1024)
	// _, err = body.Read(bytes)
	// defer body.Close()
	// if err != nil {
	// 	if err != io.EOF {
	// 		log.InfoF("addUser request error: %s", err.Error())
	// 		c.String(http.StatusBadRequest, "")
	// 		return
	// 	}
	// }
	// log.TraceF("addUser <- [%s]", string(bytes))

	err = c.BindJSON(&ui)
	if err != nil {
		log.InfoF("addUser bindjson error: %s", err.Error())
		c.JSON(http.StatusOK, ResponseMap["需要提供用户手机或者邮箱、名称和所在支部编号"])
		return
	}

	log.TraceF("Phone: %s Email: %s Name: %s Group: %s Image: %s Index: %d",
		ui.Phone, ui.Email, ui.Name, ui.Group, ui.Image, ui.Index)

	if len(ui.Phone) <= 0 {
		if tools.IsEmail(ui.Email) == false {
			c.JSON(http.StatusOK, ResponseMap["需要提供用户手机或者邮箱、名称和所在支部编号"])
			return
		}
	}

	log.TraceF("add user --> Name: %s Phone: %s Email: %s Group: %s", ui.Name, ui.Phone, ui.Email, ui.Group)

	company, err := company_manager.GetCompanyWithGroup(ui.Group)
	if err != nil {
		if err == models.Err_mongodb_error {
			c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
			return
		}
	}

	if company == nil {
		c.JSON(http.StatusOK, ResponseMap["支部编号错误"])
		return
	}

	err = user_manager.Add(ui.Phone, ui.Email, default_password, ui.Name, company.ID, ui.Group, ui.Image, ui.Index)
	if err != nil {
		log.TraceF("add user err: %s", err)
		var res *Response
		switch err {
		case models.Err_email_used:
			res = ResponseMap["Email已经被注册"]
		case models.Err_phone_used:
			res = ResponseMap["手机已经被注册"]
		default:
			res = ResponseMap["网络错误，请重试"]
		}
		c.JSON(http.StatusOK, res)
		return
	}
	log.TraceF("add user ok")
	c.JSON(http.StatusOK, NewResponse(0, "添加用户成功", nil))
}

func RemoveUser(c *gin.Context) {
	// type userInfo struct {
	// 	UID string `form:"uid" json:"uid" binding:"required"`
	// }
	// var ui userInfo

	// err := c.BindJSON(&ui)
	// if err != nil {
	// 	c.JSON(http.StatusOK, ResponseMap["需要提供用户账号"])
	// 	return
	// }
	var err error
	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(http.StatusOK, ResponseMap["需要提供用户账号"])
		return
	}

	err = user_manager.RemoveId(id)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["删除用户出错"])
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "删除用户成功", nil))
}

func ActiveUser(c *gin.Context) {

	type activeInfo struct {
		Account  string `form:"account" json:"account" binding:"required"`
		Code     string `form:"code" json:"code" binding:"required"`
		Device   string `form:"device" json:"device" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	var ai activeInfo
	var imeis models.IMEI_List

	err := c.BindJSON(&ai)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["需要用户账号、密码、短信验证码和手机deviceID"])
		return
	}

	if ai.Code != verification_code_test {
		c.JSON(http.StatusOK, ResponseMap["验证码错误"])
		return
	}

	userInfo, err := user_manager.GetAccount(ai.Account)
	if err != nil {
		if err == models.Err_not_found {
			c.JSON(http.StatusOK, ResponseMap["用户不存在"])
			return
		} else {
			c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
			return
		}
	}

	//1. the imei should not have been used; if used, it must be the user of the id
	imei, err := imei_manager.Exits(ai.Device)
	if err != nil {
		if err == models.Err_mongodb_error {
			c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
			return
		}
	}

	if imei != nil { // the device has been used
		if imei.BoundWith(ai.Account) == false {
			// if imei.UserID != ai.ID {
			c.JSON(http.StatusOK, ResponseMap["设备已经被绑定到其它账号"])
			return
		} else {
			log.TraceF("device [%s] already bound to user [%s]", ai.Device, ai.Account)
			// c.JSON(http.StatusOK, nil)
			// return
			goto boundOK
		}
	} else { // not found, means that the IMEI not bound to user yet

	}

	//2. the user with id should have less than 4 imeis
	imeis, err = imei_manager.GetAll(ai.Account)
	if err != nil {
		if err == models.Err_mongodb_error {
			c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
			return
		}
	}

	if imeis != nil && len(imeis) >= max_bound_imei {
		c.JSON(http.StatusOK, ResponseMap["绑定的设备数量超出限制"])
		return
	}

	err = user_manager.Active(userInfo, ai.Device, ai.Password)
	if err != nil {
		c.JSON(http.StatusOK, ResponseMap["网络错误，请重试"])
		return
	}

	log.TraceF("device [%s]  bound to user [%s] OK", ai.Device, ai.Account)

boundOK:
	c.JSON(http.StatusOK, NewResponse(0, "激活成功", userInfo.ToAccount()))
}

// func ResetPwd(c *gin.Context) {

// 	type resetInfo struct {
// 		ID          string `form:"id" json:"id" binding:"required"`
// 		OldPassword string `form:"oldpassword" json:"oldpassword" binding:"required"`
// 		NewPassword string `form:"newpassword" json:"newpassword" binding:"required"`
// 	}

// 	var ri resetInfo

// 	err := c.BindJSON(&ri)
// 	if err != nil {
// 		c.String(http.StatusBadRequest, "id or pwd error")
// 		return
// 	}

// 	if ri.OldPassword == ri.NewPassword {
// 		c.String(http.StatusBadRequest, "new password cannot be the same as old")
// 		return
// 	}

// 	if err := user_manager.SetNewPwd(ri.ID, ri.OldPassword, ri.NewPassword); err == nil {
// 		c.String(http.StatusOK, "OK")
// 		return
// 	} else {
// 		c.String(http.StatusBadRequest, "resetpwd error")
// 		return
// 	}
// }
