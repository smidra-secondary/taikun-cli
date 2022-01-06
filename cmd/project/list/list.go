package list

import (
	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/apiconfig"
	"github.com/itera-io/taikun-cli/cmd/cmdutils"
	"github.com/itera-io/taikun-cli/config"
	"github.com/itera-io/taikun-cli/utils/format"

	"github.com/itera-io/taikungoclient/client/projects"
	"github.com/itera-io/taikungoclient/models"
	"github.com/spf13/cobra"
)

type ListOptions struct {
	OrganizationID int32
}

func NewCmdList() *cobra.Command {
	var opts ListOptions

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRun(&opts)
		},
	}

	cmd.Flags().BoolVarP(&config.ReverseSortDirection, "reverse", "r", false, "Reverse order of results")
	cmd.Flags().Int32VarP(&opts.OrganizationID, "organization-id", "o", 0, "Organization ID (only applies for Partner role)")

	cmdutils.AddSortByFlag(cmd, &config.SortBy, models.ProjectListForUIDto{})
	cmdutils.AddLimitFlag(cmd)

	return cmd
}

func listRun(opts *ListOptions) (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	params := projects.NewProjectsListParams().WithV(apiconfig.Version)
	if opts.OrganizationID != 0 {
		params = params.WithOrganizationID(&opts.OrganizationID)
	}
	if config.ReverseSortDirection {
		apiconfig.ReverseSortDirection()
	}
	if config.SortBy != "" {
		params = params.WithSortBy(&config.SortBy).WithSortDirection(&apiconfig.SortDirection)
	}

	var projects = make([]*models.ProjectListForUIDto, 0)
	for {
		response, err := apiClient.Client.Projects.ProjectsList(params, apiClient)
		if err != nil {
			return err
		}
		projects = append(projects, response.Payload.Data...)
		projectsCount := int32(len(projects))
		if config.Limit != 0 && projectsCount >= config.Limit {
			break
		}
		if projectsCount == response.Payload.TotalCount {
			break
		}
		params = params.WithOffset(&projectsCount)
	}

	if config.Limit != 0 && int32(len(projects)) > config.Limit {
		projects = projects[:config.Limit]
	}

	format.PrintResults(projects,
		"id",
		"name",
		"organizationName",
		"status",
		"health",
		"createdAt",
		"kubernetesCurrentVersion",
		"cloudType",
		"hasKubeConfigFile",
		"quotaId",
		"expiredAt",
		"isLocked",
	)
	return
}
