FROM golang:1.18

WORKDIR /var/www/html/dysn/auth

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
#RUN go mod init github.com/DYSN-Project/auth && go mod tidy

COPY . .
#COPY ../. /var/www/html/dysn/auth

#RUN go build -v -o /usr/local/bin/app ./...

EXPOSE 8080

CMD ["go","run","main.go"]