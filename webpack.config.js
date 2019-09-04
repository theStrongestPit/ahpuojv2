const path = require('path');    //path 模块提供了一些工具函数，用于处理文件与目录的路径。
const webpack = require('webpack');       //webpack打包工具
const VueLoaderPlugin = require('vue-loader/lib/plugin');         // vue-loader 编译vue文件
const compiler = require('vue-template-compiler')            // 模板函数编译 与vue-loader配合使用
const ExtractTextPlugin = require("extract-text-webpack-plugin");
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    //mode: "development",
    entry: {       //入口
        "adminApp": "./web-admin/main.js",
        "userApp": "./web-user/main.js"
    },
    module: {            //处理项目中的不同类型的模块。
        rules: [      // rules 各种规则(数组类型) 每个规则可以分为三部分 - 条件(condition)，结果(result)和嵌套规则(nested rule)
            {
                test: /\.vue$/,
                loader: 'vue-loader',
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
                    path.resolve('/node_modules/vue-image-crop-upload'),
                ],
                use: {
                    loader: "babel-loader",
                    options: {
                        plugins: [
                            "babel-plugin-syntax-dynamic-import"
                        ]
                    }
                },
            },
            {
                test: /\.css$/,
                use: [
                    'style-loader',
                    { loader: 'css-loader', options: { sourceMap: true } },
                ]
            },
            {
                test: /\.svg$/,
                loader: 'svg-sprite-loader',
                include: [path.resolve('static/icons')],
                options: {
                    symbolId: 'icon-[name]'
                }
            },
            {
                test: /\.(ttf|eot|woff|woff2|gif)$/,
                use: [
                    'url-loader'
                ]
            },
            {
                test: /\.(sass|scss)$/,
                use: [
                    "style-loader",
                    "css-loader",
                    {
                        loader: 'postcss-loader',
                        options: {
                            sourceMap: true,
                            config: {
                                path: 'postcss.config.js'
                            }
                        }
                    },
                    { loader: 'sass-loader', options: { sourceMap: true } },
                    {
                        loader: 'sass-resources-loader',
                        options: {
                            sourceMap: true,
                            resources: [
                                path.resolve('web-common/sass/color.scss'),
                            ]
                        }
                    }
                ]
            }
        ]
    },
    plugins: [
        // 不支持独立scss文件的热更新 开发环境不使用
        // new ExtractTextPlugin({
        //     filename: "[name].min.css",
        //     allChunks: false // 注意 2
        // }),
        // new CleanWebpackPlugin(['dist']),  // 告诉清理插件的目录
        new HtmlWebpackPlugin({
            filename: 'index.html',
            template: 'web-user/index.html',
            inject: true,
            hash: false,
            minify: {
                removeComments: true,
                collapseWhitespace: true,
                removeAttributeQuotes: true
            }
        }),
        new HtmlWebpackPlugin({
            filename: 'admin_index.html',
            template: 'web-admin/admin_index.html',
            inject: true,
            hash: false,
            minify: {
                removeComments: true,
                collapseWhitespace: true,
                removeAttributeQuotes: true
            }
        }),
        new VueLoaderPlugin()                 //vue-loader插件开启
    ],
    output: {        //出口
        filename: '[name].bundle.js',    //文件名
        chunkFilename: "[name].chunk.js",
        path: path.resolve(__dirname, 'dist'),   //路径
        publicPath: "/dist/"        //srcript 引入路径
    },
    resolve: {
        //引入路径是不用写对应的后缀名
        extensions: ['.js', '.vue', '.json'],
        alias: {
            //正在使用的是vue的运行时版本，而此版本中的编译器时不可用的，我们需要把它切换成运行时 + 编译的版本
            'vue$': 'vue/dist/vue.runtime.esm.js',
            '@': path.resolve(__dirname, ''),
        }
    },
    devServer: {
        contentBase: './',
        hot: true,
        host: '127.0.0.1',
        port: '8888',
        disableHostCheck: true,
        proxy: {
            '/': {
                target: 'http://localhost',
                secure: false
            },
            '/api': {
                target: 'http://localhost:8080',
                secure: false
            }
        },
    },
};