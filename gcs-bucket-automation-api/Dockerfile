FROM pulumi/pulumi-go

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /bucket-factory  cmd/main.go

# $PROJECT_ID$PROJECT_ID$PROJECT_ID
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run
CMD ["/bucket-factory"]
