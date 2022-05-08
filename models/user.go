package models

import (
	"time"

	"gorm.io/gorm"
)

type userOrm struct {
	db *gorm.DB
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex" json:"username"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"-"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserOrmer interface {
	GetOneByID(id uint) (user User, err error)
	GetOneByUsername(username string) (user User, err error)
	InsertUser(user User) (id uint, err error)
	UpdateUser(user User) (err error)
}

func NewUserOrmer(db *gorm.DB) UserOrmer {
	//_ = db.AutoMigrate(&User{})		// builds table when enabled
	return &userOrm{db}
}

func (o *userOrm) GetOneByID(id uint) (user User, err error) {
	result := o.db.Model(&User{}).Where("id = ?", id).First(&user)
	return user, result.Error
}

func (o *userOrm) GetOneByUsername(username string) (user User, err error) {
	result := o.db.Model(&User{}).Where("username = ?", username).First(&user)
	return user, result.Error
}

func (o *userOrm) InsertUser(user User) (id uint, err error) {
	result := o.db.Model(&User{}).Create(&user)
	return user.ID, result.Error
}

func (o *userOrm) UpdateUser(user User) (err error) {
	// By default, only non-empty fields are updated. See https://gorm.io/docs/update.html#Updates-multiple-columns
	result := o.db.Model(&User{}).Model(&user).Updates(&user)
	return result.Error
}
