# API

## 目录

* [API](#api)
  * [User](#user)
    * [POST /users/login](#post-userslogin)
    * [POST /users/register](#post-usersregister)
    * [POST /user/email](#post-useremail)
    * [POST /user/valid](#post-uservalid)
    * [POST /users/logout](#post-userslogout)
    * [POST /users/Info](#post-usersinfo)
    * [GET /user/info/\{userID\}](#get-userinfouserid)
  * [Room](#room)
    * [GET /room/list/\{page\}](#get-roomlistpage)
    * [GET /room/detail/\{roomID\}](#get-roomdetailroomid)
    * [POST /room/new](#post-roomnew)
    * [POST /room/join/\{roomId\}](#post-roomjoinroomid)
    * [POST /room/ready/\{true/false\}](#post-roomreadytruefalse)
    * [POST /room/team/\{teamID\}](#post-roomteamteamid)
    * [POST /room/role/\{roleName\}](#post-roomrolerolename)
    * [POST /room/quit](#post-roomquit)
    * [POST /room/info](#post-roominfo)
    * [POST /room/own/\{userId\}](#post-roomownuserid)
    * [POST /room/out/\{userId\}](#post-roomoutuserid)
    * [POST /room/start](#post-roomstart)

## User

### POST /users/login

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
    "status": "success",
    "msg": "NikeName"
}
// 新用户
{
    "status": "success",
    "msg": "new_user"
}
// 未激活用户
{
    "status": "not_valid",
    "msg": "xxx@xx.com"
}
// 密码错误
{
    "status": "error",
    "msg": "error_pass"
}
```

### POST /users/register

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
    "status": "success",
    "Data": ""
}
// 邮箱无效 /^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/
{
    "status": "invalid_email",
    "Data": ""
}
// 用户名无效 [a-zA-Z][a-zA-Z0-9_]{0,31}
{
    "status": "invalid_name",
    "Data": ""
}
// 密码无效 >512
{
    "status": "invalid_password",
    "Data": ""
}
// 邮箱已存在
{
    "status": "exist_email",
    "Data": ""
}
// 用户名已存在
{
    "status": "exist_name",
    "Data": ""
}
// 用户名为保留字
{
    "status": "reserved_name",
    "Data": ""
}
```

### POST /user/email

获取验证码

无参数

返回

```json
{
    "status": "success",
    "Data": ""
}
{
    "status": "not_login",
    "Data": ""
}
```

### POST /user/valid

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
    "status": "success",
    "Data": ""
}
{
    "status": "not_login",
    "Data": ""
}
```

### POST /users/logout

退出登陆

无参数

返回

```json
{
    "status": "success",
    "Data": ""
}
```

### POST /users/Info

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
// 数据不合法
{
    "status": "bad_req",
    "msg": ""
}
```

### GET /user/info/{userID} 

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



## Room

以下API均需要登陆状态

### GET /room/list/{page}

获取房间列表（每页10个）page:1~10

需要登陆状态

返回：

```go
type RoomsRes struct {
	Status string              `json:"status"`
	Count  int                 `json:"count"` // 总数量
	Rooms  []services.GameRoom `json:"rooms"`
}
// GameRoom 房间数据
type GameRoom struct {
	ID        int      `json:"id"`        // 房间 ID
	OwnID     string   `json:"ownId"`     // 房主ID
	Port      int      `json:"port"`      // 房间服务器端口
	Title     string   `json:"title"`     // 标题
	GameMap   string   `json:"gameMap"`   // 游戏地图
	MaxPlayer int      `json:"maxPlayer"` // 最大人数
	Mode      string   `json:"mode"`      // 游戏模式
	Password  string   `json:"password"`  // 房间密码（如果有密码则为"password",没有就为""）
	Playing   bool     `json:"playing"`   // 是否正在玩
	Players   []Player `json:"players"`   // 玩家数据
}
// 成功
status = "success"
// 空列表
status = "null"
// 非法请求
status = "bad_req"
```



### GET /room/detail/{roomID} 

获取单个房间详情

需要已经加入了房间中

返回：

```go
type RoomRes struct {
	Status     string                `json:"status"`
	Room       services.GameRoom     `json:"room"`
	PlayerInfo []services.PlayerInfo `json:"players"`
}
// GameRoom 房间数据
type GameRoom struct {
	ID        int      `json:"id"`        // 房间 ID
	OwnID     string   `json:"ownId"`     // 房主ID
	Port      int      `json:"port"`      // 房间服务器端口
	Title     string   `json:"title"`     // 标题
	GameMap   string   `json:"gameMap"`   // 游戏地图
	MaxPlayer int      `json:"maxPlayer"` // 最大人数
	Mode      string   `json:"mode"`      // 游戏模式
	Password  string   `json:"password"`  // 房间密码（如果有密码则为"password",没有就为""）
	Playing   bool     `json:"playing"`   // 是否正在玩
	Players   []Player `json:"players"`   // 玩家数据
}
// PlayerInfo 玩家个性信息
type PlayerInfo struct {
	Player Player       `json:"player"`
	Info   UserBaseInfo `json:"info"`
}
// Player 玩家信息
type Player struct {
	UserID  string `json:"userId"`  // 玩家ID
	GameID  int    `json:"gameId"`  // 游戏内ID
	RoleID  string `json:"roleId"`  // 角色ID
	IsReady bool   `json:"isReady"` // 是否准备
	Team    int    `json:"team"`    // "1-4" - 队伍一~四
}
// UserBaseInfo 用户个性信息
type UserBaseInfo struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Gender int    `json:"gender"`
}
// 成功
status = "success"
// 没有登陆
status = "not_login"
// 非法请求
status = "bad_req"
// 房间不存在
status = "not_found"
// 玩家不在房间内
status = "not_allow"
```



###  POST /room/new 

新建并加入房间

参数：

```go
type reqNewRoom struct {
	Title     string `json:"title"`
	Password  string `json:"password"` // （没有为""）
	GameMap   string `json:"gameMap"`
	GameMode  string `json:"gameMode"`
	MaxPlayer int    `josn:"maxPlayer"`
}
```

返回值：

```json
{
    "status": "success",
    "msg": ""
}
// 成功
status = "success"
// 没有登陆
status = "not_login"
// 非法请求
status = "bad_req"
```



### POST /room/join/{roomId}

加入房间

参数：

```go
type reqJoinRoom struct {
	Password string `json:"password"`
}
```

返回值：

```json
{
    "status": "success",
    "msg": ""
}
// 成功
status = "success"
// 没有登陆
status = "not_login"
// 非法请求
status = "bad_req"
// 找不到房间
status = "not_found"
// 密码错误
status = "err_password"
// 房间玩家以达上限
status = "max_player"
```



### POST /room/ready/{true/false} 

设置准备状态

返回值

```json
{
    "status": "success",
    "msg": ""
}
// 成功
status = "success"
// 没有登陆
status = "not_login"
// 非法请求
status = "bad_req"
// 用户不在房间里面
status = "not_found"
```



### POST /room/team/{teamID} 

设置队伍

返回值：

```json
{
    "status": "success",
    "msg": ""
}
// 成功
status = "success"
// 没有登陆
status = "not_login"
// 非法请求
status = "bad_req"
// 用户不在房间里面
status = "not_found"
// 不允许的队伍/或者队伍已满
status = "not_allow"
```



### POST /room/role/{roleName} 

设置角色

返回值：

```json
{
    "status": "success",
    "msg": ""
}
// 成功
status = "success"
// 没有登陆
status = "not_login"
// 非法请求
status = "bad_req"
// 用户不在房间里面
status = "not_found"
```



### POST /room/quit 

退出房间

返回值：

```json
{
    "status": "success",
    "msg": ""
}
// 成功
status = "success"
// 没有登陆
status = "not_login"
// 用户不在房间里面
status = "not_found"
```



### POST /room/info 

设置房间信息

参数：

```go
type roomInfoReq struct {
	MaxPlayer int    `json:"maxPlayer"`// 参数为0则不更改
	GameMap   string `json:"gameMap"` // 参数为“”则不更改
	GameMode  string `json:"gameMode"`// 参数为“”则不更改
}
```

返回值：

```json
{
    "status": "success",
    "msg": ""
}
// 成功
status = "success"
// 没有登陆
status = "not_login"
// 用户不在房间里面
status = "not_found"
// 非法请求
status = "bad_req"
// 你不是房主
status = "not_allow"
```



### POST /room/own/{userId} 

设置房主

返回值：

```json
{
    "status": "success",
    "msg": ""
}
// 成功
status = "success"
// 没有登陆
status = "not_login"
// 用户不在房间里面
status = "not_found"
// 非法请求
status = "bad_req"
// 你不是房主
status = "not_allow"
```



### POST /room/out/{userId} 

踢人

返回值：

```json
{
    "status": "success",
    "msg": ""
}
// 成功
status = "success"
// 没有登陆
status = "not_login"
// 用户不在房间里面
status = "not_found"
// 非法请求
status = "bad_req"
// 你不是房主
status = "not_allow"
```



### POST /room/start 

开始游戏

返回值：

```json
{
    "status": "success",
    "msg": ""
}
// 成功
status = "success"
// 没有登陆
status = "not_login"
// 用户不在房间里面
status = "not_found"
// 你不是房主
status = "not_allow"
// 有玩家还没装准备
status = "not_ready"
// 游戏服务器已满
status = "max_server"
// 非合作模式下所有玩家都是同一队的
status = "one_team"
```

