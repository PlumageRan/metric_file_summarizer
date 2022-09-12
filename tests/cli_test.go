package tests

import (
	"AccelByteTakeHome/cli"
	"testing"
	"time"
)

func TestGenerateSummary(t *testing.T) {
	cmd := cli.CommandLine{}
	cmd.GenerateSummary("/Users/haoranzhou/Documents/AccelByteTakeHome",
		"json",
		time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.January, 3, 23, 59, 59, 0, time.UTC),
		"",
		"", "")
}
