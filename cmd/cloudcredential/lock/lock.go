package lock

import (
	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/apiconfig"
	"github.com/itera-io/taikun-cli/cmd/cmderr"
	"github.com/itera-io/taikun-cli/utils/format"
	"github.com/itera-io/taikun-cli/utils/types"

	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/models"
	"github.com/spf13/cobra"
)

func NewCmdLock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock <cloud-credential-id>",
		Short: "Lock a cloud credential",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := types.Atoi32(args[0])
			if err != nil {
				return cmderr.IDArgumentNotANumberError
			}
			return lockRun(id)
		},
	}

	return cmd
}

func lockRun(id int32) (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	body := &models.CloudLockManagerCommand{
		ID:   id,
		Mode: types.LockedMode,
	}
	params := cloud_credentials.NewCloudCredentialsLockManagerParams().WithV(apiconfig.Version).WithBody(body)
	_, err = apiClient.Client.CloudCredentials.CloudCredentialsLockManager(params, apiClient)
	if err == nil {
		format.PrintStandardSuccess()
	}

	return
}
