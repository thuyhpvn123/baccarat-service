package services

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/cmd/client"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/pkg/transaction"
	log "github.com/sirupsen/logrus"
)

type SendTransactionService interface {
	SetDeck(roomNumber uint64,encryptDeck []string)(interface{}, error)
	GetCardsService( roomNumberStr string,playerAdd string) (interface{},error)
	GetAllPlayersAdd( roomid string) (interface{},error)
	Confirm( roomid string, chosenPlayerAdd string,decodedKey []string, decodedCardsArr []string) (interface{},error)
}

type sendTransactionService struct {
	chainClient     *client.Client
	baccaratAbi     *abi.ABI
	baccaratAddress e_common.Address
}
func NewSendTransactionService(
	chainClient     *client.Client,
	baccaratAbi     *abi.ABI,
	baccaratAddress e_common.Address,
) SendTransactionService {
	return &sendTransactionService{
		chainClient:     chainClient,
		baccaratAbi:     baccaratAbi,
		baccaratAddress: baccaratAddress,
	}
}

func (h *sendTransactionService) SetDeck( roomNumber uint64,encryptDeck []string) (interface{},error) {
	log.Info("encryptDeck setDeck:",encryptDeck)
	log.Info("roomNumber la: ",roomNumber)
	var result interface{}
	result = "revert"
	input, err := h.baccaratAbi.Pack(
		"setDeck",
		uint256.NewInt(roomNumber).ToBig(),
		encryptDeck,
	)
	if err != nil {
		logger.Error("error when pack call data", err)
		return nil, err
	}
	callData := transaction.NewCallData(input)

	bData, err := callData.Marshal()
	if err != nil {
		logger.Error("error when marshal call data", err)
		return nil, err
	}

	relatedAddress := []e_common.Address{}
	maxGas := uint64(5_000_000)
	maxGasPrice := uint64(1_000_000_000)
	timeUse := uint64(0)
	receipt,err := h.chainClient.SendTransaction(
		h.baccaratAddress,
		uint256.NewInt(0),
		pb.ACTION_CALL_SMART_CONTRACT,
		bData,
		relatedAddress,
		maxGas,
		maxGasPrice,
		timeUse,
	)
	fmt.Println(receipt)
	if(receipt.Status() == pb.RECEIPT_STATUS_RETURNED){
		kq := make(map[string]interface{})
		err =  h.baccaratAbi.UnpackIntoMap(kq, "setDeck", receipt.Return())
		if err != nil {
			logger.Error("UnpackIntoMap")
		}
		result = kq[""]	
		logger.Info("SetDeck - Result - ", result)
	}else{
		result = hex.EncodeToString(receipt.Return())
		logger.Info("SetDeck - Result - ", result)
	
	}
	return result ,err

}
func (h *sendTransactionService) GetCardsService( roomNumberStr string,playerAdd string) (interface{},error) {
	var result interface{}
	result = "revert"
	roomNumber, err := strconv.ParseUint(roomNumberStr, 10, 64)
	if err != nil {
		log.Error("Error converting string to uint64:", err)
		
	}
	input, err := h.baccaratAbi.Pack(
		"getPlayerCards",
		uint256.NewInt(roomNumber).ToBig(),
		common.HexToAddress(playerAdd),
	)
	if err != nil {
		logger.Error("error when pack call data", err)
		return nil, err
	}
	callData := transaction.NewCallData(input)

	bData, err := callData.Marshal()
	if err != nil {
		logger.Error("error when marshal call data", err)
		return nil, err
	}
	fmt.Println("input: ",hex.EncodeToString(bData))
	relatedAddress := []e_common.Address{}
	maxGas := uint64(5_000_000)
	maxGasPrice := uint64(1_000_000_000)
	timeUse := uint64(0)
	receipt,err := h.chainClient.SendTransaction(
		h.baccaratAddress,
		uint256.NewInt(0),
		pb.ACTION_CALL_SMART_CONTRACT,
		bData,
		relatedAddress,
		maxGas,
		maxGasPrice,
		timeUse,
	)
	fmt.Println("GetCardsService:",receipt)
	if(receipt.Status() == pb.RECEIPT_STATUS_RETURNED){
		kq := make(map[string]interface{})
		err =  h.baccaratAbi.UnpackIntoMap(kq, "getPlayerCards", receipt.Return())
		if err != nil {
			logger.Error("UnpackIntoMap")
		}
		result = kq[""]	
		logger.Info("getPlayerCards - Result - ", result)
	}else{
		result = hex.EncodeToString(receipt.Return())
		logger.Info("getPlayerCards - Result - ", result)
	
	}
	return result ,err

}
func (h *sendTransactionService) GetAllPlayersAdd( roomid string) (interface{},error) {
	var result interface{}
	result = "revert"
	roomNumber, err := strconv.ParseUint(roomid, 10, 64)
	if err != nil {
		log.Error("Error converting string to uint64:", err)
		
	}
	input, err := h.baccaratAbi.Pack(
		"getAllPlayersAdd",
		uint256.NewInt(roomNumber).ToBig(),
	)
	if err != nil {
		logger.Error("error when pack call data GetChosenPlayerAddr", err)
		return nil, err
	}
	callData := transaction.NewCallData(input)

	bData, err := callData.Marshal()
	if err != nil {
		logger.Error("error when marshal call data GetChosenPlayerAddr", err)
		return nil, err
	}
	fmt.Println("input: ",hex.EncodeToString(bData))
	relatedAddress := []e_common.Address{}
	maxGas := uint64(5_000_000)
	maxGasPrice := uint64(1_000_000_000)
	timeUse := uint64(0)
	receipt,err := h.chainClient.SendTransaction(
		h.baccaratAddress,
		uint256.NewInt(0),
		pb.ACTION_CALL_SMART_CONTRACT,
		bData,
		relatedAddress,
		maxGas,
		maxGasPrice,
		timeUse,
	)
	fmt.Println("GetChosenPlayerAddr:",receipt)
	if(receipt.Status() == pb.RECEIPT_STATUS_RETURNED){
		kq := make(map[string]interface{})
		err =  h.baccaratAbi.UnpackIntoMap(kq, "getChosenPlayerAddr", receipt.Return())
		if err != nil {
			logger.Error("UnpackIntoMap")
		}
		result = kq[""]	
		logger.Info("getChosenPlayerAddr - Result - ", result)
	}else{
		result = hex.EncodeToString(receipt.Return())
		logger.Info("getChosenPlayerAddr - Result - ", result)
	
	}
	return result ,err

}
func (h *sendTransactionService) Confirm( roomid string, chosenPlayerAdd string,decodedKey []string, decodedCardsArr []string) (interface{},error) {
	var result interface{}
	result = "revert"
	roomNumber, err := strconv.ParseUint(roomid, 10, 64)
	if err != nil {
		log.Error("Error converting string to uint64:", err)
		
	}
	input, err := h.baccaratAbi.Pack(
		"comfirm",
		uint256.NewInt(roomNumber).ToBig(),
		common.HexToAddress(chosenPlayerAdd),
		decodedKey,
		decodedCardsArr,
	)
	if err != nil {
		logger.Error("error when pack call data Confirm", err)
		return nil, err
	}
	callData := transaction.NewCallData(input)

	bData, err := callData.Marshal()
	if err != nil {
		logger.Error("error when marshal call data Confirm", err)
		return nil, err
	}
	fmt.Println("input: ",hex.EncodeToString(bData))
	relatedAddress := []e_common.Address{}
	maxGas := uint64(5_000_000)
	maxGasPrice := uint64(1_000_000_000)
	timeUse := uint64(0)
	receipt,err := h.chainClient.SendTransaction(
		h.baccaratAddress,
		uint256.NewInt(0),
		pb.ACTION_CALL_SMART_CONTRACT,
		bData,
		relatedAddress,
		maxGas,
		maxGasPrice,
		timeUse,
	)
	fmt.Println("Confirm:",receipt)
	if(receipt.Status() == pb.RECEIPT_STATUS_RETURNED || receipt.Status() == pb.RECEIPT_STATUS_HALTED){
		// kq := make(map[string]interface{})
		// err =  h.baccaratAbi.UnpackIntoMap(kq, "comfirm", receipt.Return())
		// if err != nil {
		// 	logger.Error("UnpackIntoMap")
		// }
		result = "comfirmed"
		logger.Info("Confirm - Result - ", "comfirmed")
	}else{
		result = hex.EncodeToString(receipt.Return())
		logger.Info("Confirm - Result - ", result)
	
	}
	return result ,err

}

