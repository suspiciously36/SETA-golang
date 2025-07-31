const { ApolloServer } = require('apollo-server-express');
const express = require('express');
const cors = require('cors');
const { sequelize, testConnection } = require('./config/database');
require('dotenv').config();

const typeDefs = require('./schema/typeDefs');
const resolvers = require('./resolvers');
const { createContext } = require('./middleware/auth');

async function startServer() {
  const app = express();

  // Middleware
  app.use(
    cors({
      origin: process.env.CORS_ORIGIN || 'http://localhost:3000',
      credentials: true
    })
  );

  // Connect to PostgreSQL and sync models
  try {
    await testConnection();

    // Sync database (create tables if they don't exist)
    await sequelize.sync({
      force:
        process.env.NODE_ENV === 'development' &&
        process.env.FORCE_SYNC === 'true'
    });
    console.log('📦 Database synchronized');
  } catch (error) {
    console.error('❌ Database connection error:', error);
    process.exit(1);
  }

  // Create Apollo Server
  const server = new ApolloServer({
    typeDefs,
    resolvers,
    context: createContext,
    introspection: process.env.NODE_ENV === 'development',
    plugins: [
      // Enable GraphiQL interface
      process.env.NODE_ENV === 'development'
        ? require('apollo-server-core').ApolloServerPluginLandingPageLocalDefault(
            {
              embed: true,
              includeCookies: true
            }
          )
        : require('apollo-server-core').ApolloServerPluginLandingPageProductionDefault(
            {
              footer: false
            }
          )
    ]
  });

  await server.start();
  server.applyMiddleware({
    app,
    path: '/graphql',
    cors: {
      origin: process.env.CORS_ORIGIN || 'http://localhost:3000',
      credentials: true
    }
  });

  // Health check endpoint
  app.get('/health', (req, res) => {
    res.json({
      status: 'OK',
      service: 'user-service',
      timestamp: new Date().toISOString()
    });
  });

  const PORT = process.env.PORT || 4000;

  app.listen(PORT, () => {
    console.log(`🚀 User Service running at http://localhost:${PORT}`);
    console.log(
      `🎯 GraphQL endpoint: http://localhost:${PORT}${server.graphqlPath}`
    );
    if (process.env.NODE_ENV === 'development') {
      console.log(
        `🎨 GraphiQL UI available at: http://localhost:${PORT}${server.graphqlPath}`
      );
    }
  });
}

// Handle unhandled promise rejections
process.on('unhandledRejection', (err) => {
  console.error('❌ Unhandled Promise Rejection:', err);
  process.exit(1);
});

// Handle uncaught exceptions
process.on('uncaughtException', (err) => {
  console.error('❌ Uncaught Exception:', err);
  process.exit(1);
});

startServer().catch((error) => {
  console.error('❌ Failed to start server:', error);
  process.exit(1);
});
