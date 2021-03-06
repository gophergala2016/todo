package main

import (
	"log"

	"github.com/boj/redistore"
	"github.com/zemirco/couchdb"
)

const ip = "192.168.99.100"

var (
	db        database
	rediStore *redistore.RediStore
)

func init() {
	var err error
	// init redis for sessions
	rediStore, err = redistore.NewRediStore(10, "tcp", ip+":6379", "", []byte("secret-key"))
	if err != nil {
		panic(err)
	}
	// init couchdb
	client, err := couchdb.NewClient("http://" + ip + ":5984/")
	if err != nil {
		panic(err)
	}
	log.Println(client.Info())
	db = database{
		Client: client,
	}
}
