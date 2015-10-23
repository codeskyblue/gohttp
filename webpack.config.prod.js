var path = require('path');
var baseConfig = require('./webpack.config.base');
var config = Object.create(baseConfig);

config.entry = {
	explorer: "./public/explorer.entry.js",
	preview: "./public/preview.entry.js"
}

config.module.loaders.push(
	{
		test: /\.jsx$/, 
		loaders: ['babel'],
		include: [path.join(__dirname, 'public')]
	})

module.exports = config;
