package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jzsg/fabric-explorer/common/define"
)

//保单表
type TPolicy struct {
	define.Policy
	BlockHeight uint64    `json:"block_height"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

//func InsertService() error {
//
//}

func GetPolicyByNumber(id string) (res *TPolicy, err error) {
	err = DB.Model(&TPolicy{}).First(&res, "number = ?", id).Error
	return
}

func GetPolicies(p, c, name string) ([]TPolicy, int64, error) {
	page, _ := strconv.Atoi(p)
	if page < 1 {
		page = 1
	}
	count, _ := strconv.Atoi(c)
	if count < 5 {
		count = 10
	}
	var res []TPolicy
	var totalCount int64
	err := DB.Model(&TPolicy{}).Where("insured LIKE ?", fmt.Sprintf("'%%%s%%'", name)).
		Count(&totalCount).Order("updated_at desc").
		Limit(count).Offset((page - 1) * count).Find(&res).Error
	return res, totalCount, err
}
