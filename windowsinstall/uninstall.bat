@ECHO OFF
ECHO Uninstall StaticDeployment
PAUSE
ECHO.
DEL StaticDeployment*.exe

SET "remove_path=%~dp0"
SET "remove_path=%remove_path:~0,-1%"

SET "new_path="
for %%p in (%PATH%) do (
    if /I not "%%p"=="%remove_path%" (
        SET "new_path=!new_path!;%%p"
    )
)
ECHO Removing %remove_path% from PATH ...
SETx PATH "%new_path:~1%"
ECHO %remove_path% has been removed from PATH.

PAUSE