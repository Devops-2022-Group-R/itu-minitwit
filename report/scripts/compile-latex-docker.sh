# Be careful about line endings! If the file uses CRLF, bash might not be able to run the script
mkdir -p ./build
docker run \
  --rm \
  -i \
  --net=none \
  -v $(pwd):/data \
  kongborup/custom-latex:latest \
  pdflatex \
  --shell-escape \
  -interaction=nonstopmode \
  -file-line-error \
  -output-directory=build \
  -aux-directory=build \
  $1
