$ git clone https://github.com/wallneryan/fli-docker
$ cd ~/fli-docker/build/
$ vi Dockerfile
$ docker build --no-cache --tag clusterhq/fli-docker-build .
$ docker run --rm -it \
    -v <host-path-to-fli-docker>:/go/src/github.com/wallnerryan/fli-docker \
    -v $HOME:/output \
    clusterhq/fli-docker-build

