# Build css
FROM node:alpine AS build-css

WORKDIR /app

COPY . .

RUN npm install

RUN npx tailwindcss -i ./tailwind.css -o ./styles.css

# Build a binary
FROM golang:1.23 AS build-app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

## copy the whole source
COPY . .

## copy built css from build-css
COPY --from=build-css /app/styles.css ./handlers/assets/styles.css

## generate go code
RUN go generate ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/blogserver ./cmd/server/main.go

# Run the server
FROM scratch

COPY --from=build-app /app/blogserver /blogserver

ENTRYPOINT [ "/blogserver" ]
CMD []
