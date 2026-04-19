package backuprestore

import "fmt"

func AssessRPORTO(targetRPO, targetRTO, actualRPO, actualRTO int) (string, string) {
	if actualRPO <= targetRPO && actualRTO <= targetRTO {
		return "RPO/RTO 目标均已达成", ""
	}
	return "RPO/RTO 目标未完全达成", fmt.Sprintf("目标 RPO=%d, 实际 RPO=%d；目标 RTO=%d, 实际 RTO=%d", targetRPO, actualRPO, targetRTO, actualRTO)
}
