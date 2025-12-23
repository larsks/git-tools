package main

import (
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/larsks/git-tools/internal/gitcommand"
	"github.com/larsks/gobot/tools"
)

type (
	Config struct {
		UseCommit bool
	}
)

var config = Config{}
var git *gitcommand.GitCommand

func init() {
	flag.BoolVarP(&config.UseCommit, "commit", "c", false, "Edit files in current commit")
}

func main() {
	flag.Parse()
	git = gitcommand.NewGit()
	editor := tools.GetenvWithDefault("VISUAL", tools.GetenvWithDefault("EDITOR", "vi"))

	ref := "HEAD"
	if config.UseCommit {
		ref = "HEAD^"
	}

	output, err := git.Output("diff", "--name-only", "-z", ref)
	if err != nil {
		log.Fatalf("failed to list files: %v", err)
	}

	files_to_edit := slices.DeleteFunc(strings.Split(string(output), "\000"), func(s string) bool {
		return strings.TrimSpace(s) == ""
	})

	if len(files_to_edit) <= 0 {
		log.Fatalf("no files to edit")
	}

	cmd := exec.Command(editor, files_to_edit...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run editor: %v", err)
	}
}
