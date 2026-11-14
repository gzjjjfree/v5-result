package main

import (
	"github.com/gzjjjfree/v5-result/main/commands"
	"github.com/gzjjjfree/v5-result/main/commands/base"
	_ "github.com/gzjjjfree/v5-result/main/distro/all"
)

func main() {
	base.RootCommand.Long = "A unified platform for anti-censorship."
	base.RegisterCommand(commands.CmdRun)
	base.RegisterCommand(commands.CmdVersion)
	base.RegisterCommand(commands.CmdTest)
	base.SortLessFunc = runIsTheFirst
	base.SortCommands()
	base.Execute()
}

func runIsTheFirst(i, j *base.Command) bool {
	left := i.Name()
	right := j.Name()
	if left == "run" {
		return true
	}
	if right == "run" {
		return false
	}
	return left < right
}
