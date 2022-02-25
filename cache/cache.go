package cache

import (
	"tecace/googleapi"
)

var keys = syncKeysFromSheet()

func syncKeysFromSheet() map[string]bool {
	result := make(map[string]bool)
	data, err := googleapi.Get()
	if err != nil {
		return result
	}

	for k := range data {
		result[k] = true
	}
	return result
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
