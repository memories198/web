package user

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	// 全局 db 模式
	db *gorm.DB

	rdb *redis.Client
)

func DataBaseStart() error {
	logFile, err := os.OpenFile("./web.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer logFile.Close()

	//创建session
	dsn := "root:o?dv(%qn)Uf8*e^3AQTPB4k6L1N975xG@tcp(1.94.9.77:1763)/web?timeout=5000ms&readTimeout=5000ms&writeTimeout=5000ms&charset=utf8mb4&parseTime=true&loc=Local"
	dbSession, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(logFile, "\r\n", log.LstdFlags), // io.Writer
			logger.Config{
				SlowThreshold: time.Second, // 慢查询阈值
				LogLevel:      logger.Info, // 日志级别
				Colorful:      false,       // 禁用彩色打印
			},
		),
	})
	if err != nil {
		return err
	}

	//测试连接
	sqlDb, err := dbSession.DB()
	if err != nil {
		return err
	} else {
		err = sqlDb.Ping()
		if err != nil {
			return err
		}
	}

	db = dbSession

	err = migration()
	if err != nil {
		return err
	}

	return nil
}

var ctx = context.Background()

func MemoryStart() error {
	redisDb := redis.NewClient(&redis.Options{
		Addr: "1.94.9.77:1764", // Redis 地址
		//Addr:     "1.94.9.77:6379", // Redis 地址
		Password: "o?dv(%qn)Uf8*e^3AQTPB4k6L1N975xG", // 密码
		DB:       0,                                  // 使用默认 DB
	})
	//测试连接
	_, err := redisDb.Ping(ctx).Result()
	if err != nil {
		return err
	}

	rdb = redisDb

	return nil
}
