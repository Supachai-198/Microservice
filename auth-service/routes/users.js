const express = require('express');
const bcrypt = require('bcrypt');
const jwt = require('jsonwebtoken');
const User = require('../models/user');
const router = express.Router();
const passportJWT = require('../middleware/passport-jwt');

/* GET users listing. */
router.get('/', function(req, res, next) {
  res.send('respond with a resource');
});

router.post('/login', async function(req, res, next) {
  const { email, password } = req.body;
  const User_data = await User.findOne({ where: { email: email } });

  // check email
  if( await User_data === null )
  {
    return res.status(404).json({
      message: 'user not found'
    });
  }
  else
  {
    // check password
    const isValid = await bcrypt.compare(password, User_data.password);
    if( !isValid )
    {
      return res.status(401).json({
        message: 'invalid password'
      });
    }

    const access_token = jwt.sign(
      {
        user_id: User_data.id,
        user_role: User_data.role
      },
      process.env.JWT_SECRET,
      {
        expiresIn: '7d',
      }
    );

    return res.status(200).json({
      access_token: access_token,
      message: 'login users'
    });
  }
});

router.get('/profile', passportJWT.isLogin, async function(req, res, next) {
  const user = await User.findByPk(req.user.user_id);

  return res.status(200).json({
    id: user.id,
    fullname: user.fullname,
    email: user.email,
    role: user.role,
    message: 'user profile'
  });
});

router.post('/register', async function(req, res, next) {
  const { fullname, email, password } = req.body;

  // check if email already exists
  if( await User.findOne({ where: { email: email } }) !== null )
  {
    return res.status(409).json({
      message: 'email already exists'
    });
  }

  else
  {
    // hash password
    const salt = await bcrypt.genSalt();
    const passwordHash = await bcrypt.hash(password, salt);

    // create user
    await User.create({
      fullname: fullname,
      email: email,
      password: passwordHash
    });
    
    return res.status(200).json({
      message: 'register users'
    });
  }


  return res.status(200).json({
    message: 'register users'
  });
});

module.exports = router;
