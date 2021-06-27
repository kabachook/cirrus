package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var force bool

func NewConfig() *cobra.Command {
	v := viper.GetViper()
	rootCmd := &cobra.Command{
		Use:   "config",
		Short: "Manipulate config",
	}

	dumpCmd := &cobra.Command{
		Use:   "dump",
		Short: fmt.Sprintf("Dump config to %s", v.ConfigFileUsed()),
		Run: func(cmd *cobra.Command, args []string) {
			logger, _ := zap.NewDevelopment()
			defer logger.Sync()
			sugar := logger.Sugar()

			if force {
				err := viper.WriteConfig()
				if err != nil {
					sugar.Fatal(err)
				}
			} else {
				err := viper.SafeWriteConfig()
				if err != nil {
					sugar.Fatal(err)
				}
			}
		}}

	dumpCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite config files")

	rootCmd.AddCommand(dumpCmd)
	return rootCmd
}
