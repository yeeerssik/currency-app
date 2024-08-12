# syntax=docker/dockerfile:1

FROM golang:1.22.1

# Set destination for COPY
WORKDIR /

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /kdf_tech_job

EXPOSE 8000

# Run
CMD ["/kdf_tech_job"]