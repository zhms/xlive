function pluralize(time, label) {
	if (time === 1) {
		return time + label
	}
	return time + label + 's'
}
export function timeAgo(time) {
	const between = Date.now() / 1000 - Number(time)
	if (between < 3600) {
		return pluralize(~~(between / 60), ' minute')
	} else if (between < 86400) {
		return pluralize(~~(between / 3600), ' hour')
	} else {
		return pluralize(~~(between / 86400), ' day')
	}
}
//like 10000 => 10k
export function numberFormatter(num, digits) {
	const si = [
		{ value: 1e18, symbol: 'E' },
		{ value: 1e15, symbol: 'P' },
		{ value: 1e12, symbol: 'T' },
		{ value: 1e9, symbol: 'G' },
		{ value: 1e6, symbol: 'M' },
		{ value: 1e3, symbol: 'k' },
	]
	for (let i = 0; i < si.length; i++) {
		if (num >= si[i].value) {
			return (num / si[i].value).toFixed(digits).replace(/\.0+$|(\.[0-9]*[1-9])0+$/, '$1') + si[i].symbol
		}
	}
	return num.toString()
}
//10000 => "10,000"
export function toThousandFilter(num) {
	return (+num || 0).toString().replace(/^-?\d+/g, (m) => m.replace(/(?=(?!\b)(\d{3})+$)/g, ','))
}
//首字母大写
export function uppercaseFirst(string) {
	return string.charAt(0).toUpperCase() + string.slice(1)
}

export function shortString(orignalstr, startcount, endcount) {
	if (orignalstr.length >= startcount + endcount) {
		let start = orignalstr.substr(0, startcount)
		let end = orignalstr.substr(orignalstr.length - endcount, endcount)
		return start + '***' + end
	}
	return orignalstr
}

export function numDeal(t, n) {
	if (t == 0) return 0
	let f = t / Math.abs(t)
	n = n || 6
	if (n < 0) n = 0
	t = t * Math.pow(10, n)
	t = Math.abs(t)
	t = Math.floor(t)
	return (t / Math.pow(10, n)) * f
}
