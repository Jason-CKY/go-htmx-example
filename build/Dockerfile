FROM node:18-alpine as dev

WORKDIR /app

RUN wget https://go.dev/dl/go1.21.3.linux-amd64.tar.gz &&\
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz &&\
    rm go1.21.3.linux-amd64.tar.gz

ENV PATH="${PATH}:/usr/local/go/bin:/root/go/bin"

RUN go install github.com/a-h/templ/cmd/templ@latest &&\
    go install github.com/cosmtrek/air@latest

ENTRYPOINT ["air"]

FROM node:18-alpine as build-css

WORKDIR /app
COPY . .
RUN npm ci && npx tailwindcss -i ./build/input.css -o /output.css


FROM golang:1.21-alpine as production

RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

COPY ./src/app/go.mod ./
COPY ./src/app/go.sum ./

RUN go mod download

COPY ./src/app ./
COPY --from=build-css /output.css ./static/output.css

RUN templ generate
RUN go build -o /main
CMD [ "/main" ]
