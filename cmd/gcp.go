/*
Copyright © 2021 Danil Beltyukov <root@danil.co>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var gcpCmd = &cobra.Command{
	Use:   "gcp",
	Short: "Google Cloud Platform",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(gcpCmd)

	gcpCmd.PersistentFlags().StringP("output", "o", "text", "Output format")

	gcpCmd.PersistentFlags().String("project", "", "Project name")
	gcpCmd.PersistentFlags().String("key", "", "ServiceAccount JSON key file")
	gcpCmd.PersistentFlags().StringSlice("zones", []string{"europe-north1-a", "europe-north1-b", "europe-north1-c"}, "GCP Zones to enumerate")
	viper.BindPFlag("gcp.project", gcpCmd.PersistentFlags().Lookup("project"))
	viper.BindPFlag("gcp.key", gcpCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("gcp.zones", gcpCmd.PersistentFlags().Lookup("zones"))
}
