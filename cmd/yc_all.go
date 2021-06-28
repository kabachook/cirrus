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
	"context"

	cmdGen "github.com/kabachook/cirrus/pkg/cmd"
	"github.com/kabachook/cirrus/pkg/provider/yc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var ycAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List IPs of all resources",
	Run: func(cmd *cobra.Command, args []string) {
		v := viper.GetViper()
		ctx := context.Background()

		folderId := v.GetString(cmdGen.YcFolderId)
		token := v.GetString(cmdGen.YcToken)
		zones := v.GetStringSlice(cmdGen.YcZones)
		output, err := cmd.Parent().PersistentFlags().GetString("output")
		if err != nil {
			logger.Error(err.Error())
			return
		}

		provider, err := yc.New(ctx, yc.Config{
			FolderID: folderId,
			Token:    token,
			Zones:    zones,
			Logger:   logger.Named(yc.Name),
		})
		if err != nil {
			logger.Error(err.Error())
			return
		}

		endpoints, err := provider.All()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		logger.Info("Got endpoints")
		switch output {
		case "text":
			for _, endpoint := range endpoints {
				logger.Sugar().Infof("type: %s\tname: %s\tip: %s", endpoint.Type, endpoint.Name, endpoint.IP)
			}
		case "json":
			logger.Info("", zap.Any("endpoints", endpoints))
		}

	},
}

func init() {
	ycCmd.AddCommand(ycAllCmd)
}
