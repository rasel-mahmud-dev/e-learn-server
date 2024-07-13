FROM golang:1.22.5-bullseye

WORKDIR /usr/src/app

COPY . .




# RUN go mod download

RUN go mod tidy


RUN go install github.com/air-verse/air@latest

# RUN air init

#CMD ["go", "run", "main.go"]

CMD ["air", "-c", ".air.toml"]
