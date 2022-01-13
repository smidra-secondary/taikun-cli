package list

import (
	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/cmd/cmdutils"
	"github.com/itera-io/taikun-cli/config"
	"github.com/itera-io/taikun-cli/utils/out"

	"github.com/itera-io/taikungoclient/client/organizations"
	"github.com/itera-io/taikungoclient/models"
	"github.com/spf13/cobra"
)

func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List organizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRun()
		},
		Args:    cobra.NoArgs,
		Aliases: cmdutils.ListAliases,
	}

	cmdutils.AddSortByAndReverseFlags(cmd, models.OrganizationDetailsDto{})
	cmdutils.AddLimitFlag(cmd)

	return cmd
}

func listRun() (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	params := organizations.NewOrganizationsListParams().WithV(api.Version)
	if config.ReverseSortDirection {
		api.ReverseSortDirection()
	}
	if config.SortBy != "" {
		params = params.WithSortBy(&config.SortBy).WithSortDirection(&api.SortDirection)
	}

	var organizations = make([]*models.OrganizationDetailsDto, 0)
	for {
		response, err := apiClient.Client.Organizations.OrganizationsList(params, apiClient)
		if err != nil {
			return err
		}
		organizations = append(organizations, response.Payload.Data...)
		organizationsCount := int32(len(organizations))
		if config.Limit != 0 && organizationsCount >= config.Limit {
			break
		}
		if organizationsCount == response.Payload.TotalCount {
			break
		}
		params = params.WithOffset(&organizationsCount)
	}

	if config.Limit != 0 && int32(len(organizations)) > config.Limit {
		organizations = organizations[:config.Limit]
	}

	out.PrintResults(organizations,
		"id",
		"name",
		"fullName",
		"discountRate",
		"partnerName",
		"isEligibleUpdateSubscription",
		"isLocked",
		"isReadOnly",
		"users",
		"cloudCredentials",
		"projects",
		"servers",
	)
	return
}
