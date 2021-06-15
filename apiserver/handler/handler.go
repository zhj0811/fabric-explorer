package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jzsg/fabric-explorer/apiserver/common"
	"github.com/jzsg/fabric-explorer/apiserver/db"
	"github.com/jzsg/fabric-explorer/common/logging"
)

var logger = logging.NewLogger("debug", "handler")

func UploadPolicies(c *gin.Context) {
	fh, err := c.FormFile("file")

	if err != nil {
		logger.Errorf("Read file failed %s", err.Error())
		Response(c, err, common.RequestFormatErr, nil)
		return
	}
	policies, errCode, err := ParsePolicies(fh)
	if err != nil {
		logger.Errorf("Parse excel to policies failed %s", err.Error())
		Response(c, err, errCode, nil)
		return
	}
	err = db.DB.Create(policies).Error
	if err != nil {
		logger.Errorf("Insert policies to db failed %s", err.Error())
	}
	logger.Infof("Upload excel services success")
	Response(c, nil, common.Success, policies)
	return
}

func UploadServices(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		logger.Errorf("Read file failed %s", err.Error())
		Response(c, err, common.RequestFormatErr, nil)
		return
	}
	services, errCode, err := ParseServices(fh)
	if err != nil {
		logger.Errorf("Parse excel to services failed %s", err.Error())
		Response(c, err, errCode, nil)
		return
	}
	logger.Infof("Parse excel services info %+v", services)
	err = db.DB.Create(services).Error
	if err != nil {
		logger.Errorf("Insert services to db failed %s", err.Error())
	}
	logger.Infof("Upload excel services success")
	Response(c, nil, common.Success, services)
	return
}

func UploadPolicy(c *gin.Context) {
	policy := &db.TPolicy{}
	err := c.ShouldBindJSON(policy)
	if err != nil {
		logger.Errorf("Read request policy failed %s", err.Error())
		Response(c, err, common.RequestFormatErr, nil)
		return
	}
	errCode, err := uploadPolicy(policy)
	if err != nil {
		logger.Errorf("Upload policy %s failed", policy.Number)
		Response(c, err, errCode, nil)
		return
	}
	logger.Infof("Upload policy %s success", policy.Number)
	Response(c, nil, common.Success, nil)
	return
}

func UploadService(c *gin.Context) {
	service := &db.TService{}
	err := c.ShouldBindJSON(service)
	if err != nil {
		logger.Errorf("Read request service failed %s", err.Error())
		Response(c, err, common.RequestFormatErr, nil)
		return
	}
	errCode, err := uploadService(service)
	if err != nil {
		logger.Errorf("Upload service %s failed", service.ID)
		Response(c, err, errCode, nil)
		return
	}
	logger.Infof("Upload service %s success", service.ID)
	Response(c, nil, common.Success, nil)
	return
}

func uploadPolicy(policy *db.TPolicy) (int, error) {
	logger.Infof("Upload policy %+v", *policy)
	//var p db.TPolicy
	var count int64
	err := db.DB.Model(&db.TPolicy{}).Where("number = ?", policy.Number).Count(&count).Error
	if err != nil {
		return common.GetDBErr, err
	}
	if count == 0 {
		err = db.DB.Create(policy).Error
		if err != nil {
			return common.InsertDBErr, err
		}
	} else {
		err = db.DB.Model(&db.TPolicy{}).Save(&policy).Error
		if err != nil {
			return common.UpdateDBErr, err
		}
	}
	return common.Success, nil
}

func uploadService(service *db.TService) (int, error) {
	logger.Infof("Upload service %+v", *service)
	//var p db.TService

	if service.ID == 0 {
		err := db.DB.Create(service).Error
		if err != nil {
			return common.UpdateDBErr, err
		}
	} else {
		err := db.DB.Model(&db.TService{}).Save(&service).Error
		if err != nil {
			return common.UpdateDBErr, err
		}
	}
	return common.Success, nil
}

func GetPolicyByNumber(c *gin.Context) {
	id := c.Param("id")
	res, err := db.GetPolicyByNumber(id)
	if err != nil {
		logger.Errorf("Get policy %s failed %s", id, err.Error())
		Response(c, err, common.GetDBErr, nil)
		return
	}
	logger.Infof("Get policy %+v", *res)
	Response(c, nil, common.Success, res)
	return
}

func GetServiceById(c *gin.Context) {
	id := c.Param("id")
	res, err := db.GetServiceById(id)
	if err != nil {
		logger.Errorf("Get service %s failed %s", id, err.Error())
		Response(c, err, common.GetDBErr, nil)
		return
	}
	logger.Infof("Get service %+v", *res)
	Response(c, nil, common.Success, res)
	return
}
