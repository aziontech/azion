package manifest

import (
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

	conf := &contracts.AzionApplicationOptionsV3{
		Application: contracts.AzionJsonDataApplication{ID: 123},
		// SkipDeletion is intentionally left as nil to simulate absence in JSON
	}

	rc := newResourceContext(f, conf, &contracts.Manifest{}, "", &msgs, nil)
	if err := rc.deleteOrphanedResources(); err != nil {
		t.Fatalf("orphan removal (v3) failed with SkipDeletion absent: %v", err)
	}
}
