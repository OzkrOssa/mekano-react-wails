linux:
	@wails build -platform linux/amd64 -ldflags "-s -w" -upx -clean
windows:
	@wails build -platform windows/amd64 -ldflags "-s -w" -upx -nsis -clean