@REM Script to compile LaTeX documents using Docker in Windows
@REM Taken from https://gist.github.com/avborup/580bc0b003a702d1a84235890acd5fd7

@echo Compiling file %cd%\%1

@REM Map T:\ to the current path
subst T: "%cd%"

@REM Compile the LaTeX document via a Docker container
docker run ^
  --rm ^
  -i ^
  --net=none ^
  -v T:\:/data ^
  kongborup/custom-latex:latest ^
  pdflatex ^
  --shell-escape ^
  -interaction=nonstopmode ^
  -file-line-error ^
  -output-directory=build ^
  -aux-directory=build ^
  %1

@REM Save exit code of docker command
SET exitcode=%ERRORLEVEL%

@REM Unmap T:\ from the path
subst T: /D

@REM Exit program with exit code of docker command
exit /b %exitcode%
