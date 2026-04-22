package enterprise

import "context"

type DeliveryCatalogInput struct {
	BundleName  string
	Environment string
}

type DeliveryCatalog interface {
	BuildArtifactSummary(ctx context.Context, input DeliveryCatalogInput) string
}

type StaticDeliveryCatalog struct{}

func NewStaticDeliveryCatalog() *StaticDeliveryCatalog { return &StaticDeliveryCatalog{} }

func (c *StaticDeliveryCatalog) BuildArtifactSummary(_ context.Context, input DeliveryCatalogInput) string {
	return input.BundleName + "@" + input.Environment + " 的交付材料已汇总"
}
