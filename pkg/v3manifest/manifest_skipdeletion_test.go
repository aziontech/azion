package manifest

import (
	"context"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/testutils"
)

func TestDeleteResources_SkipDeletionAbsent_V3(t *testing.T) {
	// Empty global maps so no deletions are attempted
	CacheIds = map[string]int64{}
	RuleIds = map[string]contracts.RuleIdsStruct{}
	OriginKeys = map[string]string{}
	OriginIds = map[string]int64{}

	f, _, _ := testutils.NewFactory(nil)
	msgs := []string{}
	ctx := context.Background()

	conf := &contracts.AzionApplicationOptionsV3{
		Application: contracts.AzionJsonDataApplication{ID: 123},
		// SkipDeletion is intentionally left as nil to simulate absence in JSON
	}

	if err := deleteResources(ctx, f, conf, &msgs); err != nil {
		t.Fatalf("deleteResources (v3) failed with SkipDeletion absent: %v", err)
	}
}
