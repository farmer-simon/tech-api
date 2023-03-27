### Build 命令
```text

go env -w GOARCH=amd64
go env -w CGO_ENABLED=0

go env -w GOOS=linux
go build -o tech-build  -ldflags "-w -s"  cmd/api/main.go

go env -w GOOS=windows


```

