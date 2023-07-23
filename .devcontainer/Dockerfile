FROM mcr.microsoft.com/vscode/devcontainers/go:1.20

# allow writes for all to solve module write permission issues
RUN chmod -R a+w /go

# compatible with latest golang
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
