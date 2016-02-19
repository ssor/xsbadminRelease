#服务器地址
www.iot-top.com:8081/


#用户激活

POST api/v1/user/active

BOYD example:

{"account": "133xxxxx","code":"123456789qwertyuiop", "password": "xxxxxx", "device": "xxxxxxx"}

其中 account 表示用户账号，可能为手机号或者邮箱账号；device 为手机设备唯一码；code 是验证码，暂时使用上面的默认验证码；password 表示新设置的用户密码


#用户登录

POST api/v1/user/login

BOYD example:

{"account": "133xxxxx", "password": "xxxxxx", "device": "xxxxxxx"}

其中 account 表示用户账号，可能为手机号或者邮箱账号；device 为手机设备唯一码


