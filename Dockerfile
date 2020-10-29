# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
#FROM golang

# Copy the local package files to the container's workspace.
#ADD . /usr/src/goadmin

# Set the Current Working Directory inside the container
#WORKDIR $GOPATH/usr/src/goadmin

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
#COPY . .

# Download all the dependencies
#RUN go get -d -v ./...

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
# RUN go install github.com/golang/example/

# Install the package
#RUN go install -v ./...

# Run the outyet command by default when the container starts.
#ENTRYPOINT /go/bin/outyet


FROM golang:onbuild

# Document that the service listens on port 8080.
EXPOSE 8335
EXPOSE 8337
EXPOSE 8339

# Run the executable
#CMD ["spigot2vk_admin"]