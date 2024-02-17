/* eslint-disable */
<template>
	<div class="login-container">
		<el-form ref="loginForm" :model="form_data" class="login-form" autocomplete="on" label-position="left">
			<div class="title-container">
				<h3 class="title">演示系统</h3>
			</div>
			<el-form-item prop="Account">
				<span class="svg-container">
					<svg-icon icon-class="user" />
				</span>
				<el-input ref="Account" v-model="form_data.Account" placeholder="用户名" name="Account" type="text" tabindex="1" autocomplete="on" />
			</el-form-item>
			<el-form-item prop="Password">
				<span class="svg-container">
					<svg-icon icon-class="password" />
				</span>
				<el-input key="password" ref="Password" v-model="form_data.Password" type="password" placeholder="密码" name="Password" tabindex="2" autocomplete="on" @keyup.native="handleChangePassword" @click="handleLogin" />
			</el-form-item>
			<el-button type="primary" style="width: 100%; margin-bottom: 30px" @click.native.prevent="handleLogin">登录</el-button>
		</el-form>
	</div>
</template>

<script>
import { isStringAndNotEmpty } from '@/api/utils'
import moment from 'moment'
export default {
	name: 'Login',
	data() {
		return {
			form_data: {
				Account: localStorage.getItem('admin_account'),
				Password: localStorage.getItem('admin_password'),
			},
			password_changed: false,
		}
	},

	created() {},
	mounted() {
		if (this.form_data.Account == '') {
			this.$refs.Account.focus()
		} else if (this.form_data.Password == '') {
			this.$refs.Password.focus()
		}
	},
	destroyed() {},
	methods: {
		handleChangePassword() {
			this.password_changed = true
		},
		handleLogin() {
			if (!isStringAndNotEmpty(this.form_data.Account)) return this.$message.error('请填写账号')
			if (!isStringAndNotEmpty(this.form_data.Password)) return this.$message.error('请填写密码')
			let password = this.form_data.Password
			if (this.password_changed) password = this.$md5(this.form_data.Password)

			let logindata = {
				Account: this.form_data.Account,
				Password: password,
			}
			this.$store.dispatch('user/login', logindata).then(() => {
				localStorage.setItem('admin_account', this.form_data.Account)
				localStorage.setItem('admin_password', password)
				this.$message.success('登录成功')
				this.$store.dispatch('router/getRouters').then(() => {
					this.$router.push({ path: '/dashboard' })
				})
			})
		},
	},
}
</script>
<style lang="scss">
$bg: #283443;
$light_gray: #fff;
$cursor: #fff;
@supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
	.login-container .el-input input {
		color: $cursor;
	}
}
.login-container {
	.el-input {
		display: inline-block;
		height: 47px;
		width: 85%;
		input {
			background: transparent;
			border: 0px;
			-webkit-appearance: none;
			appearance: none;
			border-radius: 0px;
			padding: 12px 5px 12px 15px;
			color: $light_gray;
			height: 47px;
			caret-color: $cursor;
			&:-webkit-autofill {
				box-shadow: 0 0 0px 1000px $bg inset !important;
				-webkit-text-fill-color: $cursor !important;
			}
		}
	}
	.el-form-item {
		border: 1px solid rgba(255, 255, 255, 0.1);
		background: rgba(0, 0, 0, 0.1);
		border-radius: 5px;
		color: #454545;
	}
}
</style>
<style lang="scss" scoped>
$bg: #2d3a4b;
$dark_gray: #889aa4;
$light_gray: #eee;
.login-container {
	min-height: 100%;
	width: 100%;
	background-color: $bg;
	overflow: hidden;
	.login-form {
		position: relative;
		width: 520px;
		max-width: 100%;
		padding: 160px 35px 0;
		margin: 0 auto;
		overflow: hidden;
	}
	.tips {
		font-size: 14px;
		color: #fff;
		margin-bottom: 10px;
		span {
			&:first-of-type {
				margin-right: 16px;
			}
		}
	}
	.svg-container {
		padding: 6px 5px 6px 15px;
		color: $dark_gray;
		vertical-align: middle;
		width: 30px;
		display: inline-block;
	}
	.title-container {
		position: relative;
		.title {
			font-size: 26px;
			color: $light_gray;
			margin: 0px auto 40px auto;
			text-align: center;
			font-weight: bold;
		}
	}
	.show-pwd {
		position: absolute;
		right: 10px;
		top: 7px;
		font-size: 16px;
		color: $dark_gray;
		cursor: pointer;
		user-select: none;
	}
	.thirdparty-button {
		position: absolute;
		right: 0;
		bottom: 6px;
	}
	@media only screen and (max-width: 470px) {
		.thirdparty-button {
			display: none;
		}
	}
}
</style>
