# Build.
FROM golang:1.23 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o ./tmp/main cmd/compSudoku/main.go


#Deploy
FROM gcr.io/distroless/static-debian11 AS release-stage

ARG SUDOKU_API_KEY

ENV SUDOKU_API_KEY=SUDOKU_API_KEY
WORKDIR /
COPY --from=build-stage app/tmp/main /tmp/main
COPY --from=build-stage /app/web/public ./web/public
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/tmp/main"]