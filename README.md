# WhatsApp Messenger Backend

This project is a backend implementation of a WhatsApp-style messenger built with Golang, using the Echo framework. It leverages Gorm for database interactions, Go-Playground Validator for input validation, and Swagger for API documentation.

## Features

- **Echo**: High-performance, extensible, minimalistic web framework for Go.
- **Gorm**: ORM library for Go providing a powerful database abstraction.
- **Go-Playground Validator**: Struct-level and field-level validations.
- **Swagger**: Generates interactive API documentation.

## Requirements

Before you begin, ensure you have met the following requirements:

- Golang version >= 1.18
- A running instance of a database compatible with Gorm (e.g., PostgreSQL, MySQL, SQLite)
- [Swaggo](https://github.com/swaggo/swag) CLI for generating Swagger documentation

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/guvanch999/whatsapp-messenger-backend.git
   cd whatsapp-messenger-backend
   ```
   
2. **Install dependencies**:
    
   ```bash
   go mod tidy
   ```
3. **Configure environment variables:**

    Create a .env file in the project root and add the following:

    ```dotenv

    GO_ENV=development
    PORT=3030
    
    DISABLE_AUTO_MIGRATION="true"
    APP_URL=
    POSTGRES_URI=
    SUPABASE_URL=
    SUPABASE_KEY=
    BRANCH_NAME=main
    GOOGLE_CREDENTIALS=
    SECRET_KEY_FOR_HASH=

    ```
   
4. **Generate Swagger documentation:**

    ```bash
    swag init
    ```
   
5. **Run the application:**
   Start the server with the following command:
    ```bash
   go run .
    ```