SET NAME=NyarukoClipboard
SET NAMEV=%NAME%_v1.0.0_
MD bin
DEL /Q bin\*
COPY README.md bin\
SET CGO_ENABLED=0
SET GOARCH=amd64
go generate
SET GOOS=windows
go build -o bin\%NAMEV%Windows64.exe .
DEL /Q *.syso
SET GOOS=linux
go build -o bin\%NAMEV%Linux64 .
SET GOOS=darwin
SET GOARCH=amd64
go build -o bin\%NAMEV%macOSI64.
SET GOARCH=arm64
go build -o bin\%NAMEV%macOSM64.
SET GOOS=windows
go build -o bin\%NAMEV%WindowsARM64.exe .
DEL /Q *.syso
SET GOARCH=386
go generate
SET GOOS=windows
go build -o bin\%NAMEV%Windows32.exe .
DEL /Q *.syso
SET GOOS=linux
go build -o bin\%NAMEV%Linux32 .
SET NAME=
SET NAMEV=
SET CGO_ENABLED=
SET GOARCH=
SET GOOS=
