# https://gist.github.com/naesheim/18d0c0a58ee61f4674353a2f4cf71475

set -e

# latest commit
LATEST_COMMIT=$(git rev-parse HEAD)

# latest commit where path/to/folder1 was changed
REPORT_COMMIT=$(git log -1 --format=format:%H --full-diff /report)

if [ $REPORT_COMMIT = $LATEST_COMMIT ];
    then
        echo "files in Report has changed"
        #   ./scripts/compile-latex-docker.sh main.tex
        #   ./scripts/compile-latex-docker.bat main.tex
        #   .circleci/do_something.sh
else
     echo "no folders of relevance has changed"
     exit 0;
fi