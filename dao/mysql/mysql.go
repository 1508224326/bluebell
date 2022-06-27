package mysql

import (
	"fmt"
	"github.com/aiyouyo/bluebell/models"
	"github.com/aiyouyo/bluebell/settingts"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var db *sqlx.DB

func InitDB() (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"))
	if db, err = sqlx.Connect("mysql", dsn); err != nil {
		fmt.Println("数据库初始化失败", err)
		zap.L().Error("db connect failed, err: ", zap.Error(err))
		return
	}

	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_connection"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_connection"))

	return
}

func InitDBv2(cfg *settingts.MysqlConfig) (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName)
	if db, err = sqlx.Connect("mysql", dsn); err != nil {
		fmt.Println("数据库初始化失败", err)
		zap.L().Error("db connect failed, err: ", zap.Error(err))
		return
	}

	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)

	return
}

func Close() {
	_ = db.Close()
}

// GetAllPostList 查询所有帖子同步到Redis中
func GetAllPostList() (posts []*models.PostBase, err error) {
	sqlStr := `select post_id, community_id, create_time from post`
	err = db.Select(&posts, sqlStr)
	return
}
