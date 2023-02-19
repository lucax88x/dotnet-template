package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"dzor/core"
	"dzor/core/azure"
	"dzor/core/dotnet"
	"dzor/core/git"
)

func init() {
	continuosDeliveryCmd.PersistentFlags().String("buildId", "", "Current build id")
	continuosDeliveryCmd.MarkPersistentFlagRequired("buildId")
	viper.BindPFlag("buildId", continuosDeliveryCmd.PersistentFlags().Lookup("buildId"))

	continuosDeliveryCmd.PersistentFlags().String("commitMessage", "", "Commit message for gitops")
	viper.BindPFlag("commitMessage", continuosDeliveryCmd.PersistentFlags().Lookup("commitMessage"))

	rootCmd.AddCommand(continuosDeliveryCmd)
}

var continuosDeliveryCmd = &cobra.Command{
	Use:   "cd",
	Short: "runs the CD task",
	Long: `CD stands for continuos-delivery
  it will take the commit message from azure-devops env variables, if you want to override use --commitMessage`,
	Run: func(cmd *cobra.Command, args []string) {
		core.Wrap(func(ctx core.WrapContext) error {

			sdk := dotnet.CreateSdkContainer(ctx)
			sdk, solutionPath := dotnet.Restore(ctx, sdk)
			sdk = dotnet.Build(ctx, sdk, solutionPath)

			imageTag := getImageTag()
			commitMessage := getCommitMessage(ctx, imageTag)

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

func getCommitMessage(ctx core.WrapContext, imageTag string) string {
	commitMessage := azure.GetCommitMessage()

	cliCommitMessage := viper.GetString("commitMessage")

	if len(cliCommitMessage) > 0 {
		commitMessage = cliCommitMessage
	}

	if len(commitMessage) == 0 {
		ctx.Log.Fatal("empty commit message")
	}

	return fmt.Sprintf("patch: %s %s", commitMessage, imageTag)
}
