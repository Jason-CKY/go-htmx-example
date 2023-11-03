FROM cosmtrek/air:v1.49.0 as dev

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@latest

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
