# IMPORTANT: You will need to log in with Docker first:
#   echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
docker build -t kongborup/custom-latex .
docker push kongborup/custom-latex:latest
