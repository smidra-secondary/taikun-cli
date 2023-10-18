package rename

import (
	"context"
	"github.com/itera-io/taikun-cli/cmd/cmderr"
	"github.com/itera-io/taikun-cli/cmd/cmdutils"
	"github.com/itera-io/taikun-cli/utils/out"
	"github.com/itera-io/taikun-cli/utils/types"
	tk "github.com/itera-io/taikungoclient"
	taikuncore "github.com/itera-io/taikungoclient/client"
	"github.com/spf13/cobra"
)

type RenameOptions struct {
	ID   int32
	Name string
}

func NewCmdRename() *cobra.Command {
	var opts RenameOptions

	cmd := cobra.Command{
		Use:   "rename <standalone-profile-id>",
		Short: "Rename a standalone profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.ID, err = types.Atoi32(args[0])
			if err != nil {
				return cmderr.ErrIDArgumentNotANumber
			}
			return renameRun(&opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "New name (required)")
	cmdutils.MarkFlagRequired(&cmd, "name")

	return &cmd
}

func renameRun(opts *RenameOptions) (err error) {
	// Create and authenticated client to the Taikun API
	myApiClient := tk.NewClient()

	// Prepare the arguments for the query
	body := taikuncore.StandAloneProfileUpdateCommand{
		Id:   &opts.ID,
		Name: *taikuncore.NewNullableString(&opts.Name),
	}

	// Execute a query into the API + graceful exit
	response, err := myApiClient.Client.StandaloneProfileAPI.StandaloneprofileEdit(context.TODO()).StandAloneProfileUpdateCommand(body).Execute()
	if err != nil {
		return tk.CreateError(response, err)
	}

	out.PrintStandardSuccess()
	return
	/*
		apiClient, err := taikungoclient.NewClient()
		if err != nil {
			return
		}

		body := models.StandAloneProfileUpdateCommand{
			ID:   opts.ID,
			Name: opts.Name,
		}

		params := stand_alone_profile.NewStandAloneProfileEditParams().WithV(taikungoclient.Version)
		params = params.WithBody(&body)

		_, err = apiClient.Client.StandAloneProfile.StandAloneProfileEdit(params, apiClient)
		if err == nil {
			out.PrintStandardSuccess()
		}

		return
	*/
}
