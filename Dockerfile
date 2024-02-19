FROM golang:1.21.4 AS build
WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o /bin/app

FROM gcr.io/distroless/static-debian11
COPY --from=build /bin/app /bin/app
COPY .env.prod /bin/.env.prod

EXPOSE 8080

CMD ["/bin/app", "/bin/.env.prod" ]