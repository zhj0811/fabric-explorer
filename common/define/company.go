package define

type Company struct {
	ID           string `json:"id"`               //全局唯一标识
	Name         string `json:"name"`             //公司名称
	PreviousName string `json:"previous_name"`    //曾用名
	USCC         string `json:"uscc" gorm:"uscc"` //统一社会信用代码
	Province     string `json:"province"`         //省
	City         string `json:"city"`             //市
	County       string `json:"county"`           //县区
	Address      string `json:"address"`          //详细地址
	Longitude    string `json:"longitude"`        //经度
	Latitude     string `json:"latitude"`         //纬度
	Type         string `json:"type"`             //行业类别
	TxID         string `json:"tx_id" gorm:"column:tx_id; index"`
	Extension1   string `json:"extension1"` //预留字段1
	Extension2   string `json:"extension2"` //预留字段2
	Extension3   string `json:"extension3"` //预留字段3
}

type CompanyInfo struct {
	Company
	BlockHeight uint64 `json:"block_height"`
	Timestamp   int64  `json:"timestamp"`
}
