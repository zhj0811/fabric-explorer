package handler

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/channelconfig"
	"github.com/jzsg/fabric-explorer/apiserver/db"
	"github.com/jzsg/fabric-explorer/apiserver/sdk"
	"github.com/zhj0811/gohfc/pkg/parseBlock"
	"gorm.io/gorm"
)

func Listening() error {
	block, err := db.GetLatestBlock()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		block = &db.Block{BlockHeight: -1}
	}
	logger.Infof("Get block height %d", block.BlockHeight)
	blockChan := make(chan parseBlock.Block)
	go handleEvent(blockChan)
	blockChan, err = sdk.ListenEventFullBlock(block.BlockHeight + 1)
	if err != nil {
		return err
	}
	handleEvent(blockChan)
	return nil
}

func handleEvent(blockChan chan parseBlock.Block) {
	for {
		block := <-blockChan
		handleBlock(block)
		handleTx(block)
	}
}

func handleBlock(event parseBlock.Block) error {
	block := db.Block{
		BlockHeight:    int64(event.Header.Number),
		BlockHash:      fmt.Sprintf("%x", event.Header.DataHash),
		BlockPreHash:   fmt.Sprintf("%x", event.Header.PreviousHash),
		BlockTimeStamp: event.BlockTimeStamp,
		BlockTxCount:   len(event.Transactions),
		BlockSize:      int64(event.Size),
		//BlockConfirmTime: event.BlockTimeStamp.Unix()-event.FirstTxTime.Unix(),
	}
	return db.DB.Create(&block).Error
}

func handleTx(event parseBlock.Block) {
	txCount := len(event.Transactions)
	if txCount <= 0 {
		return
	}

	for txIndex, txTmp := range event.Transactions {
		// insert data to tx table
		//channel := txTmp.ChannelHeader.ChannelId
		tx := db.Tx{
			TxID: txTmp.ChannelHeader.TxId,
		}

		tx.TxHash = fmt.Sprintf("%x", txTmp.ProposalHash)
		tx.TxBlockHeight = int64(event.Header.Number)
		startTime := time.Unix(txTmp.ChannelHeader.Timestamp.Seconds, int64(txTmp.ChannelHeader.Timestamp.Nanos)).UTC()
		tx.TxStartTime = startTime
		tx.TxConfirmTime = int64(event.BlockTimeStamp.Unix() - tx.TxStartTime.Unix())
		tx.TxType = common.HeaderType_name[txTmp.ChannelHeader.Type]
		tx.TxValidStatus = txTmp.ValidationCode == 0
		tx.TxID = txTmp.ChannelHeader.TxId
		tx.TxIndex = int64(txIndex)
		tx.TxSize = int64(txTmp.Size)
		tx.TxBlockHash = fmt.Sprintf("%x", event.Header.DataHash)
		// tx.TxContent = event.TxContent
		tx.LastConfigBlockNum = int64(event.LastConfigBlockNumber.LastConfigBlockNum)

		tx.TxFinishTime = event.BlockTimeStamp
		//chaincode := db.Chaincode{}
		if tx.TxType == common.HeaderType_name[int32(common.HeaderType_ENDORSER_TRANSACTION)] {
			if len(txTmp.ChaincodeSpec.Input.Args) > 0 {
				tx.Function = txTmp.ChaincodeSpec.Input.Args[0]
			}

			if len(txTmp.ChaincodeSpec.Input.Args) > 1 {
				args, _ := json.Marshal(txTmp.ChaincodeSpec.Input.Args[1:])
				tx.Args = string(args)
			}
			proposalResponse, err := proto.Marshal(txTmp.Response)
			if err != nil {
				logger.Error("proto marshal ProposalResponse err : ", err)
				continue
			}
			tx.ProposalResponse = base64.StdEncoding.EncodeToString(proposalResponse)
			tx.ChaincodeID = txTmp.ChannelHeader.ChaincodeId.Name
			tx.ChaincodeVersion = txTmp.ChannelHeader.ChaincodeId.Version
			tx.ChaincodeType = peer.ChaincodeSpec_Type_name[int32(txTmp.ChaincodeSpec.Type)]
			tx.ChaincodePath = txTmp.ChannelHeader.ChaincodeId.Path
			if (tx.Function == "deploy" || tx.Function == "upgrade") && tx.ChaincodeID == "lscc" {
				//chaincode.ChannelUUID = event.ChannelID
				//chaincode.UUID = tx.ChaincodeID
				//depc, err := utils.GetChaincodeDeploymentSpec([]byte(txTmp.ChaincodeSpec.Input.Args[2]))
				//if err == nil {
				//	var args []string
				//	for _, arg := range depc.ChaincodeSpec.Input.Args {
				//		args = append(args, string(arg))
				//	}
				//	initParam, _ := json.Marshal(args)
				//	chaincode.InitParam = string(initParam)
				//	chaincode.Version = depc.ChaincodeSpec.ChaincodeId.Version
				//	policy := &common.SignaturePolicyEnvelope{}
				//	err := proto.Unmarshal([]byte(txTmp.ChaincodeSpec.Input.Args[3]), policy)
				//	if err == nil {
				//		buf := &bytes.Buffer{}
				//		if err := protolator.DeepMarshalJSON(buf, policy); err == nil {
				//			chaincode.EndorsePolicy = buf.String()
				//		}
				//	}
				//
				//}
			} else {
				//chaincode.GetNewestChaincodeByUUID(tx.ChaincodeID)
				//tx.ChaincodePath = chaincode.ChaincodePath
				//tx.InitParam = chaincode.InitParam
				//tx.EndorsePolicy = chaincode.EndorsePolicy
				//tx.ChaincodeVersion = chaincode.Version
				//ccInfo, err := fabric_db.GetChaincodeByUUID(chaincode.UUID)
				//if err == nil {
				//	tx.ChaincodeName = ccInfo.Name
				//}

			}
			tx.MspId = txTmp.SignatureHeader.MspId
			var certs []*x509.Certificate
			for _, endorsement := range txTmp.Endorsements {
				certs = append(certs, endorsement.SignatureHeader.Certificate)
				if tx.EndorseSignature != "" {
					tx.EndorseSignature += ","
				}
				tx.EndorseSignature += base64.StdEncoding.EncodeToString(endorsement.Signature)
			}
			jsonByte, _ := json.Marshal(certs)
			tx.EndorseCerts = string(jsonByte)

			//if tx.Function == "deploy" && tx.ChaincodeID == "lscc" {
			//	chaincode.Insert()
			//} else if tx.Function == "upgrade" && tx.ChaincodeID == "lscc" {
			//	chaincode.Update()
			//}

		} else if tx.TxType == common.HeaderType_name[int32(common.HeaderType_CONFIG)] {
			config := &parseBlock.Config{}
			if err := json.Unmarshal([]byte(event.Config), config); err != nil {
				logger.Error("failed to json unmarshal config block : ", event.Config)
				continue
			}

			policy, _ := json.Marshal(config.ChannelGroup.Groups[channelconfig.ApplicationGroupKey].Policies)
			tx.Policy = string(policy)
			tx.MspId = txTmp.SignatureHeader.MspId

			// get peer org info ca cert and anchors
			//for mspid, org := range config.ChannelGroup.Groups[channelconfig.ApplicationGroupKey].Groups {
			//	if _, ok := org.Values[channelconfig.AnchorPeersKey]; ok {
			//		// get anchors info
			//		for _, anchor := range org.Values[channelconfig.AnchorPeersKey].Value.AnchorPeers {
			//			if tx.AnchorInfo != "" {
			//				tx.AnchorInfo += ","
			//			}
			//			tx.AnchorInfo += anchor.Host
			//		}
			//	}
			//
			//	//orgCa := db.NewOrgCaTable()
			//	//if err := orgCa.GetCaByMspID(mspid, txTmp.ChannelHeader.ChannelId); err == nil {
			//	//	continue
			//	//}
			//	cas := org.Values[channelconfig.MSPKey].Value.Config.RootCerts
			//	for _, ca := range cas {
			//		caCert, err := base64.StdEncoding.DecodeString(ca)
			//		if err != nil {
			//			logger.Errorf("failed to decode ca use base64, ca = %s, err : %s", string(caCert), err)
			//			continue
			//		}
			//	//	cert, err := ParseCerterficate(caCert)
			//	//	if err != nil {
			//	//		logger.Error("failed to parse root ca, err : ", err)
			//	//		continue
			//	//	}
			//	//
			//	//	data, err := json.Marshal(cert)
			//	//	if err != nil {
			//	//		continue
			//	//	}
			//	//
			//	//	if orgCa.Ca != "" {
			//	//		orgCa.Ca += ","
			//	//	}
			//	//	orgCa.Ca += string(data)
			//	//}
			//	//orgCa.MspId = mspid
			//	//orgCa.ChannelUUID = txTmp.ChannelHeader.ChannelId
			//	//orgCa.Insert()
			//}
			//// get orderer org ca cert
			////for mspid, org := range config.ChannelGroup.Groups[channelconfig.OrdererGroupKey].Groups {
			////	orgCa := db.NewOrgCaTable()
			////	if err := orgCa.GetCaByMspID(mspid, txTmp.ChannelHeader.ChannelId); err == nil {
			////		continue
			////	}
			////	cas := org.Values[channelconfig.MSPKey].Value.Config.RootCerts
			////	for _, ca := range cas {
			////		caCert, err := base64.StdEncoding.DecodeString(ca)
			////		if err != nil {
			////			logger.Errorf("failed to decode ca use base64, ca = %s, err : %s", string(caCert), err)
			////			continue
			////		}
			////		cert, err := ParseCerterficate(caCert)
			////		if err != nil {
			////			logger.Error("failed to parse root ca, err : ", err)
			////			continue
			////		}
			////
			////		data, err := json.Marshal(cert)
			////		if err != nil {
			////			continue
			////		}
			////
			////		if orgCa.Ca != "" {
			////			orgCa.Ca += ","
			////		}
			////		orgCa.Ca += string(data)
			////	}
			////	orgCa.MspId = mspid
			////	orgCa.ChannelUUID = txTmp.ChannelHeader.ChannelId
			////	orgCa.Insert()
			//}
		}

		//orgCa := db.NewOrgCaTable()
		//if err := orgCa.GetCaByMspID(tx.MspId, txTmp.ChannelHeader.ChannelId); err != nil {
		//	logger.Error("failed to get ca cert by mspid, err : ", err)
		//} else {
		//	tx.CaCert = orgCa.Ca
		//}

		jsonByte, _ := json.Marshal(txTmp.TxActionSignatureHeader.Certificate)
		tx.TxActionCert = string(jsonByte)

		tx.Signature = base64.StdEncoding.EncodeToString(txTmp.Signature)
		if event.BlockCreatorSignature != nil {
			if event.BlockCreatorSignature.SignatureHeader != nil {
				jsonByte, _ = json.Marshal(event.BlockCreatorSignature.SignatureHeader.Certificate)
				tx.BlockCreatorCert = string(jsonByte)
			}
		}
		if tx.TxID == "" {
			tx.TxID = tx.TxHash
		}
		if err := db.DB.Create(&tx); err != nil {
			return
		}
	}
}
