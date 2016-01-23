package main

import "github.com/boj/redistore"

var (
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
}
