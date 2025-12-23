COMMANDS = git-ptt git-resume

all: $(COMMANDS)

tidy:
	go mod tidy

git-ptt: tidy
	go build ./cmd/git-ptt

git-resume: tidy
	go build ./cmd/git-resume
