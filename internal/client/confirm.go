package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/meta-node-blockchain/meta-node-baccarat/internal/database"
)

// admin call after bidder choose players to compare
func (c *Client) ComfirmAll(roomid string)(map[string]interface{} ,error){
	var result map[string]interface{}
	// read key array from levelDb
	callmap1 := map[string]interface{}{
		"key": roomid ,
	}
	result1 := database.ReadValueStorage(callmap1, c.DB)
	var keysArray []string
	err := json.Unmarshal(result1["value"].([]byte), &keysArray)
	fmt.Println("encrypted-keys", keysArray)
	//read deck of roomid from levelDb
	callmap2 := map[string]interface{}{
		"key": roomid + "EncryptDeck",
	}

	result2 := database.ReadValueStorage(callmap2, c.DB)
	var deckKq []string
	err = json.Unmarshal(result2["value"].([]byte), &deckKq)
	if err != nil {
		fmt.Println("error Unmarshal:", err)
	}
	fmt.Println("encrypted-deck", deckKq)
	//get chosen addresses
	kq,err := c.getAllPlayersAdd(roomid)
	chosenAddr := kq.([]common.Address)
	var players []string
	var kqs []string
	for i := 0; i < len(chosenAddr); i++ {
		chosenPlayerAdd := fmt.Sprintf("%s",chosenAddr[i])

		//getCards for each chosen address
		cardArr ,err := c.getCards(roomid, chosenPlayerAdd)
		if err != nil {
			log.Error("err:",err)
			return nil, err
		}
		//store encoded cards of each player
		call := map[string]interface{}{
			"encrypted-cards": cardArr,
			"encrypted-deck":  deckKq,
			"encrypted-keys":  keysArray,
		}
		decodedKey := FindKeys(call)
		// //store keys of each player in db
		// barrayKey, _ := json.Marshal(decodedKey)
		// strbarrayKey := string(barrayKey)
		// callmap := map[string]interface{}{
		// 	"key":  roomid + chosenPlayerAdd+ "key",
		// 	"data": strbarrayKey,
		// }
		// result := WriteValueStorage(callmap, levelDb)
		// fmt.Println("write Keys of a player to storage success:", result["success"])
		newcardArr := make([]interface{}, len(cardArr.([]string)))
		for i, v := range cardArr.([]string) {
			newcardArr[i] = v
		}
		newkeysArray := make([]interface{}, len(decodedKey))
		for i, v := range decodedKey {
			newkeysArray[i] = v
		}
		call1 := map[string]interface{}{
			"encrypted-cards": newcardArr,
			"encrypted-keys":  newkeysArray,
		}
		fmt.Println("encrypted-cards:",newcardArr)
		//decoded cards
		kq,msg :=c.DecryptDeck(call1)
		if kq == nil{
			fmt.Println("error DecryptDeck:", msg)
			return nil, errors.New(msg)
		}
		decodedCardsArr,ok:=kq["message"].([]string)
		if !ok {
			msg:= "error input message comfirm"
			log.Error(msg)
		return nil,errors.New(msg)
		}
		//store decoded cards of each player
		// barrayCards, _ := json.Marshal(decodedCardsArr)
		// strbarrayCards := string(barrayCards)
		// callmap1 := map[string]interface{}{
		// 	"key":  roomid + chosenPlayerAdd+ "cards",
		// 	"data": strbarrayCards,
		// }
		// result1 := WriteValueStorage(callmap1, levelDb)
		// fmt.Println("write decoded cards of a player to storage success:", result1["success"])

		comfirmKq,err:=c.comfirm(roomid, chosenPlayerAdd,decodedKey, decodedCardsArr )
		if err != nil {
			fmt.Println("error comfirm:", err)
			return nil, err
		}
		players = append(players, chosenPlayerAdd)
		kqs = append(kqs,comfirmKq.(string) )
		// go cli.sentToClient("comfirm", result)
	}
	result =map[string]interface{}{
			"address":players,
			"message":kqs,
		}
		fmt.Println(result)
	return result,nil

}
func (c *Client) getAllPlayersAdd(roomid string) (interface{},error) {
	chosenAddrs,err := c.Serv.GetAllPlayersAdd(roomid)
	if err != nil {
		log.Error("Error GetChosenPlayerAddr", err)	
		return nil, err
	}
	return chosenAddrs,nil
}

func (c *Client) comfirm(roomid string, chosenPlayerAdd string,decodedKey []string, decodedCardsArr []string) (interface{},error){
	result,err := c.Serv.Confirm(roomid, chosenPlayerAdd,decodedKey, decodedCardsArr)
	if err != nil {
		log.Error("Error Confirm", err)	
		return nil, err
	}
	return result,nil
}
