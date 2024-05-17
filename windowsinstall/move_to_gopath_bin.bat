DEL "%GOPATH%\bin\StaticDeployment*.exe"
MOVE "StaticDeployment*.exe" "%GOPATH%\bin"
DEL *.md
START CMD /C DEL move_to_gopath_bin.bat
EXIT