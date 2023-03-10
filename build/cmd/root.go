package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "Use debug mode")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

var rootCmd = &cobra.Command{
	Use:   "dzor",
	Short: "dzor is a builder CLI",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
