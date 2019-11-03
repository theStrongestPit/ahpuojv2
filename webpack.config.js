const path = require('path'); // path 模块提供了一些工具函数，用于处理文件与目录的路径。

const HtmlWebpackPlugin = require('html-webpack-plugin');
const merge = require('webpack-merge');
const webpackBaseConfig = require('./webpack.base.config');

const webpackDevConfig = {
  mode: 'development',
  devtool: 'source-map',
  output: {
    // 出口
    filename: 'js/[name].bundle.js', // 文件名
    chunkFilename: 'js/[name].chunk.js',
    path: path.resolve(__dirname, 'dist'), // 路径
    publicPath: '/' // srcript 引入路径
  },
  module: {
    // 处理项目中的不同类型的模块。
    rules: [
      {
        test: /\.(js|vue)$/,
        loader: 'eslint-loader',
        enforce: 'pre',
        include: [
          path.resolve(__dirname, 'web-admin'),
          path.resolve(__dirname, 'web-common'),
          path.resolve(__dirname, 'web-user')
        ], // 指定检查的目录
        options: {
          // 这里的配置项参数将会被传递到 eslint 的 CLIEngine
          formatter: require('eslint-friendly-formatter') // 指定错误报告的格式规范
        }
      }
    ]
  },
  plugins: [
    // 多入口 链式调用
    new HtmlWebpackPlugin({
      filename: 'index.html',
      template: 'web-common/index.html',
      inject: true,
      hash: false,
      chunks: ['userApp'],
      minify: {
        removeComments: true,
        collapseWhitespace: true,
        removeAttributeQuotes: true
      }
    }),
    new HtmlWebpackPlugin({
      filename: 'admin/index.html',
      template: 'web-common/index.html',
      inject: true,
      hash: false,
      chunks: ['adminApp'],
      minify: {
        removeComments: true,
        collapseWhitespace: true,
        removeAttributeQuotes: true
      }
    })
  ],
  devServer: {
    open: true,
    hot: true,
    host: '127.0.0.1',
    port: '8888',
    disableHostCheck: true,
    // 这个配置 真是找了好久才知道 这样可以起到和生产环境配置nginx一样的效果
    historyApiFallback: {
      rewrites: [
        {from: /^\/$/, to: '/index.html'},
        {from: /^\/admin/, to: '/admin/index.html'}
      ]
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        secure: false
      }
    }
  }
};

module.exports = merge(webpackBaseConfig, webpackDevConfig);
