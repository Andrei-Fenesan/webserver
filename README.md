# Web Server in Go

This project contains a simple web server written in Go. It allows you to return statics HTML files.
The server listens on port 80 (by default for HTTP) and 443 (by default for HTTPS) and serves files from the `resources`
directory. The CGI scripts should be
in the `cgi-bin` directory inside the directory root directory

## Prerequisites

- Go 1.16 or later

## Features

- Serve static HTML files
- Serve dynamic content using CGI
- Allow HTTPS traffic
- Protect against directory traversal attacks
- Support for command-line flags to configure the server
- Redirects HTTP traffic to HTTPS

## ⚙️ Command-line Flags

| Flag              | Type   | Default       | Description                                        |
|-------------------|--------|---------------|----------------------------------------------------|
| `-src`            | string | `./resources` | Root directory to serve files from                 |
| `-httpPort`       | int    | `80`          | Port to run the HTTP server on                     |
| `-httpsPort`      | int    | `443`         | Port to run the HTTPS server on                    |
| `-certPath`       | string | `""`          | Path to the SSL certificate (required for HTTPS)   |
| `-privateKeyPath` | string | `""`          | Path to the SSL private key (required for HTTPS)   |
| `-httpRedirectTo` | string | `"localhost"` | Target host to redirect HTTP traffic to over HTTPS |

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

## Compiling

1. Compile the program
   ```sh
   go build
   ```

## 🛠 Usage

### 1. Basic HTTP Server (no TLS)

Runs a simple HTTP server that serves static files from the `./static` directory on port `8080`.

```bash
go run main.go -src ./static -httpPort 8080
```

### 2. HTTPS Usage with HTTP Redirection

To run the server with **HTTPS enabled** and automatic **HTTP to HTTPS redirection**, use the following command:

```bash
go run main.go \
  -src ./static \
  -httpPort 80 \
  -httpsPort 443 \
  -certPath ./certs/cert.pem \
  -privateKeyPath ./certs/key.pem \
  -httpRedirectTo example.com
  ```

## Run the webserver as a container

All the command are relative to the webserver directory

1. Build the image

```sh
docker build --tag {insert_desired_image_neme}
```

2. Create the container

```sh
docker create --name {insert_desired_container_neme} -p {desired_port}:8080 {insert_above_image_name}
```

3. Copy the resources in the container from the host filesystem

```sh
docker cp ./{insert_resource_directory} {insert_above_container_name}:/app
```

4. Start the container

```sh
docker start {insert_above_container_name}
```

Example of a single line command that does the same as the above:

```sh
docker build --tag websvimg . && docker create --name websv -p 9991:8080 websvimg && docker cp ./resources websv:/app && docker start websv
```

## Future Improvements

- Add support for basic authentication.
- Add support for other HTTP methods (e.g., POST, PUT, DELETE).