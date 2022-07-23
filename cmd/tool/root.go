package tool

import (
	"github.com/spf13/cobra"
)

// NewRootCommand returns the root of all tool commands.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tool",
		Short: "Tool for nmea",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
		Hidden: true,
	}

	cmd.AddCommand(NewCmdRecord())
	cmd.AddCommand(NewCmdPlayback())

	return cmd
}
