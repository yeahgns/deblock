@echo off
cd /d %~dp0
java -Xmx4G -Xms4G -jar server.jar nogui
pause
