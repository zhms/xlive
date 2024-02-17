import moment from 'moment'
export function getQueryObject(url) {
	url = url == null ? window.location.href : url
	const search = url.substring(url.lastIndexOf('?') + 1)
	const obj = {}
	const reg = /([^?&=]+)=([^?&=]*)/g
	search.replace(reg, (rs, $1, $2) => {
		const name = decodeURIComponent($1)
		let val = decodeURIComponent($2)
		val = String(val)
		obj[name] = val
		return rs
	})
	return obj
}

export function byteLength(str) {
	let s = str.length
	for (var i = str.length - 1; i >= 0; i--) {
		const code = str.charCodeAt(i)
		if (code > 0x7f && code <= 0x7ff) s++
		else if (code > 0x7ff && code <= 0xffff) s += 2
		if (code >= 0xdc00 && code <= 0xdfff) i--
	}
	return s
}

export function cleanArray(actual) {
	const newArray = []
	for (let i = 0; i < actual.length; i++) {
		if (actual[i]) {
			newArray.push(actual[i])
		}
	}
	return newArray
}

export function param(json) {
	if (!json) return ''
	return cleanArray(
		Object.keys(json).map((key) => {
			if (json[key] === undefined) return ''
			return encodeURIComponent(key) + '=' + encodeURIComponent(json[key])
		})
	).join('&')
}

export function param2Obj(url) {
	const search = decodeURIComponent(url.split('?')[1]).replace(/\+/g, ' ')
	if (!search) {
		return {}
	}
	const obj = {}
	const searchArr = search.split('&')
	searchArr.forEach((v) => {
		const index = v.indexOf('=')
		if (index !== -1) {
			const name = v.substring(0, index)
			const val = v.substring(index + 1, v.length)
			obj[name] = val
		}
	})
	return obj
}

export function objectMerge(target, source) {
	if (typeof target !== 'object') {
		target = {}
	}
	if (Array.isArray(source)) {
		return source.slice()
	}
	Object.keys(source).forEach((property) => {
		const sourceProperty = source[property]
		if (typeof sourceProperty === 'object') {
			target[property] = objectMerge(target[property], sourceProperty)
		} else {
			target[property] = sourceProperty
		}
	})
	return target
}

export function uniqueArr(arr) {
	return Array.from(new Set(arr))
}

export function createUniqueString() {
	const timestamp = +new Date() + ''
	const randomNum = parseInt((1 + Math.random()) * 65536) + ''
	return (+(randomNum + timestamp)).toString(32)
}

export function hasClass(ele, cls) {
	return !!ele.className.match(new RegExp('(\\s|^)' + cls + '(\\s|$)'))
}

export function addClass(ele, cls) {
	if (!hasClass(ele, cls)) ele.className += ' ' + cls
}

export function removeClass(ele, cls) {
	if (hasClass(ele, cls)) {
		const reg = new RegExp('(\\s|^)' + cls + '(\\s|$)')
		ele.className = ele.className.replace(reg, ' ')
	}
}

export function isStringAndNotEmpty(str) {
	if (typeof str === 'string' || str instanceof String) {
		if (str.length > 0) return true
		return false
	}
	return false
}

export function isExternal(path) {
	return /^(https?:|mailto:|tel:)/.test(path)
}

export function validURL(url) {
	const reg =
		/^(https?|ftp):\/\/([a-zA-Z0-9.-]+(:[a-zA-Z0-9.&%$-]+)*@)*((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])){3}|([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.(com|edu|gov|int|mil|net|org|biz|arpa|info|name|pro|aero|coop|museum|[a-zA-Z]{2}))(:[0-9]+)*(\/($|[a-zA-Z0-9.,?'\\+&%$#=~_-]+))*$/
	return reg.test(url)
}

export function validLowerCase(str) {
	const reg = /^[a-z]+$/
	return reg.test(str)
}

export function validUpperCase(str) {
	const reg = /^[A-Z]+$/
	return reg.test(str)
}

export function validEmail(email) {
	const reg = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
	return reg.test(email)
}

export function isArray(arg) {
	if (typeof Array.isArray === 'undefined') {
		return Object.prototype.toString.call(arg) === '[object Array]'
	}
	return Array.isArray(arg)
}

export function isString(str) {
	if (typeof str === 'string' || str instanceof String) {
		return true
	}
	return false
}

export function getTimeStamp(t) {
	if (!t || t == 0) return 0
	return moment.utc(moment(t).format('YYYY-MM-DD HH:mm:ss')).valueOf() - 28800000
}
