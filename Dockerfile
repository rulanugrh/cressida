# Get docker Image
FROM golang:1.22-alpine

# Default Argument
ARG APP_PORT

# Workdir Default Config
WORKDIR /usr/src/app

# COPY all source code to local docker
COPY . .

# Running Build
RUN go mod tidy
RUN go build -o main

# EXPOSE port
EXPOSE ${APP_PORT}

# Running Main File
CMD [ "./main" ]