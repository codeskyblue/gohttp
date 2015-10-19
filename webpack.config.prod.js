var baseConfig = require('./webpack.config.base');
var config = Object.create(baseConfig);

config.entry = [
	"./public/entry.js"
];

module.exports = config;
