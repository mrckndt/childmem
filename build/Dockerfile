FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY childmem.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o childmem .


FROM alpine:latest

ENV INTERVAL=30
ENV PNAME="mattermost"

WORKDIR /app

COPY --from=builder /build/childmem .
COPY build/loop.sh .
RUN chmod +x loop.sh childmem
RUN mkdir -p /data

CMD ./loop.sh ${INTERVAL} "/app/childmem \
    -pname ${PNAME} \
    -output \"/data/child_mem.csv\" \
    -includeParent"
