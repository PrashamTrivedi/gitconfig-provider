/*
Copyright © 2021 NAME HERE prash2488@gmail.com

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
	"path/filepath"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gitconfig-provider",
	Short: "Adds various git configs based on your git origin",
}

func init() {
	home, _ := os.UserHomeDir()

	dbPath := filepath.Join(home, "gitProviders.json")
	fmt.Println(dbPath)
	defaultFilesData := "[{\"name\":\"github\",\"url\":\"https://github.com/\"},{\"name\":\"bitbucket\",\"url\":\"https://bitbucket.org/\"},{\"name\":\"gitlab\",\"url\":\"https://gitlab.com/\"}]"
	dbFile, errorData := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0600)
	if errorData != nil {
		fmt.Println("Error in processing file:", errorData.Error())
		os.Exit(1)
	}
	dbFile.WriteString(defaultFilesData)
}
