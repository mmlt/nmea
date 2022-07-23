package tool

import (
	"github.com/mmlt/nmea/pkg/record"
	"github.com/spf13/cobra"
)

// NewCmdRecord returns a command to record NMEA sentences.
func NewCmdRecord() *cobra.Command {
	// flags
	var (
		host      string
		filename  string
		timestamp bool
	)

	cmd := cobra.Command{
		Use:   "record --host address --file name [--timestamp]",
		Short: "Record NMEA sencentences send by host to file",
		Long:  `Record NMEA sencentences send by host in a file for diagnostics or playback.`,
		Run: func(c *cobra.Command, args []string) {
			rr, err := record.Open(host, filename, timestamp)
			exitOnError(err)

			defer rr.Close()
			// close network connection to make blocking read in Run exit
			go func() {
				select {
				case <-c.Context().Done():
					rr.Close()
				}
			}()
			err = rr.Run(c.Context())
			exitOnError(err)
		},
	}

	cmd.Flags().StringVar(&host, "host", "localhost:10110", "The address of the host to connect to.")
	must(cmd.MarkFlagRequired("host"))
	cmd.Flags().StringVar(&filename, "file", "", "The name of output file.")
	must(cmd.MarkFlagRequired("file"))
	cmd.Flags().BoolVar(&timestamp, "timestamp", true, "Add timestamp to output.")

	return &cmd
}
