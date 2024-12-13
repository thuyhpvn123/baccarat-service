package client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	// "time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/meta-node-blockchain/meta-node-baccarat/internal/database"
	"github.com/meta-node-blockchain/meta-node/pkg/bls"
	cm "github.com/meta-node-blockchain/meta-node/pkg/common"
	log "github.com/sirupsen/logrus"
	// "golang.org/x/exp/rand"
)

/***
Example:
	testPubkey1 = common.FromHex("a2702ce6bbfb2e013935781bac50a0e168732bd957861e6fbf185d688c82ade34c9f33fead179decb5953b3382b061df")
	testSign1   = common.FromHex("a507c03ab7ebb69a4b3adc22a0347bb2466788e6a3baa174a62bd74cdff60dfd6d6ba9ec6237098f1ceef6013bfeff1d0c8be716266710e1493c422293a676e7f168007324a23435d4590896f97f8e3686cf0c280240b9406800c1cec6bafb5d")
	testHash1   = common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
	sign1 := Sign(cm.PrivateKeyFromBytes(common.FromHex(testSecret1)), testHash1.Bytes())
	fmt.Printf("Sign1: %v\n", common.Bytes2Hex(sign1.Bytes()))
***/

func (c *Client) VerifySign(call map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	hash, ok := call["hash"].(string)
	if !ok {
		result = (map[string]interface{}{
			"success": true,
			"data":    false,
		})
		return result
	}
	bhash := common.FromHex(hash)
	sign, ok := call["sign"].(string)
	if !ok {
		result = (map[string]interface{}{
			"success": true,
			"data":    false,
		})
		return result
	}
	bsign := common.FromHex(sign)

	pubKey, ok := call["pubKey"].(string)
	if !ok {
		result = (map[string]interface{}{
			"success": true,
			"data":    false,
		})
		return result
	}
	bpubKey := common.FromHex(pubKey)
	pubKeyCm := cm.PubkeyFromBytes(bpubKey)
	signCm := cm.SignFromBytes(bsign)
	success := bls.VerifySign(pubKeyCm, signCm, bhash)
	address := crypto.Keccak256(bpubKey)[12:]
	result = (map[string]interface{}{
		"success": true,
		"data":    success,
		"address": hex.EncodeToString(address),
	})
	log.Info("verifySign:", success)
	return result

}
func (c *Client) GetKeyForPlayer(roomNumber string,playerAdd string) {
	cardArr,err := c.getCards(roomNumber,playerAdd)
	if err != nil {
		log.Error("err:",err)
		c.sentToClient("get-key-for-player",err)
	}
	log.Info("Get card done")
	log.Info("encrypted-cards", cardArr)
	//
	mapdeckKq := map[string]interface{}{
        "key": roomNumber + "EncryptDeck",
	}
	kq := database.ReadValueStorage(mapdeckKq, c.DB)
	var deckKq []string
	if kq["value"] != nil {
		err := json.Unmarshal(kq["value"].([]byte),&deckKq)
		if err != nil {
			log.Error("err:",err)
			c.sentToClient("get-key-for-player",err)
		}
		log.Info("encrypted-deck", deckKq)
	}else{
		c.sentToClient("get-key-for-player", "no encrypt deck")

	}
	//
	mapkeysArray := map[string]interface{}{
        "key": roomNumber ,
	}
	kq = database.ReadValueStorage(mapkeysArray, c.DB)
	var keysArray []string
	if kq["value"] != nil {
		err := json.Unmarshal(kq["value"].([]byte),&keysArray)
		if err != nil {
			log.Error("err:",err)
		}
		fmt.Println("encrypted-keys", keysArray)
	}else{
		c.sentToClient("get-key-for-player", "no encrypt key")
	}
	//
	call := map[string]interface{}{
		"encrypted-cards": cardArr,
		"encrypted-deck":  deckKq,
		"encrypted-keys":  keysArray,
	}
	fmt.Println("call la:",call)
	keyArr := FindKeys(call)
	// indices := findIndices(cardArr.([]string), deckKq)
	// keyArr := findArray(keysArray, indices)

	c.sentToClient("get-key-for-player", keyArr)
}
func (c *Client) getCards(roomNumber string,playerAdd string) (interface{},error) {
	cards, err := c.Serv.GetCardsService(roomNumber,playerAdd)
	if err != nil {
		log.Error("Error GetCardService", err)	
	}
	return cards,err
}
// find array keys
func FindKeys(callMap map[string]interface{}) []string {
	fmt.Println("Find Keys")
	encryptedDeck,ok := callMap["encrypted-deck"].([]string)
	if !ok {
		fmt.Println("Error: encryptedDeck is not a []string or does not exist.")
		return []string{}
	}
	fmt.Println("encryptedDeck la:",callMap["encrypted-deck"])
	encryptedCardArr,ok := callMap["encrypted-cards"].([]string)
	if !ok {
		fmt.Println("Error: encryptedCardArr is not a []string or does not exist.")
		return []string{}
	}
	fmt.Println("encryptedCardArr la:",callMap["encrypted-cards"])

	encryptedKey := callMap["encrypted-keys"].([]string)
	if !ok {
		fmt.Println("Error: encryptedKey is not a []string or does not exist.")
		return []string{}
	}
	fmt.Println("encryptedKey la:",callMap["encrypted-keys"])

	indices := findIndices(encryptedDeck, encryptedCardArr)
	fmt.Println("indices la:",indices)
	result := findArray(encryptedKey, indices)
	return result
}

func findIndices(firstArray []string, secondArray []string) []int {
	fmt.Println("firstArray la:",firstArray)
	fmt.Println("secondArray la:",secondArray)
	indices := make([]int, len(secondArray))
	indexMap := make(map[string]bool)

	// Create a map of elements from the first array for efficient lookup
	for _, num := range firstArray {
		indexMap[num] = true
	}

	// Find the indices of elements from the second array in the first array
	for i, num := range secondArray {
		if indexMap[num] {
			indices[i] = findIndex(firstArray, num)
		} else {
			indices[i] = -1 // Element not found in the first array
		}
	}

	return indices
}

func findIndex(arr []string, target string) int {
	for i, num := range arr {
		if num == target {
			return i
		}
	}
	return -1 // Element not found
}
func findArray(firstArray []string, secondArray []int) []string {
	result := make([]string, len(secondArray))

	for i, num := range secondArray {
		if num >= 0 && num < len(firstArray) {
			result[i] = firstArray[num]
		} else {
			result[i] = "nil" // Invalid index, assign -1 as the value
		}
	}

	return result
}

func (c *Client) GetSign(callMap map[string]interface{}) cm.Sign {
	privateKey := callMap["privateKey"].(string)
	addressForSign := callMap["address"].(string)
	// Initialize the random number generator
	// rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 100
	// randomNumber := rand.Intn(101)

	message := common.FromHex(addressForSign + strconv.Itoa(0))
	//vd addressForSign= "0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2"
	keyPair := bls.NewKeyPair(common.FromHex(privateKey))
	//vd privateKey="36e1aa979f98c7154fb2491491ec044ccac099651209ccfbe2561746dbe29ebb"
	hash := crypto.Keccak256(message)
	prikey := keyPair.PrivateKey()

	sign := bls.Sign(prikey, hash)
	pubkey := keyPair.BytesPublicKey()
	add := keyPair.Address()
	address := crypto.Keccak256(pubkey)[12:]
	fmt.Println("sign:", sign)
	fmt.Println("pubkey:", hex.EncodeToString(pubkey))
	fmt.Println("address:", add)
	fmt.Println("hash:", hex.EncodeToString(hash))
	fmt.Println("address tu publickey:", hex.EncodeToString(address))

	call := map[string]interface{}{
		"sign": fmt.Sprint(sign),
		"hash":  hex.EncodeToString(hash),
	}
	go c.sentToClient("get-sign", call)
	return sign
}
/***
sign: a0ba415b563556862a0f379e4930e62f89178d080ed520825912c866a07917d3e08b1941fc0fca19fe6163894befcb8801161f71a2bb050e4e5beef9c422a4ef237b96d0d77830977a293536a2476021885a59b4891279877ff393995912f189
pubkey: 8d23b1e4bd15581b30633660b042d5107df085b57c2a12bccb0dde2e444f708397899110ee2dc1d98eea952c340f93e9
address: 0xE730d4572f20A4d701EBb80b8b5aFA99b36d5e49
hash: 15958c00decd9bf7ba1a448d27cf104a1ab5fa7e4c0339b691f9d35638af289d
address tu publickey: e730d4572f20a4d701ebb80b8b5afa99b36d5e49
***/

// func (c *Client) GetSign(callMap map[string]interface{}) cm.Sign {
// 	fmt.Println("111111111111111")
// 	privateKey := "068ab5318764afcaa90089c3bbe54bfab7c472dffc9cf4abf5ab3024b2111228"
// 	publickey:= "80fbb74c4d6a803ae42f3bf86c72756adb4a40eebfaab15fa4a6db4d90aaeca31006fac317987ec076cc047010918f29"
// 	nonce := 27
	
// 	message := []byte(publickey + strconv.Itoa(nonce))
// 	keyPair := bls.NewKeyPair(common.FromHex(privateKey))
// 	prikey := keyPair.PrivateKey()
  
// 	sign := bls.Sign(prikey, crypto.Keccak256(message))
// 	fmt.Println("sign:", sign)
  
// 	return sign
// }
