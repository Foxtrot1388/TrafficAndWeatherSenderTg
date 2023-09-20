FROM golang:1.20.7-bullseye AS builder

WORKDIR /usr/local/go/src/
RUN apt update
RUN apt -y install wget
RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-2/wkhtmltox_0.12.6.1-2.bullseye_amd64.deb
RUN apt install -f -y ./wkhtmltox_0.12.6.1-2.bullseye_amd64.deb
ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .
RUN go build -mod=mod -o app.exe cmd/main.go

FROM debian:11

RUN apt update
RUN apt -y install wget
RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
RUN apt install -f -y ./google-chrome-stable_current_amd64.deb
RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-2/wkhtmltox_0.12.6.1-2.bullseye_amd64.deb
RUN apt install -f -y ./wkhtmltox_0.12.6.1-2.bullseye_amd64.deb

COPY --from=builder /usr/local/go/src/app.exe /
COPY --from=builder /usr/local/go/src/app.yaml /

CMD ["/app.exe"]