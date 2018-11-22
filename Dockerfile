## BUILD ##
FROM golang:alpine as builder

# copy the app folder and set the working directory
COPY . $GOPATH/src/github.com/popcore/verisart_exercise
WORKDIR $GOPATH/src/github.com/popcore/verisart_exercise

# build executable
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/verisart

## RELEASE ##
FROM alpine:3.8
COPY --from=builder build/verisart build/verisart
USER appuser
EXPOSE 9091
ENTRYPOINT ["build/verisart"]


