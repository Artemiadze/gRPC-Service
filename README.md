# gRPC Service on Golang #
---

## The general scheme of the project — how authorization will work

### The characters:

**User** — a person who is forced to log in because he wants to use our URL Shortener.
**URL Shortener** — the service that will be the SSO client
**An authorization server (SSO)** is a service that can authorize, provide information about user rights, etc.

### How it will work:

The user (or the application used by him) sends a request to the SSO to receive a JWT authorization token
With this token, it goes to the URL Shortener to perform various useful queries — create short links, delete them, etc.
The URL Shortener receives a request from the client, extracts a token from it, by which it understands who came and what it is allowed to do.

## 🛠️ Installation

### Prerequisites
- Docker 🐳

### Steps to Install
1. Clone the repository:
   ```bash
   git clone 
   ```
2. Navigate to the project directory:
   ```
   cd gRPC-Service
   ```
3. Start the application using Docker:
   ```
   docker compose up --build
   ```
