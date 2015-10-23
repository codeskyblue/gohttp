var webpack = require('webpack');
var path = require('path');
var baseConfig = require('./webpack.config.base');

var config = Object.create(baseConfig);

config.evtool = 'eval';
config.entry = {
	explorer: [
		'webpack-dev-server/client?http://localhost:3000',
		'webpack/hot/only-dev-server',
		"./public/explorer.entry.js"
	],
	preview: [
		'webpack-dev-server/client?http://localhost:3000',
		'webpack/hot/only-dev-server',
		"./public/preview.entry.js"
	]
}

config.plugins.push(
	new webpack.HotModuleReplacementPlugin(),
	new webpack.NoErrorsPlugin())

config.module.loaders.push(
	{
		test: /\.jsx$/, 
		loaders: ['react-hot', 'babel'],
		include: [path.join(__dirname, 'public')]
	})
module.exports = config;
