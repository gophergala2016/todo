package main

import (
	"log"

	"github.com/boj/redistore"
	"github.com/zemirco/couchdb"
)

var (
	db        Database
	RediStore *redistore.RediStore
)

func init() {
	var err error
	// init redis for sessions
	// fix hard coded ip to docker machine
	RediStore, err = redistore.NewRediStore(10, "tcp", "192.168.99.100:6379", "", []byte("secret-key"))
	if err != nil {
		panic(err)
	}
	// init couchdb
	client, err := couchdb.NewClient("http://192.168.99.100:5984/")
	if err != nil {
		panic(err)
	}
	log.Println(client.Info())
	db = Database{
		Client: client,
	}
}
