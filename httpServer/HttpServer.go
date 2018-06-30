package httpServer

import (
	"time"

	"github.com/XMatrixStudio/Violet.SDK.Go"
	"github.com/ZhenlyChen/BugServer/httpServer/controllers"
	"github.com/ZhenlyChen/BugServer/httpServer/models"
	"github.com/ZhenlyChen/BugServer/httpServer/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"github.com/ZhenlyChen/BugServer/gameServer"
)

type ServerConfig struct {
	Host string                  `yaml:"Host"` // 服务器监听地址
	Port string                  `yaml:"Port"` // 服务器监听端口
	Dev  bool                    `yaml:"Dev"`  // 是否开发环境
	Game gameServer.ServerConfig `yaml:"Game"` // 游戏服务器
}

// Config 配置文件
type Config struct {
	Mongo  models.Mongo     `yaml:"Mongo"`  // mongoDB配置
	Server ServerConfig     `yaml:"Server"` // iris配置
	Violet violetSdk.Config `yaml:"Violet"` // Violet配置
}

func RunServer(c Config) {
	// 初始化数据库
	Model, err := models.NewModel(c.Mongo)
	if err != nil {
		panic(err)
	}
	// 初始化服务
	Service := services.NewService(Model)

	// 启动服务器
	app := iris.New()
	if c.Server.Dev {
		app.Logger().SetLevel("debug")
	}

	sessionManager := sessions.New(sessions.Config{
		Cookie:  "sessionBug",
		Expires: 24 * time.Hour,
	})

	users := mvc.New(app.Party("/user"))
	userService := Service.GetUserService()
	userService.InitViolet(c.Violet)
	users.Register(userService, sessionManager.Start)
	users.Handle(new(controllers.UsersController))

	roomService := Service.GetRoomService()
	roomService.InitGameServer(c.Server.Game)
	roomService.NewRoom()

	app.Run(
		// Starts the web server
		iris.Addr(c.Server.Host+":"+c.Server.Port),
		// Disables the updater.
		iris.WithoutVersionChecker,
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
	)
}
