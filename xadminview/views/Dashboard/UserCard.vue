<template>
	<el-card style="margin-bottom: 20px; width: 400px">
		<div slot="header" class="clearfix">
			<span>个人资料</span>
		</div>
		<div class="user-profile">
			<div class="box-center">
				<pan-thumb image="avatar.gif" :height="'100px'" :width="'100px'" :hoverable="false">
					<div style="margin-top: -5px">小兔子乖乖,把门儿打开</div>
				</pan-thumb>
			</div>
			<div class="box-center">
				<div class="user-name text-center">{{ userinfo.Account }}</div>
			</div>
		</div>

		<div class="user-bio" style="margin-top: -15px">
			<div class="user-skills user-bio-section">
				<div class="user-bio-section-header"></div>
				<div class="user-bio-section-body">
					<div class="progress-item">
						<span style="margin-left: 80px">{{ `登录次数 ：${userinfo.LoginCount}` }}</span>
					</div>
					<div class="progress-item" style="margin-top: 10px">
						<span style="margin-left: 96px">{{ `登录Ip ： ${userinfo.Ip}` }}</span>
					</div>
					<div class="progress-item" style="margin-top: 10px">
						<span style="margin-left: 52px">{{ `上次登录时间 ： ${userinfo.LoginTime}` }}</span>
					</div>
					<div style="margin-top: 20px">
						<el-button type="danger" style="width: 200px; margin-left: 85px" @click="logout">退出登录</el-button>
					</div>
				</div>
			</div>
		</div>
	</el-card>
</template>

<script>
import { mapGetters } from 'vuex'
import PanThumb from './PanThumb.vue'
export default {
	components: { PanThumb },
	props: {
		user: {
			type: Object,
			default: () => {
				return {
					name: '',
					email: '',
					avatar: '',
					role: '',
				}
			},
		},
	},
	computed: {
		...mapGetters(['userinfo']),
	},
	methods: {
		logout() {
			this.$post('/v1/admin_user_logout', { Token: sessionStorage.getItem('x-token') }).then(() => {})
			sessionStorage.removeItem('userinfo')
			sessionStorage.removeItem('x-token')
			this.$router.push({ path: '/' })
		},
	},
}
</script>

<style lang="scss" scoped>
.box-center {
	margin: 0 auto;
	display: table;
}

.text-muted {
	color: #777;
}

.user-profile {
	.user-name {
		font-weight: bold;
	}

	.box-center {
		padding-top: 10px;
	}

	.user-role {
		padding-top: 10px;
		font-weight: 400;
		font-size: 14px;
	}

	.box-social {
		padding-top: 30px;

		.el-table {
			border-top: 1px solid #dfe6ec;
		}
	}

	.user-follow {
		padding-top: 20px;
	}
}

.user-bio {
	margin-top: 20px;
	color: #606266;

	span {
		padding-left: 4px;
	}

	.user-bio-section {
		font-size: 14px;
		padding: 15px 0;

		.user-bio-section-header {
			border-bottom: 1px solid #dfe6ec;
			padding-bottom: 10px;
			margin-bottom: 10px;
			font-weight: bold;
		}
	}
}
</style>
