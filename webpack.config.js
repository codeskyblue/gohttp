var webpack = require('webpack');
var path = require('path');
var baseConfig = require('./webpack.config.base');

var config = Object.create(baseConfig);

config.evtool = 'eval';
config.entry = [
	'webpack-dev-server/client?http://localhost:3000',
	'webpack/hot/only-dev-server',
	"./public/entry.js"
];
config.plugins.push(
	new webpack.HotModuleReplacementPlugin(),
	new webpack.NoErrorsPlugin())

module.exports = config;
