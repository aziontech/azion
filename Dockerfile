FROM golang:1.25.9

WORKDIR /go/src/app

COPY . .

CMD ["make", "build", "lint"]

ENV PATH="/go/bin:${PATH}"
