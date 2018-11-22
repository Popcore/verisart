## BUILD ##
FROM golang:alpine as build-env

# copy the app folder and set the working directory
WORKDIR /go/src/github.com/popcore/verisart_exercise
COPY . .

# build executable
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o verisart_app

## RELEASE ##
FROM alpine:3.8
COPY --from=build-env /go/src/github.com/popcore/verisart_exercise .
EXPOSE 9091
ENTRYPOINT ./verisart_app


