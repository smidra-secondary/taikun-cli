package create

import (
	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/apiconfig"
	"github.com/itera-io/taikun-cli/cmd/cmdutils"
	"github.com/itera-io/taikun-cli/config"
	"github.com/itera-io/taikun-cli/utils/format"

	"github.com/itera-io/taikungoclient/client/showback"
	"github.com/itera-io/taikungoclient/models"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	Name           string
	OrganizationID int32
	Password       string
	URL            string
	Username       string
	IDOnly         bool
}

func NewCmdCreate() *cobra.Command {
	var opts CreateOptions

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

	cmdutils.AddIdOnlyFlag(&cmd, &opts.IDOnly)

	return &cmd
}

func printResult(resource interface{}) {
	if config.OutputFormat == config.OutputFormatJson {
		format.PrettyPrintJson(resource)
	} else if config.OutputFormat == config.OutputFormatTable {
		format.PrettyPrintApiResponseTable(resource,
			"id",
			"name",
			"organizationName",
			"url",
			"createdAt",
			"isLocked",
		)
	}
}

func createRun(opts *CreateOptions) (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	body := models.CreateShowbackCredentialCommand{
		Name:           opts.Name,
		OrganizationID: opts.OrganizationID,
		Password:       opts.Password,
		URL:            opts.URL,
		Username:       opts.Username,
	}

	params := showback.NewShowbackCreateCredentialParams().WithV(apiconfig.Version)
	params = params.WithBody(&body)

	response, err := apiClient.Client.Showback.ShowbackCreateCredential(params, apiClient)
	if err == nil {
		if opts.IDOnly {
			format.PrintResourceID(response.Payload)
		} else {
			printResult(response.Payload)
		}
	}

	return
}
