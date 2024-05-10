START /WAIT CMD /C buildclean.bat
CD Minify
START CMD /C build.bat
CD ..
CD ZhCodeConv
START CMD /C build.bat
CD ..
CD PluginDemo
START CMD /C build.bat
CD ..
START /WAIT CMD /C build.bat
PAUSE
CD bin
for /d %%D in (*) do (
    7z a -tzip -mx=9 "%%D.zip" "%%D"
    RD /S /Q "%%D"
)
CD ..
