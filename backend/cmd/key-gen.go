package cmd

import (
	linkservice "github.com/adewoleadenigbagbe/url-shortner-service/key-generator-service"
	"github.com/spf13/cobra"
)

func linkserviceCommand() *cobra.Command {
	var shortlinkCmd = &cobra.Command{
		Use:   "generatelink",
		Short: "Generate shortlink",
		Long:  `Generate shortlink for the future`,
		Run: func(cmd *cobra.Command, args []string) {
			linkservice.Run()
		},
	}

	return shortlinkCmd
}
