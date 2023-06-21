const { DataTypes } = require('sequelize');
const sequelize = require('../db/mysql');

const User = sequelize.define(
    'User', {
        fullname: { type: DataTypes.STRING(255) },
        email: {
            type: DataTypes.STRING(255),
            unique: true,
        },
        password: { type: DataTypes.STRING(255) },
        role: {
            type: DataTypes.STRING(100),
            defaultValue: 'member',
        },
    },
    {
        tableName: 'users',
        timestamps: false,
    }
);

module.exports = User;