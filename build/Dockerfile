FROM golang:latest

VOLUME /output

RUN go get gopkg.in/yaml.v2

ENTRYPOINT ["go", "build", "-o", "/output/fli-docker", "github.com/wallnerryan/fli-docker/cmd/fli-docker"]

