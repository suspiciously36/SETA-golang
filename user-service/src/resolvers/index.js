const User = require('../models/User');
const {
  generateToken,
  requireAuth,
  requireRole
} = require('../middleware/auth');
const {
  UserInputError,
  AuthenticationError
} = require('apollo-server-express');

const resolvers = {
  Query: {
    fetchUsers: async (_, __, { user }) => {
      // Only admins and managers can fetch all users
      requireRole(user, ['ADMIN', 'MANAGER']);

      try {
        return await User.findAll({
          order: [['createdAt', 'DESC']],
          attributes: { exclude: ['passwordHash'] }
        });
      } catch (error) {
        throw new Error('Failed to fetch users');
      }
    },

    me: (_, __, { user }) => {
      requireAuth(user);
      return user;
    }
  },

  Mutation: {
    createUser: async (_, { input }) => {
      const { username, email, password, role } = input;

      try {
        // Check if user already exists
        const existingUser = await User.findOne({
          where: {
            [require('sequelize').Op.or]: [{ email }, { username }]
          }
        });

        if (existingUser) {
          throw new UserInputError(
            'User with this email or username already exists'
          );
        }

        // Validate password strength
        if (password.length < 6) {
          throw new UserInputError(
            'Password must be at least 6 characters long'
          );
        }

        const user = await User.create({
          username,
          email,
          passwordHash: password, // Will be hashed by the model hook
          role: role || 'USER'
        });

        // Return user without password hash
        const { passwordHash, ...userWithoutPassword } = user.toJSON();
        return userWithoutPassword;
      } catch (error) {
        if (error.name === 'SequelizeUniqueConstraintError') {
          throw new UserInputError(
            'User with this email or username already exists'
          );
        }
        throw new Error(error.message || 'Failed to create user');
      }
    },

    login: async (_, { input }) => {
      const { email, password } = input;

      try {
        // Find user by email
        const user = await User.findOne({ where: { email } });
        if (!user) {
          throw new AuthenticationError('Invalid email or password');
        }

        // Check password
        const isValidPassword = await user.comparePassword(password);
        if (!isValidPassword) {
          throw new AuthenticationError('Invalid email or password');
        }

        // Generate token
        const token = generateToken(user.userId);

        // Return user without password hash
        const { passwordHash, ...userWithoutPassword } = user.toJSON();

        return {
          token,
          user: userWithoutPassword
        };
      } catch (error) {
        if (error instanceof AuthenticationError) {
          throw error;
        }
        throw new Error('Login failed');
      }
    },

    logout: (_, __, { user }) => {
      requireAuth(user);
      // In a real-world scenario, you might want to blacklist the token
      return 'Successfully logged out';
    },

    updateUser: async (_, { id, input }, { user }) => {
      requireAuth(user);

      // Users can only update their own profile unless they're admin
      if (user.userId !== parseInt(id) && user.role !== 'ADMIN') {
        throw new Error('You can only update your own profile');
      }

      try {
        const [updatedRowsCount] = await User.update(input, {
          where: { userId: id },
          returning: true
        });

        if (updatedRowsCount === 0) {
          throw new Error('User not found');
        }

        const updatedUser = await User.findByPk(id, {
          attributes: { exclude: ['passwordHash'] }
        });

        return updatedUser;
      } catch (error) {
        throw new Error(error.message || 'Failed to update user');
      }
    },

    deleteUser: async (_, { id }, { user }) => {
      // Only admins can delete users
      requireRole(user, ['ADMIN']);

      try {
        const userToDelete = await User.findByPk(id);

        if (!userToDelete) {
          throw new Error('User not found');
        }

        await userToDelete.destroy();

        return `User ${userToDelete.username} deleted successfully`;
      } catch (error) {
        throw new Error(error.message || 'Failed to delete user');
      }
    }
  }
};

module.exports = resolvers;
