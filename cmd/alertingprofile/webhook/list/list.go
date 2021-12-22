package list

import (
	"fmt"

	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/apiconfig"
	"github.com/itera-io/taikun-cli/cmd/cmderr"
	"github.com/itera-io/taikun-cli/utils/format"
	"github.com/itera-io/taikun-cli/utils/types"

	"github.com/itera-io/taikungoclient/client/alerting_profiles"
	"github.com/spf13/cobra"
)

type ListOptions struct {
	AlertingProfileID int32
	Limit             int32
}

func NewCmdList() *cobra.Command {
	var opts ListOptions

	cmd := &cobra.Command{
		Use:   "list <alerting-profile-id>",
		Short: "List an alerting profile's webhooks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Limit < 0 {
				return cmderr.NegativeLimitFlagError
			}
			alertingProfileID, err := types.Atoi32(args[0])
			if err != nil {
				return cmderr.IDArgumentNotANumberError
			}
			opts.AlertingProfileID = alertingProfileID
			return listRun(&opts)
		},
	}

	cmd.Flags().Int32VarP(&opts.Limit, "limit", "l", 0, "Limit number of results (limitless by default)")

	return cmd
}

func listRun(opts *ListOptions) (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	params := alerting_profiles.NewAlertingProfilesListParams().WithV(apiconfig.Version)
	params = params.WithID(&opts.AlertingProfileID)

	response, err := apiClient.Client.AlertingProfiles.AlertingProfilesList(params, apiClient)
	if err != nil {
		return err
	}
	if len(response.Payload.Data) != 1 {
		return fmt.Errorf("Alerting profile with ID %d not found.", opts.AlertingProfileID)
	}
	alertingWebhooks := response.Payload.Data[0].Webhooks

	if opts.Limit != 0 && int32(len(alertingWebhooks)) > opts.Limit {
		alertingWebhooks = alertingWebhooks[:opts.Limit]
	}

	format.PrintResults(alertingWebhooks,
		"url",
	)
	return
}
