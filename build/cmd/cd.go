package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"dzor/core"
	"dzor/core/dotnet"
	"dzor/core/git"
)

func init() {
	continuosDeliveryCmd.PersistentFlags().String("buildId", "", "Current build id")
	continuosDeliveryCmd.MarkPersistentFlagRequired("buildId")
	viper.BindPFlag("buildId", continuosDeliveryCmd.PersistentFlags().Lookup("buildId"))

	continuosDeliveryCmd.PersistentFlags().String("commitMessage", "", "Commit message for gitops")
	continuosDeliveryCmd.MarkPersistentFlagRequired("commitMessage")
	viper.BindPFlag("commitMessage", continuosDeliveryCmd.PersistentFlags().Lookup("commitMessage"))

	rootCmd.AddCommand(continuosDeliveryCmd)
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

			imageTag := getImageTag()
			commitMessage := getCommitMessage(imageTag)

			ctx.Log.Infof("using image-tag %v", imageTag)

			dotnet.Dockerize(ctx, sdk, imageTag)
			git.PatchGitOps(ctx, sdk, imageTag, commitMessage)

			return nil
		})
	},
}

func getImageTag() string {
	return fmt.Sprintf("v%s.%s", time.Now().Format("2006.0102"), viper.GetString("buildId"))
}

func getCommitMessage(imageTag string) string {
	return fmt.Sprintf("patch: %s %s", viper.GetString("commitMessage"), imageTag)
}
