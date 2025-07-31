# User Service

A GraphQL microservice for user management built with Node.js, Apollo Server, and PostgreSQL.

## Features

- GraphQL API for user management
- JWT-based authentication
- Role-based access control (ADMIN, USER, MANAGER)
- Password hashing with bcrypt
- Input validation and error handling
- PostgreSQL integration with Sequelize ORM

## API Operations

### Queries

- `fetchUsers`: List all users (Admin/Manager only)
- `me`: Get current user profile

### Mutations

- `createUser(input: CreateUserInput!)`: Create a new user
- `login(input: LoginInput!)`: Login and receive JWT token
- `logout()`: Logout current user
- `updateUser(id: ID!, input: UpdateUserInput!)`: Update user profile
- `deleteUser(id: ID!)`: Delete user (Admin only)

## Setup

1. Install dependencies:

```bash
npm install
```

2. Copy environment variables:

```bash
cp .env.example .env
```

3. Update the `.env` file with your configuration:

```
NODE_ENV=development
PORT=4000

# PostgreSQL Database Configuration
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=goapi_db
DATABASE_USER=goapi_user
DATABASE_PASSWORD=goapi_password

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRES_IN=7d

# CORS Configuration
CORS_ORIGIN=http://localhost:3000

# Development options
FORCE_SYNC=false
```

4. Start PostgreSQL (make sure PostgreSQL is running with the same database as go-apis-service)

5. Start the development server:

```bash
npm run dev
```

## Usage

The GraphQL playground will be available at `http://localhost:4000/graphql` in development mode.

### Example Queries

**Create User:**

```graphql
mutation {
  createUser(
    input: {
      username: "johndoe"
      email: "john@example.com"
      password: "password123"
      role: USER
    }
  ) {
    userId
    username
    email
    role
  }
}
```

**Login:**

```graphql
mutation {
  login(input: { email: "john@example.com", password: "password123" }) {
    token
    user {
      userId
      username
      email
      role
    }
  }
}
```

**Fetch Users (with Authorization header):**

```graphql
query {
  fetchUsers {
    userId
    username
    email
    role
    createdAt
  }
}
```

## Authentication

Include the JWT token in the Authorization header:

```
Authorization: Bearer your-jwt-token-here
```

## Project Structure

```
src/
├── index.js              # Main server file
├── config/
│   └── database.js       # PostgreSQL connection
├── models/
│   └── User.js           # User model (Sequelize)
├── schema/
│   └── typeDefs.js       # GraphQL schema
├── resolvers/
│   └── index.js          # GraphQL resolvers
└── middleware/
    └── auth.js           # Authentication middleware
```
