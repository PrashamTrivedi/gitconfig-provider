package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//create two string variables to hold providerName and providerURL
var providerName, providerURL string
var doForcefully bool

//create cobra command
var addProviderCmd = &cobra.Command{
	Use:   "addProvider",
	Short: "Adds a git provider for you",
	Long:  `This command will add a git provider for you`,
	Run: func(cmd *cobra.Command, args []string) {
		//call add provider with the providerName and providerURL keep properties blank
		provider, index := AddProvider(providerName, providerURL,doForcefully)
		if index == -1 {
			fmt.Println("Provider", provider.Name, "added successfully")
		} else {
			fmt.Println("Recognized existing provider with:", provider.Name, "& url:", provider.Url)
		}
	},
}

func init() {
	addProviderCmd.Flags().StringVarP(&providerName, "name", "n", "", "Name of the provider")
	addProviderCmd.Flags().StringVarP(&providerURL, "url", "u", "", "URL of the provider")
	addProviderCmd.Flags().BoolVarP(&doForcefully, "force", "f", false, "HIGHLY EXPERIMENTAL: Add the provider forcefully")

	RootCmd.AddCommand(addProviderCmd)
}
