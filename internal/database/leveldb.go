package database

import (
	"log"
	"github.com/syndtr/goleveldb/leveldb"
)
func Open( path string) (*leveldb.DB,error){
	db, err := leveldb.OpenFile(path, nil)
    if err != nil {
        return nil, err
    }
    // defer db.Close()
	return db,err
}
func WriteValueStorage(call map[string]interface{},db *leveldb.DB) map[string]interface{} {
	key, ok := call["key"].(string)
	data, ok := call["data"].(string)
	if !ok || key == "" || data == "" {
	
	return map[string]interface{}{
	"error": "EWVS-001",
	"message": "Key and data is required",
	}
	}

	// securyDb.Write(key, data)
	err := db.Put([]byte(key), []byte(data), nil)
	if err!=nil{
		log.Fatal()
	}
	return map[string]interface{}{
	"success": true,
	}
}
 func ReadValueStorage(call map[string]interface{},db *leveldb.DB) map[string]interface{} {
	key, ok := call["key"].(string)
	if !ok || key == "" {
		return map[string]interface{}{
			"success": false,
			"msg": "Key not found",
		}
	} else {
		value,err := db.Get([]byte(key), nil)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"msg": "",
			}
		
		} else {
			return map[string]interface{}{
				"value": value,
				"success": true,
				"msg": "",
			}
			
		}
	}
}
func DeleteKeyStorage(call map[string]interface{},db *leveldb.DB) map[string]interface{} {
	key, ok := call["key"].(string)
	if !ok || key == "" {
	return map[string]interface{}{
	"error": "EDKS-001",
	"message": "Key is required",
	}
	}
	// securyDb.Delete(key)
	err := db.Delete([]byte(key), nil)
	if err!=nil{
		log.Fatal()
	}
	return map[string]interface{}{
	"success": true,
	}
}