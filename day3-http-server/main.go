package main

import (
	"Pigcache/day3-http-server/pigcache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string {
	"Tom": "630",
	"Jack": "589",
	"Sam": "567",
}

func main()  {
	pigcache.NewGroup("scores", 2<<10, pigcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := pigcache.NewHTTPPool(addr)
	log.Println("pigcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
