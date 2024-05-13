@ECHO OFF
ECHO Install StaticDeployment
ECHO.
ECHO Please enter the installation folder and press Enter.
ECHO Default folder: GOPATH\bin
ECHO GOPATH=%GOPATH%
ECHO.
set /p targetFolder="Install To :"
if "%targetFolder%"=="" set targetFolder=%GOPATH%\bin
ECHO.
ECHO Installing...
DEL "%targetFolder%"\StaticDeployment*.exe
MOVE StaticDeployment*.exe "%targetFolder%"\
setlocal
echo %PATH% | find /I "%targetFolder%" >nul
ECHO Adding %targetFolder% to PATH ...
if errorlevel 1 (
    setx PATH "%PATH%;%targetFolder%"
)
ECHO %targetFolder% has been added to PATH.
ECHO OK
PAUSE