START /WAIT CMD /C buildclean.bat
CD Minify
START CMD /C build.bat
CD ..
CD ZhCodeConv
START CMD /C build.bat
CD ..
CD Join
START CMD /C build.bat
CD ..
CD PluginDemo
START CMD /C build.bat
CD ..
START /WAIT CMD /C build.bat
PAUSE
XCOPY /S /Y Join\bin bin
RD /S /Q Join\bin
XCOPY /S /Y Minify\bin bin
RD /S /Q Minify\bin
XCOPY /S /Y PluginDemo\bin bin
RD /S /Q PluginDemo\bin
XCOPY /S /Y ZhCodeConv\bin bin
RD /S /Q ZhCodeConv\bin
CD bin
for /d %%D in (*) do (
    7z a -tzip -mx=9 "%%D.zip" "%%D"
    RD /S /Q "%%D"
)
CD ..
ECHO %GOPATH%
