/*
Copyright © 2025 Daniel Soler dsolerh.cinter95@gmail.com

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
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-wrkpub",
	Short: "Publish packages on a workspace with individual tags",
	Long:  `Publish individual packages in a workspace each with their own version`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// add config file flag
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		".go-workpublish.yaml",
		"config file (default is $PWD/.go-workpublish.yaml)",
	)
}
