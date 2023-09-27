FROM golang:1.20

WORKDIR /server

COPY go.mod ./
RUN go mod download

COPY . .
RUN apt-get -y update && apt-get -y upgrade && apt-get install -y --no-install-recommends ffmpeg

RUN go build -o /cmd

EXPOSE 8055

CMD ["/cmd"]