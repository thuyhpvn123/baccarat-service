package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	random "math/rand"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/database"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/constant"
	log "github.com/sirupsen/logrus"

)
func (c *Client)GeneratePlayerKeys(call map[string]interface{})map[string]interface{} {
	result:=make(map[string]interface{})
	roomid,ok := call["roomNumber"].(string)
	if !ok {
		log.Error("error input roomNumber-GeneratePlayerKeys")
	}
	keysArr:=make([]string, 52)

	for i := 0; i < 52; i++ {
		key := make([]byte, 16) // AES-256 requires a 32-byte key
		_, err := rand.Read(key)

		if err != nil {
			result=(map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})	
			return result
		}

		keysArr[i] = hex.EncodeToString(key)
	}
	// keysArray=keysArr
	// leveldb, err := leveldb.OpenFile("./db/device_info", nil)
	// if err != nil {
	// 	panic(err)
	// }
	barray,_:=json.Marshal(keysArr)
    var strbarray string
    strbarray=string(barray)
	
	callmap :=map[string]interface{}{
		"key":roomid,
		"data":strbarray,

	}
	//save key array to roomid in leveldb
	kq:= database.WriteValueStorage(callmap,c.DB)
	log.Info("write player keys to storage success:",kq["success"])
	fmt.Println("keysArr:",keysArr)
	result=(map[string]interface{}{
		"success": true,
		"message": keysArr,
	})	

	return result
}

// createDeck creates a standard deck of 52 cards
func (c *Client)CreateDeck() map[string]interface{} {
	const deckSize = 52
	result:=make(map[string]interface{})
	ranks := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	suits := []string{"S", "H", "D", "C"}

	deck := make([]string, deckSize)

	for i, suit := range suits {
		for j, rank := range ranks {
			card := rank + suit
			deck[i*len(ranks)+j] = card
		}
	}
	result=(map[string]interface{}{
		// "success": true,
		"deck": deck,
	})	
	fmt.Println("Create deck success")
	return result
}

// shuffleDeck shuffles the given deck of cards
func(c *Client)ShuffleDeck(call map[string]interface{})[]string {
	// result:=make(map[string]interface{})
	deck := call["deck"].([]string)
	random.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	fmt.Println("Shuffle deck success:",deck)
	return deck
}
// encryptDeck encrypts each card in the deck using the private keys of the players
func(c *Client) EncryptDeck(deck []string, arrmap []string, roomid string) []string {
	fmt.Println("begin encrypt deck ")
	encryptedDeck := make([]string, len(deck))

		for i, priKey := range arrmap{
			card:=deck[i]
			encryptedCard:= encryption(card,priKey)
			card =encryptedCard
			encryptedDeck[i] = card

		}
		// deckKq=encryptedDeck
		fmt.Println("EncryptDeck khi setDeck laf:",encryptedDeck)
		//save encryptDeck in leveldb
		barray,_:=json.Marshal(encryptedDeck)
		var strbarray string
		strbarray=string(barray)
		callmap :=map[string]interface{}{
			"key":roomid + "EncryptDeck",
			"data":strbarray,
	
		}
		result:= database.WriteValueStorage(callmap,c.DB)
		log.Info("write encrypted deck to storage success:",result["success"])
	return encryptedDeck
}
func createCipher(key string) cipher.Block {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatalf("Failed to create the AES cipher: %s", err)
	}
	return c
}
func encryption(plainText string,key string) string{
	bytes := []byte(plainText)
	blockCipher := createCipher(key)
	stream := cipher.NewCTR(blockCipher, constant.IV)
	// Buffer for storing decrypted data
	encryptedData := make([]byte, len(bytes))
	stream.XORKeyStream(encryptedData, bytes)
	result:=base64.StdEncoding.EncodeToString(encryptedData)
	return result
}
func(c *Client) DecryptDeck(call map[string]interface{}) (map[string]interface{},string)  {
	fmt.Println("-----------")
	result:=make(map[string]interface{})
	fmt.Println("cal encrytedDeck:",call["encrypted-cards"])
	fmt.Println("cal encrytedKeys:",call["encrypted-keys"])
	encrytedDeck,ok := call["encrypted-cards"].([]interface{})
	if !ok {
		msg:= "error input encrypted-cards"
		log.Error(msg)
		return nil,msg
	}
	playerKeys,ok := call["encrypted-keys"].([]interface{})
	if !ok {
		msg:= "error input encrypted-keys"
		log.Error(msg)
		return nil,msg
	}
	decryptedDeck := make([]string, len(encrytedDeck))
	for i, encryptedcard := range encrytedDeck {
		decryptedBlockBytes:= decryption(encryptedcard.(string),playerKeys[i].(string))
		encryptedcard =string(decryptedBlockBytes)
		decryptedDeck[i] = encryptedcard.(string)

	}
	fmt.Println("decryptedDeck:",decryptedDeck)
	fmt.Printf("decryptedDeck type:%T",decryptedDeck[0])
	fmt.Printf("decryptedDeck value:%v",decryptedDeck[0])

	result=(map[string]interface{}{
		"success": true,
		"message": decryptedDeck,
	})	
	return result,""
} 

func decryption(encrypted string,key string) []byte {
	bytes,err:=base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		fmt.Println("Error decoding the encrypted string:", err)
		return nil
	}
	blockCipher := createCipher(key)
	stream := cipher.NewCTR(blockCipher, constant.IV)
	stream.XORKeyStream(bytes, bytes)
	return bytes
}
