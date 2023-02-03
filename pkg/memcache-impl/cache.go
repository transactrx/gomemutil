package memcache_impl

import (
	"github.com/bradfitz/gomemcache/memcache"
	"log"
	"os"
	"strings"
)

var logger = log.New(os.Stdout, "memcached", log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
var memcachedClient *memcache.Client
var serverList string
var seed string

// Add writes the given item, if no value already exists for its key
func AddCacheValue(memKey string, value []byte, expirationInSecond int) error {

	item := &memcache.Item{
		Key:        getCacheKey(memKey),
		Value:      value,
		Expiration: int32(expirationInSecond),
	}
	logger.Printf("AddCacheValue key:%s, expirationInSecond:%d", memKey, expirationInSecond)
	// Add the item to the memcache, if the key does not already exist
	if err := getMemcachedClient().Add(item); err == memcache.ErrNotStored {
		logger.Printf("item with key %q already exists", item.Key)
		return nil
	} else if err != nil {
		logger.Printf("error adding item: %v", err)
		return err
	}

	return nil
}

// Get gets the item for the given key.
func GetCachedObject(memKey string) ([]byte, error) {
	logger.Printf("GetCachedObject key:%s ", memKey)
	if item, err := getMemcachedClient().Get(getCacheKey(memKey)); err == memcache.ErrCacheMiss {
		logger.Printf("item %s not in the cache", memKey)
		return nil, nil
	} else if err != nil {
		logger.Printf("error getting item: %v", err)
		return nil, err
	} else {
		return item.Value, nil
	}

}

// Delete deletes the item with the provided key.
func DeleteCachedObject(memKey string) error {
	logger.Printf("DeleteCachedObject key:%s ", memKey)
	// Get the item from the memcache
	if err := getMemcachedClient().Delete(getCacheKey(memKey)); err == memcache.ErrCacheMiss {
		logger.Printf("item %s not in the cache", memKey)
	} else if err != nil {
		logger.Printf("error deleting item: %v", err)
		return err
	}

	return nil
}

func getCacheKey(memKey string) string {
	return fetchSeed() + memKey
}

func fetchSeed() string {
	if (len(strings.TrimSpace(seed)) > 0) == false {
		seed = "defaultseed"
		confSeed := os.Getenv("MEMCACHED_SEED")
		if len(strings.TrimSpace(confSeed)) > 0 {
			seed = confSeed
		}
	}
	return seed
}

func getMemcachedClient() *memcache.Client {
	if memcachedClient == nil {
		serverList = os.Getenv("MEMCACHED_SERVERLIST")

		logger.Printf("******************************************* MEMCACHED CONFIG *******************************************")
		logger.Printf("********************************************************************************************************")
		logger.Printf("MEMCACHED_SERVERLIST: %s", serverList)
		logger.Printf("MEMCACHED_SEED: %s", fetchSeed())
		logger.Printf("********************************************************************************************************")
		logger.Printf("********************************************************************************************************")
		memcachedClient = memcache.New(getMemcachedServerList(serverList)...)
	}

	return memcachedClient
}

func getMemcachedServerList(conn string) []string {

	serverlist := []string{}
	var serverSplited []string = strings.Split(conn, ";")

	for _, v := range serverSplited {
		if len(v) > 0 {
			serverlist = append(serverlist, v)
		}
	}

	if len(serverlist) == 0 {
		serverlist = append(serverlist, "127.0.0.1:11211")
	}

	return serverlist
}
