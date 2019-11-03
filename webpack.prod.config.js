require('babel-polyfill')
const path = require('path') // path 模块提供了一些工具函数，用于处理文件与目录的路径。
const merge = require('webpack-merge')
const webpackBaseConfig = require('./webpack.base.config')
const CompressionPlugin = require('compression-webpack-plugin')
const TerserPlugin = require('terser-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const OptimizeCssAssetsPlugin = require('optimize-css-assets-webpack-plugin') // 压缩css插件

const webpackProdConfig = {
  output: {
    // 出口
    filename: 'js/[name].bundle.[chunkhash].js', // 文件名
    chunkFilename: 'js/[name].chunk.[chunkhash].js',
    path: path.resolve(__dirname, 'dist'), // 路径
    publicPath: '/' // srcript 引入路径
  },
  optimization: {
    splitChunks: {
      cacheGroups: {
        vendor: {
          test: /node_modules/,
          chunks: 'initial',
          priority: 10,
          enforce: true
        },
        styles: {
          test: /\.s?css$/,
          chunks: 'all',
          enforce: true,
          priority: 20
        }
      }
    },
    minimizer: [
      new TerserPlugin({
        extractComments: true,
        cache: true,
        parallel: true,
        sourceMap: true,
        terserOptions: {
          extractComments: 'all',
          compress: {
            drop_console: true
          }
        }
      })
    ]
  },
  plugins: [
    new MiniCssExtractPlugin({
      filename: 'css/[name].[contenthash].css',
      chunkFilename: 'css/[name].chunk.[contenthash].css'
    }),
    new OptimizeCssAssetsPlugin(),
    // 多入口 链式调用
    new HtmlWebpackPlugin({
      filename: 'index.html',
      template: 'web-common/index.html',
      inject: true,
      hash: false,
      chunks: [
        'userApp',
        'vendor~userApp',
        'vendor~adminApp~userApp',
        'styles~userApp',
        'styles~adminApp~userApp'
      ],
      minify: {
        removeComments: true,
        collapseWhitespace: true,
        removeAttributeQuotes: true
      }
    }),
    new HtmlWebpackPlugin({
      filename: 'admin_index.html',
      template: 'web-common/index.html',
      inject: true,
      hash: false,
      chunks: [
        'adminApp',
        'vendor~adminApp',
        'vendor~adminApp~userApp',
        'styles~adminApp',
        'styles~adminApp~userApp'
      ],
      minify: {
        removeComments: true,
        collapseWhitespace: true,
        removeAttributeQuotes: true
      }
    }),
    new CompressionPlugin({
      test: /\.js(\?.*)?$/i
    }),
    new CleanWebpackPlugin()
  ]
}
module.exports = merge(webpackBaseConfig, webpackProdConfig)
