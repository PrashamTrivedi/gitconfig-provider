package cmd

import (
	"github.com/spf13/cobra"
)

//create two string variables to hold providerName and providerURL
var providerNameToRemove, providerURLToRemove string

//create cobra command
var removeProvider = &cobra.Command{
	Use:   "removeProvider",
	Short: "Removes a git provider ",
	Long:  `This command will remove a git provider`,
	Run: func(cmd *cobra.Command, args []string) {
		//call add provider with the providerName and providerURL keep properties blank
		RemoveProvider(providerNameToRemove, providerURLToRemove)
	},
}

func init() {
	removeProvider.Flags().StringVarP(&providerNameToRemove, "name", "n", "", "Name of the provider")
	removeProvider.Flags().StringVarP(&providerURLToRemove, "url", "u", "", "URL of the provider")

	RootCmd.AddCommand(removeProvider)
}
