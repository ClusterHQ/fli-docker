FROM golang:latest

VOLUME /output

RUN go get gopkg.in/yaml.v2
RUN go get github.com/docker/libcompose
RUN go get golang.org/x/net/context

ENTRYPOINT ["go", "build", "-o", "/output/fli-docker", "github.com/ClusterHQ/fli-docker/cmd/fli-docker"]

