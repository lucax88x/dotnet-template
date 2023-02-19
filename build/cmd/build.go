package cmd

import (
	"dzor/core"
	"dzor/core/dotnet"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build the projects",
	Long:  `build the projects of the solution`,
	Run: func(cmd *cobra.Command, args []string) {
		core.Wrap(func(ctx core.WrapContext) error {

			sdk := dotnet.CreateSdkContainer(ctx)
			sdk, solutionPath := dotnet.Restore(ctx, sdk)
			sdk = dotnet.Build(ctx, sdk, solutionPath)

			return nil
		})
	},
}
