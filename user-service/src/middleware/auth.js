const jwt = require('jsonwebtoken');
const User = require('../models/User');

const generateToken = (userId) => {
  return jwt.sign({ userId }, process.env.JWT_SECRET, {
    expiresIn: process.env.JWT_EXPIRES_IN || '7d'
  });
};

const verifyToken = (token) => {
  return jwt.verify(token, process.env.JWT_SECRET);
};

const createContext = async ({ req }) => {
  let user = null;

  const token = req.headers.authorization?.replace('Bearer ', '');

  if (token) {
    try {
      const decoded = verifyToken(token);
      user = await User.findByPk(decoded.userId, {
        attributes: { exclude: ['passwordHash'] }
      });
    } catch (error) {
      console.error('Invalid token:', error.message);
    }
  }

  return { user, token };
};

const requireAuth = (user) => {
  if (!user) {
    throw new Error('Authentication required');
  }
  return user;
};

const requireRole = (user, allowedRoles) => {
  requireAuth(user);

  if (!allowedRoles.includes(user.role)) {
    throw new Error('Insufficient permissions');
  }

  return user;
};

module.exports = {
  generateToken,
  verifyToken,
  createContext,
  requireAuth,
  requireRole
};
