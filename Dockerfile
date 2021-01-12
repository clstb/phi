FROM golang@sha256:00967cc62fe598f62aec4b63c9f53f1de00f838384ab39e18ea02e5d2853e5ba as builder
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates make protoc && update-ca-certificates
# Create appuser
ENV USER=appuser
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
# Set working directory
WORKDIR $GOPATH/src/github.com/clstb/phi
# Fetch dependencies.
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify
RUN go get google.golang.org/protobuf/cmd/protoc-gen-go
RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
# Copy source
COPY . .
# Generate proto
RUN make gen-proto
# Build the binary
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/phi
############################
# STEP 2 build a small image
############################
FROM scratch
# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Copy our static executable
COPY --from=builder /go/bin/phi /go/bin/phi
# Use an unprivileged user.
USER appuser:appuser
# Run the hello binary.
ENTRYPOINT ["/go/bin/phi"]
