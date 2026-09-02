FROM golang:1.26.1

WORKDIR /go/src/app

COPY . .

CMD ["make", "build", "lint"]

ENV PATH="/go/bin:${PATH}"
