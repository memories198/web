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

func AddRepository(username, repoUsername, repoPassword, repository string) error {
	var user User
	err := db.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return err
	}

	var repo Repository
	err = db.Model(&Repository{}).Unscoped().Where("user_id = ? and repo_username = ? and server_address = ? ",
		user.ID, repoUsername, repository).First(&repo).Error
	if err == nil {
		//如果存在已经被软删除的相同行，将它的deleted_at设为nil就可以恢复出来，从而完成添加仓库信息
		if repo.DeletedAt.Valid == true {
			db.Model(&Repository{}).Unscoped().Where("user_id = ? and repo_username = ? and server_address = ? ",
				user.ID, repoUsername, repository).Update("deleted_at", nil)
		} else {
			return errors.New("已添加该仓库")
		}
		return nil
	}

	err = db.Model(&Repository{}).Create(&Repository{
		RepoUsername:  repoUsername,
		RepoPassword:  repoPassword,
		UserID:        strconv.Itoa(int(user.ID)),
		ServerAddress: repository,
	}).Error
	if err != nil {
		return errors.New("将仓库信息保存至数据库失败：" + err.Error())
	}
	return nil
}
func RemoveRepository(username, repoUsername, repository string) error {
	var user User
	err := db.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return err
	}

	err = db.Model(&Repository{}).Where("user_id = ? and repo_username = ? and server_address = ? ",
		user.ID, repoUsername, repository).First(nil).Error
	if err != nil {
		return errors.New("该仓库不存在：" + err.Error())
	}

	err = db.Model(&Repository{}).Where("user_id = ? and repo_username = ? and server_address = ?",
		user.ID, repoUsername, repository).Delete(nil).Error
	if err != nil {
		return errors.New("从数据库删除仓库信息失败：" + err.Error())
	}
	return nil
}
func GetRepository(username, repository string) ([]Repository, error) {
	var user User
	err := db.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	var repo []Repository
	err = db.Model(&Repository{}).Where("user_id = ? and server_address = ?",
		user.ID, repository).First(&repo).Error
	if err != nil {
		return nil, err
	}
	return repo, nil
}
func GetAllRepositories(username string) ([]*Repository, error) {
	var user User
	err := db.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	var repo []*Repository
	err = db.Model(&Repository{}).Where("user_id = ?",
		user.ID).Find(&repo).Error
	if err != nil {
		return nil, err
	}
	return repo, nil
}
