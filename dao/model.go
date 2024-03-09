package dao

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null;type:varchar(20);unique"`
	Password string `gorm:"not null;type:varchar(20)"`
}
type Server struct {
	gorm.Model
	ServerAddress string `gorm:"not null;type:varchar(256)"`

	User   User `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID string
}
type Repository struct {
	gorm.Model
	RepoUsername string `gorm:"not null;type:varchar(256);uniqueIndex:idx_repoUsername_userID"`
	RepoPassword string `gorm:"not null;type:varchar(256)"`

	User   User   `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID string `gorm:"uniqueIndex:idx_repoUsername_userID"`
}
