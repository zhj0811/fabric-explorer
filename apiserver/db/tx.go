package db

import (
	"fmt"
	"time"
)

// Tx tx info
type Tx struct {
	TxID               string    `json:"tx_id" gorm:"primary_key;column:tx_id"`
	TxHash             string    `json:"tx_hash" gorm:"index;column:tx_hash"`
	TxIndex            int64     `json:"tx_index" gorm:"column:tx_index"`
	TxBlockHeight      int64     `json:"tx_block_height" gorm:"not null,index;column:tx_block_height"`
	TxBlockHash        string    `json:"tx_block_hash" gorm:"not null,index;column:tx_block_hash"`
	TxStartTime        time.Time `json:"tx_start_time" gorm:"column:tx_start_time"` //交易开始时间（发起交易的时间）
	TxFinishTime       time.Time `json:"tx_finish_time" gorm:"column:tx_finish_time"`
	TxConfirmTime      int64     `json:"tx_confirm_time" gorm:"column:tx_confirm_time"` //交易完成时间（交易从发起到完成的时间）
	TxType             string    `json:"tx_type" gorm:"column:tx_type"`                 //普通交易和配置交易
	TxValidStatus      bool      `json:"tx_valid_status" gorm:"column:tx_valid_status"` //交易是否合法
	TxSize             int64     `json:"tx_size" gorm:"column:tx_size"`                 //交易大小
	LastConfigBlockNum int64     `json:"last_config_block_num"`                         //上一次修改该channel配置的区块号
	Function           string    `json:"function"`                                      //调用函数名称
	Args               string    `json:"args" gorm:"type:text"`                         //参数列表
	ProposalResponse   string    `json:"proposal_response" gorm:"type:text"`            //响应
	ChaincodeID        string    `json:"chaincode_id"`                                  //合约id
	ChaincodeName      string    `json:"chaincode_name"`
	ChaincodeVersion   string    `json:"chaincode_version"`                  //合约版本
	ChaincodeType      string    `json:"chaincode_type"`                     //合约类型（golang or others）
	ChaincodePath      string    `json:"chaincode_path"`                     //合约安装路径
	EndorsePolicy      string    `json:"endorse_policy" gorm:"type:text"`    //合约背书策略
	InitParam          string    `json:"init_param"`                         //合约初始化参数
	Signature          string    `json:"signature"`                          //交易发起者签名
	EndorseSignature   string    `json:"endorse_signature" gorm:"type:text"` //peer节点背书签名(数组，以逗号隔开)
	AnchorInfo         string    `json:"anchor_info"`                        //锚节点信息(配置交易中显示)
	Policy             string    `json:"policy" gorm:"type:text"`            //策略（配置交易中显示）
	MspId              string    `json:"msp_id"`                             //交易发起者msp id（配置交易中显示）
	OrgName            string    `json:"org_name"`
	CaCert             string    `json:"ca_cert" gorm:"type:text"`            //组织根证书证书
	TxActionCert       string    `json:"tx_action_cert" gorm:"type:text"`     //客户端证书（交易发起者证书）
	EndorseCerts       string    `json:"endorse_certs" gorm:"type:text"`      //Endorse证书
	BlockCreatorCert   string    `json:"block_creator_cert" gorm:"type:text"` //order证书
	//TxContent []byte `json:"tx_content" gorm:"column:tx_content"`//交易内容
	CreatedAt time.Time `json:"created_at"`
}

func GetTxsByBlockHeight(height string) (res []*Tx, err error) {
	err = DB.Debug().Model(&Tx{}).Where(fmt.Sprintf("tx_block_height = %s", height)).Find(&res).Error
	return
}

type Chaincode struct {
	TableName     string    `json:"-" gorm:"-"`
	ID            uint      `json:"id" gorm:"primary_key"`           //需要做唯一索引,所以必须存在。
	UUID          string    `json:"uuid" gorm:"not null,index"`      //后端识别名
	Version       string    `json:"version"`                         //版本
	ChannelUUID   string    `json:"channel_uuid"`                    //对应channel
	EndorsePolicy string    `json:"endorse_policy" gorm:"type:text"` //背书策略
	InitParam     string    `json:"init_param"`                      //初始化参数
	ChaincodePath string    `json:"chaincode_path"`                  //合约路径
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
