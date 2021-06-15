package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jzsg/fabric-explorer/apiserver/common"
	"github.com/jzsg/fabric-explorer/apiserver/db"
)

type ListInfo struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

type PolicyFullInfo struct {
	Policy   db.TPolicy     `json:"policy"`
	Services []*db.TService `json:"services"`
}

func Response(c *gin.Context, err error, errCode int, data interface{}) {
	res := &common.ResponseInfo{
		Data: data,
	}
	if err != nil {
		res.ErrCode = errCode
		res.ErrMsg = err.Error()
	}
	//ret, _ := json.Marshal(res)
	//c.Writer.Write(ret)
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, res)
	return
}
