FROM golang:1.20

WORKDIR /server

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /cmd

EXPOSE 8055

CMD ["/cmd"]