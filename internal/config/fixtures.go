package config

import (
	"os"
)

func GetFixtureDirectories() []string {
	// used for the vhs tape recording stuff
	if os.Getenv("JUMPER_FIXTURES") == "" {
		return nil
	}

	return []string{
		"~/dev/blog",
		"~/dev/habit-tracker",
		"~/dev/todo-list",
		"~/dev/notes-service",
		"~/dev/sass-project",
		"~/dev/sandbox/leetcode-practice",
		"~/dev/website",
		"~/work/internal-tools",
		"~/work/billing-service",
		"~/work/platform-infra",
		"~/work/reporting-dashboard",
	}
}
