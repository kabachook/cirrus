package cmd

import (
	"fmt"

	"github.com/kabachook/cirrus/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var force bool

// NewConfig generates `config dump` command
func NewConfig() *cobra.Command {
	v := viper.GetViper()
	logger := config.Logger
	rootCmd := &cobra.Command{
		Use:   "config",
		Short: "Manipulate config",
	}

	dumpCmd := &cobra.Command{
		Use:   "dump",
		Short: fmt.Sprintf("Dump config to %s", v.ConfigFileUsed()),
		Run: func(cmd *cobra.Command, args []string) {
			if force {
				err := viper.WriteConfig()
				if err != nil {
					logger.Error(err.Error())
				}
			} else {
				err := viper.SafeWriteConfig()
				if err != nil {
					logger.Fatal(err.Error())
				}
			}
		}}

	dumpCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite config files")
	rootCmd.AddCommand(dumpCmd)

	return rootCmd
}
