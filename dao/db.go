package dao

import (
	"errors"
	"strconv"
)

func migration() error {
	err := db.AutoMigrate(User{}, Server{}, Repository{})
	if err != nil {
		return err
	}
	return nil
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
func GetUserAllServers(username string) (serversAddress []string, err error) {
	var u User
	err = db.Model(&User{}).Where("username = ?", username).First(&u).Error
	if err != nil {
		return
	}

	var servers []Server
	err = db.Model(&Server{}).Where("user_id = ?", u.ID).Find(&servers).Error
	if err != nil {
		return
	}

	for _, server := range servers {
		serversAddress = append(serversAddress, server.ServerAddress)
	}
	return
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

func AddRepository(username, repoUsername, repoPassword string) error {
	//db.Model(&Repository{}).Where()
	return nil
}
