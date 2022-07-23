package tool

import (
	"github.com/mmlt/nmea/pkg/record"
	"github.com/spf13/cobra"
)

// NewCmdPlayback returns a command to playback NMEA sentences.
func NewCmdPlayback() *cobra.Command {
	// flags
	var (
		host     string
		filename string
	)

	cmd := cobra.Command{
		Use:   "playback --host address --file name [--timestamp]",
		Short: "Playback NMEA sencentences from file and send them to host",
		Long: `Playback NMEA sencentences from file and send them to host.
If the file contains timestamps the same interval will be used.`,
		Run: func(c *cobra.Command, args []string) {
			rr, err := record.OpenPlayback(host, filename)
			exitOnError(err)

			defer rr.Close()

			err = rr.Run(c.Context())
			exitOnError(err)
		},
	}

	cmd.Flags().StringVar(&host, "host", "localhost:10110", "The address of the host to connect to.")
	must(cmd.MarkFlagRequired("host"))
	cmd.Flags().StringVar(&filename, "file", "", "The name of input file.")
	must(cmd.MarkFlagRequired("file"))

	return &cmd
}
