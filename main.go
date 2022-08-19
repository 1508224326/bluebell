package main

import (
	"context"
	"fmt"
	"github.com/aiyouyo/bluebell/controller"
	"github.com/aiyouyo/bluebell/dao/mysql"
	"github.com/aiyouyo/bluebell/dao/redis"
	"github.com/aiyouyo/bluebell/logger"
	"github.com/aiyouyo/bluebell/pkg/snowflake"
	"github.com/aiyouyo/bluebell/routes"
	"github.com/aiyouyo/bluebell/settingts"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"go.uber.org/zap"
)

//var srv *http.Server

// @title bluebell项目
// @version 1.0
// @description 帖子发表展示 投票等功能
// @termsOfService http://swagger.io/terms/

// @contact.name yongdeng
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:9000
// @BasePath 127.0.0.1:9000/api/v1
func main() {

	// 1. 初始化配置
	if err := settingts.InitV2(); err != nil {
		fmt.Println("init settings failed, err: ", err)
		return
	}

	// 2. 初始化日志
	if err := logger.InitV2(settingts.Conf.LogConfig, settingts.Conf.Mode); err != nil {
		fmt.Println("logger init failed, err: ", err)
		return
	}
	defer zap.L().Sync()

	// 3. 初始化 mysql 连接
	if err := mysql.InitDBv2(settingts.Conf.MysqlConfig); err != nil {
		zap.L().Error(settingts.Conf.MysqlConfig.Host)
		zap.L().Error("mysql connection failed", zap.Error(err))
		return
	}
	defer mysql.Close()

	//4. 初始化redis
	if err := redis.InitV2(settingts.Conf.RedisConfig); err != nil {
		fmt.Println("redis connection failed, err: ", err)
		zap.L().Error("redis connection failed, err", zap.Error(err))
		return
	}
	defer redis.CLose()

	//同步mysql数据到redis
	//list, err := mysql.GetAllPostList()
	//if err != nil {
	//	zap.L().Error("Mysql同步数据出错", zap.Error(err))
	//	fmt.Println("mysql.GetAllPostList() err : ", err)
	//	return
	//}
	//err = redis.SetAllPosts(list)
	//if err != nil {
	//	zap.L().Error("Redis同步数据出错", zap.Error(err))
	//	fmt.Println("redis.SetAllPosts(list) err : ", err)
	//	return
	//}

	// 初始化 ID 生成器
	if err := snowflake.Init(settingts.Conf.StartTime, settingts.Conf.MachineID); err != nil {
		fmt.Println("snowflake init failed, err: ", err)
		return
	}

	// 初始化gin框架的翻译器 翻译错误信息
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Println("controller.InitTrans init failed, err: ", err)
		return
	}

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	// 5. 注册路由
	r := routes.Setup()

	// 6. 启动服务（优雅关机）
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d",
			settingts.Conf.Host,
			settingts.Conf.Port,
		),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen err", zap.Error(err))
		}
		fmt.Println("========================================================")
	}()

	quit := make(chan os.Signal, 1)

	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 会发送 syscall.SIGINT 信号，我们常用的 Ctrl + C 就是触发的这个信号
	// kill -9 会发送 syscall.SIGKILL 信号，但不能被捕获到

	// signal.Notify 把收到的 syscall.SIGINT 或 syscall.SIGTERM 信号转给 quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server")

	// 创建一个5秒超时的context
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFunc()
	// 5 五秒内优雅的关机
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown, err: ", err)
	}
	zap.L().Info("Server exiting")

}
