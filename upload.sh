set -e
docker build . -t rofleksey/shantaram-cms:latest --platform linux/amd64
docker push rofleksey/shantaram-cms:latest
