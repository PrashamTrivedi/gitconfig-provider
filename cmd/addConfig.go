/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var providerNameToAdd string
var keyToAdd string
var value string

// addConfigCmd represents the addConfig command
var addConfigCmd = &cobra.Command{
	Use:   "addConfig",
	Short: "Adds config to your git provider",

	Run: func(cmd *cobra.Command, args []string) {
		AddProviderPropertyFromName(providerNameToAdd, keyToAdd, value)
	},
}

func init() {
	addConfigCmd.Flags().StringVarP(&providerNameToAdd, "provider", "p", "", "Name of provider to store property")
	addConfigCmd.Flags().StringVarP(&keyToAdd, "key", "k", "", "Key of property")
	addConfigCmd.Flags().StringVarP(&value, "value", "v", "", "Value of property")

	addConfigCmd.MarkFlagRequired("provider")
	addConfigCmd.MarkFlagRequired("key")
	addConfigCmd.MarkFlagRequired("value")

	RootCmd.AddCommand(addConfigCmd)
}
