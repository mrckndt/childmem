FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY childmem.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o childmem .


FROM alpine:latest

ENV INTERVAL=30
ENV PNAME="mattermost"

# create non-root user
RUN addgroup -g 1000 childmem && \
    adduser -D -u 1000 -G childmem childmem

WORKDIR /app

COPY --from=builder /build/childmem .
COPY build/loop.sh .
RUN chmod +x loop.sh childmem && \
    chown -R childmem:childmem /app

USER childmem

RUN mkdir -p /app/data

CMD ["/bin/sh", "-c", "./loop.sh ${INTERVAL} \"/app/childmem -pname ${PNAME} -output \\\"/app/data/child_mem.csv\\\" -includeParent\""]
