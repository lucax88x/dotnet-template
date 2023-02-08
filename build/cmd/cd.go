package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"dzor/core"
	"dzor/core/dotnet"
)

func init() {
	// viper.BindPFlag("buildId", continuosDeliveryCmd.PersistentFlags().Lookup("buildId"))

	rootCmd.AddCommand(continuosDeliveryCmd)

	continuosDeliveryCmd.PersistentFlags().String("buildId", "", "Current build id")
	continuosDeliveryCmd.MarkPersistentFlagRequired("buildId")
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

			imageTag := core.GetImageTag()

			dotnet.Dockerize(ctx, sdk, imageTag)
			dotnet.PatchGitOps(ctx, sdk, imageTag)

			return nil
		})
	},
}
