package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
	"web/config"
)

var (

	// 全局 db 模式
	db *gorm.DB

	rdb *redis.Client
)

func DataBaseStart() error {

	//创建session
	dsn := config.MysqlConfig.User + ":" + config.MysqlConfig.Password + "@tcp(" + config.MysqlConfig.Host + ")/" +
		config.MysqlConfig.Database +
		"?timeout=3000ms&readTimeout=5000ms&writeTimeout=5000ms&charset=utf8mb4&parseTime=true&loc=Local"
	dbSession, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(config.GormLogFile, "\r\n", log.LstdFlags), // io.Writer
			logger.Config{
				SlowThreshold: time.Second, // 慢查询阈值
				LogLevel:      logger.Info, // 日志级别
				Colorful:      false,       // 禁用彩色打印
			},
		),
		PrepareStmt: true, //使用软删除
	})
	if err != nil {
		return err
	}

	//测试连接,由于设置了超时时间，所以不用测试连接

	db = dbSession

	err = migration()
	if err != nil {
		return err
	}

	return nil
}

func MemoryStart() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RedisConfig.Timeout)*time.Millisecond)
	defer cancel()

	redisDb := redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Host,     // Redis 地址
		Password: config.RedisConfig.Password, // 密码
		DB:       config.RedisConfig.Database, // 使用默认 DB
	})

	//只是声明了客户端，需要测试客户端和服务器的连接
	_, err := redisDb.Ping(ctx).Result()
	if err != nil {
		return err
	}

	rdb = redisDb
	return nil
}
