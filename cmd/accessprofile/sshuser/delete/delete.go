package delete

import (
	"taikun-cli/api"
	"taikun-cli/utils"
	"taikun-cli/utils/types"

	"github.com/itera-io/taikungoclient/client/ssh_users"
	"github.com/itera-io/taikungoclient/models"
	"github.com/spf13/cobra"
)

func NewCmdDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <ssh-user-id>",
		Short: "Delete SSH user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sshUserID, err := types.Atoi32(args[0])
			if err != nil {
				return utils.WrongIDArgumentFormatError
			}
			return deleteRun(sshUserID)
		},
	}
	return cmd
}

func deleteRun(id int32) (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	body := models.DeleteSSHUserCommand{
		ID: id,
	}
	params := ssh_users.NewSSHUsersDeleteParams().WithV(utils.ApiVersion).WithBody(&body)
	_, err = apiClient.Client.SSHUsers.SSHUsersDelete(params, apiClient)

	if err == nil {
		utils.PrintDeleteSuccess("SSH user", id)
	}

	return
}
