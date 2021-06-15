package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jzsg/fabric-explorer/apiserver/common"
	"github.com/jzsg/fabric-explorer/apiserver/db"
	"github.com/jzsg/fabric-explorer/apiserver/sdk"
	"github.com/jzsg/fabric-explorer/common/define"
)

func GetServices(c *gin.Context) {
	page := c.Query("page")
	count := c.Query("count")
	insured := c.Query("insured")
	number := c.Query("number")
	list, totalCount, err := db.GetServices(page, count, number, insured)
	if err != nil {
		logger.Errorf("Get service list failed %s", err.Error())
		Response(c, err, common.GetDBErr, nil)
		return
	}
	logger.Infof("Get service %+v", list)
	res := ListInfo{Total: totalCount, List: list}
	Response(c, nil, common.Success, res)
	return
}

func UploadInvokeService(c *gin.Context) {
	service := &define.Service{}
	err := c.ShouldBindJSON(service)
	if err != nil {
		logger.Errorf("Read service policy failed %s", err.Error())
		Response(c, err, common.RequestFormatErr, nil)
		return
	}
	txId, errCode, err := invokeService(service)
	if err != nil {
		logger.Errorf("Invoke service failed %s", err.Error())
		Response(c, err, errCode, nil)
		return
	}
	logger.Infof("Invoke policy %s success, tx id: %s", service.ID, txId)
	Response(c, nil, common.Success, txId)
	return
}

func invokeService(service *define.Service) (string, int, error) {
	args, err := json.Marshal(service)
	if err != nil {
		//logger.Errorf("Marshal policy failed %s", err.Error())
		//Response(c, err, common.RequestFormatErr, nil)
		return "", common.RequestFormatErr, err
	}
	req := []string{define.SaveService, string(args)}
	res, err := sdk.Invoke(req)
	if err != nil {
		return "", common.InvokeErr, err
	}
	return res.TxID, common.Success, nil
}

func QueryService(c *gin.Context) {
	id := c.Param("id")
	res, errCode, err := queryService(id)
	if err != nil {
		logger.Errorf("Fabric query service %s failed %s", id, err.Error())
		Response(c, err, errCode, nil)
		return
	}
	logger.Infof("Fabric query service success %+v", res)
	Response(c, nil, common.Success, res)
	return
}

func queryService(id string) (*define.ServiceInfo, int, error) {
	bytes, err := queryByKey(id)
	if err != nil {
		return nil, common.QueryErr, err
	}
	if bytes == nil {
		return nil, common.Success, nil
	}
	res := &define.ServiceInfo{}
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, common.UnmarshalJSONErr, err
	}
	filterTx, err := sdk.GetFilterTxByTxID(res.TxID)
	if err != nil {
		return nil, common.QueryErr, err
	}
	res.BlockHeight = filterTx.BlockNum
	res.Timestamp = filterTx.Timestamp
	return res, common.Success, nil
}
