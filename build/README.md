How to build the project

```
$ git clone https://github.com/ClusterHQ/fli-docker

$ cd ~/fli-docker/build/

$ docker build --no-cache --tag clusterhq/fli-docker-build .

$ docker run --rm -it \
    -v <host-path-to-fli-docker>:/go/src/github.com/ClusterHQ/fli-docker \
    -v /tmp/:/output \
    clusterhq/fli-docker-build
```

The binary will be available at  `/tmp/fli-docker`
