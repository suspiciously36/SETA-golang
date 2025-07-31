const { gql } = require('apollo-server-express');

const typeDefs = gql`
  type User {
    userId: ID!
    username: String!
    email: String!
    role: UserRole!
    createdAt: String!
    updatedAt: String!
  }

  enum UserRole {
    ADMIN
    USER
    MANAGER
  }

  type AuthPayload {
    token: String!
    user: User!
  }

  type Query {
    fetchUsers: [User!]!
    me: User
  }

  type Mutation {
    createUser(input: CreateUserInput!): User!
    login(input: LoginInput!): AuthPayload!
    logout: String!
    updateUser(id: ID!, input: UpdateUserInput!): User!
    deleteUser(id: ID!): String!
  }

  input CreateUserInput {
    username: String!
    email: String!
    password: String!
    role: UserRole = USER
  }

  input LoginInput {
    email: String!
    password: String!
  }

  input UpdateUserInput {
    username: String
    email: String
    role: UserRole
  }
`;

module.exports = typeDefs;
