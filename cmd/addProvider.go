package cmd

import (
	"github.com/spf13/cobra"
)

//create two string variables to hold providerName and providerURL
var providerName, providerURL string

//create cobra command
var addProviderCmd = &cobra.Command{
	Use:   "addProvider",
	Short: "Adds a git provider for you",
	Long:  `This command will add a git provider for you`,
	Run: func(cmd *cobra.Command, args []string) {
		//call add provider with the providerName and providerURL keep properties blank
		AddProvider(providerName, providerURL)
	},
}

func init() {
	addProviderCmd.Flags().StringVarP(&providerName, "name", "n", "", "Name of the provider")
	addProviderCmd.Flags().StringVarP(&providerURL, "url", "u", "", "URL of the provider")

	RootCmd.AddCommand(addProviderCmd)
}
