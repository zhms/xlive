<template>
	<div class="dialogBox">
		<el-dialog :title="title" :visible.sync="visable" width="500px" center @open="handleOpen">
			<el-form label-width="90px" :inline="true" :model="itemdata">
				<el-form-item label="角色名:" style="margin-left: 30px">
					<el-input v-model="itemdata.role_name" :disabled="title == '编辑角色'"></el-input>
				</el-form-item>
				<el-form-item label="备注:" style="margin-left: 30px">
					<el-input v-model="itemdata.memo"></el-input>
				</el-form-item>
				<el-form-item label="上级角色:" style="margin-left: 30px">
					<el-select v-model="itemdata.parent" placeholder="请选择" style="width: 150px" @change="showAuthTree" clearable>
						<el-option v-for="item in roles" :key="item.role_name" :label="item.role_name" :value="item.role_name"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item label="权限:" style="margin-left: 30px" v-show="tree_show">
					<el-tree style="margin-left: 85px; margin-top: -30px; width: 250px" ref="authtree" :default-checked-keys="tree_selected" node-key="path" :props="tree_data" show-checkbox> </el-tree>
				</el-form-item>
			</el-form>
			<span slot="footer" class="dialog-footer">
				<el-button type="primary" @click="handleCommit">确定</el-button>
			</span>
		</el-dialog>
	</div>
</template>
<script>
import dlgbase from '@/api/dlgbase'
export default {
	extends: dlgbase,
	data() {
		return {
			roles: [],
			tree_show: false,
			tree_data: {},
			tree_selected: [],
			parentroledata: null,
			superroledata: null,
			roledata: null,
		}
	},
	methods: {
		commitData(next) {
			if (!this.itemdata.role_name) return this.$message.error('请填写角色名')
			if (!this.itemdata.parent) return this.$message.error('请选择上级角色')
			let setdisable = (node) => {
				for (let n in node) {
					if (typeof node[n] == 'object') {
						setdisable(node[n])
					} else {
						node[n] = 0
					}
				}
			}
			let newroledata = JSON.parse(JSON.stringify(this.superroledata))
			setdisable(newroledata)
			let selected = this.$refs.authtree.getCheckedNodes()
			for (let i = 0; i < selected.length; i++) {
				if (!selected[i].leaf) continue
				let path = selected[i].path.split('.')
				let pn = newroledata
				for (let i = 0; i < path.length - 1; i++) {
					pn = pn[path[i]]
				}
				pn[path[path.length - 1]] = 1
			}
			if (this.title == '编辑角色') {
				let data = {
					seller_id: this.itemdata.seller_id,
					role_name: this.itemdata.role_name,
					role_data: JSON.stringify(newroledata),
					memo: this.itemdata.memo,
					parent: this.itemdata.parent,
				}
				this.$post('/v1/admin_role/update_admin_role', data, { google: true }).then(() => {
					this.$message.success('编辑成功')
					next(true)
				})
			}
			if (this.title == '添加角色') {
				let data = {
					seller_id: this.itemdata.seller_id || this.seller_id,
					parent: this.itemdata.parent,
					role_name: this.itemdata.role_name,
					role_data: JSON.stringify(newroledata),
					memo: this.itemdata.memo,
					state: 1,
				}
				this.$post('/v1/admin_role/create_admin_role', data, { google: true }).then(() => {
					this.$message.success('添加成功')
					next(true)
				})
			}
		},
		onOpen() {
			this.tree_show = false
			this.superroledata = null
			this.roles = []
			if (this.title == '编辑角色') {
				this.getRoles()
				this.showAuthTree()
			} else {
				this.getRoles()
			}
		},
		getRoles() {
			setTimeout(() => {
				let data = {
					seller_id: this.itemdata.seller_id || this.seller_id,
				}
				this.$post('/v1/admin_role/get_admin_role', data).then((roledata) => {
					this.roles = roledata.data
				})
			}, 200)
		},
		showAuthTree() {
			if (!this.itemdata.parent) {
				this.tree_show = false
				return
			}
			let data1 = {
				seller_id: this.itemdata.seller_id || this.seller_id,
				role_name: this.itemdata.parent,
			}
			if (data1.role_name == 'god') data1.role_name = '运营商超管'
			this.$post('/v1/admin_role/get_admin_role', data1).then((parentrole) => {
				this.parentroledata = JSON.parse(parentrole.data[0].role_data)
				let data2 = {
					seller_id: this.itemdata.seller_id || this.seller_id,
					role_name: '运营商超管',
				}
				this.$post('/v1/admin_role/get_admin_role', data2).then((superrole) => {
					this.superroledata = JSON.parse(superrole.data[0].role_data)
					let data3 = {
						seller_id: this.itemdata.seller_id || this.seller_id,
						role_name: this.itemdata.role_name || '-',
					}
					this.$post('/v1/admin_role/get_admin_role', data3).then((roledata) => {
						console.log(roledata)
						if (roledata.data.length == 0) {
							roledata.data = [
								{
									role_data: '{"系统首页":{"查":1}}',
								},
							]
						}
						this.roledata = JSON.parse(roledata.data[0].role_data)
						let treedata = this.getTreeData()
						this.$refs.authtree.root.setData(treedata.menu)
						this.tree_selected = treedata.selected
						this.tree_show = true
					})
				})
			})
		},
		getTreeData() {
			let setdisable = (node) => {
				for (let n in node) {
					if (typeof node[n] == 'object') {
						setdisable(node[n])
					} else {
						node[n] = 0
					}
				}
			}
			setdisable(this.superroledata)
			let setenable = (parent, node) => {
				for (let n in node) {
					if (typeof node[n] == 'object') {
						let p = parent + `.${n}`
						setenable(p, node[n])
					} else {
						if (node[n] == 1) {
							let p = parent.split('.')
							let pn = this.superroledata
							for (let j = 0; j < p.length; j++) {
								pn = pn[p[j]]
							}
							pn[n] = 1
						}
					}
				}
			}
			for (let n in this.parentroledata) {
				let parent = `${n}`
				setenable(parent, this.parentroledata[n])
			}
			let menu = []
			let submenu = (node, root) => {
				for (let n in root) {
					if (typeof root[n] == 'object') {
						let subnode = {
							path: node.path + '.' + n,
							label: n,
							children: [],
						}
						node.children.push(subnode)
						submenu(subnode, root[n])
					} else {
						let path = node.path + '.' + n
						let p = path.split('.')
						let pr = this.parentroledata
						for (let i = 0; i < p.length; i++) {
							pr = pr[p[i]]
						}
						if (pr == 1) {
							let subnode = {
								path: path,
								label: n,
								leaf: true,
							}
							node.children.push(subnode)
						}
					}
				}
			}
			for (let n in this.superroledata) {
				let node = {
					path: n,
					label: n,
					children: [],
				}
				if (n == '系统首页') continue
				menu.push(node)
				submenu(node, this.superroledata[n])
			}
			let selected = []
			let getselected = (parent, node) => {
				for (let n in node) {
					if (typeof node[n] == 'object') {
						let p = parent + `.${n}`
						getselected(p, node[n])
					} else {
						if (node[n] == 1) {
							selected.push(`${parent}.${n}`)
						}
					}
				}
			}
			for (let n in this.roledata) {
				let parent = `${n}`
				getselected(parent, this.roledata[n])
			}
			for (let i = 0; i < menu.length; i++) {
				if (!menu[i].children) continue
				for (let j = 0; j < menu[i].children.length; j++) {
					if (!menu[i].children[j].children) continue
					for (let k = 0; k < menu[i].children[j].children.length; k++) {
						if (!menu[i].children[j].children[k].children) continue
						if (menu[i].children[j].children[k].children.length == 0) {
							menu[i].children[j].children.splice(k, 1)
							k--
						}
					}
				}
			}
			for (let i = 0; i < menu.length; i++) {
				if (!menu[i].children) continue
				for (let j = 0; j < menu[i].children.length; j++) {
					if (!menu[i].children[j].children) continue
					if (menu[i].children[j].children.length == 0) {
						menu[i].children.splice(j, 1)
						j--
					}
				}
			}
			for (let i = 0; i < menu.length; i++) {
				if (!menu[i].children) continue
				if (menu[i].children.length == 0) {
					menu.splice(i, 1)
					i--
				}
			}
			return { menu, selected }
		},
	},
}
</script>
