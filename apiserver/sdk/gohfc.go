package sdk

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	gohfc "github.com/zhj0811/gohfc/pkg"
	"github.com/zhj0811/gohfc/pkg/parseBlock"
)

var (
	hfcSdk gohfc.Sdk
)

//InitSDKs 初始化sdk
func InitSDKs(configFilePath string) error {
	viper.SetEnvPrefix("core")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigFile(configFilePath)

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error when initializing SDK config: %s", err)
	}
	fmt.Println("viper read in config success.")
	hfcSdk, err = gohfc.New(configFilePath)
	if err != nil {
		return err
	}
	return nil
}

//Invoke invoke
func Invoke(in []string) (*gohfc.InvokeResponse, error) {
	res, err := hfcSdk.Invoke(in, nil, "", "")
	if err != nil && res == nil {
		return nil, err
	}
	return res, err
}

//Query 查询
func Query(in []string) ([]byte, error) {
	res, err := hfcSdk.Query(in, nil, "", "")
	if err != nil {
		return nil, err
	}
	if res[0].Error != nil {
		return nil, err
	}
	return res[0].Response.Response.Payload, nil
}

//GetUserId 获取userid
func GetUserId() string {
	return viper.GetString("user.id")
}

func GetBlockHeightByTxID(txId string) (uint64, error) {
	block, err := hfcSdk.GetBlockByTxID(txId, "")
	if err != nil {
		return 0, err
	}
	return block.BlockNum, nil
}

func GetFilterTxByTxID(txId string) (*parseBlock.FilterTx, error) {
	tx, err := hfcSdk.GetFilterTxByTxID(txId, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func ListenEventFullBlock(startNum int64) (chan parseBlock.Block, error) {
	return hfcSdk.ListenEventFullBlock("", startNum)
}
