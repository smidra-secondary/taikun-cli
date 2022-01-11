package commit

import (
	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/apiconfig"
	"github.com/itera-io/taikun-cli/utils/format"
	"github.com/itera-io/taikun-cli/utils/types"
	"github.com/itera-io/taikungoclient/client/projects"
	"github.com/spf13/cobra"
)

type CommitOptions struct {
	ProjectID int32
}

func NewCmdCommit() *cobra.Command {
	var opts CommitOptions

	cmd := cobra.Command{
		Use:   "commit <project-id>",
		Short: "Commit the project's changes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.ProjectID, err = types.Atoi32(args[0])
			if err != nil {
				return
			}
			return commitRun(&opts)
		},
	}

	return &cmd
}

func commitRun(opts *CommitOptions) (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	params := projects.NewProjectsCommitParams().WithV(apiconfig.Version)
	params = params.WithProjectID(opts.ProjectID)

	_, err = apiClient.Client.Projects.ProjectsCommit(params, apiClient)
	if err == nil {
		format.PrintStandardSuccess()
	}

	return
}
