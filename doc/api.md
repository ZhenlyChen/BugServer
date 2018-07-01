# API

## 目录

* [API](#api)
  * [目录](#%E7%9B%AE%E5%BD%95)
  * [User](#user)
    * [POST /users/login 登陆](#post-userslogin-%E7%99%BB%E9%99%86)
    * [POST /users/register 注册](#post-usersregister-%E6%B3%A8%E5%86%8C)
    * [POST /user/email 获取邮箱验证码](#post-useremail-%E8%8E%B7%E5%8F%96%E9%82%AE%E7%AE%B1%E9%AA%8C%E8%AF%81%E7%A0%81)
    * [POST /user/valid 验证邮箱验证码](#post-uservalid-%E9%AA%8C%E8%AF%81%E9%82%AE%E7%AE%B1%E9%AA%8C%E8%AF%81%E7%A0%81)
    * [POST /users/logout 退出登陆](#post-userslogout-%E9%80%80%E5%87%BA%E7%99%BB%E9%99%86)
    * [POST /users/Info 设置用户信息](#post-usersinfo-%E8%AE%BE%E7%BD%AE%E7%94%A8%E6%88%B7%E4%BF%A1%E6%81%AF)
    * [GET /user/info/\{userID\} 获取用户信息](#get-userinfouserid-%E8%8E%B7%E5%8F%96%E7%94%A8%E6%88%B7%E4%BF%A1%E6%81%AF)
  * [Room](#room)
    * [GET /room/list/\{page\} 获取房间列表](#get-roomlistpage-%E8%8E%B7%E5%8F%96%E6%88%BF%E9%97%B4%E5%88%97%E8%A1%A8)
    * [GET /room/detail 获取自己所在房间详情](#get-roomdetail-%E8%8E%B7%E5%8F%96%E8%87%AA%E5%B7%B1%E6%89%80%E5%9C%A8%E6%88%BF%E9%97%B4%E8%AF%A6%E6%83%85)
    * [POST /room/new 新建并加入房间](#post-roomnew-%E6%96%B0%E5%BB%BA%E5%B9%B6%E5%8A%A0%E5%85%A5%E6%88%BF%E9%97%B4)
    * [POST /room/join/\{roomId\} 加入房间](#post-roomjoinroomid-%E5%8A%A0%E5%85%A5%E6%88%BF%E9%97%B4)
    * [POST /room/ready/\{true/false\} 设置准备状态](#post-roomreadytruefalse-%E8%AE%BE%E7%BD%AE%E5%87%86%E5%A4%87%E7%8A%B6%E6%80%81)
    * [POST /room/team/\{teamID\} 设置队伍](#post-roomteamteamid-%E8%AE%BE%E7%BD%AE%E9%98%9F%E4%BC%8D)
    * [POST /room/role/\{roleName\} 设置角色](#post-roomrolerolename-%E8%AE%BE%E7%BD%AE%E8%A7%92%E8%89%B2)
    * [POST /room/quit 退出房间](#post-roomquit-%E9%80%80%E5%87%BA%E6%88%BF%E9%97%B4)
    * [POST /room/info 设置房间信息](#post-roominfo-%E8%AE%BE%E7%BD%AE%E6%88%BF%E9%97%B4%E4%BF%A1%E6%81%AF)
    * [POST /room/own/\{userId\} 设置房主](#post-roomownuserid-%E8%AE%BE%E7%BD%AE%E6%88%BF%E4%B8%BB)
    * [POST /room/out/\{userId\} 踢人](#post-roomoutuserid-%E8%B8%A2%E4%BA%BA)
    * [POST /room/start 开始游戏](#post-roomstart-%E5%BC%80%E5%A7%8B%E6%B8%B8%E6%88%8F)

## User

### POST /users/login 登陆

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

### POST /users/register 注册

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

### POST /user/email 获取邮箱验证码

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

### POST /user/valid 验证邮箱验证码

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

### POST /users/logout 退出登陆

无参数

返回

```json
{
    "status": "success",
    "Data": ""
}
```

### POST /users/Info 设置用户信息

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

### GET /user/info/{userID} 获取用户信息

userID为空时候获取自身信息

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

### GET /room/list/{page} 获取房间列表

（每页10个）page:1~10

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



### GET /room/detail 获取自己所在房间详情

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
// 玩家不在房间内/或者房间不存在
status = "not_found"
```



###  POST /room/new 新建并加入房间

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



### POST /room/join/{roomId} 加入房间

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



### POST /room/ready/{true/false} 设置准备状态

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



### POST /room/team/{teamID} 设置队伍

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



### POST /room/role/{roleName} 设置角色

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



### POST /room/quit 退出房间

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



---

以下为房主专用API

### POST /room/info 设置房间信息

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



### POST /room/own/{userId} 设置房主

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



### POST /room/out/{userId} 踢人

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



### POST /room/start 开始游戏

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



## GameServer

基于帧同步队列

以下为UDP请求

### 请求

#### 加入对局

```go
0
// UserComeIn ...
type UserComeIn struct {
	ID int `json:"id"`
}
```

#### 退出对局

```go
3
```

#### 发送命令

```go
1
type UserData struct {
	ID    int `json:"id"`
	Input int `json:"input"`
}
```

#### 设置当前帧数

```go
2
type UserBack struct {
	ID    int `json:"id"`
	Frame int `json:"frame"`
}
```



### 返回

#### 当前数据

一次性最多返回10帧数据

```go
// ResData ...
type ResData struct {
	Data []FrameState `json:"data"`
}

// FrameState ...
type FrameState struct {
	FrameID  int       `json:"frameID"`
	Commends []Commend `json:"commends"`
}

// Commend ...
type Commend struct {
	UserID int `json:"id"`
	Input  int `json:"input"`
}
```

