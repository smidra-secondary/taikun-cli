package info

import (
	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/apiconfig"
	"github.com/itera-io/taikun-cli/utils/format"
	"github.com/itera-io/taikungoclient/client/users"
	"github.com/spf13/cobra"
)

func NewCmdInfo() *cobra.Command {
	cmd := cobra.Command{
		Use:   "info",
		Short: "Retrieve your information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return infoRun()
		},
	}

	return &cmd
}

func infoRun() (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	params := users.NewUsersDetailsParams().WithV(apiconfig.Version)

	response, err := apiClient.Client.Users.UsersDetails(params, apiClient)
	if err == nil {
		format.PrintResultVertical(response.Payload.Data,
			"id",
			"username",
			"role",
			"organizationName",
			"email",
			"displayName",
			"isEmailConfirmed",
			"isEmailNotificationEnabled",
			"isApprovedByPartner",
			"owner",
			"isLocked",
			"createdAt",
			"lastLoginAt",
		)
	}

	return
}
