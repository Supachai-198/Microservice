const express = require('express');
const path = require('path');
const cookieParser = require('cookie-parser');
const logger = require('morgan');
const cors = require('cors');
const dotenv = require("dotenv");
const mongoose = require('mongoose');
dotenv.config();

const indexRouter = require('./routes/index');
const productRouter = require('./routes/product');

const app = express();

app.use(cors());
app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

mongoose.connect(process.env.MONGODB_URI)

const api_version = 'api/v1';
app.use('/', indexRouter);
app.use(`/${api_version}/product`, productRouter);

module.exports = app;
