package main

import (
	"io/ioutil"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type consulInstance struct {
	client *api.Client

	cacheFilePath string
}

const LOCK_KEY = "REDIS_LOCK_KEY"

func newConsulInstance(cacheFile string) (*consulInstance, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &consulInstance{
		client:        client,
		cacheFilePath: cacheFile,
	}, nil
}

func (c *consulInstance) getSession(sessionName string) (string, error) {
	content, err := ioutil.ReadFile(c.cacheFilePath)
	if err == nil && len(content) > 0 {
		log.Debugf("Retrieved consul session id from cache file: %v\n", string(content))
		return string(content), nil
	}

	session, _, err := c.client.Session().Create(&api.SessionEntry{
		Name: sessionName,
		TTL:  "30s",
	}, nil)

	// cache session id for future invocations
	ioutil.WriteFile(c.cacheFilePath, []byte(session), 0666)

	log.Debugf("Created new consul session: %v\n", session)

	return session, nil
}

func (c *consulInstance) renewSession(id string) error {
	_, _, err := c.client.Session().Renew(id, nil)
	log.Debug("Renewed consul session.")
	return err
}

func (c *consulInstance) isConsulReady() error {
	// Test that we can set and read a test key/value
	kv := c.client.KV()
	p := &api.KVPair{Key: "TestKey", Value: []byte("bla")}
	_, err := kv.Put(p, nil)

	if err != nil {
		return err
	}

	// Do lookup
	_, _, err = kv.Get("TestKey", nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *consulInstance) getMaster() (string, error) {
	pair, _, err := c.client.KV().Get(LOCK_KEY, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve master")
	}

	if pair != nil{
		if pair.Session == "" {
			log.Infoln("Found LOCK_KEY without lock but with value.")
			return "", nil
		}

		return string(pair.Value), nil
	}

	return "", nil
}

func (c *consulInstance) tryLockMaster(sessionId, redisAddress string) (bool, error) {
	acquired, _, err := c.client.KV().Acquire(&api.KVPair{
		Key:     LOCK_KEY,
		Session: sessionId,
		Value:   []byte(redisAddress),
	}, nil)

	return acquired, err
}
