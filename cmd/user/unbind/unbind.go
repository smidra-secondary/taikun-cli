package unbind

import (
	"fmt"
	"taikun-cli/api"
	"taikun-cli/cmd/cmdutils"

	"github.com/itera-io/taikungoclient/client/user_projects"
	"github.com/itera-io/taikungoclient/models"
	"github.com/spf13/cobra"
)

type UnbindOptions struct {
	Username  string
	ProjectID int
}

func NewCmdUnbind() *cobra.Command {
	var opts UnbindOptions

	cmd := &cobra.Command{
		Use:   "unbind",
		Short: "Unbind a user from a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return unbindRun(&opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Username, "username", "u", "", "Username (required)")
	cmd.MarkFlagRequired("username")

	cmd.Flags().IntVarP(&opts.ProjectID, "project-id", "p", 0, "Project ID (required)")
	cmd.MarkFlagRequired("project-id")

	return cmd
}

func unbindRun(opts *UnbindOptions) (err error) {
	apiClient, err := api.NewClient()
	if err != nil {
		return
	}

	body := &models.BindProjectsCommand{
		UserName: opts.Username,
		Projects: []*models.UpdateUserProjectDto{
			{
				ProjectID: int32(opts.ProjectID),
				IsBound:   false,
			},
		},
	}

	params := user_projects.NewUserProjectsBindProjectsParams().WithV(cmdutils.ApiVersion).WithBody(body)
	response, err := apiClient.Client.UserProjects.UserProjectsBindProjects(params, apiClient)
	if err == nil {
		fmt.Println(response.Payload)
	}

	return
}
