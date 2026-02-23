FROM golang:1.25.6

WORKDIR /go/src/app

COPY . .

CMD ["make", "build", "lint"]

ENV PATH="/go/bin:${PATH}"
