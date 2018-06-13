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
)

type HttpConfig struct {
	Host          string `yaml:"Host"`          // 服务器监听地址
	Port          string `yaml:"Port"`          // 服务器监听端口
	Dev           bool   `yaml:"Dev"`           // 是否开发环境
	PortPoolBegin int64  `yaml:"PortPoolBegin"` // 游戏服务器地址池开始
	PortPoolSize  int64  `yaml:"PortPoolSize"`  // 最大负载
}

// Config 配置文件
type Config struct {
	Mongo      models.Mongo     `yaml:"Mongo"`  // mongoDB配置
	HttpServer HttpConfig       `yaml:"Server"` // iris配置
	Violet     violetSdk.Config `yaml:"Violet"` // Violet配置
}

func RunServer(c Config) {
	// 初始化数据库
	Model, err := models.NewModel(c.Mongo)
	if err != nil {
		panic(err)
	}
	// 初始化服务
	Service := services.NewService(Model)
	userService := Service.NewUserService()
	userService.InitViolet(c.Violet)

	// 启动服务器
	app := iris.New()
	if c.HttpServer.Dev {
		app.Logger().SetLevel("debug")
	}

	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessionBug",
		Expires: 24 * time.Hour,
	})
	// "/users" based mvc application.
	users := mvc.New(app.Party("/users"))
	// Bind the "userService" to the UserController's Service (interface) field.
	users.Register(userService, sessManager.Start)
	users.Handle(new(controllers.UsersController))

	app.Run(
		// Starts the web server
		iris.Addr(c.HttpServer.Host+":"+c.HttpServer.Port),
		// Disables the updater.
		iris.WithoutVersionChecker,
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
	)
}
