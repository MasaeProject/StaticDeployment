SET NAME=StaticDeployment
SET CGO_ENABLED=1
SET GOFLAGS=-buildvcs=false
DEL %NAME%*.exe
CD Join
go build -o ..\%NAME%_Join.exe
CD ..
CD minify
go build -o ..\%NAME%_Minify.exe
CD ..
CD zhcodeconv
go build -o ..\%NAME%_ZhCodeConv.exe
CD ..
go build -o %NAME%.exe
SET CGO_ENABLED=
SET GOFLAGS=
%NAME%.exe test\testconfig.yaml
SET NAME=
