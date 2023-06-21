const passport = require("passport");
const JwtStrategy = require("passport-jwt").Strategy,
  ExtractJwt = require("passport-jwt").ExtractJwt;

  const opts = {};
  opts.jwtFromRequest = ExtractJwt.fromAuthHeaderAsBearerToken();
  opts.secretOrKey = process.env.JWT_SECRET;
  // opts.issuer = 'accounts.examplesoft.com';
  // opts.audience = 'yoursite.net';

  passport.use(
    new JwtStrategy(opts, async function (jwt_payload, done) {
      try {
        const user = { user_id: jwt_payload.user_id, user_role: jwt_payload.user_role };
        if (user) {
          return done(null, user); // ส่งข้อมูล user ไปกับ req
        }
      } catch (error) {
        done(error);
      }
    })
  );



// check authen
module.exports.isLogin = passport.authenticate("jwt", { session: false });