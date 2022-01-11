package etc

import (
	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/apiconfig"
	"github.com/itera-io/taikun-cli/utils/format"
	"github.com/itera-io/taikun-cli/utils/types"
	"github.com/itera-io/taikungoclient/client/notifications"
	"github.com/itera-io/taikungoclient/models"
	"github.com/spf13/cobra"
)

type EtcOptions struct {
	ProjectID int32
}

func NewCmdEtc() *cobra.Command {
	var opts EtcOptions

	cmd := cobra.Command{
		Use:   "etc <project-id>",
		Short: "Get estimated time of completion for project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.ProjectID, err = types.Atoi32(args[0])
			if err != nil {
				return
			}
			return etcRun(&opts)
		},
	}

	return &cmd
}

func etcRun(opts *EtcOptions) (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	body := models.GetProjectOperationCommand{
		ProjectID: opts.ProjectID,
	}

	params := notifications.NewNotificationsGetProjectOperationMessagesParams().WithV(apiconfig.Version)
	params = params.WithBody(&body)

	response, err := apiClient.Client.Notifications.NotificationsGetProjectOperationMessages(params, apiClient)
	if err == nil {
		format.PrintResult(response, "operation", "estimatedTime")
	}

	return
}
