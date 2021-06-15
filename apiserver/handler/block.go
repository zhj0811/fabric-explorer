package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jzsg/fabric-explorer/apiserver/common"
	"github.com/jzsg/fabric-explorer/apiserver/db"
)

func GetBlocks(ctx *gin.Context) {
	page := ctx.Query("page")
	count := ctx.Query("count")
	blocks, totalCount, err := db.PagingQueryBlocks(page, count)
	if err != nil {
		logger.Errorf("Get blocks failed %s", err.Error())
		Response(ctx, err, common.GetDBErr, nil)
		return
	}
	logger.Infof("Get blocks %+v", blocks)
	res := ListInfo{List: blocks, Total: totalCount}
	Response(ctx, nil, common.Success, res)
	return
}

func GetBlockFullInfo(ctx *gin.Context) {
	height := ctx.Param("id")
	res, err := db.GetTxsByBlockHeight(height)
	if err != nil {
		logger.Errorf("Get transactions failed %s", err.Error())
		Response(ctx, err, common.GetDBErr, nil)
		return
	}
	logger.Infof("Get transactions %+v", res)
	//res := ListInfo{List: blocks}
	Response(ctx, nil, common.Success, res)
	return
}
