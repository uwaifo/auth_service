FROM golang:alpine as build-env

RUN apk add --update --no-cache ca-certificates git

WORKDIR /app

# set src directory
ENV SRC_DIR=/app_src

# set cgo enabled to 0
ENV CGO_ENABLED 0

# set the port of the application
ENV PORT ":80"

# create the SRC directory
RUN mkdir $SRC_DIR

# add source code to directory
ADD . $SRC_DIR

# install added dependency in the src directory
RUN cd $SRC_DIR; go mod download

# Build and compile the go application
RUN cd $SRC_DIR; go build -o app-auth

# copy the executable and webapp dir into workdir
RUN cd $SRC_DIR; cp app-auth jwtRS256.key jwtRS256.key.pub /app/

# expose port 8000
EXPOSE 80

ENTRYPOINT ./app-auth
