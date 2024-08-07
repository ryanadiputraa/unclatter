FROM golang:1.22-alpine AS build

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Run config script
RUN chmod +x ./config/config.sh

ARG PORT
ARG FE_URL
ARG POSTGRES_HOST
ARG POSTGRES_PORT
ARG POSTGRES_USER
ARG POSTGRES_PASSWORD
ARG POSTGRES_DB
ARG JWT_SECRET
ARG GOOGLE_REDIRECT_URL
ARG GOOGLE_CLIENT_ID
ARG GOOGLE_CLIENT_SECRET
ARG GOOGLE_STATE

RUN sh config/config.sh ${PORT} ${FE_URL} ${POSTGRES_HOST} ${POSTGRES_PORT} ${POSTGRES_USER} ${POSTGRES_PASSWORD} ${POSTGRES_DB} ${JWT_SECRET} ${GOOGLE_REDIRECT_URL} ${GOOGLE_CLIENT_ID} ${GOOGLE_CLIENT_SECRET} ${GOOGLE_STATE}

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o unclatter cmd/api/main.go

FROM alpine:3.20

WORKDIR /app

# Copy the config file from the build stage
COPY --from=build /app/config/config.yml /app/config/config.yml

# Copy app from build stage
COPY --from=build /app/unclatter /app/unclatter

EXPOSE 80

CMD [ "./unclatter" ]
