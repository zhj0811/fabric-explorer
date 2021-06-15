package handler

import (
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jzsg/fabric-explorer/apiserver/common"
	"github.com/jzsg/fabric-explorer/apiserver/db"
)

func Login(c *gin.Context) {
	var name = c.GetHeader("user")
	var passwd = c.GetHeader("passwd")
	token, errCode, err := login(name, passwd)
	if err != nil {
		logger.Errorf("User %s login failed %s", name, err.Error())
		Response(c, err, errCode, nil)
		return
	}
	logger.Infof("User %s login success", name)
	Response(c, nil, common.Success, token)
	return
}

func login(name, passwd string) (string, int, error) {
	user, err := db.GetUserByName(name)
	if err != nil {
		return "", common.UserInvalidErr, errors.New("name invalid")
	}
	if user.Passwd != PasswdEncryMD5(passwd) {
		return "", common.UserNameOrPasswdErr, errors.New("name invalid")
	}
	return GenerateToken(name), common.Success, nil
}

// User password MD5 encryption
func PasswdEncryMD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func Register(c *gin.Context) {
	var user db.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		logger.Errorf("Read request service failed %s", err.Error())
		Response(c, err, common.RequestFormatErr, nil)
		return
	}
	errCode, err := register(&user)
	if err != nil {
		logger.Errorf("Register new user %s failed %s", user.Name, err.Error())
		Response(c, err, errCode, nil)
		return
	}
	logger.Infof("Register new user %s success", user.Name)
	Response(c, nil, common.Success, nil)
	return
}

func register(user *db.User) (int, error) {
	var count int64
	err := db.DB.Model(&db.User{}).Where("name = ?", user.Name).Count(&count).Error
	if err != nil {
		return common.GetDBErr, err
	}
	err = db.DB.Model(&db.User{}).Create(&user).Error
	return common.InsertDBErr, err
}
