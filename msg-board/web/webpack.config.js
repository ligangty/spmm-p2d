const path = require('path');
const webpack = require('webpack');
const CleanWebpackPlugin = require('clean-webpack-plugin');

const outputDirectory = 'build';

module.exports = {
  entry: './src/client/index.js',
  output: {
    path: path.resolve(__dirname, outputDirectory),
    filename: 'index_bundle.js'
  },
  mode: 'development',
  devtool: 'inline-source-map',
  devServer: {
    port: 3000,
    contentBase: './build',
    historyApiFallback: {
      rewrites: [
        { from: /\/browse\/.*/, to: '/' },
      ]
    },
    proxy: {
      '/api': 'http://localhost:8080'
    }
  },
  // plugins: [
  //   new CleanWebpackPlugin([outputDirectory])
  // ],
  module: {
    rules: [
      { test: /\.js$/, use: 'babel-loader', exclude: /node_modules/ },
      { test: /\.jsx?$/, use: 'babel-loader', exclude: /node_modules/ },
      { test: /\.css$/, use: ['style-loader', 'css-loader'] }
    ]
  }

}
