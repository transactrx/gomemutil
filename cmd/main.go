package main

import (
	"fmt"
	memcache_impl "github.com/transactrx/gomemutil/pkg/memcache-impl"
)

func main() {

	memKey := "test1"
	value := "peterlaanguila"
	expirationInSecond := 60

	memcache_impl.AddCacheValue(memKey, []byte(value), expirationInSecond)
	memValue, err := memcache_impl.GetCachedObject(memKey)
	if err != nil {
		fmt.Printf("se formo %v", err)
	}

	fmt.Printf("Results key:%s, value:%s", memKey, string(memValue))
}
