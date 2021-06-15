package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jzsg/fabric-explorer/apiserver/common"
	"github.com/jzsg/fabric-explorer/apiserver/db"
	"github.com/jzsg/fabric-explorer/apiserver/sdk"
	"github.com/jzsg/fabric-explorer/common/define"
)

func GetPolicies(c *gin.Context) {
	page := c.Query("page")
	count := c.Query("count")
	insured := c.Query("insured")
	list, totalCount, err := db.GetPolicies(page, count, insured)
	if err != nil {
		logger.Errorf("Get policy list failed %s", err.Error())
		Response(c, err, common.GetDBErr, nil)
		return
	}
	logger.Infof("Get policies %+v", list)
	res := ListInfo{Total: totalCount, List: list}
	Response(c, nil, common.Success, res)
	return
}

func UploadInvokePolicy(c *gin.Context) {
	policy := &define.Policy{}
	err := c.ShouldBindJSON(policy)
	if err != nil {
		logger.Errorf("Read request policy failed %s", err.Error())
		Response(c, err, common.RequestFormatErr, nil)
		return
	}
	txId, errCode, err := invokePolicy(policy)
	if err != nil {
		logger.Errorf("Invoke policy failed %s", err.Error())
		Response(c, err, errCode, nil)
		return
	}
	logger.Infof("Invoke policy %s success, tx id: %s", policy.ID, txId)
	Response(c, nil, common.Success, txId)
	return
}

func invokePolicy(policy *define.Policy) (string, int, error) {
	args, err := json.Marshal(policy)
	if err != nil {
		//logger.Errorf("Marshal policy failed %s", err.Error())
		//Response(c, err, common.RequestFormatErr, nil)
		return "", common.RequestFormatErr, err
	}
	req := []string{define.SavePolicy, string(args)}
	res, err := sdk.Invoke(req)
	if err != nil {
		return "", common.InvokeErr, err
	}
	return res.TxID, common.Success, nil
}
