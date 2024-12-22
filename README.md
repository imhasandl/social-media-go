[![CI](https://github.com/imhasandl/go-restapi/actions/workflows/ci.yml/badge.svg)](https://github.com/imhasandl/go-restapi/actions/workflows/ci.yml)

# Go RESTful API

This is a simple RESTful API written in Golang.

## Documentation

There is 18 types of requests, but here in documentation are only essential ones and you can check the rest of them in the main file.

WARNING: The main function contains request that can delete all user from database and posts, so be careful. 

## Functionality

- User Management:
    - Login and register users using JWT tokens.
    - Change user email and password (authenticated access required).
    - List all users (optional admin privileges).
    - Get user information by email or ID (authenticated access or admin privileges).
- Post Management:
    - Create, list, retrieve, update, and delete posts (authenticated access required for most operations).
- Status Check:
    - Verifies if the API server is running.

### API Configuration

The API configuration is defined by the `apiConfig` struct. It holds the following fields:

* `db`: Database connection pool
* `status`: Function to check the status of the server
* `jwtSecret`: Secret key used for JWT authentication
* `webhookKey`: Secret key used for validating webhooks

### Endpoints

#### Status Check

* **Method:** GET
* **URL:** `/status`
* **Description:** Checks if the server is running.
* **Response:** JSON object with a `message` field indicating the server status (e.g., "OK").

#### User Management

* **Login**
    * **Method:** POST
    * **URL:** `/api/users/login`
    * **Description:** Login a user using a JWT token in the header.
    * **Request Body:** JSON object with the following fields:
        * `email`: User's email address
        * `password`: User's password
    * **Response:** JSON object with a `token` field containing the JWT token on success, or an error message on failure.
* **Register**
    * **Method:** POST
    * **URL:** `/api/users/register`
    * **Description:** Register a new user.
    * **Request Body:** JSON object with the following fields:
        * `email`: User's email address
        * `password`: User's password
    * **Response:** JSON object with a message indicating success or failure.
* **Change User Information**
    * **Method:** PUT
    * **URL:** `/api/users/change`
    * **Description:** Change a user's email and password. Requires a JWT token in the header.
    * **Request Body:** JSON object with the following fields (optional):
        * `email`: New email address
        * `password`: New password
    * **Response:** JSON object with a message indicating success or failure.
* **Get All Users**
    * **Method:** GET
    * **URL:** `/api/users`
    * **Description:** Retrieves a list of all users. Requires admin privileges.
    * **Response:** JSON array containing user objects. Each user object has the following fields:
        * `id`: User's unique identifier
        * `email`: User's email address
* **Get User by Email**
    * **Method:** GET
    * **URL:** `/api/users/{email}`
    * **Description:** Retrieves a user by their email address.
    * **Response:** JSON object containing the user information or an error message if the user is not found.
* **Get User by ID**
    * **Method:** GET
    * **URL:** `/api/users/{user_id}`
    * **Description:** Retrieves a user by their ID.
    * **Response:** JSON object containing the user information or an error message if the user is not found.

#### Post Management

* **List All Posts**
    * **Method:** GET
    * **URL:** `/api/posts`
    * **Description:** Retrieves a list of all posts.
    * **Response:** JSON array containing post objects. Each post object has the following fields:
        * `id`: Post's unique identifier
        * `title`: Post title
        * `content`: Post content (optional)
        * `author_id`: ID of the user who created the post (optional)
* **Create Post**
    * **Method:** POST
    * **URL:** `/api/posts`
    * **Description:** Creates a new post. Requires a JWT token in the header.
    * **Request Body:** JSON object with the following fields:
        * `title`: Post title
        * `content`: Post content (optional)
    * **Response:** JSON object containing the newly created post information or an error message on failure.
* **Get Post by ID**
    * **Method:** GET
    * **URL:** `/api/posts/{post_id}`
    * **Description:** Retrieves a post by its ID.
    * **Response:** JSON object containing the post information or an error message if the post is not found.
* **Change Post**
    * **Method:** PUT
    * **URL:** `/api/posts/{post_id}`
    * **Description:** Changes a post. Requires a JWT token in the header and ownership of the post.
    * **Request Body:**
