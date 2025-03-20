FROM golang:1.24.1

WORKDIR /go/src/app

COPY . .

CMD ["make", "build", "lint"]

ENV PATH="/go/bin:${PATH}"
