const mongoose = require('mongoose');

const modelSchema = new mongoose.Schema(
    {
        product_id: {
            type: Number,
            require: true,
            unique: true
        },
        product_name: {
            type: String,
        },
        product_price: {
            type: Number,
        },
        user_id: {
            type: Number,
        }
    },
    {
        collation: 'products',
        timestamps: false,
        versionKey: false
    }
);

const Product = mongoose.model('Product', modelSchema);

module.exports = Product;