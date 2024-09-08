# Build css
FROM node:alpine AS build-css

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY tailwind.config.js ./
COPY tailwind.css ./
COPY **/*.templ ./

RUN mkdir static
RUN npx tailwindcss -i ./tailwind.css -o ./static/styles.css

# Build a binary
FROM golang:1.22 AS build-app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN go install github.com/a-h/templ/cmd/templ@latest

## copy the whole source
COPY . .

## copy built css from build-css
COPY --from=build-css /app/static/styles.css ./static/styles.css

## generate go code from .templ files
RUN go generate ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/blogserver ./cmd/server/main.go

# Run the server
FROM scratch

COPY --from=build-app /app/blogserver /blogserver

CMD ["/blogserver"]
