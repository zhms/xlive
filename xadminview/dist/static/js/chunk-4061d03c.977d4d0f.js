(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-4061d03c"],{"206e":function(t,e,a){"use strict";a.r(e);var i=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"container"},[a("el-form",{attrs:{inline:!0,model:t.filters}},[a("el-form-item",{attrs:{label:""}},[a("el-input",{staticStyle:{width:"150px"},attrs:{placeholder:"管理员",clearable:!0},model:{value:t.filters.account,callback:function(e){t.$set(t.filters,"account",e)},expression:"filters.account"}})],1),a("el-form-item",{attrs:{label:""}},[a("el-input",{staticStyle:{width:"150px"},attrs:{placeholder:"操作",clearable:!0},model:{value:t.filters.opt_name,callback:function(e){t.$set(t.filters,"opt_name",e)},expression:"filters.opt_name"}})],1),a("el-form-item",{attrs:{label:""}},[a("el-date-picker",{attrs:{type:"daterange","range-separator":"至","start-placeholder":"开始日期","end-placeholder":"结束日期",clearable:""},model:{value:t.filters.DateRange,callback:function(e){t.$set(t.filters,"DateRange",e)},expression:"filters.DateRange"}})],1),a("el-form-item",[a("el-button",{attrs:{type:"primary",icon:"el-icon-refresh"},on:{click:t.handleQuery}},[t._v("查询")])],1)],1),a("el-table",{staticClass:"table",staticStyle:{"margin-top":"-13px"},attrs:{data:t.table_data,border:"","max-height":"670px","cell-style":{padding:"3px"},"highlight-current-row":!0}},[a("el-table-column",{attrs:{align:"center",prop:"id",label:"序号",width:"80"}}),a("el-table-column",{attrs:{align:"center",prop:"account",label:"管理员",width:"100"}}),a("el-table-column",{attrs:{align:"center",prop:"opt_name",label:"操作",width:"200"}}),a("el-table-column",{attrs:{align:"center",prop:"req_ip",label:"ip",width:"130"}}),a("el-table-column",{attrs:{align:"center",prop:"create_time",label:"时间",width:"200"},scopedSlots:t._u([{key:"default",fn:function(e){return[a("span",[t._v(t._s(t._f("北京时间")(e.row.create_time)))])]}}])}),a("el-table-column",{attrs:{label:"内容","show-overflow-tooltip":""},scopedSlots:t._u([{key:"default",fn:function(e){return[a("span",{staticStyle:{cursor:"pointer"},attrs:{type:"text"},on:{click:function(a){return t.copy(e.row.req_data)}}},[t._v(" "+t._s(e.row.req_data))])]}}])})],1),a("div",{staticClass:"pagination"},[a("el-pagination",{staticStyle:{"margin-right":"10px"},attrs:{background:"",layout:"sizes,total, prev, pager, next, jumper","current-page":t.page,"page-sizes":t.page_sizes,total:t.total,"page-size":t.page_size},on:{"size-change":t.page_sizeChange,"current-change":t.pageChange}})],1)],1)},n=[],l=(a("a9e3"),a("7e79"),a("a417")),o=a("c1df"),s=a.n(o),r={extends:l["a"],data:function(){return{filters:{DateRange:[s()().format("YYYY-MM-DD"),s()().format("YYYY-MM-DD")]}}},created:function(){this.getTableData()},methods:{getTableData:function(){var t=this,e=this.getQueryData();e.channel_id=Number(e.channel_id),this.$post("/v1/admin_get_opt_log",e).then((function(e){t.table_data=e.data,t.total=e.total}))}}},c=r,d=a("2877"),u=Object(d["a"])(c,i,n,!1,null,null,null);e["default"]=u.exports},"7e79":function(t,e,a){},a417:function(t,e,a){"use strict";var i=a("5530"),n=(a("e9c4"),a("b64b"),a("d3b7"),a("159b"),a("5c96")),l=(a("0c6d"),a("2f62")),o=(a("73f5"),a("c1df")),s=a.n(o),r=a("f62d"),c=a("b76c");e["a"]={extends:c["a"],data:function(){return{filters:{},page:1,page_size:15,page_sizes:[15,50,100,200,500,1e3,1500],total:0,table_data:[],itemdata:null,Export:0,columnpicker:!1,column:{},dialog0:{show:!1,title:"",itemdata:{}},dialog1:{show:!1,title:"",itemdata:{}},dialog2:{show:!1,title:"",itemdata:{}},dialog3:{show:!1,title:"",itemdata:{}},dialog4:{show:!1,title:"",itemdata:{}},dialog5:{show:!1,title:"",itemdata:{}},dialog6:{show:!1,title:"",itemdata:{}},dialog7:{show:!1,title:"",itemdata:{}},dialog8:{show:!1,title:"",itemdata:{}},dialog9:{show:!1,title:"",itemdata:{}},symbols:[{id:"trx",name:"TRX"},{id:"usdt",name:"USDT"}]}},computed:Object(i["a"])({},Object(l["b"])(["sellers","channels","games","rooms","states"])),created:function(){if(this.columns){var t="column_".concat(this.$vnode.key),e=localStorage.getItem(t);if(e)this.column=JSON.parse(e);else{for(var a=0;a<this.columns.length;a++)this.column[this.columns[a]]=!0;localStorage.setItem(t,JSON.stringify(this.column))}}},methods:{pageChange:function(t){this.page=t,this.getTableData()},page_sizeChange:function(t){this.page=1,this.page_size=t,this.getTableData()},handleQuery:function(t){var e=this;1===t?(this.Export=1,this.getTableData((function(t){e.Export=0,window.open("/sapi"+t.filename)}))):(this.Export=0,this.page=1,this.getTableData())},handleAdd:function(t){var e=this;this.AddItem(t,(function(a,i){i=i||{},e["dialog".concat(t)].itemdata=i,e["dialog".concat(t)].title=a,e["dialog".concat(t)].show=!0}))},handleEdit:function(t,e){var a=this;this["dialog".concat(e)].itemdata=JSON.parse(JSON.stringify(t)),this.ModifyItem(e,(function(t){a["dialog".concat(e)].title=t,a["dialog".concat(e)].show=!0}),this["dialog".concat(e)].itemdata)},handleDelete:function(t,e){this.DeleteItem(t,e)},getQueryData:function(){var t={};t.page=this.page,t.page_size=this.page_size;var e=JSON.parse(JSON.stringify(this.filters));for(var a in e.channel_id=e.channel_id||0,e)null!=e[a]&&("DateRange"==a&&Object(r["a"])(e[a])?(t["start_time"]=s()(e[a][0]).format("YYYY-MM-DD HH:mm:ss"),t["end_time"]=s()(s()(e[a][1]).valueOf()+864e5).format("YYYY-MM-DD HH:mm:ss")):"DateTimeRange"==a&&Object(r["a"])(e[a])?(t["start_time"]=s()(e[a][0]).format("YYYY-MM-DD HH:mm:ss"),t["end_time"]=s()(s()(e[a][1]).valueOf()+1e3).format("YYYY-MM-DD HH:mm:ss")):(t[a]=e[a],"symbol"==a.toLowerCase()&&(t[a]=e[a].toLowerCase())));return 1==this.Export&&(t.Export=this.Export),t},copy:function(t){var e=document.createElement("textarea");e.value="".concat(t),document.body.appendChild(e),e.select(),document.execCommand("copy"),e.remove(),Object(n["Message"])({type:"success",message:"复制成功",center:!0})},getSummaries:function(t){var e=this,a=t.columns,i=[];return a.forEach((function(t,a){if(0!=a)e.getTotal&&(i[a]=e.getTotal(t.property));else{if(e.getTotal){var n=e.getTotal(0);if(n)return void(i[a]=n)}i[a]="总计"}})),i},selectColumn:function(){this.columnpicker=!0},setColumn:function(t){for(var e in t)this.column[e]=t[e];var a="column_".concat(this.$vnode.key);localStorage.setItem(a,JSON.stringify(this.column))}}}},b76c:function(t,e,a){"use strict";e["a"]={data:function(){return{MapYesNo:{1:"是",2:"否"},ListYesNo:[{id:1,name:"是"},{id:2,name:"否"}]}}}}}]);