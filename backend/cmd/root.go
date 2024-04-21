package cmd

import (
	"github.com/spf13/cobra"
)

type UrlShortner struct {
	rootCmd *cobra.Command
}

func NewUrlShortner() *UrlShortner {
	tc := &UrlShortner{
		rootCmd: &cobra.Command{
			Use:   "urlshortner",
			Short: "UrlShortner CLI",
			// no need to provide the default cobra completion command
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
	}

	return tc
}

func (urlshortner *UrlShortner) Start() error {
	urlshortner.rootCmd.AddCommand(serveApiCommand())
	urlshortner.rootCmd.AddCommand(linkserviceCommand())
	urlshortner.rootCmd.AddCommand(billingserviceCommand())
	return urlshortner.execute()
}

func (urlshortner *UrlShortner) execute() error {
	err := urlshortner.rootCmd.Execute()
	if err != nil {
		return err
	}

	return nil
}
