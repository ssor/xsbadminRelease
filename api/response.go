package api

var (
	ResponseMap map[string]*Response
)

func init() {
	ResponseMap = make(map[string]*Response)

	ResponseMap["操作成功"] = NewResponse(0, "", nil)
	ResponseMap["网络错误，请重试"] = NewResponse(1, "网络错误，请重试", nil)
	ResponseMap["需要提供用户账号、新密码和设备唯一码"] = NewResponse(2, "需要提供用户账号、新密码和设备唯一码", nil)
	ResponseMap["用户不存在"] = NewResponse(3, "用户不存在", nil)
	ResponseMap["密码错误"] = NewResponse(4, "密码错误", nil)
	ResponseMap["用户尚未激活"] = NewResponse(5, "用户尚未激活", nil)
	ResponseMap["设备尚未绑定"] = NewResponse(6, "设备尚未与账号绑定，请先激活该设备", nil)
	ResponseMap["用户已存在"] = NewResponse(7, "用户不存在", nil)
	ResponseMap["需要提供用户账号和新密码"] = NewResponse(8, "需要提供用户账号和新密码", nil)
	ResponseMap["需要提供用户账号"] = NewResponse(9, "需要提供用户账号", nil)
	ResponseMap["删除用户出错"] = NewResponse(10, "删除用户出错", nil)
	ResponseMap["需要用户账号、密码、短信验证码和手机deviceID"] = NewResponse(11, "需要用户账号、密码、短信验证码和手机deviceID", nil)
	ResponseMap["验证码错误"] = NewResponse(12, "验证码错误", nil)
	ResponseMap["设备已经被绑定到其它账号"] = NewResponse(13, "设备已经被绑定到其它账号", nil)
	ResponseMap["绑定的设备数量超出限制"] = NewResponse(14, "绑定的设备数量超出限制", nil)
	ResponseMap["重置密码出错，请重试"] = NewResponse(15, "重置密码出错，请重试", nil)
	ResponseMap["需要提供用户手机或者邮箱、名称和所在支部编号"] = NewResponse(16, "需要提供用户手机或者邮箱、名称和所在支部编号", nil)
	ResponseMap["Email已经被注册"] = NewResponse(17, "Email已经被注册", nil)
	ResponseMap["手机已经被注册"] = NewResponse(18, "手机已经被注册", nil)
	ResponseMap["缺少支部参数信息"] = NewResponse(19, "缺少支部参数信息", nil)
	ResponseMap["该单位不存在或者单位编号错误"] = NewResponse(20, "该单位不存在或者单位编号错误", nil)
	ResponseMap["需要提供单位编号"] = NewResponse(21, "需要提供单位编号", nil)
	ResponseMap["需要提供单位名称"] = NewResponse(22, "需要提供单位名称", nil)
	ResponseMap["需要提供支部名称和所在单位编号"] = NewResponse(23, "需要提供支部名称和所在单位编号", nil)
	ResponseMap["需要支部和单位编号"] = NewResponse(24, "需要支部和单位编号", nil)
	ResponseMap["支部不属于该单位"] = NewResponse(25, "支部不属于该单位", nil)
	ResponseMap["支部编号错误"] = NewResponse(26, "支部编号错误", nil)
	ResponseMap["上传出错"] = NewResponse(27, "上传出错", nil)
	ResponseMap["需要提供图书服务器编号和图书名称"] = NewResponse(28, "需要提供图书服务器编号和图书名称", nil)
	ResponseMap["需要提供栏目和图书编号"] = NewResponse(29, "需要提供栏目和图书编号", nil)
	ResponseMap["图书不存在"] = NewResponse(30, "图书不存在", nil)
}

func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
