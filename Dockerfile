FROM golang:1.18-alpine AS build
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build .

FROM alpine:3.16.0
WORKDIR /app

RUN apk add tzdata
RUN adduser --disabled-password --no-create-home reisen

COPY --from=build /app/reisen .
COPY --from=build /app/templates templates
COPY --from=build /app/static static

USER reisen
CMD ./reisen
