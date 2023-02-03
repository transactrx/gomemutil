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

func AddCacheValue(memKey string, value []byte, expirationInSecond int) error {

	item := &memcache.Item{
		Key:        fetchSeed() + memKey,
		Value:      value,
		Expiration: int32(expirationInSecond),
	}

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

func GetCachedObject(memKey string) ([]byte, error) {

	// Get the item from the memcache
	if item, err := getMemcachedClient().Get(fetchSeed() + memKey); err == memcache.ErrCacheMiss {
		logger.Print("item not in the cache")
		return nil, nil
	} else if err != nil {
		logger.Printf("error getting item: %v", err)
		return nil, err
	} else {
		return item.Value, nil
	}

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

		logger.Printf("******************************************* NATS CONFIG *******************************************")
		logger.Printf("***************************************************************************************************")
		logger.Printf("MEMCACHED_SERVERLIST: %s", serverList)
		logger.Printf("MEMCACHED_SEED: %s", fetchSeed())
		logger.Printf("***************************************************************************************************")
		logger.Printf("***************************************************************************************************")
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
