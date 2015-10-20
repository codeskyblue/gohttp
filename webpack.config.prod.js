var baseConfig = require('./webpack.config.base');
var config = Object.create(baseConfig);

config.entry = [
	"./public/entry.js"
];

config.module.loaders.push(
	{
		test: /\.jsx$/, 
		loaders: ['babel'],
		include: [path.join(__dirname, 'public')]
	})

module.exports = config;
