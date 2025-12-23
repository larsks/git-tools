package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/larsks/gobot/tools"
	flag "github.com/spf13/pflag"

	"github.com/larsks/git-tools/internal/gitcommand"
)

type (
	Config struct {
		RemoteName  string
		DryRun      bool
		Force       bool
		TrailerName string
	}
)

var config = Config{}
var git *gitcommand.GitCommand

func init() {
	flag.StringVarP(&config.RemoteName, "remote", "r", tools.GetenvWithDefault("GIT_PTT_REMOTE", "origin"), "Name of target remote")
	flag.StringVarP(&config.TrailerName, "trailer-name", "t", tools.GetenvWithDefault("GIT_PTT_TRAILER", "x-branch"), "Git trailer containing branch name")
	flag.BoolVarP(&config.DryRun, "dry-run", "n", false, "Do not actually push anything")
	flag.BoolVarP(&config.Force, "force", "f", false, "Use --force instead of --force-with-lease")
}

func main() {
	flag.Parse()
	git = gitcommand.NewGit()

	var from, to string
	if flag.NArg() < 1 {
		from = os.Getenv("GIT_PTT_FROM")
		if from == "" {
			log.Fatalf("Missing required <from> argument")
		}
	} else {
		from = flag.Arg(0)
	}

	if strings.Contains(from, "..") {
		parts := strings.SplitN(from, "..", 2)
		from, to = parts[0], parts[1]
	} else {
		to = "HEAD"
	}

	for cid, err := range git.RevList(from, to, true) {
		if err != nil {
			log.Fatalf("failed to get revision list: %v", err)
		}

		trailers, err := git.GetTrailers(cid)
		if err != nil {
			log.Fatalf("failed to get trailed for %s: %v", cid, err)
		}

		if remote_branch, exists := trailers["x-branch"]; exists {
			log.Printf("push commit %s to %s %s", cid, config.RemoteName, remote_branch)

			if !config.DryRun {
				var force_arg string
				if config.Force {
					force_arg = "--force"
				} else {
					force_arg = "--force-with-lease"
				}

				if err := git.Run("push", force_arg, config.RemoteName, fmt.Sprintf("%s:%s", cid, remote_branch)); err != nil {
					log.Fatalf("failed to push %s to %s %s", cid, config.RemoteName, remote_branch)
				}
			}
		}
	}
}
