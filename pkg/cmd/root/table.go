package root

import (
	buildCmd "github.com/aziontech/azion-cli/pkg/cmd/build"
	"github.com/aziontech/azion-cli/pkg/cmd/clone"
	"github.com/aziontech/azion-cli/pkg/cmd/completion"
	"github.com/aziontech/azion-cli/pkg/cmd/config"
	"github.com/aziontech/azion-cli/pkg/cmd/create"
	"github.com/aziontech/azion-cli/pkg/cmd/delete"
	deploycmd "github.com/aziontech/azion-cli/pkg/cmd/deploy"
	"github.com/aziontech/azion-cli/pkg/cmd/describe"
	devcmd "github.com/aziontech/azion-cli/pkg/cmd/dev"
	initcmd "github.com/aziontech/azion-cli/pkg/cmd/init"
	linkcmd "github.com/aziontech/azion-cli/pkg/cmd/link"
	"github.com/aziontech/azion-cli/pkg/cmd/list"
	"github.com/aziontech/azion-cli/pkg/cmd/login"
	"github.com/aziontech/azion-cli/pkg/cmd/logout"
	logcmd "github.com/aziontech/azion-cli/pkg/cmd/logs"
	"github.com/aziontech/azion-cli/pkg/cmd/profiles"
	"github.com/aziontech/azion-cli/pkg/cmd/purge"
	"github.com/aziontech/azion-cli/pkg/cmd/reset"
	"github.com/aziontech/azion-cli/pkg/cmd/rollback"
	"github.com/aziontech/azion-cli/pkg/cmd/sync"
	"github.com/aziontech/azion-cli/pkg/cmd/unlink"
	"github.com/aziontech/azion-cli/pkg/cmd/update"
	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmd/warmup"
	"github.com/aziontech/azion-cli/pkg/cmd/whoami"
	"github.com/aziontech/azion-cli/pkg/cmdregistry"
	v3buildCmd "github.com/aziontech/azion-cli/pkg/v3commands/build"
	v3deploycmd "github.com/aziontech/azion-cli/pkg/v3commands/deploy"
	v3devcmd "github.com/aziontech/azion-cli/pkg/v3commands/dev"
	v3initcmd "github.com/aziontech/azion-cli/pkg/v3commands/init"
	v3linkcmd "github.com/aziontech/azion-cli/pkg/v3commands/link"
	v3login "github.com/aziontech/azion-cli/pkg/v3commands/login"
	v3purge "github.com/aziontech/azion-cli/pkg/v3commands/purge"
	v3sync "github.com/aziontech/azion-cli/pkg/v3commands/sync"
	v3unlink "github.com/aziontech/azion-cli/pkg/v3commands/unlink"
)

// children lists every top-level command, with the API generation it belongs
// to. The root builds this list once and lets the table decide what the
// account's generation sees, instead of keeping one registration function per
// generation.
//
// The five verb dispatchers are version-agnostic here: each one is a single
// shell that reads its own table (pkg/cmd/<verb>/table.go) to pick its
// children. The commands still listed per generation are the ones whose whole
// implementation differs — the deploy pipeline and the project-scaffolding
// commands around it, which phase 3 addresses.
var children = []cmdregistry.Entry{
	cmdregistry.Any(create.NewCmd),
	cmdregistry.Any(delete.NewCmd),
	cmdregistry.Any(describe.NewCmd),
	cmdregistry.Any(list.NewCmd),
	cmdregistry.Any(update.NewCmd),

	cmdregistry.Any(completion.NewCmd),
	cmdregistry.Any(logcmd.NewCmd),
	cmdregistry.Any(logout.NewCmd),
	cmdregistry.Any(profiles.NewCmd),
	cmdregistry.Any(reset.NewCmd),
	cmdregistry.Any(rollback.NewCmd),
	cmdregistry.Any(version.NewCmd),
	cmdregistry.Any(whoami.NewCmd),

	cmdregistry.V4Only(buildCmd.NewCmd),
	cmdregistry.V4Only(deploycmd.NewCmd),
	cmdregistry.V4Only(devcmd.NewCmd),
	cmdregistry.V4Only(initcmd.NewCmd),
	cmdregistry.V4Only(linkcmd.NewCmd),
	cmdregistry.V4Only(login.New),
	cmdregistry.V4Only(purge.NewCmd),
	cmdregistry.V4Only(sync.NewCmd),
	cmdregistry.V4Only(unlink.NewCmd),

	// Offered only on v4: these have no v3 counterpart.
	cmdregistry.V4Only(clone.NewCmd),
	cmdregistry.V4Only(config.NewCmd),
	cmdregistry.V4Only(warmup.NewCmd),

	cmdregistry.V3Only(v3buildCmd.NewCmd),
	cmdregistry.V3Only(v3deploycmd.NewCmd),
	cmdregistry.V3Only(v3devcmd.NewCmd),
	cmdregistry.V3Only(v3initcmd.NewCmd),
	cmdregistry.V3Only(v3linkcmd.NewCmd),
	cmdregistry.V3Only(v3login.New),
	cmdregistry.V3Only(v3purge.NewCmd),
	cmdregistry.V3Only(v3sync.NewCmd),
	cmdregistry.V3Only(v3unlink.NewCmd),
}
