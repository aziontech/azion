package create

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/create"
	"github.com/aziontech/azion-cli/pkg/apiversion"
	origin "github.com/aziontech/azion-cli/pkg/cmd/create/origin"
	personalToken "github.com/aziontech/azion-cli/pkg/cmd/create/personal_token"
	profile "github.com/aziontech/azion-cli/pkg/cmd/create/profile"
	v3CacheSetting "github.com/aziontech/azion-cli/pkg/cmd/create/v3/cache_setting"
	v3Domain "github.com/aziontech/azion-cli/pkg/cmd/create/v3/domain"
	v3EdgeApplications "github.com/aziontech/azion-cli/pkg/cmd/create/v3/edge_applications"
	v3EdgeFunction "github.com/aziontech/azion-cli/pkg/cmd/create/v3/edge_function"
	v3EdgeStorage "github.com/aziontech/azion-cli/pkg/cmd/create/v3/edge_storage"
	v3RulesEngine "github.com/aziontech/azion-cli/pkg/cmd/create/v3/rules_engine"
	v4Applications "github.com/aziontech/azion-cli/pkg/cmd/create/v4/applications"
	v4CacheSetting "github.com/aziontech/azion-cli/pkg/cmd/create/v4/cache_setting"
	v4Connector "github.com/aziontech/azion-cli/pkg/cmd/create/v4/connector"
	v4Crl "github.com/aziontech/azion-cli/pkg/cmd/create/v4/crl"
	v4Csr "github.com/aziontech/azion-cli/pkg/cmd/create/v4/csr"
	v4CustomPages "github.com/aziontech/azion-cli/pkg/cmd/create/v4/custom_pages"
	v4DataStream "github.com/aziontech/azion-cli/pkg/cmd/create/v4/data_stream"
	v4DeviceGroups "github.com/aziontech/azion-cli/pkg/cmd/create/v4/device_groups"
	v4DigitalCertificate "github.com/aziontech/azion-cli/pkg/cmd/create/v4/digital_certificate"
	v4DnsRecord "github.com/aziontech/azion-cli/pkg/cmd/create/v4/dns_record"
	v4DnsZone "github.com/aziontech/azion-cli/pkg/cmd/create/v4/dns_zone"
	// v4Domain "github.com/aziontech/azion-cli/pkg/cmd/create/v4/domain"
	v4Firewall "github.com/aziontech/azion-cli/pkg/cmd/create/v4/firewall"
	v4FirewallInstance "github.com/aziontech/azion-cli/pkg/cmd/create/v4/firewall_instance"
	v4FirewallRules "github.com/aziontech/azion-cli/pkg/cmd/create/v4/firewall_rules"
	v4Function "github.com/aziontech/azion-cli/pkg/cmd/create/v4/function"
	v4FunctionInstance "github.com/aziontech/azion-cli/pkg/cmd/create/v4/function_instance"
	// v4Kv "github.com/aziontech/azion-cli/pkg/cmd/create/v4/kv"
	v4NetworkList "github.com/aziontech/azion-cli/pkg/cmd/create/v4/network_list"
	v4RulesEngine "github.com/aziontech/azion-cli/pkg/cmd/create/v4/rules_engine"
	v4Storage "github.com/aziontech/azion-cli/pkg/cmd/create/v4/storage"
	v4Waf "github.com/aziontech/azion-cli/pkg/cmd/create/v4/waf"
	v4WafExceptions "github.com/aziontech/azion-cli/pkg/cmd/create/v4/waf_exceptions"
	v4WorkloadDeployment "github.com/aziontech/azion-cli/pkg/cmd/create/v4/workload_deployment"
	v4Workloads "github.com/aziontech/azion-cli/pkg/cmd/create/v4/workloads"
	variables "github.com/aziontech/azion-cli/pkg/cmd/create/variables"
	"github.com/aziontech/azion-cli/pkg/cmdregistry"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmd builds the `azion create` dispatcher. One shell serves every API
// generation: which sub-commands it offers, and which examples it shows, are
// read from the children and examples tables below, filtered by the
// generation carried on the factory.
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     msg.Usage,
		Short:   msg.ShortDescription,
		Long:    msg.LongDescription,
		Example: heredoc.Doc(cmdregistry.Example(examples, f)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(cmdregistry.Children(children, f)...)

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}

// children lists every sub-command `azion create` can offer. Entries carry the
// API generation they belong to, so the single dispatcher above serves both
// generations from this one table.
//
// Commented entries are resources whose command exists but is not yet wired up.
var children = []cmdregistry.Entry{
	cmdregistry.V4Only(v4Applications.NewCmd),
	cmdregistry.V4Only(v4RulesEngine.NewCmd),
	cmdregistry.Any(personalToken.NewCmd),
	cmdregistry.Any(origin.NewCmd),
	cmdregistry.V4Only(v4CacheSetting.NewCmd),
	cmdregistry.V4Only(v4DeviceGroups.NewCmd),
	cmdregistry.V4Only(v4DnsZone.NewCmd),
	cmdregistry.V4Only(v4DnsRecord.NewCmd),
	cmdregistry.V4Only(v4Function.NewCmd),
	cmdregistry.Any(variables.NewCmd),
	cmdregistry.V4Only(v4Storage.NewCmd),
	cmdregistry.V4Only(v4Workloads.NewCmd),
	cmdregistry.V4Only(v4WorkloadDeployment.NewCmd),
	cmdregistry.V4Only(v4Connector.NewCmd),
	cmdregistry.V4Only(v4CustomPages.NewCmd),
	cmdregistry.V4Only(v4DataStream.NewCmd),
	cmdregistry.V4Only(v4FunctionInstance.NewCmd),
	cmdregistry.Any(profile.NewCmd),
	cmdregistry.V4Only(v4NetworkList.NewCmd),
	// cmdregistry.V4Only(v4Kv.NewCmd),
	cmdregistry.V4Only(v4Firewall.NewCmd),
	cmdregistry.V4Only(v4FirewallInstance.NewCmd),
	cmdregistry.V4Only(v4FirewallRules.NewCmd),
	cmdregistry.V4Only(v4Waf.NewCmd),
	cmdregistry.V4Only(v4WafExceptions.NewCmd),
	cmdregistry.V4Only(v4DigitalCertificate.NewCmd),
	cmdregistry.V4Only(v4Csr.NewCmd),
	cmdregistry.V4Only(v4Crl.NewCmd),
	// cmdregistry.V4Only(v4Domain.NewCmd),
	cmdregistry.V3Only(v3EdgeApplications.NewCmd),
	cmdregistry.V3Only(v3RulesEngine.NewCmd),
	cmdregistry.V3Only(v3Domain.NewCmd),
	cmdregistry.V3Only(v3CacheSetting.NewCmd),
	cmdregistry.V3Only(v3EdgeFunction.NewCmd),
	cmdregistry.V3Only(v3EdgeStorage.NewCmd),
}

// examples is the curated example block for `azion create --help`. It is a
// hand-picked highlight rather than the full child list, and it differs per
// generation because the resources do, so it is data here like the children.
var examples = map[apiversion.Version]string{
	apiversion.V4: `
		$ azion create --help
		$ azion create application -h
		$ azion create connector -h
		$ azion create workload -h
		$ azion create network-list -h
		$ azion create firewall -h		
        `,
	apiversion.V3: `
		$ azion create --help
		$ azion create edge-application -h
		$ azion create rules-engine -h
		$ azion create domain -h
		$ azion create personal-token -h
		$ azion create origin -h
		$ azion create cache-setting -h
		$ azion create edge-function -h
		$ azion create variables -h
		$ azion create edge-storage -h
		$ azion create profile -h
        `,
}
