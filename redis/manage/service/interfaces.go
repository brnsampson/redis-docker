package service

// These are mostly notes for implementing the service.go definitions.
type ServiceMasterFinder interface {
	FindServiceMaster() error
}

type ServiceMasterAssigner interface {
	AssignServiceMaster() error
}

type ServiceMasterFinderAssigner interface {
	ServiceMasterFinder
	ServiceMasterAssigner
}

// Defining requirements for this package to work
type RemoteAddrsReader interface {
	ReadRemoteAddrs() ([]string, error)
}

type RemoteKeyReader interface {
	ReadRemoteKey(string) (string, error)
}

type RemoteServiceLocker interface {
	LockRemoteService() error
}

type RemoteServiceUnlocker interface {
	UnlockRemoteService() error
}

type RemoteServiceLockUnlocker interface {
	RemoteServiceLocker
	RemoteServiceUnlocker
}

type InstanceReadyChecker interface {
	IsInstanceReady() (bool, error)
}

type InstanceMasterReader interface {
	ReadInstanceMaster() (string, error)
}

type InstanceMasterUpdater interface {
	UpdateInstanceMaster(addr string) error
}

type InstanceMasterClaimer interface {
	ClaimInstanceMaster() error
}

type InstanceMasterReaderUpdater interface {
	InstanceMasterReader
	InstanceMasterUpdater
	InstanceMasterClaimer
}
