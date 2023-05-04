FROM golang:1.19-alpine AS build
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build .

FROM alpine:3.17.3
WORKDIR /app

RUN apk add tzdata
RUN adduser --disabled-password --no-create-home reisen

COPY --from=build /app/reisen .
COPY --from=build /app/templates templates
COPY --from=build /app/static static
RUN chown -R reisen:reisen /app/static

USER reisen

CMD ./reisen
