package cmd

import (
	"github.com/adewoleadenigbagbe/url-shortner-service/server"
	"github.com/spf13/cobra"
)

func serveApiCommand() *cobra.Command {
	var apiCmd = &cobra.Command{
		Use:   "serveapi",
		Short: "Serve the API on the Specified host",
		Run: func(cmd *cobra.Command, args []string) {
			server.InitializeAPI()
		},
	}

	return apiCmd
}
