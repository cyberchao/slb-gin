<template>
  <div>
    <el-table :data="tableData" border stripe>
      <el-table-column label="uuid" min-width="250" prop="uuid"></el-table-column>
      <el-table-column label="用户名" min-width="150" prop="userName"></el-table-column>
      <el-table-column label="昵称" min-width="150" prop="nickName"></el-table-column>
      <el-table-column label="用户角色" min-width="150">
        <template slot-scope="scope">
          <el-cascader
            @change="changerole(scope.row)"
            v-model="scope.row.role.roleId"
            :options="authOptions"
            :show-all-levels="false"
            :props="{ checkStrictly: true,label:'roleName',value:'roleId',disabled:'disabled',emitPath:false}"
            filterable
          ></el-cascader>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      :style="{float:'right',padding:'20px'}"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
      layout="total, sizes, prev, pager, next, jumper"
    ></el-pagination>
  </div>
</template>


<script>
// 获取列表内容封装在mixins内部  getTableData方法 初始化已封装完成
const path = process.env.VUE_APP_BASE_API;
import {
  getUserList,
  setUserRole,
} from "@/api/user";
import { getRoleList } from "@/api/role";
import infoList from "@/mixins/infoList";
import { mapGetters } from "vuex";
export default {
  name: "Api",
  mixins: [infoList],
  data() {
    return {
      listApi: getUserList,
      path: path,
      authOptions: [],
      addUserDialog: false,
      userInfo: {
        username: "",
        password: "",
        nickName: "",
        roleId: ""
      },
    };
  },
  computed: {
    ...mapGetters("user", ["token"])
  },
  methods: {
    setOptions(authData) {
      this.authOptions = [];
      this.setroleOptions(authData, this.authOptions);
    },
    setroleOptions(roleData, optionsData) {
      roleData &&
        roleData.map(item => {
          if (item.children && item.children.length) {
            const option = {
              roleId: item.roleId,
              roleName: item.roleName,
              children: []
            };
            this.setroleOptions(item.children, option.children);
            optionsData.push(option);
          } else {
            const option = {
              roleId: item.roleId,
              roleName: item.roleName
            };
            optionsData.push(option);
          }
        });
    },
    async changeRole(row) {
      const res = await setUserRole({
        uuid: row.uuid,
        roleId: row.role.roleId
      });
      if (res.code == 0) {
        this.$message({ type: "success", message: "角色设置成功" });
      }
    }
  },
  async created() {
    this.getTableData();
    const res = await getRoleList({ page: 1, pageSize: 999 });
    this.setOptions(res.data.list);
  }
};
</script>
<style lang="scss">

.button-box {
  padding: 10px 20px;
  .el-button {
    float: right;
  }
}

.user-dialog {
  .header-img-box {
  width: 200px;
  height: 200px;
  border: 1px dashed #ccc;
  border-radius: 20px;
  text-align: center;
  line-height: 200px;
  cursor: pointer;
}
}
</style>