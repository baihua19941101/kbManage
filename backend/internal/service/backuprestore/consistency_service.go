package backuprestore

func BuildConsistencyNotice(jobType string) string {
	switch jobType {
	case "cross-cluster-restore":
		return "跨集群恢复需要额外确认网络、存储和身份材料已预置"
	case "environment-migration":
		return "环境迁移前需要人工确认切换窗口和回滚步骤"
	default:
		return "恢复结果需结合业务探针与关键事务校验"
	}
}
