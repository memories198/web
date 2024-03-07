package dao

import (
	"context"
	"time"
	"web/config"
)

func MemoryGetKey(k string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RedisConfig.Timeout)*time.Millisecond)
	defer cancel()

	v, err := rdb.Get(ctx, k).Result()
	if err != nil {
		return "", err
	}
	return v, nil
}
func MemorySetKey(k, v string, t int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RedisConfig.Timeout)*time.Millisecond)
	defer cancel()

	err := rdb.Set(ctx, k, v, time.Duration(t)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}
func MemoryDelKey(k string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RedisConfig.Timeout)*time.Millisecond)
	defer cancel()

	err := rdb.Del(ctx, k).Err()
	if err != nil {
		return err
	}
	return nil
}
func MemoryLPush(username, cookie string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RedisConfig.Timeout)*time.Millisecond)
	defer cancel()

	err := rdb.LPush(ctx, "cookie:"+username, cookie).Err()
	if err != nil {
		return err
	}
	return nil
}
func MemoryLPop(username string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RedisConfig.Timeout)*time.Millisecond)
	defer cancel()

	value, err := rdb.LPop(ctx, "cookie:"+username).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
func MemorySetHash(u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RedisConfig.Timeout)*time.Millisecond)
	defer cancel()

	err := rdb.HMSet(ctx, "user:"+u.Username, "Username", u.Username, "Password", u.Password).Err()
	if err != nil {
		return err
	}
	return nil
}

func MemoryGetHash(key string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RedisConfig.Timeout)*time.Millisecond)
	defer cancel()

	user, err := rdb.HGetAll(ctx, "user:"+key).Result()
	if err != nil {
		return nil, err
	}
	return &User{
		Username: user["Username"],
		Password: user["Password"],
	}, nil
}
func MemorySetExpire(key string, expireTime int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RedisConfig.Timeout)*time.Millisecond)
	defer cancel()
	rdb.Expire(ctx, key, time.Duration(expireTime))
	return nil
}
