package define

type Service struct {
	ID           string `json:"id" gorm:"column:id;primary_key;"`
	Number       string `json:"number" binding:"required" gorm:"not null;index"` //保单号
	Insured      string `json:"insured"`                                         //被保险人
	Date         string `json:"date"`                                            //服务时间
	Organization string `json:"organization"`                                    //服务机构
	Type         string `json:"type"`                                            //服务类型
	TxID         string `json:"tx_id" gorm:"column:tx_id;"`
	Extension1   string `json:"extension1"` //预留字段1
	Extension2   string `json:"extension2"` //预留字段2
	Extension3   string `json:"extension3"` //预留字段3
}

type ServiceInfo struct {
	Service
	BlockHeight uint64 `json:"block_height"`
	Timestamp   int64  `json:"timestamp"`
}
