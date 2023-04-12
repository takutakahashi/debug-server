FROM golang:1.19 as builder

ADD . /
RUN go mod download
RUN go build -o cmd ./main.go

FROM ubuntu
COPY --from=builder cmd /
CMD ["/cmd"]
