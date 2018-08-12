package main

import (
	"github.com/hashicorp/consul/api"
)

func isConsulReady() (bool, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return false, err
	}

	// Test that we can set and read a test key/value
	kv := client.KV()
	p := &api.KVPair{Key: "TestKey", Value: []byte("bla")}
	_, err = kv.Put(p, nil)

	if err != nil {
		return false, err
	}

	// Do lookup
	_, _, err = kv.Get("TestKey", nil)
	if err != nil {
		return false, err
	}
	return true, nil
}
