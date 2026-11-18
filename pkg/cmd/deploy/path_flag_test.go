package deploy

import (
	"testing"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPathFlagIsDeprecated(t *testing.T) {
	c := NewCobraCmd(NewDeployCmd(&cmdutil.Factory{}))
	fl := c.Flags().Lookup("path")
	require.NotNil(t, fl)
	assert.Equal(t, msg.DeployPathFlagDeprecated, fl.Deprecated)
	assert.True(t, fl.Hidden)
	assert.NotContains(t, c.Example, "--path")
}
