FROM golang:1.26.9

WORKDIR /go/src/app

COPY . .

CMD ["make", "build", "lint"]

ENV PATH="/go/bin:${PATH}"
