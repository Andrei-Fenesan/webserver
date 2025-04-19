# Web Server in Go

This project contains a simple web server written in Go. It allows you to return statics HTML files.
The server listens on port 8080 (by default) and serves files from the `resources` directory. The CGI scripts should be in the `cgi-bin` directory inside the directory root directory
The port and the directory can be configured using command-line flags:
- `-port` specifies the port number
- `-src` specifies the directory to serve files from

## Prerequisites

- Go 1.16 or later

## Features
- Serve static HTML files
- Serve dynamic content using CGI
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
## Compiling
1. Compile the program
   ```sh
   go build
   ```

## Run the webserver as a container
1. Build the image
```sh
docker build --tag {insert_desired_image_neme}
```
2. Create the container
```sh
docker create --name {insert_desired_container_neme} -p {desired_port}:8080 {insert_above_image_name}
```
3. Copy the resources in the container from the host filesytem
```sh
docker cp ./{insert_resource_directory} {insert_above_container_name}:/app
```
4. Start the container
```sh
docker start {insert_above_container_name}
```
Example of a single line command that do the same as the above:
```sh
docker build --tag websvimg . && docker create --name websv -p 9991:8080 websvimg && docker cp ./resources websv:/app && docker start websv
```

## Future Improvements
 - Add support for serving files over HTTPS.
 - Add support for basic authentication.
 - Add support for other HTTP methods (e.g., POST, PUT, DELETE).
 - Add support for all platforms (Windows, Linux, MacOS). Currently the server works on Linux and MacOS