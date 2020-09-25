FROM golang
ENV REPO=github.com/Jakatut/NAD-A3
RUN apt-get update && apt-get install -y ca-certificates git-core ssh
RUN  mkdir -p /go/src/$REPO \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
ENV APP=$GOPATH/src/$REPO
ENV PORT=8080
ADD ./app $APP
ADD ./app/messages $APP
ADD ./app/routes $APP
ADD ./app/handlers $APP

WORKDIR $APP
EXPOSE 8080

RUN go build -o NAD-A3 .
CMD ["./docker-go-starter-kit"]