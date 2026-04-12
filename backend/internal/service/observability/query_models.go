package observability

type LogsQueryRequest struct {
	ClusterIDs []uint64
	Namespace  string
	Workload   string
	Pod        string
	Container  string
	Keyword    string
	StartAt    string
	EndAt      string
	Limit      int
}

type EventsQueryRequest struct {
	ClusterID    string
	Namespace    string
	ResourceKind string
	ResourceName string
	EventType    string
	StartAt      string
	EndAt        string
}

type MetricQueryRequest struct {
	ClusterIDs  []uint64
	SubjectType string
	SubjectRef  string
	MetricKey   string
	StartAt     string
	EndAt       string
	Step        string
}

type ResourceContextQuery struct {
	ClusterID    string
	Namespace    string
	ResourceKind string
	ResourceName string
	Keyword      string
	StartAt      string
	EndAt        string
}
