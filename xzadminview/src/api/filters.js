import moment from 'moment'
//格式化时间为北京时间
export function 北京时间(value) {
	if (!value || value == '') return ''
	return moment(value).format('YYYY-MM-DD HH:mm:ss')
}
//格式化时间为utc时间
export function 世界时间(value) {
	if (!value || value == '') return ''
	return moment(value).utc().format('YYYY-MM-DD HH:mm:ss')
}

export function 是否(value) {
	if (value == 1) return '是'
	return '否'
}

export function 启用禁用(value) {
	if (value == 1) return '启用'
	return '禁用'
}

export function amount2(value) {
	if (!value) return '0.00'
	value = parseFloat(value)
	value = Math.floor(value * 100) / 100
	return value
}
