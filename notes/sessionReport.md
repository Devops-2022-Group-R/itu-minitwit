The plan so far
- [x] Find changes in the report directory by using git diff between latest commit and previous 
- [x] build main.LaTeX by using the script 
- [x] add changes and commit with [skip ci] 
    - [x] ensure using the correct github credentials
    - [x] ensure the branch is correct use the env variables correct
- [x] push the changes to github



Issue with files not mounting into docker volume [CircleCi docs](https://support.circleci.com/hc/en-us/articles/360007324514-How-can-I-use-Docker-volume-mounting-on-CircleCI-).

Fixed by using a machine executor instead [CircleCi Executor types](https://circleci.com/docs/2.0/executor-types/#using-machine)
