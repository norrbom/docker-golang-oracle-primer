FROM golang:1.16

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y libaio1 wget unzip
RUN wget -O /tmp/instantclient-basic-linux-x64.zip https://download.oracle.com/otn_software/linux/instantclient/193000/instantclient-basic-linux.x64-19.3.0.0.0dbru.zip
RUN mkdir -p /usr/lib/oracle && unzip /tmp/instantclient-basic-linux-x64.zip -d /usr/lib/oracle
RUN ldconfig -v /usr/lib/oracle/instantclient_19_3
RUN ldd /usr/lib/oracle/instantclient_19_3/libclntsh.so

WORKDIR $GOPATH/src/app/
COPY src .
RUN GO111MODULE=on go mod init zero-to-prod.norrbom.org
RUN cat go.mod
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o /go/bin/app

ENTRYPOINT ["/go/bin/app"]

