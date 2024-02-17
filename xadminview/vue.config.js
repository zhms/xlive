'use strict'
const path = require('path')
function resolve(dir) {
	return path.join(__dirname, dir)
}
const name = '哈希后台管理'
module.exports = {
	publicPath: '/',
	outputDir: 'dist',
	assetsDir: 'static',
	productionSourceMap: false,
	devServer: {
		proxy: {
			'/api': {
				target: 'http://127.0.0.1:6500', //https://admin.f6921k.live  http://192.168.2.33:4534
				changeOrigin: true,
				pathRewrite: {
					'^/api': '/api',
				},
			},
		},
	},
	configureWebpack: {
		name: name,
		resolve: {
			alias: {
				'@': resolve('src'),
			},
		},
	},
	chainWebpack(config) {
		config.plugin('preload').tap(() => [
			{
				rel: 'preload',
				fileBlacklist: [/\.map$/, /hot-update\.js$/, /runtime\..*\.js$/],
				include: 'initial',
			},
		])
		config.plugins.delete('prefetch')
		config.module.rule('svg').exclude.add(resolve('src/icons')).end()
		config.module
			.rule('icons')
			.test(/\.svg$/)
			.include.add(resolve('src/icons'))
			.end()
			.use('svg-sprite-loader')
			.loader('svg-sprite-loader')
			.options({
				symbolId: 'icon-[name]',
			})
			.end()
		config.when(process.env.NODE_ENV !== 'development', (config) => {
			config
				.plugin('ScriptExtHtmlWebpackPlugin')
				.after('html')
				.use('script-ext-html-webpack-plugin', [
					{
						inline: /runtime\..*\.js$/,
					},
				])
				.end()
			config.optimization.splitChunks({
				chunks: 'all',
				cacheGroups: {
					libs: {
						name: 'chunk-libs',
						test: /[\\/]node_modules[\\/]/,
						priority: 10,
						chunks: 'initial',
					},
					elementUI: {
						name: 'chunk-elementUI',
						priority: 20,
						test: /[\\/]node_modules[\\/]_?element-ui(.*)/,
					},
					commons: {
						name: 'chunk-commons',
						test: resolve('src/components'),
						minChunks: 3,
						priority: 5,
						reuseExistingChunk: true,
					},
				},
			})
			config.optimization.runtimeChunk('single')
		})
	},
}
