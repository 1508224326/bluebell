package settingts

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Host          string `mapstructure:"host"`
	Name          string `mapstructure:"name"`
	Mode          string `mapstructure:"mode"`
	Version       string `mapstructure:"version"`
	StartTime     string `mapstructure:"start_time"`
	Port          int    `mapstructure:"port"`
	MachineID     int64  `mapstructure:"machine_id"`
	*LogConfig    `mapstructure:"log"`
	*MysqlConfig  `mapstructure:"mysql"`
	*RedisConfig  `mapstructure:"redis"`
	*TokenConfig  `mapstructure:"token"`
	*ExpireConfig `mapstructure:"expire"`
}

type LogConfig struct {
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
}

type MysqlConfig struct {
	Port        int    `mapstructure:"port"`
	MaxOpenConn int    `mapstructure:"max_open_connection"`
	MaxIdleConn int    `mapstructure:"max_idle_connection"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DBName      string `mapstructure:"dbname"`
	Host        string `mapstructure:"host"`
}

type RedisConfig struct {
	PoolSize int    `mapstructure:"pool_size"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
}

type TokenConfig struct {
	Duration int64 `mapstructure:"duration"` // token过期时间
}

type ExpireConfig struct {
	PostExpire  int64 `mapstructure:"post_expire"`
	VotedExpire int64 `mapstructure:"voted_expire"`
}

func InitV2() (err error) {
	viper.SetConfigFile("config.yaml")
	//viper.SetConfigName("config") // 1. 设置配置文件名字
	//viper.SetConfigType("yaml")   // 2. 设置文件类型
	viper.AddConfigPath(".") // 3. 配置文件路径
	// 4. 读取加载
	if err = viper.ReadInConfig(); err != nil {
		fmt.Println("配置文件加载失败")
		return
	}
	// 5. 将读取到的信息反序列化到结构体中
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Println("配置文件加载失败")
		return
	}

	// 设置监听
	viper.WatchConfig()

	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件已被修改")
		// 当配置文件被修改时，将配置重新载入到conf中
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Println("配置文件加载失败")
			return
		}
	})
	return

}
