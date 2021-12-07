package root

import (
	"github.com/spf13/cobra"
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "taikun",
		Short: "Taikun CLI",
		Long:  `Manage Taikun resources from the command line.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	return cmd
}
