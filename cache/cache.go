package cache

import (
	"fmt"
	"tecace/googleapi"
)

var keys = make(map[string]bool)

func SyncKeysFromSheet() {

	data, err := googleapi.Get()
	if err != nil {
		return
	}

	for k := range data {
		fmt.Printf("Loading key: %s\n", k)
		keys[k] = true
	}
}

func AddKey(key string) {
	keys[key] = true
}

func HasKey(key string) bool {
	_, exists := keys[key]
	return exists
}

func RemoveKey(key string) {
	delete(keys, key)
}
