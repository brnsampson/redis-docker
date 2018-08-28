package remote

type RemoteShimmer interface {
	Lock(string) error
	Unlock(string) error
	Readkey(string) (string, error)
	ReadRemoteNodes(string) (*RemoteNodes, error)
}
