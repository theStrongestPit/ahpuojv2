const path = require('path');
const HappyPack = require('happypack');
const VueLoaderPlugin = require('vue-loader/lib/plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
module.exports = {
  entry: {
    // 入口
    adminApp: './web-admin/main.js',
    userApp: './web-user/main.js'
  },
  module: {
    // 处理项目中的不同类型的模块。
    rules: [
      // rules 各种规则(数组类型) 每个规则可以分为三部分 - 条件(condition)，结果(result)和嵌套规则(nested rule)
      {
        test: /\.vue$/,
        loader: 'vue-loader'
      },
      {
        test: /\.pug$/,
        loader: 'pug-plain-loader'
      },
      {
        test: /\.js$/,
        include: [
          path.resolve('web-user'),
          path.resolve('web-admin'),
          path.resolve('web-common'),
          path.resolve('/node_modules/vue-image-crop-upload')
        ],
        loader: 'happypack/loader?id=happybabel'
      },
      {
        test: /\.css$/,
        use: [
          'style-loader',
          {
            loader: 'css-loader',
            options: {
              sourceMap: true
            }
          }
        ]
      },
      {
        test: /\.(png|jpe?g|gif|webp)(\?.*)?$/,
        use: [
          {
            loader: 'url-loader',
            options: {
              limit: 4096,
              fallback: {
                loader: 'file-loader',
                options: {
                  name: 'static/img/[name].[hash:8].[ext]'
                }
              }
            }
          }
        ]
      },
      {
        test: /\.svg$/,
        loader: 'svg-sprite-loader',
        include: [path.resolve('web-common/assets/icons')],
        options: {
          symbolId: 'icon-[name]'
        }
      },
      {
        test: /\.(ttf|eot|woff|woff2|gif)$/,
        use: ['url-loader']
      },
      {
        test: /\.(sass|scss)$/,
        use: [
          'style-loader',
          'css-loader',
          {
            loader: 'postcss-loader',
            options: {
              sourceMap: true,
              config: {
                path: 'postcss.config.js'
              }
            }
          },
          {
            loader: 'sass-loader',
            options: {
              sourceMap: true
            }
          },
          {
            loader: 'sass-resources-loader',
            options: {
              sourceMap: true,
              resources: [path.resolve('web-common/sass/color.scss')]
            }
          }
        ]
      }
    ]
  },
  plugins: [
    new HappyPack({
      id: 'happybabel',
      loaders: [
        {
          loader: 'babel-loader',
          options: {
            plugins: ['babel-plugin-syntax-dynamic-import']
          }
        }
      ],
      // 开启 4 个线程
      threads: 4
    }),
    new VueLoaderPlugin(), // vue-loader插件开启,
    new CopyWebpackPlugin([
      {
        from: path.resolve('web-common/static'),
        to: 'static/'
      }
    ])
  ],
  resolve: {
    // 引入路径是不用写对应的后缀名
    extensions: ['.js', '.vue', '.json'],
    alias: {
      // 正在使用的是vue的运行时版本，而此版本中的编译器时不可用的，我们需要把它切换成运行时 + 编译的版本
      vue$: 'vue/dist/vue.runtime.esm.js',
      '@': path.resolve(__dirname, '')
    }
  }
};
