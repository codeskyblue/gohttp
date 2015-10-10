var webpack = require('webpack');
var path = require('path');

module.exports = {
	evtool: 'eval',
	entry: [
		'webpack-dev-server/client?http://localhost:3000',
		'webpack/hot/only-dev-server',
		"./public/entry.js"
	],
	output: {
		path: path.join(__dirname, 'public'),
		filename: "bundle.js",
		publicPath: "/-/"
	},
	module: {
		loaders: [
			{test: /\.css$/, loader: "style!css"},
			{
				test: /\.jsx$/, 
				loaders: ['react-hot', 'babel'],
				include: [path.join(__dirname, 'public')]
			}
		]
	},
	plugins: [
		new webpack.HotModuleReplacementPlugin()
    ],
	resolve: {
		extensions: ['', '.js', 'jsx']
	}
}
