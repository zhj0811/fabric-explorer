package define

//保单
type Policy struct {
	ID         string `json:"id"`                                                   //全局唯一标识
	Number     string `json:"number" binding:"required" gorm:"not null;primaryKey"` //保单号
	Type       string `json:"type" gorm:"not null"`                                 //保单类型
	Insured    string `json:"insured"`                                              //被保险人
	USCC       string `json:"uscc" gorm:"uscc"`                                     //统一社会信用代码
	StartAt    string `json:"start_at" gorm:"start_at"`                             //起保时间
	ExpireAt   string `json:"expire_at" gorm:"expire_at"`                           //终保时间
	Pooled     bool   `json:"pooled"`                                               //是否共保
	Insurer    string `json:"insurer"`                                              //承保人
	Amount     string `json:"amount"`                                               //保险金额（元）
	Premium    string `json:"premium"`                                              //保险费（元）
	Rate       string `json:"rate"`                                                 //保险费率
	Content    string `json:"content"`                                              //保障内容
	TxID       string `json:"tx_id" gorm:"column:tx_id; index"`
	Extension1 string `json:"extension1"` //预留字段1
	Extension2 string `json:"extension2"` //预留字段2
	Extension3 string `json:"extension3"` //预留字段3
}

type PolicyInfo struct {
	Policy
	BlockHeight uint64 `json:"block_height"`
	Timestamp   int64  `json:"timestamp"`
}
