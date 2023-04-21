build:
	GOOS=linux go build -o  build/bingchat_linux_amd64 cmd/main.go
	GOOS=darwin go build -o  build/bingchat_drawin_amd64 cmd/main.go
	GOOS=windows go build -o  build/bingchat_windows_amd64.exe cmd/main.go