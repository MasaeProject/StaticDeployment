CD Minify
START build.bat
CD ..
CD ZhCodeConv
START build.bat
CD ..
build.bat
PAUSE
CD bin
for /d %%D in (*) do (
    7z a -tzip -mx=9 "%%D.zip" "%%D"
    RD /S /Q "%%D"
)
CD ..