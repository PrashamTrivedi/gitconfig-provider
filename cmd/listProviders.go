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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// listProvidersCmd represents the listProviders command
var listProvidersCmd = &cobra.Command{
	Use:   "listProviders",
	Short: "Lists git providers",

	Run: func(cmd *cobra.Command, args []string) {
		gitProviders, err := GetProviders()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		for i, provider := range gitProviders {
			fmt.Printf("%d. %s (%s)\n", i+1, provider.Name, provider.Url)
			if len(provider.Properties) > 0 {
				fmt.Println("   Properties")
				for key, val := range provider.Properties {
					fmt.Printf("\t%s: %s\n", key, val)
				}
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(listProvidersCmd)
}
