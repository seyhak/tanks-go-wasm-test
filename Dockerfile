FROM golang:1.22-alpine

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify

# COPY . .
# RUN go build -v -o /usr/local/bin/app ./...

CMD ["sh"]

# docker build -t gotanks .
# docker image ls
# docker run --mount type=bind,source=.,target=/usr/src/app -t -i gotanks
