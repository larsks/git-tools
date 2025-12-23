COMMANDS = git-ptt

all: $(COMMANDS)

tidy:
	go mod tidy

git-ptt: tidy
	go build ./cmd/git-ptt
