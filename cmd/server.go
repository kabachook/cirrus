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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kabachook/cirrus/pkg/config"
	"github.com/kabachook/cirrus/pkg/provider"
	"github.com/kabachook/cirrus/pkg/provider/gcp"
	"github.com/kabachook/cirrus/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var logger = config.Logger

var (
	listen           string
	providersEnabled []string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run as server",
	Run: func(cmd *cobra.Command, args []string) {
		v := viper.GetViper()
		ctx := context.Background()

		var providers []provider.Provider

		if len(providersEnabled) == 0 {
			logger.Error("No providers enabled")
			return
		}

		for _, name := range providersEnabled {
			switch name {
			case "gcp":
				p, err := gcp.New(ctx, gcp.Config{
					Project: v.GetString("gcp.project"),
					Options: []option.ClientOption{
						option.WithCredentialsFile(v.GetString("gcp.key")),
					},
					Zones:  v.GetStringSlice("gcp.zones"),
					Logger: logger.Named("gcp"),
				})
				if err != nil {
					logger.Error(err.Error())
					return
				}
				providers = append(providers, p)
			}
		}

		server, err := server.New(ctx, server.Config{
			Logger: logger.Named("server"),
			Server: &http.Server{
				Addr: listen,
			},
			Providers: providers,
		})
		if err != nil {
			logger.Error(err.Error())
			return
		}

		go listenToSystemSignals(ctx, server)

		if err = server.Run(); err != nil {
			logger.Error(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().StringVar(&listen, "listen", ":3232", "Address to listen to")
	serverCmd.PersistentFlags().StringSliceVar(&providersEnabled, "providers", []string{"gcp"}, "Providers enabled")
	viper.BindPFlag("server.listen", serverCmd.PersistentFlags().Lookup("listen"))
	viper.BindPFlag("server.providers", serverCmd.PersistentFlags().Lookup("providers"))
}

func listenToSystemSignals(ctx context.Context, s *server.Server) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	//lint:ignore S1000 side-effects
	for {
		select {
		case <-signalChan:
			logger.Info("Shutting down...")
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			if err := s.Shutdown(ctx); err != nil {
				logger.Error("Timed out waiting for server to shut down")
			}
			return
		}
	}
}
