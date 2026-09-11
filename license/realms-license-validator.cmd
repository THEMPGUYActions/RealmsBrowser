@echo off
setlocal
if "%REALMS_LICENSE_SSH_AUTH_FILE%"=="" exit /b 2
ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=accept-new -i "%REALMS_LICENSE_SSH_AUTH_FILE%" realmsnetwork@ssh-realmsnetwork.alwaysdata.net %*
exit /b %ERRORLEVEL%
