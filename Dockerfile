FROM golang:1.19 AS build

RUN apt-get update && apt-get install --no-install-recommends -y \
    upx \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/aws-ecr-cleaner .

RUN upx --best --lzma /usr/local/bin/aws-ecr-cleaner

FROM ubuntu:22.04

RUN apt-get update && apt-get install --no-install-recommends -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /usr/local/bin/aws-ecr-cleaner /usr/local/bin/aws-ecr-cleaner

CMD [ "/usr/local/bin/aws-ecr-cleaner" ]
