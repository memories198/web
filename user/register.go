package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null varchar(20) unique"`
	Password string `gorm:"not null varchar(20)"`
}
type Server struct {
	gorm.Model
	ServerAddress string `gorm:"varchar(256)"`

	User   User `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID string
}
