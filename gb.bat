@echo off
set d=%~dp0
set d=%d:~0,-1%
for %%f in (%d%) do set myfolder=%%~nxf
echo %myfolder%
FOR /F "delims=" %%i IN ('git config user.name') DO set user=%%i
echo %user%
set ghproject=github.com\%user%\%myfolder%
if not "%GOPATH%" == "" (
	if not exist "%GOPATH%\src\%ghproject%" (
		mkdir "%GOPATH%\src\github.com\%user%" 2>NUL
		mklink /J "%GOPATH%\src\%ghproject%" "%~dp0\deps\src\%ghproject%"
	)
	if not exist "%GOPATH%\pkg\windows_amd64\%ghproject%" (
		mkdir "%GOPATH%\pkg\windows_amd64\github.com\%user%" 2>NUL
		mklink /J "%GOPATH%\pkg\windows_amd64\%ghproject%" "%~dp0\deps\pkg\windows_amd64\%ghproject%"
	)
)
setlocal
if not exist "%~dp0\deps\src\%ghproject%" (
	mkdir %~dp0\deps\src\github.com\%user% 2>NUL
	mklink /J "%~dp0\deps\src\%ghproject%" "%~dp0"
)
set GOPATH=%~dp0deps
set GOBIN=%~dp0bin
cd "%~dp0\deps\src\%ghproject%"
call go build . github.com/%user%/%myfolder% || goto eof
go test -coverprofile=coverage.out github.com/%user%/%myfolder%|grep -v -e "^\.\.*$"|grep -v "^$"|grep -v "thus far"
set msg=doskey %myfolder%=%~dp0bin\%myfolder%.exe $*
doskey /macros:all|grep %myfolder% 1>NUL || echo %msg% |clip && echo %msg%
endlocal
set d=
set myfolder=
set user=
