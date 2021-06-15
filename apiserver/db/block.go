package db

import (
	"strconv"
	"time"
)

type Block struct {
	BlockHeight      int64     `json:"block_height" gorm:"column:block_height"`
	BlockHash        string    `json:"block_hash" gorm:"column:block_hash"`
	BlockPreHash     string    `json:"block_pre_hash" gorm:"column:block_pre_hash"`
	BlockTimeStamp   time.Time `json:"block_time_stamp" gorm:"column:block_time_stamp"`
	BlockConfirmTime int64     `json:"block_confirm_time" gorm:"column:block_confirm_time"`
	BlockTxCount     int       `json:"block_tx_count" gorm:"column:block_tx_count"`
	BlockSize        int64     `json:"block_size" gorm:"column:block_size"`
	CreatedAt        time.Time
}

func GetLatestBlock() (res *Block, err error) {
	err = DB.Model(&Block{}).Order("block_height desc").First(&res).Error
	return
}

func GetBlocks(limit int) (res []*Block, err error) {
	err = DB.Model(&Block{}).Order("block_height desc").Limit(limit).Find(&res).Error
	return
}

func PagingQueryBlocks(page, count string) (res []*Block, totalCount int64, err error) {
	p, _ := strconv.Atoi(page)
	if p < 1 {
		p = 1
	}
	limit, _ := strconv.Atoi(count)
	if limit < 1 {
		limit = 10
	}

	err = DB.Model(&Block{}).Order("block_height desc").Count(&totalCount).Limit(limit).Offset((p - 1) * limit).Find(&res).Error
	return
}
