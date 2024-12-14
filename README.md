# SafePass Server

SafePass is a secure password manager application. This repository contains the server-side code for SafePass, which provides APIs for user authentication, password management, and vault operations.

## Features

- **User Authentication**: Users can register and log in using their email and master password.
- **Password Vault**: Each user has a vault where they can store and manage their passwords.
- **Password Management**: Users can create, update, retrieve, and delete passwords in their vault.
- **JWT Authentication**: Secure authentication using JSON Web Tokens (JWT).
- **Supabase Integration**: Uses Supabase as the backend database for storing user and vault information.
- **Logging**: Logs requests and responses for debugging and monitoring purposes.

## Project Structure

```
.
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── testing/
│       └── testing.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── auth_handlers.go
│   │   │   └── vault_handlers.go
│   │   ├── middlewares/
│   │   │   ├── auth_middleware.go
│   │   │   └── log_middleware.go
│   │   └── routes/
│   │       └── routes.go
│   ├── config/
│   │   └── config.go
│   ├── consts/
│   │   └── consts.go
│   ├── database/
│   │   └── database.go
│   ├── logging/
│   │   └── logging.go
│   ├── repositories/
│   │   ├── password_repository.go
│   │   ├── user_repository.go
│   │   └── vault_repository.go
│   └── services/
│       ├── auth_services.go
│       ├── user_services.go
│       └── vault_services.go
├── pkg/
│   ├── crypto/
│   │   └── crypto.go
│   ├── dotenv/
│   │   └── dotenv.go
│   ├── dtos/
│   │   ├── password/
│   │   │   ├── create_password.go
│   │   │   └── create_password_request.go
│   │   ├── user/
│   │   │   ├── create_user.go
│   │   │   ├── create_user_request.go
│   │   │   ├── login_request.go
│   │   │   └── update_user_request.go
│   │   └── vault/
│   │       └── create_vault.go
│   └── models/
│       ├── error.go
│       ├── identity.go
│       ├── password.go
│       ├── response.go
│       ├── token.go
│       └── user.go
├── .env
├── .gitignore
├── config.yaml
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/safepass-server.git
    cd safepass-server
    ```

2. Install dependencies:
    ```sh
    go mod download
    ```

3. Set up environment variables:
    Create a .env file in the root directory and add the following variables:

    ```env
    SUPABASE_REST_URL=your_supabase_rest_url
    SUPABASE_API_KEY=your_supabase_api_key
    JWT_SECRET_KEY=your_jwt_secret_key
    ```

4. Configure the application:
    Modify the config.yaml file to set your server, JWT, and logging configurations.

## Running the Server

To start the server, run:
```sh
go run cmd/server/main.go
```

The server will start on the host and port specified in the config.yaml file.

## API Endpoints

### Authentication

- **POST /api/v1/auth/login**: Log in a user.
- **POST /api/v1/auth/register**: Register a new user.

### Vault

- **GET /api/v1/vault/@me**: Get the user's vault.
- **GET /api/v1/vault/passwords**: Get all passwords in the user's vault.
- **GET /api/v1/vault/password**: Get a specific password by ID.
- **POST /api/v1/vault/password/create**: Create a new password in the vault.
- **POST /api/v1/vault/password/update/{id}**: Update an existing password.
- **POST /api/v1/vault/password/delete/{id}**: Delete a password.

## Logging

Logs are written to log.txt by default.

## License

This project is licensed under the GNU Affero General Public License v3.0. See the LICENSE file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.