package gitcommand

import (
	"bufio"
	"fmt"
	"iter"
	"os"
	"os/exec"
	"strings"
)

type (
	GitCommand struct {
		gitPath string
	}
)

func NewGit() *GitCommand {
	return &GitCommand{
		gitPath: "git",
	}
}

func (g *GitCommand) WithGitPath(path string) *GitCommand {
	g.gitPath = path
	return g
}

func (g *GitCommand) Command(args ...string) *exec.Cmd {
	return exec.Command(g.gitPath, args...)
}

func (g *GitCommand) GetTrailers(commitSHA string) (map[string]string, error) {
	// Get the commit message body using git log
	cmdLog := g.Command("log", "-1", commitSHA, "--pretty=format:%B")
	output, err := cmdLog.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit message: %w", err)
	}
	commitMsg := string(output)

	// Use git interpret-trailers to parse the output
	// --parse option outputs each trailer on its own line
	cmdInterpret := g.Command("interpret-trailers", "--parse")
	cmdInterpret.Stdin = strings.NewReader(commitMsg)
	trailersOutput, err := cmdInterpret.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to interpret trailers: %w", err)
	}

	trailers := make(map[string]string)
	lines := strings.SplitSeq(strings.TrimSpace(string(trailersOutput)), "\n")

	for line := range lines {
		// Trailers use "Key: Value" format by default
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" {
				trailers[key] = value
			}
		}
	}
	return trailers, nil
}

func (g *GitCommand) Run(args ...string) error {
	cmd := g.Command(args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (g *GitCommand) RevList(from, to string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		// Get the commit message body using git log
		cmd := g.Command("rev-list", "--reverse", fmt.Sprintf("%s..%s", from, to))
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			if !yield("", err) {
				return
			}
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			if !yield("", err) {
				return
			}
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if !yield(scanner.Text(), nil) {
				return
			}
		}

		if err := cmd.Wait(); err != nil {
			if !yield("", err) {
				return
			}
		}
	}
}
