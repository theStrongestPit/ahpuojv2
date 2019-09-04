require('babel-polyfill')
const path = require('path'); //path 模块提供了一些工具函数，用于处理文件与目录的路径。
const webpack = require('webpack'); //webpack打包工具
const VueLoaderPlugin = require('vue-loader/lib/plugin'); // vue-loader 编译vue文件
const compiler = require('vue-template-compiler') // 模板函数编译 与vue-loader配合使用
// const ExtractTextPlugin = require("extract-text-webpack-plugin");
const CompressionPlugin = require('compression-webpack-plugin');
const TerserPlugin = require('terser-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const OptimizeCssAssetsPlugin = require('optimize-css-assets-webpack-plugin');//压缩css插件

module.exports = {
    //mode: "development",
    entry: { //入口
        "adminApp": "./web-admin/main.js",
        "userApp": "./web-user/main.js"
    },
    module: { //处理项目中的不同类型的模块。
        rules: [ // rules 各种规则(数组类型) 每个规则可以分为三部分 - 条件(condition)，结果(result)和嵌套规则(nested rule)
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
                    {
                        loader: "style-loader",
                        options: {
                            singleton: true
                        }
                    },
                    MiniCssExtractPlugin.loader,
                    {
                        loader: 'css-loader',
                        options: {
                            sourceMap: true,
                        }
                    },
                    {
                        loader: 'postcss-loader',
                        options: {
                            sourceMap: true,
                            config: {
                                path: 'postcss.config.js'
                            }
                        }
                    },
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
                    MiniCssExtractPlugin.loader,
                    {
                        loader: 'css-loader',
                        options: {
                            sourceMap: true,
                        }
                    },
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
            },
            {
                test: /\.vue$/,
                use: [
                    'vue-loader',
                ],
            },
            {
                test: /\.pug$/,
                use: [
                    'pug-plain-loader'
                ]
            }
        ]
    },
    optimization: {
        splitChunks: {
            cacheGroups: {
                vendor: {
                    test: /node_modules/,
                    chunks: "initial",
                    priority: 10,
                    enforce: true
                },
                styles: {
                    test: /\.s?css$/,
                    chunks: 'all',
                    enforce: true,
                    priority: 20,
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
                        drop_console: true,
                    },
                }
            }),
        ],
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: "[name].[contenthash].css",
            chunkFilename: "[name].chunk.[contenthash].css",
        }),
        new OptimizeCssAssetsPlugin(),
        // 多入口 链式调用
        new HtmlWebpackPlugin({
            filename: 'index.html',
            template: 'web-user/index.html',
            inject: true,
            hash: false,
            chunks: ['userApp', 'vendor~userApp', 'vendor~adminApp~userApp', 'styles~userApp', 'styles~adminApp~userApp'],
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
            chunks: ['adminApp', 'vendor~adminApp', 'vendor~adminApp~userApp', 'styles~adminApp', 'styles~adminApp~userApp'],
            minify: {
                removeComments: true,
                collapseWhitespace: true,
                removeAttributeQuotes: true
            }
        }),
        new VueLoaderPlugin(), //vue-loader插件开启
        new CompressionPlugin({
            test: /\.js(\?.*)?$/i,
        }),
        new CleanWebpackPlugin(
            {
                cleanOnceBeforeBuildPatterns: ['dist/[name].bundle.*.js', 'dist/[name].chunk.*.js', '[name].*.css'],　 //匹配删除的文件　　　　　　　　　//根目录
                verbose: true,        　　　　　　　　　　//开启在控制台输出信息
                dry: false        　　　　　　　　　　//启用删除文件
            }
        )
    ],
    output: { //出口
        filename: '[name].bundle.[chunkhash].js', //文件名
        chunkFilename: '[name].chunk.[chunkhash].js',
        path: path.resolve(__dirname, 'dist'), //路径
        publicPath: '/dist/' //srcript 引入路径
    },
    resolve: {
        //引入路径是不用写对应的后缀名
        extensions: ['.js', '.vue', '.json'],
        alias: {
            // 线上环境使用运行时版本
            'vue$': 'vue/dist/vue.runtime.esm.js',
            '@': path.resolve(__dirname, ''),
        }
    }
};