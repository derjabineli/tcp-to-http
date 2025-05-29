# TCP-to-HTTP

## Description
TCP-to-HTTP is a minimalist HTTP/1.1 server implementation built from scratch in Go. It parses incoming TCP connections, processes valid HTTP requests as per RFC 9110, 9112, and generates corresponding HTTP responses.

This project was designed for learning and exploring the world of HTTP and HTTP servers and, at this time, is not intended for production use. 

## Features
- HTTP/1.1 Compliance – Supports HTTP/1.1 request/response parsing and response
- TCP Socket Handling – Manages TCP connections directly.
- References RFC 9110 & 9112

## Installation
Clone this repository 
***git clone https://github.com/derjabineli/tcp-to-http.git***

## Usage
Run `go run ./cmd/httpserver` to start up a HTTP server on port 42010


