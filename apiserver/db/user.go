package db

import (
	"time"
)

type User struct {
	//ID        string    `json:"id" gorm:"primary_key;column:id"`                               //user唯一标识
	Name      string    `gorm:"primaryKey;column:name" form:"name" json:"name" binding:"required"` //名称唯一 - 可修改
	Nickname  string    `json:"nickname" gorm:"nickname" binding:"required"`                       //昵称名
	Role      uint32    `json:"role"`                                                              // 1 创建链权限 0 无创建链权限
	Admin     bool      `json:"admin"`                                                             //是否已删除
	Passwd    string    `json:"passwd" gorm:"not null" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetUserByName(name string) (user User, err error) {
	err = DB.Model(&User{}).First(&user, "name = ?", name).Error
	return
}
