
echo off

echo 'build drawin'
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o ./a_build_dist/drawin/Andex Andex.go

echo '========================'

echo 'build linux_amd64'
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o  ./a_build_dist/linux/Andex Andex.go

echo '========================'

echo 'build windows_amd64'
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o ./a_build_dist/win/Andex.exe Andex.go

echo '========================'