package create

import (
	"taikun-cli/api"
	"taikun-cli/cmd/cmdutils"

	"github.com/itera-io/taikungoclient/client/azure"
	"github.com/itera-io/taikungoclient/models"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	Name                  string
	AzureSubscriptionId   string
	AzureClientId         string
	AzureClientSecret     string
	AzureTenantId         string
	AzureLocation         string
	AzureAvailabilityZone string
	OrganizationID        int32
}

func NewCmdCreate() *cobra.Command {
	var opts CreateOptions

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create an azure cloud credential",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			return createRun(&opts)
		},
	}

	cmd.Flags().StringVarP(&opts.AzureSubscriptionId, "subscription-id", "s", "", "Azure Subscription ID (required)")
	cmdutils.MarkFlagRequired(cmd, "subscription-id")

	cmd.Flags().StringVarP(&opts.AzureClientId, "client-id", "c", "", "Azure Client ID (required)")
	cmdutils.MarkFlagRequired(cmd, "client-id")

	cmd.Flags().StringVarP(&opts.AzureClientSecret, "client-secret", "p", "", "Azure Client Secret (required)")
	cmdutils.MarkFlagRequired(cmd, "client-secret")

	cmd.Flags().StringVarP(&opts.AzureTenantId, "tenant-id", "t", "", "Azure Tenant ID (required)")
	cmdutils.MarkFlagRequired(cmd, "tenant-id")

	cmd.Flags().StringVarP(&opts.AzureLocation, "location", "l", "", "Azure Location (required)")
	cmdutils.MarkFlagRequired(cmd, "location")

	cmd.Flags().StringVarP(&opts.AzureAvailabilityZone, "availability-zone", "a", "", "Azure Availability Zone (required)")
	cmdutils.MarkFlagRequired(cmd, "availability-zone")

	cmd.Flags().Int32VarP(&opts.OrganizationID, "organization-id", "o", 0, "Organization ID")

	return cmd
}

func createRun(opts *CreateOptions) (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	body := &models.CreateAzureCloudCommand{
		Name:                  opts.Name,
		AzureSubscriptionID:   opts.AzureSubscriptionId,
		AzureClientID:         opts.AzureClientId,
		AzureClientSecret:     opts.AzureClientSecret,
		AzureTenantID:         opts.AzureTenantId,
		AzureLocation:         opts.AzureLocation,
		AzureAvailabilityZone: opts.AzureAvailabilityZone,
		OrganizationID:        opts.OrganizationID,
	}

	params := azure.NewAzureCreateParams().WithV(cmdutils.ApiVersion).WithBody(body)
	response, err := apiClient.Client.Azure.AzureCreate(params, apiClient)
	if err == nil {
		cmdutils.PrettyPrint(response.Payload)
	}

	return
}
