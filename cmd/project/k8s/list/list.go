package list

import (
	"context"
	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/cmd/cmderr"
	"github.com/itera-io/taikun-cli/cmd/cmdutils"
	"github.com/itera-io/taikun-cli/config"
	"github.com/itera-io/taikun-cli/utils/out"
	"github.com/itera-io/taikun-cli/utils/out/field"
	"github.com/itera-io/taikun-cli/utils/out/fields"
	"github.com/itera-io/taikun-cli/utils/types"
	tk "github.com/itera-io/taikungoclient"
	taikuncore "github.com/itera-io/taikungoclient/client"
	"github.com/spf13/cobra"
)

var listFields = fields.New(
	[]*field.Field{
		field.NewVisible(
			"ID", "id",
		),
		field.NewVisible(
			"NAME", "name",
		),
		field.NewVisible(
			"IP", "ipAddress",
		),
		field.NewHidden(
			"CLOUD", "cloudType",
		),
		field.NewVisibleWithToStringFunc(
			"AVAILABILITY-ZONE", "availabilityZone", out.FormatAvailabilityZones,
		),
		field.NewVisible(
			"FLAVOR", "",
			// JSON property name is set in the listRun function
			// as it depends on the server's cloud type
		),
		field.NewVisible(
			"CPU", "cpu",
		),
		field.NewVisibleWithToStringFunc(
			"RAM", "ram", out.FormatBToGiB,
		),
		field.NewVisibleWithToStringFunc(
			"DISK", "diskSize", out.FormatBToGiB,
		),
		field.NewVisible(
			"ROLE", "role",
		),
		field.NewVisible(
			"STATUS", "status",
		),
		field.NewHidden(
			"PROJECT", "projectName",
		),
		field.NewHidden(
			"PROJECT-ID", "projectId",
		),
		field.NewHidden(
			"ORG", "organizationName",
		),
		field.NewHidden(
			"ORG-ID", "organizationId",
		),
		field.NewVisibleWithToStringFunc(
			"CREATED-AT", "createdAt", out.FormatDateTimeString,
		),
		field.NewHidden(
			"CREATED-BY", "createdBy",
		),
		field.NewHiddenWithToStringFunc(
			"LAST-MODIFIED", "lastModified", out.FormatDateTimeString,
		),
		field.NewHidden(
			"LAST-MODIFIED-BY", "lastModifiedBy",
		),
		field.NewVisible(
			"WASM", "wasmEnabled",
		),
	},
)

type ListOptions struct {
	ProjectID int32
}

func NewCmdList() *cobra.Command {
	var opts ListOptions

	cmd := cobra.Command{
		Use:   "list <project-id>",
		Short: "List a project's Kubernetes servers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.ProjectID, err = types.Atoi32(args[0])
			if err != nil {
				return
			}
			return listRun(&opts)
		},
		Aliases: cmdutils.ListAliases,
	}

	cmdutils.AddSortByAndReverseFlags(&cmd, "projects-k8s", listFields)
	cmdutils.AddColumnsFlag(&cmd, listFields)

	return &cmd
}

func listRun(opts *ListOptions) (err error) {
	projectServers, err := ListServers(opts)
	if err == nil {
		flavorJsonPropertyName, err := getFlavorField(projectServers)
		if err != nil {
			return err
		}

		if err := listFields.SetFieldJsonPropertyName("FLAVOR", flavorJsonPropertyName); err != nil {
			return err
		}

		return out.PrintResults(projectServers, listFields)
	}

	return

}

func ListServers(opts *ListOptions) (projectServers []taikuncore.ServerListDto, err error) {
	myApiClient := tk.NewClient()
	myApiRequest := myApiClient.Client.ServersAPI.ServersDetails(context.TODO(), opts.ProjectID)
	if config.SortBy != "" {
		myApiRequest = myApiRequest.SortBy(config.SortBy).SortDirection(*api.GetSortDirection())
	}
	data, response, err := myApiRequest.Execute()
	if err != nil {
		err = tk.CreateError(response, err)
		return
	}
	projectServers = data.GetData()
	return

}

func getFlavorField(servers []taikuncore.ServerListDto) (string, error) {
	if len(servers) == 0 {
		return "flavor", nil
	}

	if servers[0].GetAwsInstanceType() != "" {
		return "awsInstanceType", nil
	}

	if servers[0].GetAzureVmSize() != "" {
		return "azureVmSize", nil
	}

	if servers[0].GetOpenstackFlavor() != "" {
		return "openstackFlavor", nil
	}

	if servers[0].GetGoogleMachineType() != "" {
		return "googleMachineType", nil
	}

	return "", cmderr.ErrServerHasNoFlavors
}
