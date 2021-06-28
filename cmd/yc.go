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
	cmdGen "github.com/kabachook/cirrus/pkg/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ycCmd = &cobra.Command{
	Use:   "yc",
	Short: "Yandex Cloud",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(ycCmd)

	ycCmd.PersistentFlags().StringP("output", "o", "text", "Output format")

	ycCmd.PersistentFlags().String("folder-id", "", "folderId")
	ycCmd.PersistentFlags().String("token", "", "OAuth token")
	ycCmd.PersistentFlags().StringSlice("zones", []string{"ru-central1-a", "ru-central1-b", "ru-central1-c"}, "Zones to enumerate")
	viper.BindPFlag(cmdGen.YcFolderId, ycCmd.PersistentFlags().Lookup("folder-id"))
	viper.BindPFlag(cmdGen.YcToken, ycCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag(cmdGen.YcZones, ycCmd.PersistentFlags().Lookup("zones"))
}
