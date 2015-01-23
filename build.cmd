if EXIST ccTray2Slack.exe (
	del ccTray2Slack.exe
)

go build

@if NOT "%ERRORLEVEL%" == "0" (
	@ECHO It's gone wrong :-( %ERRORLEVEL%
	@GOTO END
)

@SET XMLURL=http://ci.internal.comparethemarket.local:8153/go/cctray.xml
@REM @SET XMLURL=http://localhost/ccTray.xml

ccTray2Slack.exe -url %XMLURL% -config watch.json -username %CC_USR% -password %CC_PWD%

:END