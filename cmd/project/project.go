package project

import (
	"github.com/itera-io/taikun-cli/cmd/project/create"
	"github.com/itera-io/taikun-cli/cmd/project/delete"
	"github.com/itera-io/taikun-cli/cmd/project/list"
	"github.com/itera-io/taikun-cli/cmd/project/lock"
	"github.com/itera-io/taikun-cli/cmd/project/quotas"
	"github.com/itera-io/taikun-cli/cmd/project/unlock"

	"github.com/spf13/cobra"
)

func NewCmdProject() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project <command>",
		Short: "Manage projects",
	}

	cmd.AddCommand(create.NewCmdCreate())
	cmd.AddCommand(delete.NewCmdDelete())
	cmd.AddCommand(list.NewCmdList())
	cmd.AddCommand(lock.NewCmdLock())
	cmd.AddCommand(quotas.NewCmdQuotas())
	cmd.AddCommand(unlock.NewCmdUnlock())

	return cmd
}
