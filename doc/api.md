# API

##  /users/login

登陆

参数

```go
type LoginReq struct {
	Name     string
	Password string
}
```

返回

```json
{
    "State": "success",
    "Data": ""
}
{
    "State": "not_valid",
    "Data": "xxx@xx.com"
}
{
    "State": "error",
    "Data": "error_pass"
}
```

## /users/register

注册

参数

```go
type RegisterReq struct {
	Name     string
	Email    string
	Password string
}
```

## /users/email

获取验证码

无参数

## /users/valid

验证邮箱

参数

```go
type ValidReq struct {
	VCode string
}
```

## /users/logout

退出登陆

无参数