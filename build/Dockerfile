FROM golang:1.21-alpine as build

WORKDIR /app

COPY ./src/app/go.mod ./
COPY ./src/app/go.sum ./

RUN go mod download

COPY ./src/app ./

RUN go build -o /htmx-todo-app

FROM scratch as production

WORKDIR /

COPY --from=build /htmx-todo-app /htmx-todo-app

EXPOSE 8080
CMD [ "/htmx-todo-app" ]