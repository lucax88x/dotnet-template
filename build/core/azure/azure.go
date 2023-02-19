package azure

import "os"

func GetCommitMessage() string {
	value := os.Getenv("Build.SourceVersionMessage")

	return value
}
