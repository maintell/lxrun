SET GOOS=windows
SET GOARCH=amd64
go build -ldflags "-w -s" -o step_amd64.exe
SET GOARCH=386
go build -ldflags "-w -s" -o step_x86.exe
SET GOARCH=arm
go build -ldflags "-w -s" -o step_arm.exe

SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-w -s" -o step_linux_amd64
SET GOARCH=386
go build -ldflags "-w -s" -o step_linux_x86
SET GOARCH=arm
go build -ldflags "-w -s" -o step_linux_arm