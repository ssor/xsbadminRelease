#服务器地址
www.iot-top.com:8081/


#用户激活

POST api/v1/user/active

BOYD example:

{"account": "133xxxxx","code":"123456789qwertyuiop", "password": "xxxxxx", "device": "xxxxxxx"}

其中 account 表示用户账号，可能为手机号或者邮箱账号；device 为手机设备唯一码；code 是验证码，暂时使用上面的默认验证码；password 表示新设置的用户密码

#发送验证码
GET api/v1/user/sendVericationCode/:phone



#用户登录

POST api/v1/user/login

BOYD example:

{"account": "133xxxxx", "password": "xxxxxx", "device": "xxxxxxx"}

其中 account 表示用户账号，可能为手机号或者邮箱账号；device 为手机设备唯一码


#公告

GET api/v1/notice/newestNotices/:company?last=:para

获取最新的公告
company: 单位的 ID
para: 上次请求的最新的公告的 ID，默认值 0


#新闻

GET api/v1/news/newestNews/:company?last=:para

获取最新的新闻
company: 单位的 ID
para: 上次请求的最新的新闻的 ID，默认值 0

#任务

GET api/v1/task/allTasksForCompany/:company?expired=true&limited=true
获取一个单位的所有任务，包括上级布置的任务

当 expired 设置为空，不限制过期与否
当 expired 设置为 true，获取所有过期任务
当 expired 设置为 false，获取所有未过期任务

当 limited 设置为空，不限制限时和非限时任务
当 limited 设置为 true，获取的任务均为限时任务
当 limited 设置为 false，获取的任务均为非限时任务

wc -l

sed -i '1,Nd'