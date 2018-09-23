package main

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Manager interface {
	handleHealthCheck() error
	handlePreStart() error
	handleChange() error
}

type singleMasterManager struct {
	redis   *redisInstance
	consul  *consulInstance
	session string
}

func newSingleMasterManager(redis *redisInstance, consul *consulInstance, sessionName string) (Manager, error) {
	return &singleMasterManager{
		redis:   redis,
		consul:  consul,
		session: sessionName,
	}, nil
}

func (m *singleMasterManager) handleHealthCheck() error {
	sessionId, err := m.consul.getSession(m.session)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve consul session")
	}

	if err := m.redis.isRedisReady(); err != nil {
		return errors.Wrap(err, "redis is not healthy")
	}

	if err := m.consul.renewSession(sessionId); err != nil {
		return errors.Wrap(err, "failed to renew consul session")
	}

	expectedMaster, err := m.consul.getMaster()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve master from consul")
	}

	isMaster, err := m.redis.isRedisMaster()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve local master status")
	}

	if expectedMaster == "" {
		log.Warnln("We somehow lost our master. Starting promoting process.")
		return m.evaluateMaster()
	}

	if expectedMaster != m.session && isMaster {
		return errors.Errorf("we are master but %q should be", expectedMaster)
	}

	if expectedMaster == m.session && !isMaster {
		return errors.New("We should be master but are not.")
	}

	return nil
}

func (m *singleMasterManager) handlePreStart() error {
	return m.evaluateMaster()
}

func (m *singleMasterManager) handleChange() error {
	return m.evaluateMaster()
}

func (m *singleMasterManager) evaluateMaster() error {
	sessionId, err := m.consul.getSession(m.session)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve consul session")
	}

	currentMaster, err := m.consul.getMaster()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve current global master")
	}

	if currentMaster != "" {
		if err := m.redis.makeSlave(currentMaster); err == nil {
			log.Infof("We are not enslaved to the master %v\n", currentMaster)
			return nil
		} else {
			return errors.Wrap(err, "failed to enslave to current master")
		}
		return nil
	}

	newMaster, err := m.consul.tryLockMaster(sessionId, m.session)
	if err != nil {
		return errors.Wrap(err, "failed to try to acquire lock")
	}

	if !newMaster {
		fmt.Println("failed to retrieve lock, trying again")
		return m.evaluateMaster()
	}

	m.redis.makeMaster()
	log.Infof("We promoted ourselves to master")

	return nil
}
