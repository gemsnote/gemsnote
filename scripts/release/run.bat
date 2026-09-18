@echo off
setlocal

set "ROOT=%~dp0"
set "APP_LINK=%ROOT%runtime\github.com\gemsnote\gemsnote"

if not exist "%APP_LINK%" mklink /J "%APP_LINK%" "%ROOT%" >nul

"%ROOT%bin\gemsnote.exe" -importPath github.com/gemsnote/gemsnote -srcPath "%ROOT%runtime" -runMode prod %*
