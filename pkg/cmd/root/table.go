package root

import (
	buildCmd "github.com/aziontech/azion-cli/pkg/cmd/build"
	v3buildCmd "github.com/aziontech/azion-cli/pkg/cmd/build/v3"
	"github.com/aziontech/azion-cli/pkg/cmd/clone"
	"github.com/aziontech/azion-cli/pkg/cmd/completion"
	"github.com/aziontech/azion-cli/pkg/cmd/config"
	"github.com/aziontech/azion-cli/pkg/cmd/create"
	"github.com/aziontech/azion-cli/pkg/cmd/delete"
	deploycmd "github.com/aziontech/azion-cli/pkg/cmd/deploy"
	v3deploycmd "github.com/aziontech/azion-cli/pkg/cmd/deploy/v3"
	v6deploycmd "github.com/aziontech/azion-cli/pkg/cmd/deploy/v6"
	"github.com/aziontech/azion-cli/pkg/cmd/describe"
	devcmd "github.com/aziontech/azion-cli/pkg/cmd/dev"
	v3devcmd "github.com/aziontech/azion-cli/pkg/cmd/dev/v3"
	initcmd "github.com/aziontech/azion-cli/pkg/cmd/init"
	v3initcmd "github.com/aziontech/azion-cli/pkg/cmd/init/v3"
	linkcmd "github.com/aziontech/azion-cli/pkg/cmd/link"
	v3linkcmd "github.com/aziontech/azion-cli/pkg/cmd/link/v3"
	"github.com/aziontech/azion-cli/pkg/cmd/list"
	"github.com/aziontech/azion-cli/pkg/cmd/login"
	v3login "github.com/aziontech/azion-cli/pkg/cmd/login/v3"
	"github.com/aziontech/azion-cli/pkg/cmd/logout"
	logcmd "github.com/aziontech/azion-cli/pkg/cmd/logs"
	"github.com/aziontech/azion-cli/pkg/cmd/profiles"
	"github.com/aziontech/azion-cli/pkg/cmd/purge"
	v3purge "github.com/aziontech/azion-cli/pkg/cmd/purge/v3"
	"github.com/aziontech/azion-cli/pkg/cmd/reset"
	"github.com/aziontech/azion-cli/pkg/cmd/rollback"
	"github.com/aziontech/azion-cli/pkg/cmd/sync"
	v3sync "github.com/aziontech/azion-cli/pkg/cmd/sync/v3"
	"github.com/aziontech/azion-cli/pkg/cmd/unlink"
	v3unlink "github.com/aziontech/azion-cli/pkg/cmd/unlink/v3"
	"github.com/aziontech/azion-cli/pkg/cmd/update"
	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmd/warmup"
	"github.com/aziontech/azion-cli/pkg/cmd/whoami"
	"github.com/aziontech/azion-cli/pkg/cmdregistry"
)

// children lists every top-level command, with the API generation it belongs
// to. The root builds this list once and lets the table decide what the
// account's generation sees, instead of keeping one registration function per
// generation.
//
// v6 accounts get a short list for now: the project commands (init, link,
// build, dev, deploy) and the ones that touch no resource API (login, logout,
// profiles, version, whoami, completion, reset). The resource commands, the
// CRUD verbs among them, stay on v3 and v4 until they have a v6
// implementation.
var children = []cmdregistry.Entry{
	// Every generation.
	cmdregistry.Any(completion.NewCmd),
	cmdregistry.Any(logout.NewCmd),
	cmdregistry.Any(profiles.NewCmd),
	cmdregistry.Any(reset.NewCmd),
	cmdregistry.Any(version.NewCmd),
	cmdregistry.Any(whoami.NewCmd),

	// v3 and v4. The five verb dispatchers are one shell each that reads the
	// children table in its own file.
	cmdregistry.V3V4(create.NewCmd),
	cmdregistry.V3V4(delete.NewCmd),
	cmdregistry.V3V4(describe.NewCmd),
	cmdregistry.V3V4(list.NewCmd),
	cmdregistry.V3V4(update.NewCmd),
	cmdregistry.V3V4(logcmd.NewCmd),
	cmdregistry.V3V4(rollback.NewCmd),

	// v4 and v6: one implementation that reads the generation from the factory.
	cmdregistry.V4V6(buildCmd.NewCmd),
	cmdregistry.V4V6(devcmd.NewCmd),
	cmdregistry.V4V6(initcmd.NewCmd),
	cmdregistry.V4V6(linkcmd.NewCmd),
	cmdregistry.V4V6(login.New),

	cmdregistry.V4Only(deploycmd.NewCmd),
	cmdregistry.V4Only(purge.NewCmd),
	cmdregistry.V4Only(sync.NewCmd),
	cmdregistry.V4Only(unlink.NewCmd),
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

	cmdregistry.V6Only(v6deploycmd.NewCmd),
}
