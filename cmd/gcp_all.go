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

	"github.com/kabachook/cirrus/pkg/config"
	"github.com/kabachook/cirrus/pkg/provider/gcp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "List IPs of all resources",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewDevelopment()
		defer logger.Sync()
		sugar := logger.Sugar()

		project, err := cmd.Parent().PersistentFlags().GetString("project")
		if err != nil {
			sugar.Fatal(err)
		}
		key, err := cmd.Parent().PersistentFlags().GetString("key")
		if err != nil {
			sugar.Fatal(err)
		}

		ctx := context.Background()
		provider, err := gcp.New(ctx, config.ConfigGCP{
			Project: project,
			Options: []option.ClientOption{
				option.WithCredentialsFile(key),
			},
			Logger: sugar.Named("gcp").Desugar(),
		})
		if err != nil {
			sugar.Fatal(err)
		}

		endpoints, err := provider.All()
		if err != nil {
			sugar.Fatal(err)
		}

		sugar.Infow("Got endpoints", zap.Any("endpoints", endpoints))

	},
}

func init() {
	gcpCmd.AddCommand(allCmd)
}