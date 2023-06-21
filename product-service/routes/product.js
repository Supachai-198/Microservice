const express = require('express');
const router = express.Router();
const passportJWT = require('../middleware/passport-jwt');
const Product = require('../models/product');

/* GET users listing. */
router.post('/insert', passportJWT.isLogin, async function(req, res, next) {
  const { product_id, product_name, product_price } = req.body;
  const newProduct = await Product.create({
    product_id: product_id,
    product_name: product_name,
    product_price: product_price,
    user_id: req.user.user_id
  });

  return res.status(201).json({
    message: 'product added',
    data: newProduct
  });
});

router.get('/get', async function(req, res, next) {
  const products = await Product.find({})
  .collation({ locale: 'en', strength: 1 });

  return res.status(200).json({
    message: 'product',
    data: products
  });
});

module.exports = router;
