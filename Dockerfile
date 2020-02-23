FROM golang:1.12.7-alpine

# install os dependencies & utilitiess
RUN apk update && apk add --no-cache git=2.22.2-r0

RUN git clone https://github.com/burubur/helloworld.git /service/helloworld
WORKDIR "/service/helloworld"
RUN COMMITHASH=$(git rev-parse --short HEAD) && go build -o helloworld -ldflags "-X main.CommitHash=$COMMITHASH" .
EXPOSE 8080
ENTRYPOINT ["/service/helloworld/helloworld"]
