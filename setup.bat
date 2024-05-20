@echo off
go run . setup

echo Press any key to continue...
pause > nul

echo Building...

call build.bat

powershell -ExecutionPolicy Bypass -File .\StartReg.ps1 .\bytebrother.exe

echo Now just run bytebrother.exe
pause > nul
