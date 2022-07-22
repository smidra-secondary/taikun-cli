package info

import (
	"github.com/itera-io/taikun-cli/cmd/cmderr"
	"github.com/itera-io/taikun-cli/cmd/cmdutils"
	"github.com/itera-io/taikun-cli/utils/out"
	"github.com/itera-io/taikun-cli/utils/out/field"
	"github.com/itera-io/taikun-cli/utils/out/fields"
	"github.com/itera-io/taikun-cli/utils/types"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/servers"
	"github.com/spf13/cobra"
)

var infoFields = fields.New(
	[]*field.Field{
		field.NewHidden(
			"ID", "projectId",
		),
		field.NewVisible(
			"NAME", "projectName",
		),
		field.NewVisible(
			"ORG", "organizationName",
		),
		field.NewVisible(
			"ORG-ID", "organizationId",
		),
		field.NewVisibleWithToStringFunc(
			"HEALTH", "projectHealth", out.FormatProjectHealth,
		),
		field.NewVisible(
			"STATUS", "projectStatus",
		),
		field.NewVisible(
			"ACCESS-PROFILE", "accessProfileName",
		),
		field.NewVisible(
			"ACCESS-PROFILE-ID", "accessProfileId",
		),
		field.NewVisibleWithToStringFunc(
			"CLOUD", "cloudType", out.FormatCloudType,
		),
		field.NewVisible(
			"CLOUD-CREDENTIAL", "cloudName",
		),
		field.NewVisible(
			"CLOUD-CREDENTIAL-ID", "cloudId",
		),
		field.NewVisible(
			"HAS-ALERTING-PROFILE", "hasAlertingProfile",
		),
		field.NewVisible(
			"AUTO-UPGRADES", "isAutoUpgrade",
		),
		field.NewVisible(
			"UPGRADABLE", "hasNextVersion",
		),
		field.NewVisible(
			"HAS-FLAVORS", "hasSelectedFlavors",
		),
		field.NewVisible(
			"BACKUPS", "isBackupEnabled",
		),
		field.NewVisible(
			"MAINTENANCE", "isMaintenanceModeEnabled",
		),
		field.NewVisible(
			"MONITORING", "isMonitoringEnabled",
		),
		field.NewVisible(
			"POLICY-PROFILE", "opaProfileName",
		),
		field.NewVisible(
			"POLICY-PROFILE-ID", "opaProfileId",
		),
		field.NewVisible(
			"HAS-KUBECONFIG", "hasKubeConfigFile",
		),
		field.NewVisible(
			"K8S-VERSION", "kubernetesCurrentVersion",
		),
		field.NewHidden(
			"KUBESPRAY-VERSION", "kubeCurrentVersion",
		),
		field.NewVisible(
			"K8S-PROFILE-ID", "kubernetesProfileId",
		),
		field.NewVisible(
			"K8S-PROFILE", "kubernetesProfileName",
		),
		field.NewVisible(
			"QUOTA-ID", "quotaId",
		),
		field.NewVisible(
			"REVISIONS", "projectRevision",
		),
		field.NewVisible(
			"SERVERS", "totalCount",
		),
		field.NewVisible(
			"BASTIONS", "bastion",
		),
		field.NewVisible(
			"KUBEMASTERS", "master",
		),
		field.NewVisible(
			"KUBEWORKERS", "worker",
		),
		field.NewVisible(
			"TOTAL-CPU", "usedCpu",
		),
		field.NewVisibleWithToStringFunc(
			"TOTAL-DISK", "usedDiskSize", out.FormatBToGiB,
		),
		field.NewVisibleWithToStringFunc(
			"TOTAL-RAM", "usedRam", out.FormatBToGiB,
		),
		field.NewVisibleWithToStringFunc(
			"LOCK", "isLocked", out.FormatLockStatus,
		),
	},
)

type InfoOptions struct {
	ProjectID int32
}

func NewCmdInfo() *cobra.Command {
	var opts InfoOptions

	cmd := cobra.Command{
		Use:   "info <project-id>",
		Short: "Get detailed information on a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.ProjectID, err = types.Atoi32(args[0])
			if err != nil {
				return cmderr.ErrIDArgumentNotANumber
			}
			return infoRun(&opts)
		},
	}

	cmdutils.AddColumnsFlag(&cmd, infoFields)

	return &cmd
}

func infoRun(opts *InfoOptions) (err error) {
	apiClient, err := taikungoclient.NewClient()
	if err != nil {
		return
	}

	params := servers.NewServersDetailsParams().WithV(taikungoclient.Version)
	params = params.WithProjectID(opts.ProjectID)

	response, err := apiClient.Client.Servers.ServersDetails(params, apiClient)
	if err == nil {
		return out.PrintResult(response.Payload.Project, infoFields)
	}

	return
}
