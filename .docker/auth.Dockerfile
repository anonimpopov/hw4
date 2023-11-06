FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/auth/app ./app

CMD ["/app/app", "-c", "config.yaml"]
