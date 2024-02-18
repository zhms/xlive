import path from 'path'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const INVALID_CHAR_REGEX = /[\x00-\x1F\x7F<>*#"{}|^[\]`;?:&=+$,]/g
const DRIVE_LETTER_REGEX = /^[a-z]:/i

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [vue()],
	base: '',
	resolve: {
		alias: {
			'@': fileURLToPath(new URL('./src', import.meta.url)),
			'vue-i18n': 'vue-i18n/dist/vue-i18n.cjs.js',
		},
	},
	css: {
		preprocessorOptions: {
			scss: {
				additionalData: '@import "./src/assets/global.scss";',
			},
		},
	},
	server: {
		host: '0.0.0.0',
		proxy: {
			'/api': {
				// target: "https://blackrock.bochats.com",
				//target: 'http://blackrock-13.232.26.141.nip.io ',
				target: 'http://localhost:3400',
				changeOrigin: true,
			},
			'/ws': {
				target: 'ws://qa-13.232.26.141.nip.io',
				changeOrigin: true,
				ws: true,
			},
		},
	},
	build: {
		// 修复github pages不支持_plugin-vue_export-helper问题
		rollupOptions: {
			output: {
				// https://github.com/rollup/rollup/blob/master/src/utils/sanitizeFileName.ts
				sanitizeFileName(name) {
					const match = DRIVE_LETTER_REGEX.exec(name)
					const driveLetter = match ? match[0] : ''
					// substr 是被淘汰語法，因此要改 slice
					return driveLetter + name.slice(driveLetter.length).replace(INVALID_CHAR_REGEX, '')
				},
			},
		},
	},
})
