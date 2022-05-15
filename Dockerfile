FROM golang:1.18-alpine3.15

# install os dependencies & utilitiess
RUN apk update && apk add --no-cache git

# clone a fresh public repository
RUN git clone https://github.com/burubur/helloworld.git /opt/helloworld
WORKDIR "/opt/helloworld"

# build the image with commithash as tag
RUN COMMITHASH=$(git rev-parse --short HEAD) && go build -o helloworld -ldflags "-X main.CommitHash=$COMMITHASH" .

# open 8080 port
EXPOSE 8080

# run the helloworld API Server
ENTRYPOINT ["/opt/helloworld/helloworld"]
