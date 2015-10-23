var webpack = require('webpack');
var path = require('path');

module.exports = {
	output: {
		path: path.join(__dirname, 'public'),
		filename: "[name].bundle.js",
		publicPath: "/-/"
	},
	module: {
		loaders: [
			{test: /\.css$/, loader: "style!css"}
		]
	},
	plugins: [
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
