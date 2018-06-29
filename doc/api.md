# API

## POST /users/login

登陆

参数

```go
type LoginReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
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

## POST /users/register

注册

参数

```go
type RegisterReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
```

返回

```json
{
    "State": "success",
    "Data": ""
}
// 邮箱无效 /^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/
{
    "State": "invalid_email",
    "Data": ""
}
// 用户名无效 [a-zA-Z][a-zA-Z0-9_]{0,31}
{
    "State": "invalid_name",
    "Data": ""
}
// 密码无效 >512
{
    "State": "invalid_password",
    "Data": ""
}
// 邮箱已存在
{
    "State": "exist_email",
    "Data": ""
}
// 用户名已存在
{
    "State": "exist_name",
    "Data": ""
}
// 用户名为保留字
{
    "State": "reserved_name",
    "Data": ""
}
```

## POST /user/email

获取验证码

无参数

返回

```json
{
    "State": "success",
    "Data": ""
}
{
    "State": "not_login",
    "Data": ""
}
```

## POST /user/valid

验证邮箱

参数

```go
type ValidReq struct {
	VCode string `json:"vCode"`
}
```

返回

```json
{
    "State": "success",
    "Data": ""
}
{
    "State": "not_login",
    "Data": ""
}
```

## POST /users/logout

退出登陆

无参数

返回

```json
{
    "State": "success",
    "Data": ""
}
```

## POST /users/Info

设置用户信息

参数

```go
type InfoReq struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Gender int    `json:"gender"`
}
```

返回

```json
{
    "status": "success",
    "msg": ""
}
```

## GET /user/info/{userID} 

获取用户信息，userID为空时候获取自身信息

返回

```go
type UserRes struct {
	Status   string `json:"status"`
	NikeName string `json:"nikeName"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
	Level    int    `json:"level"`
}
// 成功
status = "success"
// 未登录
status = "not_login"
// 非法请求
status = "bad_req"
// 内部错误
status = "error"
```


