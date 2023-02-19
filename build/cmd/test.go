package cmd

import (
	"github.com/spf13/cobra"

	"dzor/core"
	"dzor/core/dotnet"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "run tests",
	Long:  `run tests of the solution`,
	Run: func(cmd *cobra.Command, args []string) {
		core.Wrap(func(ctx core.WrapContext) error {
			sdk := dotnet.CreateSdkContainer(ctx)
			sdk, solutionPath := dotnet.Restore(ctx, sdk)
			sdk = dotnet.Build(ctx, sdk, solutionPath)
			sdk = dotnet.RunUnitTests(ctx, sdk, solutionPath)
			sdk = dotnet.RunIntegrationTests(ctx, sdk, solutionPath)

			return nil
		})
	},
}
