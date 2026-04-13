package exec

type SessionTarget struct {
	ClusterID     uint64
	Namespace     string
	PodName       string
	ContainerName string
}
