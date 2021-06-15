package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jzsg/fabric-explorer/common/define"
)

//服务表
type TService struct {
	ID uint64 `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	define.Service
	TxID        string    `json:"tx_id" gorm:"column:tx_id"`
	BlockHeight uint64    `json:"block_height"`
	UploadAt    time.Time `json:"upload_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetServiceById(id string) (res *TService, err error) {
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return
	}
	err = DB.Model(&TService{}).First(&res, "id = ?", uid).Error
	return
}

func GetServices(p, c, number, insured string) ([]TService, int64, error) {
	page, _ := strconv.Atoi(p)
	if page < 1 {
		page = 1
	}
	count, _ := strconv.Atoi(c)
	if count < 5 {
		count = 10
	}
	var res []TService
	var totalCount int64
	err := DB.Model(&TPolicy{}).Where(fmt.Sprintf("number LIKE '%%%s%%' AND insured LIKE '%%%s%%'", number, insured)).
		Count(&totalCount).Order("updated_at desc").
		Limit(count).Offset((page - 1) * count).Find(&res).Error
	return res, totalCount, err
}

func GetServicesByNumber(number string) (res []*TService, err error) {
	err = DB.Model(&TService{}).Where("number = ?", number).Find(&res).Error
	return
}
