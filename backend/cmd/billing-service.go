package cmd

import (
	billing "github.com/adewoleadenigbagbe/url-shortner-service/billing-service"
	"github.com/spf13/cobra"
)

func billingserviceCommand() *cobra.Command {
	var shortlinkCmd = &cobra.Command{
		Use:   "billing",
		Short: "Generate revenue",
		Long:  `Generate revenue`,
		Run: func(cmd *cobra.Command, args []string) {
			billing.Run()
		},
	}

	return shortlinkCmd
}
