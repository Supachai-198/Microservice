const express = require('express');
const bcrypt = require('bcrypt');
const jwt = require('jsonwebtoken');
const User = require('../models/user');
const router = express.Router();
const passportJWT = require('../middleware/passport-jwt');
const connectRabbitMQ = require('../config/rabbitmq')

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
    const NewUser = await User.create({
      fullname: fullname,
      email: email,
      password: passwordHash
    });

    // connect to rabbitmq
    const channel = await connectRabbitMQ();

    // create exchange
    await channel.assertExchange('ex.supachai.direct', 'direct', { durable: true });

    // create queue
    await channel.assertQueue('q.supachai.direct.product.service', { durable: true });

    // publish message
    await channel.publish(
      'ex.supachai.direct',
      'rk.supachai.product.service',
      Buffer.from(JSON.stringify({
        id: NewUser.id,
        fullname: NewUser.fullname,
        role: NewUser.role
      })),
      {
        contentType: 'application/json',
        contentEncoding: 'utf-8',
        type: 'UserCreated',
        persistent: true
      }
    );
    
    //! fanout
    await channel.assertExchange('ex.supachai.fanout', 'fanout', { durable: true });

    // create queue
    await channel.assertQueue('q.supachai.product.service', { durable: true });

    // publish message
    await channel.publish(
      'ex.supachai.fanout',
      '',
      Buffer.from(JSON.stringify({
        message: 'auth-service fanout message'
      })),
      {
        contentType: 'application/json',
        contentEncoding: 'utf-8',
        type: 'UserFanout',
        persistent: true
      }
    );

    return res.status(200).json({
      message: 'register users'
    });
  }
});

module.exports = router;
