FROM golang:1.20

WORKDIR /server

COPY go.mod ./
RUN go mod download

COPY . .
# install ffmpeg 
RUN apt-get -y update && apt-get -y upgrade && apt-get install -y --no-install-recommends ffmpeg

RUN go build -o ./bin/server ./cmd/main.go 

EXPOSE 8055

CMD ["./bin/server"]