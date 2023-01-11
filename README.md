# My E-commerce Web Application

An e-commerce web application that allows users to browse and purchase products, and also allows admin to add and remove products and check statistics.

## Features
- User authentication and registration
- Browse and purchase products
- Add products to cart
- Checkout and Payment gateway integration
- Admin dashboard to manage products and view statistics

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- [Go](https://golang.org/doc/install)
- [Node.js](https://nodejs.org/en/download/)
- [MySQL](https://www.mysql.com/downloads/)
- [Git](https://git-scm.com/downloads)

### Installation

# Clone the repository
$ git clone https://github.com/username/myproject.git

# Navigate to the project directory
$ cd myproject

# Install the Go dependencies
$ go get -v -t -d ./...

# Create a .env file in the root of the project and add the following environment variables:
#  DB_USER=<your_db_user>
#  DB_PASSWORD=<your_db_password>
#  DB_NAME=<your_db_name>
#  JWT_SECRET=<your_secret_key>

# Start the backend server
$ go run main.go

# Navigate to the frontend directory
$ cd frontend

# Install the npm dependencies
$ npm install

# Start the frontend server
$ npm start

# Access the application
$ open http://localhost:3000

## Built With

# Backend:
- Go
- Gorm
- MySQL

# Frontend:
- React
- React-Bootstrap
- Stripe

# Author
- Similadayo

# License
- MIT
