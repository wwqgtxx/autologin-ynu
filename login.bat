@echo off
cd %~dp0
if %processor_architecture%==x86 (
set AUTOLOGIN="%~dp0/autologin-windows-386.exe"
) else (
set AUTOLOGIN="%~dp0/autologin-windows-amd64.exe"
)
%AUTOLOGIN%
pause