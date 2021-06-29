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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// listProvidersCmd represents the listProviders command
var listProvidersCmd = &cobra.Command{
	Use:   "listProviders",
	Short: "Lists git providers",

	Run: func(cmd *cobra.Command, args []string) {
		home, _ := homedir.Dir()
		dbPath := filepath.Join(home, "gitProviders.json")
		var gitProviders []GitProvider
		fileBytes, err := ioutil.ReadFile(dbPath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if err := json.Unmarshal(fileBytes, &gitProviders); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		for i, provider := range gitProviders {
			fmt.Printf("%d. %s (%s)\n", i+1, provider.Name, provider.Url)
		}

	},
}

func init() {
	RootCmd.AddCommand(listProvidersCmd)
}
