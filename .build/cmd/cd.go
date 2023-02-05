package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"dzor/core"
	"dzor/core/dotnet"
)

func init() {
	rootCmd.AddCommand(continuosDeliveryCmd)
	continuosDeliveryCmd.PersistentFlags().String("version", "1.0.0", "Version of the delivery")
	viper.BindPFlag("version", continuosDeliveryCmd.PersistentFlags().Lookup("version"))
}

var continuosDeliveryCmd = &cobra.Command{
	Use:   "cd",
	Short: "runs the CD task",
	Long:  `CD stands for continuos-delivery`,
	Run: func(cmd *cobra.Command, args []string) {
		core.Wrap(func(ctx core.WrapContext) error {
			sdk := dotnet.CreateSdkContainer(ctx)
			sdk, solutionPath := dotnet.Restore(ctx, sdk)
			sdk = dotnet.Build(ctx, sdk, solutionPath)
			dotnet.Dockerize(ctx, sdk)

			return nil
		})
	},
}
