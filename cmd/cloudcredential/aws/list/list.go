package list

import (
	"context"
	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/cmd/cmdutils"
	"github.com/itera-io/taikun-cli/config"
	"github.com/itera-io/taikun-cli/utils/out"
	"github.com/itera-io/taikun-cli/utils/out/field"
	"github.com/itera-io/taikun-cli/utils/out/fields"
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
			"ORG", "organizationName",
		),
		field.NewHidden(
			"ORG-ID", "organizationId",
		),
		field.NewVisible(
			"REGION", "region",
		),
		field.NewVisible(
			"AVAILABILITY-ZONES", "availabilityZones",
		),
		field.NewVisible(
			"DEFAULT", "isDefault",
		),
		field.NewVisibleWithToStringFunc(
			"LOCK", "isLocked", out.FormatLockStatus,
		),
		field.NewHidden(
			"CREATED-BY", "createdBy",
		),
	},
)

type ListOptions struct {
	OrganizationID int32
	Limit          int32
}

func NewCmdList() *cobra.Command {
	var opts ListOptions

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List AWS cloud credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRun(&opts)
		},
		Args:    cobra.NoArgs,
		Aliases: cmdutils.ListAliases,
	}

	cmd.Flags().Int32VarP(&opts.OrganizationID, "organization-id", "o", 0, "Organization ID (only applies for Partner role)")

	cmdutils.AddLimitFlag(cmd, &opts.Limit)
	cmdutils.AddSortByAndReverseFlags(cmd, "cloud-credentials", listFields)
	cmdutils.AddColumnsFlag(cmd, listFields)

	return cmd
}

func listRun(opts *ListOptions) error {
	amazonCloudCredentials, err := ListCloudCredentialsAws(opts)
	if err != nil {
		return err
	}

	return out.PrintResults(amazonCloudCredentials, listFields)
}

func ListCloudCredentialsAws(opts *ListOptions) (credentials []interface{}, err error) {
	myApiClient := tk.NewClient()
	myRequest := myApiClient.Client.CloudCredentialAPI.CloudcredentialsDashboardList(context.TODO())
	if opts.OrganizationID != 0 {
		myRequest = myRequest.OrganizationId(opts.OrganizationID)
	}

	if config.SortBy != "" {
		myRequest = myRequest.SortBy(config.SortBy).SortDirection(*api.GetSortDirection())
	}

	var amazonCloudCredentials = make([]taikuncore.AmazonCredentialsListDto, 0)

	for {
		data, response, newError := myRequest.Execute()
		if newError != nil {
			err = tk.CreateError(response, err)
			return
		}

		amazonCloudCredentials = append(amazonCloudCredentials, data.Amazon...)

		count := int32(len(amazonCloudCredentials))
		if opts.Limit != 0 && count >= opts.Limit {
			break
		}

		if count == data.GetTotalCountAws() {
			break
		}

		myRequest = myRequest.Offset(count)
	}

	if opts.Limit != 0 && int32(len(amazonCloudCredentials)) > opts.Limit {
		amazonCloudCredentials = amazonCloudCredentials[:opts.Limit]
	}

	credentials = make([]interface{}, len(amazonCloudCredentials))
	for i, credential := range amazonCloudCredentials {
		credentials[i] = credential
	}

	return credentials, nil

}
