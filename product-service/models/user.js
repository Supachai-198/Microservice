const mongoose = require('mongoose');

const modelSchema = new mongoose.Schema(
    {
        user_id: {
            type: Number,
            require: true,
            unique: true
        },
        fullname: {
            type: String,
        },
        role: {
            type: String,
        }
    },
    {
        collation: 'users',
        timestamps: false,
        versionKey: false
    }
);

const User = mongoose.model('User', modelSchema);

module.exports = User;