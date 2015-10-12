var webpack = require('webpack');
var path = require('path');

module.exports = {
	evtool: 'eval',
	entry: [
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
				loaders: ['babel'],
				include: [path.join(__dirname, 'public')]
			}
		]
	},
	plugins: [
		new webpack.HotModuleReplacementPlugin(),
		new webpack.ProvidePlugin({
			$: "jquery",
			jQuery: "jquery",
			"window.jQuery": "jquery"
		})
    ],
	resolve: {
		extensions: ['', '.js', 'jsx']
	}
}
