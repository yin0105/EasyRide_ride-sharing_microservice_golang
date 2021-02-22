FROM golang:apline

WORKDIR /app
ADD . /app

RUN cd /app && go build -o easyride

EXPOSE 8088

ENTRYPOINT ["./easyride"]