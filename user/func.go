package user

import (
	"errors"
	"strconv"
	"time"
)

func migration() error {
	err := db.AutoMigrate(User{}, Server{})
	if err != nil {
		return err
	}
	return nil
}

func readUserData() {

}

func GetUser(username string) (*User, error) {
	var user User
	err := db.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user *User) error {
	err := db.Model(&User{}).Where("username = ?", user.Username).Updates(user).Error
	if err != nil {
		return err
	}
	return nil
}

func RegisterUser(user *User) error {
	err := db.Model(&User{}).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}
func DeleteUser(user *User) error {
	err := db.Model(&User{}).Delete(user).Error
	if err != nil {
		return err
	}
	return nil
}

func AddServer(username, ipAndPort string) error {
	var u User
	db.Model(&User{}).Where("username = ?", username).First(&u)

	var server Server
	err := db.Model(&Server{}).Where("server_address = ? and user_id = ?", ipAndPort, u.ID).First(&server).Error
	if err == nil {
		return errors.New("已添加该服务器")
	}

	err = db.Model(&Server{}).Create(&Server{
		ServerAddress: ipAndPort,
		UserID:        strconv.Itoa(int(u.ID)),
	}).Error
	if err != nil {
		return err
	}
	return nil
}
func RemoveServer(username, ipAndPort string) error {
	var u User
	db.Model(&User{}).Where("username = ?", username).First(&u)

	var server Server
	err := db.Model(&Server{}).Where("server_address = ? and user_id = ?", ipAndPort, u.ID).First(&server).Error
	if err != nil {
		return errors.New("不存在该服务器")
	}

	err = db.Where("server_address = ? AND user_id = ?", ipAndPort, strconv.Itoa(int(u.ID))).Delete(&Server{}).Error

	if err != nil {
		return err
	}
	return nil
}

func GetUserServers(username string) []string {
	var u User
	db.Model(&User{}).Where("username = ?", username).First(&u)

	var servers []Server
	db.Model(&Server{}).Where("user_id = ?", u.ID).Find(&servers)

	var serversNames []string
	for _, server := range servers {
		serversNames = append(serversNames, server.ServerAddress)
	}
	return serversNames
}
func MemoryGetKey(k string) (string, error) {
	v, err := rdb.Get(ctx, k).Result()
	if err != nil {
		return "", err
	}
	return v, nil
}
func MemorySetKey(k, v string, t int) error {
	err := rdb.Set(ctx, k, v, time.Duration(t)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}
func MemoryDelKey(k string) error {
	err := rdb.Del(ctx, k).Err()
	if err != nil {
		return err
	}
	return nil
}
func MemoryLPush(username, cookie string) error {
	err := rdb.LPush(ctx, "cookie:"+username, cookie).Err()
	if err != nil {
		return err
	}
	return nil
}
func MemoryLPop(username string) (string, error) {
	value, err := rdb.LPop(ctx, "cookie:"+username).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
func MemorySetHash(u *User) error {
	err := rdb.HMSet(ctx, "user:"+u.Username, "Username", u.Username, "Password", u.Password).Err()
	if err != nil {
		return err
	}
	return nil
}

func MemoryGetHash(key string) (*User, error) {
	user, err := rdb.HGetAll(ctx, "user:"+key).Result()
	if err != nil {
		return nil, err
	}
	return &User{
		Username: user["Username"],
		Password: user["Password"],
	}, nil
}
