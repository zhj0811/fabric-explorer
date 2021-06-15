package handler

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jzsg/fabric-explorer/apiserver/common"
)

const (
	issuer = "dbzl"
	expire = 3 * time.Hour
	key    = "dbzlpasswd"
)

// TokenAuthMiddleware check token
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			Response(c, errors.New("token is nil"), common.TokenNilErr, nil)
			c.Abort()
			return
		}
		t, err := ParseToken(token)
		if err != nil {
			Response(c, err, common.TokenInvalidErr, nil)
			c.Abort()
			return
		}
		c.Request.Header.Set("id", t)
		c.Next()
	}
}

//生成jwt signed token
func GenerateToken(id string) string {
	//jwfConf := config.JWTInfo()
	claim := jwt.StandardClaims{
		//Audience:
		ExpiresAt: time.Now().Add(expire).Unix(),
		Id:        id,
		Issuer:    issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		panic(err)
	}
	return ss
}

//parse jwt token
func ParseToken(ss string) (string, error) {
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(ss, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	return claims.Id, nil
}
