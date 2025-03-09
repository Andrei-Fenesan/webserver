# Web Server in Go

This project contains a simple web server written in Go. It allows you to return statics HTML files.
The server listens on port 8080 (by default) and serves files from the `resources` directory.
The port and the directory can be configured using command-line flags:
- `-port` specifies the port number
- `-src` specifies the directory to serve files from

## Prerequisites

- Go 1.16 or later

## Features
- Serve static HTML files
- Protect against directory traversal attacks
- Support for command-line flags to configure the server

## Installation

1. Clone the repository:
   ```sh
   git clone <repository-url>
   cd <repository-directory>
   ```
2. Start the webserver:
   ```sh
   go run .
   ```
   Alternatively, you can start the server on the desired port and serve files from a specific directory:
   ```sh
      go run . -port 8081 -src /path/to/directory
   ```
## Future Improvements
 - Add support for serving dynamic content using CGI.