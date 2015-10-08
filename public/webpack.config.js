var webpack = require('webpack');
var path = require('path');

module.exports = {
	entry: [
		'webpack-dev-server/client?http://localhost:3000',
		'webpack/hot/only-dev-server',
		"./entry.js"
	],
	output: {
		path: __dirname,
		filename: "bundle.js"
	},
	module: {
		loaders: [
			{test: /\.css$/, loader: "style!css"},
			{
				test: /\.jsx$/, 
				loaders: ['react-hot', 'jsx?harmony'],
				include: [path.join(__dirname, '.')]
			}
		]
	},
	plugins: [
		new webpack.HotModuleReplacementPlugin(),
		new webpack.NoErrorsPlugin()
    ],
	resolve: {
		extensions: ['', '.js', 'jsx']
	}
}
