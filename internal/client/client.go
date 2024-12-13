package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/model"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/services"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

type Client struct {
	Conn *websocket.Conn
	MsgChan chan model.Message
	DB *leveldb.DB
	Serv services.SendTransactionService
}

func (c *Client) ReadMessage() {
	for {
		var msg map[string]interface{}
		err := c.Conn.ReadJSON(&msg)
		if err != nil || c.Conn == nil {
			log.Info("error: %v", err)
			c.Conn.Close()
			break
		}
		
		c.handleCallChain(msg)
	}
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		message, ok := <-c.MsgChan
		if !ok {
			return
		}
		c.Conn.WriteJSON(message)
	}
}

func (c *Client) handleCallChain(msg map[string]interface{}) {
	fmt.Println(msg)
	switch msg["command"] {
	case "test":
		if testValue, ok := msg["value"].(string) ; ok && strings.ToLower(testValue) == "ping" {
			c.sentToClient("message", "pong")
		}
	case "deal-cards":
		call:=msg["value"].(map[string]interface{}) 
		kq := c.GeneratePlayerKeys(call)
		callCreateDeck := c.CreateDeck()
		deck := c.ShuffleDeck(callCreateDeck)	
		playerkeys,ok:=kq["message"].([]string)
		if !ok {
			log.Error("error parse playerkeys")
		}
		roomNumberStr,ok:=call["roomNumber"].(string)
		if !ok {
			log.Error("error input roomNumber-deal cards")
		}
		encryptDeck := c.EncryptDeck(deck,playerkeys,roomNumberStr)
		fmt.Println("encryptDeck:", encryptDeck)

		roomNumber, err := strconv.ParseUint(roomNumberStr, 10, 64)
		if err != nil {
			log.Error("Error converting string to uint64:", err)
			
		}
		result,err:= c.Serv.SetDeck(roomNumber,encryptDeck)
		if err != nil {
			log.Error("error client call SetDeck",err)
		}
		go c.sentToClient("set-deck", result)
	case "get-key-for-player":
		fmt.Println("map la:",msg["value"])
		call:=msg["value"].(map[string]interface{}) 
		roomNumberStr,ok:=call["roomNumber"].(string)
		if !ok {
			log.Error("error input roomNumber-getkey")
		}
		//verify first
		verifyKq:=c.VerifySign(call)
		if(verifyKq["data"].(bool)==true){
			c.GetKeyForPlayer(roomNumberStr,verifyKq["address"].(string))
		}else{
			c.sentToClient("get-key-for-player", "Not Authorised Address")
		}
		
	case "decrypt-cards":
		call:=msg["value"].(map[string]interface{})
		kq,msg:=c.DecryptDeck(call)
		if kq != nil {
			c.sentToClient("decrypt-cards",kq)
		}else {
			c.sentToClient("decrypt-cards",msg)
		}
	case "get-sign":
			call:=msg["value"].(map[string]interface{})
			c.GetSign(call)
	case "comfirm":	
		call,ok := msg["value"].(map[string]interface{})
		if !ok{
			log.Error(fmt.Sprintf("error input comfirm"))
			return
		}	
		roomid,ok:= call["roomid"].(string)
		if !ok{
			log.Error(fmt.Sprintf("error input roomid-comfirm"))
			return
		}	
		result,err := c.ComfirmAll(roomid)
		if err!= nil{
			c.sentToClient("comfirm fail : ",err)
		}else{
			c.sentToClient("comfirm success : ",result)
		}
		

	default:
		log.Warn("Require call not match: ", msg)
	}
}

func (c *Client) sentToClient(command string, data interface{}) {
	c.MsgChan <- model.Message{command, data}
}