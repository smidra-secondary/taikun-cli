package list

import (
	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/cmd/cmdutils"
	"github.com/itera-io/taikun-cli/utils/out"
	"github.com/itera-io/taikun-cli/utils/out/field"
	"github.com/itera-io/taikun-cli/utils/out/fieldnames"
	"github.com/itera-io/taikun-cli/utils/out/fields"

	"github.com/itera-io/taikungoclient/client/organizations"
	"github.com/itera-io/taikungoclient/models"
	"github.com/spf13/cobra"
)

var ListFields = fields.New(
	[]*field.Field{
		field.NewVisible(
			"ID", "id",
		),
		field.NewVisible(
			"NAME", "name",
		),
		field.NewVisible(
			"FULL-NAME", "fullName",
		),
		field.NewVisible(
			"DISCOUNT-RATE", "discountRate",
		),
		field.NewVisible(
			"PARTNER", "partnerName",
		),
		field.NewVisibleWithToStringFunc(
			fieldnames.IsLocked, "isLocked", out.FormatLockStatus,
		),
		field.NewVisible(
			"READ-ONLY", "isReadOnly",
		),
		field.NewHidden(
			"EMAIL", "email",
		),
		field.NewHidden(
			"BILLING-EMAIL", "billingEmail",
		),
		field.NewHidden(
			"CITY", "city",
		),
		field.NewHidden(
			"COUNTRY", "country",
		),
		field.NewHidden(
			"PHONE", "phone",
		),
		field.NewHidden(
			"VAT", "vatNumber",
		),
		field.NewHidden(
			"SUBSCRIPTION-UPDATES", "isEligibleUpdateSubscription",
		),
		field.NewHidden(
			"PARTNER-ID", "partnerId",
		),
		field.NewHidden(
			"CLOUD-CREDENTIALS", "cloudCredentials",
		),
		field.NewHidden(
			"PROJECTS", "projects",
		),
		field.NewHidden(
			"SERVERS", "servers",
		),
		field.NewHidden(
			"USERS", "users",
		),
		field.NewHiddenWithToStringFunc(
			"CREATED-AT", "createdAt", out.FormatDateTimeString,
		),
	},
)

func NewCmdList() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "List organizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRun()
		},
		Args:    cobra.NoArgs,
		Aliases: cmdutils.ListAliases,
	}

	cmdutils.AddSortByAndReverseFlags(&cmd, ListFields)
	cmdutils.AddColumnsFlag(&cmd, ListFields)
	cmdutils.AddLimitFlag(&cmd)

	return &cmd
}

func listRun() (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	params := organizations.NewOrganizationsListParams().WithV(api.Version)
	if opts.SortBy != "" {
		params = params.WithSortBy(opts.GetSortByParam(ListFields)).WithSortDirection(api.GetSortDirection())
	}

	var organizations = make([]*models.OrganizationDetailsDto, 0)
	for {
		response, err := apiClient.Client.Organizations.OrganizationsList(params, apiClient)
		if err != nil {
			return err
		}
		organizations = append(organizations, response.Payload.Data...)
		organizationsCount := int32(len(organizations))
		if opts.Limit != 0 && organizationsCount >= opts.Limit {
			break
		}
		if organizationsCount == response.Payload.TotalCount {
			break
		}
		params = params.WithOffset(&organizationsCount)
	}

	if opts.Limit != 0 && int32(len(organizations)) > opts.Limit {
		organizations = organizations[:opts.Limit]
	}

	out.PrintResults(organizations, ListFields)
	return
}
