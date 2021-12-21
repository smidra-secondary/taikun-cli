package create

import (
	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/apiconfig"
	"github.com/itera-io/taikun-cli/cmd/cmdutils"
	"github.com/itera-io/taikun-cli/utils/format"

	"github.com/itera-io/taikungoclient/client/showback"
	"github.com/itera-io/taikungoclient/models"
	"github.com/spf13/cobra"
)

func NewCmdCreate() *cobra.Command {
	var opts models.CreateShowbackCredentialCommand

	cmd := cobra.Command{
		Use:   "create <name>",
		Short: "Create a showback credential",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			return createRun(&opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Password, "password", "p", "", "Password (Prometheus or other) (required)")
	cmdutils.MarkFlagRequired(&cmd, "password")

	cmd.Flags().StringVarP(&opts.Username, "username", "l", "", "Username (Prometheus or other) (required)")
	cmdutils.MarkFlagRequired(&cmd, "username")

	cmd.Flags().StringVarP(&opts.URL, "url", "u", "", "URL of the source (required)")
	cmdutils.MarkFlagRequired(&cmd, "url")

	cmd.Flags().Int32VarP(&opts.OrganizationID, "organization-id", "o", 0, "Organization ID (only applies for Partner role)")

	return &cmd
}

func createRun(opts *models.CreateShowbackCredentialCommand) (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	params := showback.NewShowbackCreateCredentialParams().WithV(apiconfig.Version)
	params = params.WithBody(opts)

	response, err := apiClient.Client.Showback.ShowbackCreateCredential(params, apiClient)
	if err == nil {
		format.PrettyPrintJson(response.Payload)
	}

	return
}
