APP_NAME="kloroform"

# Windows build
GOOS=windows GOARCH=amd64 go build -o ${APP_NAME}_windows_amd64.exe

# Linux build
GOOS=linux GOARCH=amd64 go build -o ${APP_NAME}_linux_amd64
