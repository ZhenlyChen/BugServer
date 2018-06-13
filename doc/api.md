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
// 老用户登陆成功
{
    "State": "success",
    "Data": "NikeName"
}
// 新用户
{
    "State": "success",
    "Data": "new_user"
}
// 未激活用户
{
    "State": "not_valid",
    "Data": "xxx@xx.com"
}
// 密码错误
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

返回
```json
{
    "State": "success"
    "Data": ""
}
{
    "State": "error"
    "Data": "exist_email"
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

返回
```json
{
    "State": "success"
    "Data": ""
}
{
    "State": "error"
    "Data": "not_login"
}
```

## /users/logout

退出登陆

无参数

返回
```json
{
    "State": "success"
    "Data": ""
}
```

## /users/userName

设置用户昵称

参数

```go
type Req struct {
	Name string
}
```

返回
```json
{
    "State": "success"
    "Data": ""
}
```

## GET /users/userBaseInfo

获取用户状态

返回
```go
type UserRes struct {
	State string
	NikeName string
	Avatar string
	Gender int
	Level int
}
```
state = "success" / "not_login" / "error"