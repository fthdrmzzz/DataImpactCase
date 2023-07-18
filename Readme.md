# Data Impact CRUD Rest Api 

This app is a REST API written in Go that interacts with a MongoDB database. The application allows you to perform CRUD operations on user data.

In this app, following programming principles are satisfied:
- **Modularity**: The code is organized into separate functions, methods, and packages, allowing for better code organization and reusability.
- **Separation of Concerns**: The code separates different responsibilities into separate functions, methods, and packages, promoting a clear and modular codebase.
- **Encapsulation**: Data and behavior are encapsulated within appropriate structs and functions, ensuring data integrity and preventing direct access to internal components.
- **Error Handling**: Errors are properly handled using Go's built-in error handling mechanisms, such as checking for specific error types and returning appropriate error messages to the client.
- **Documentation**: The code is accompanied by documentation in the form of comments and a separate Markdown file, providing clear explanations of the purpose and usage of each component.
- **Containerization**: Docker is used to containerize the application and MongoDB, promoting isolation, reproducibility, and ease of deployment.

## Prerequisites

Before running MyApp, ensure that you have the latest versions of the following prerequisites installed on your system:
- **Go** 
- **Docker**
- **MongoDB** 

## Running the Application

First clone the repository:
```
git clone https://github.com/fthdrmzzz/DataImpactCase.git
```
Move into the repository directory:
```
cd DataImpactCase
```

### Running via Docker

Inside the repository directory, run the following commands:
```$
make run
```
It will run `docker-compose up` in order to start MongoDB and App containers.

### Running Directly

Inside the repository directory run:
```
go run .\main.go
```
It will build and run the application.

## Using the Application

- `POST /api/users`: Create a new user by providing the necessary details in the request body.

- `POST /api/users/login`: Login to a user's profile by providing the user ID and password.

- `GET /api/users`: Get a list of all users.

- `GET /api/users/{id}`: Get a user by their ID.

- `PUT /user/{id}`: Update a user by their ID.

- `DELETE /api/users/{id}`: Delete a user by their ID.

I have used Postman app for querying and testing the endpoints. You can find the link to api queries [here](https://winter-space-661038.postman.co/workspace/MyWorkspace~d6566199-be2c-4685-92b6-1c67e005e5e7/collection/17115256-29a39a7b-3fe5-40a0-9742-a457ad601847?action=share&creator=17115256).

