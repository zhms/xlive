<template>
	<div class="container">
		<el-row :gutter="20">
			<el-col :span="8">
				<el-card shadow="hover" class="mgb20" style="height: 330px; width: 400px">
					<div class="user-info">
						<img src="../../public/avatar.jpg" class="user-avator" alt />
						<div class="user-info-cont">
							<div class="user-info-name">{{ userinfo.account }}</div>
							<div>{{ userinfo.role }}</div>
						</div>
					</div>
					<div class="user-info-list"><span style="margin-left: 59px">账号：</span> {{ userinfo.account }}</div>
					<div class="user-info-list"><span style="margin-left: 43px">登录IP：</span> {{ userinfo.login_ip }}</div>
					<div class="user-info-list"><span style="margin-left: 27px">登录次数：</span> {{ userinfo.login_count }}</div>
					<div class="user-info-list"><span>上次登录时间：</span> {{ userinfo.login_time }}</div>
				</el-card>
			</el-col>
		</el-row>
	</div>
</template>

<script>
import '@/styles/main.css'
import { mapGetters } from 'vuex'
import CountTo from 'vue-count-to'
import UserCard from './UserCard.vue'
import base from '@/api/base'
export default {
	name: 'Dashboard',
	extends: base,
	components: { CountTo, UserCard },
	data() {
		return {}
	},
	computed: {
		...mapGetters(['userinfo']),
	},
	created() {},
}
</script>

<style lang="scss" scoped>
.panel-group {
	margin-top: 18px;

	.card-panel-col {
		margin-bottom: 32px;
	}

	.card-panel {
		height: 108px;
		cursor: pointer;
		font-size: 12px;
		position: relative;
		overflow: hidden;
		color: #666;
		background: rgb(255, 255, 255);
		box-shadow: 4px 4px 40px rgba(0, 0, 0, 0.05);
		border-color: rgba(0, 0, 0, 0.05);

		&:hover {
			.card-panel-icon-wrapper {
				color: #fff;
			}

			.icon-people {
				background: #40c9c6;
			}

			.icon-message {
				background: #36a3f7;
			}

			.icon-money {
				background: #f4516c;
			}

			.icon-shopping {
				background: #34bfa3;
			}
		}

		.icon-people {
			color: #40c9c6;
		}

		.icon-message {
			color: #36a3f7;
		}

		.icon-money {
			color: #f4516c;
		}

		.icon-shopping {
			color: #34bfa3;
		}

		.card-panel-icon-wrapper {
			float: left;
			margin: 14px 0 0 14px;
			padding: 16px;
			transition: all 0.38s ease-out;
			border-radius: 6px;
		}

		.card-panel-icon {
			float: left;
			font-size: 48px;
		}

		.card-panel-description {
			float: right;
			font-weight: bold;
			margin: 26px;
			margin-left: 0px;

			.card-panel-text {
				line-height: 18px;
				color: rgba(0, 0, 0, 0.45);
				font-size: 16px;
				margin-bottom: 12px;
			}

			.card-panel-num {
				font-size: 20px;
			}
		}
	}
}

@media (max-width: 550px) {
	.card-panel-description {
		display: none;
	}

	.card-panel-icon-wrapper {
		float: none !important;
		width: 100%;
		height: 100%;
		margin: 0 !important;

		.svg-icon {
			display: block;
			margin: 14px auto !important;
			float: none !important;
		}
	}
}
.el-row {
	margin-bottom: 20px;
}

.grid-content {
	display: flex;
	align-items: center;
	height: 100px;
}

.grid-cont-right {
	flex: 1;
	text-align: center;
	font-size: 14px;
	color: #999;
}

.grid-num {
	font-size: 30px;
	font-weight: bold;
}

.grid-con-icon {
	font-size: 50px;
	width: 100px;
	height: 100px;
	text-align: center;
	line-height: 100px;
	color: #fff;
}

.grid-con-1 .grid-con-icon {
	background: rgb(45, 140, 240);
}

.grid-con-1 .grid-num {
	color: rgb(45, 140, 240);
}

.grid-con-2 .grid-con-icon {
	background: rgb(100, 213, 114);
}

.grid-con-2 .grid-num {
	color: rgb(45, 140, 240);
}

.grid-con-3 .grid-con-icon {
	background: rgb(242, 94, 67);
}

.grid-con-3 .grid-num {
	color: rgb(242, 94, 67);
}

.user-info {
	display: flex;
	align-items: center;
	padding-bottom: 20px;
	border-bottom: 2px solid #ccc;
	margin-bottom: 20px;
}

.user-avator {
	width: 120px;
	height: 120px;
	border-radius: 50%;
}

.user-info-cont {
	padding-left: 50px;
	flex: 1;
	font-size: 14px;
	color: #999;
}

.user-info-cont div:first-child {
	font-size: 30px;
	color: #222;
}

.user-info-list {
	font-size: 14px;
	color: #999;
	line-height: 25px;
}

.user-info-list span {
	margin-left: 1px;
}

.mgb20 {
	margin-bottom: 20px;
}

.todo-item {
	font-size: 14px;
}

.todo-item-del {
	text-decoration: line-through;
	color: #999;
}

.schart {
	width: 100%;
	height: 300px;
}
</style>
