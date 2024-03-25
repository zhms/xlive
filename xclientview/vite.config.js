import path from 'path'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const INVALID_CHAR_REGEX = /[\x00-\x1F\x7F<>*#"{}|^[\]`;?:&=+$,]/g
const DRIVE_LETTER_REGEX = /^[a-z]:/i

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
				additionalData: '@import "./src/css/global.scss";',
			},
		},
	},
	server: {
		host: '0.0.0.0',
		proxy: {
			'/api': {
				// target: "https://blackrock.bochats.com",
				//target: 'http://blackrock-13.232.26.141.nip.io ',
				target: 'http://10.10.234.82:3400',
				changeOrigin: true,
				ws: true,
			},
			'/ws': {
				//target: 'ws://qa-13.232.26.141.nip.io',
				target: 'ws://10.10.234.82:2740',
				changeOrigin: true,
				ws: true,
			},
		},
	},
	build: {
		rollupOptions: {
			output: {
				sanitizeFileName(name) {
					const match = DRIVE_LETTER_REGEX.exec(name)
					const driveLetter = match ? match[0] : ''
					return driveLetter + name.slice(driveLetter.length).replace(INVALID_CHAR_REGEX, '')
				},
			},
		},
	},
})
