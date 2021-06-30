<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="集群">
          <el-select v-model="searchInfo.cluster" multiple placeholder="选择">
            <el-option
              v-for="item in nginx_clusters"
              :key="item.value"
              :label="item.value"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="环境">
          <el-select v-model="searchInfo.env" multiple placeholder="选择">
            <el-option
              v-for="item in envs"
              :key="item.value"
              :label="item.value"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="ip">
          <el-input placeholder="ip" v-model="searchInfo.ip"></el-input>
        </el-form-item>

        <el-form-item>
          <el-button @click="onSubmit" type="primary">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="addhost()" type="primary">新增</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="nginx_t()" type="primary">Check</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="nginx_s()" type="primary">Reload</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      :data="tableData"
      border
      stripe
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55"> </el-table-column>
      <el-table-column
        label="集群"
        min-width="50"
        prop="cluster"
      ></el-table-column>
      <el-table-column label="环境" min-width="50" prop="env"></el-table-column>
      <el-table-column label="IP" min-width="50" prop="ip"></el-table-column>
      <el-table-column
        label="reload时间"
        min-width="150"
        prop="reload_time"
      ></el-table-column>
      <el-table-column
        label="描述"
        min-width="150"
        prop="description"
      ></el-table-column>

      <el-table-column fixed="right" label="操作" width="300">
        <template slot-scope="scope">
          <el-button
            @click="deleteHost(scope.row)"
            size="small"
            type="danger"
            icon="el-icon-delete"
          ></el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      :style="{ float: 'right', padding: '20px' }"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
      layout="total, sizes, prev, pager, next, jumper"
    ></el-pagination>
    <el-dialog :title="dialogTitle" :visible.sync="dialogFormVisible">
      <el-form ref="form" :rules="rules" :model="form" class="demo-form-inline">
        <el-form-item label="集群" prop="cluster">
          <el-select v-model="form.cluster" placeholder="选择">
            <el-option
              v-for="item in nginx_clusters"
              :key="item.value"
              :label="item.value"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="环境" prop="env">
          <el-select v-model="form.env" placeholder="选择">
            <el-option
              v-for="item in envs"
              :key="item.value"
              :label="item.value"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="ip" prop="ip">
          <el-input placeholder="ip" v-model="form.ip"></el-input>
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            placeholder="description"
            v-model="form.description"
          ></el-input>
        </el-form-item>
      </el-form>
      <div class="dialog-footer" slot="footer">
        <el-button @click="closeDialog">取 消</el-button>
        <el-button @click="save('form')" type="primary">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
// 获取列表内容封装在mixins内部  getTableData方法 初始化已封装完成 条件搜索时候 请把条件安好后台定制的结构体字段 放到 this.searchInfo 中即可实现条件搜索

import {
  getHostList,
  createHost,
  checkHost,
  reloadHost,
  deleteHost,
} from "@/api/host";
import infoList from "@/mixins/infoList";

export default {
  name: "Host",
  mixins: [infoList],
  data() {
    return {
      page: 1,
      total: 10,
      pageSize: 10,
      tableData: [],
      searchInfo: {},
      lineNumbers: true,
      dialogFormVisible: false,
      dialogTitle: "新增Host",
      listApi: getHostList,
      form: {
        env: "",
        cluster: "",
        description: "",
        ip: "",
      },
      multipleSelection: [],
      hostList: [],
      nginx_clusters: this.GLOBAL.nginx_clusters,
      envs: this.GLOBAL.envs,
      rules: {
        cluster: [
          { required: true, message: "请选择nginx集群", trigger: "change" },
        ],
        env: [{ required: true, message: "请选择环境", trigger: "change" }],
        ip: [{ required: true, message: "请填写ip", trigger: "change" }],
      },
    };
  },
  methods: {
    handleSizeChange(val) {
      this.pageSize = val;
      this.getTableData();
    },
    handleCurrentChange(val) {
      this.page = val;
      this.getTableData();
    },
    //条件搜索前端看此方法
    onSubmit() {
      this.page = 1;
      this.pageSize = 10;
      this.getTableData();
    },
    closeDialog() {
      this.dialogFormVisible = false;
    },
    addhost() {
      this.dialogFormVisible = true;
    },
    checkhost(row) {
      checkHost({ id: row.id }).then((res) => {
        if (res.code == 0) {
          this.$message({
            type: "success",
            message: "保存成功!",
          });
        }
      });
    },
    reloadhost(row) {
      reloadHost({ id: row.id }).then((res) => {
        if (res.code == 0) {
          this.$message({
            type: "success",
            message: "重载成功!",
          });
        }
      });
    },
    openDialog(row) {
      this.dialogTitlethis = "新增Host";
      this.code = row.src;
      this.dialogFormVisible = true;
    },
    async deleteHost(row) {
      this.$confirm("确认是否删除该server", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(async () => {
          const res = await deleteHost({ id: row.id });
          if (res.code == 0) {
            this.$message({
              type: "success",
              message: "删除成功!",
            });
            if (this.tableData.length == 1) {
              this.page--;
            }
            this.getTableData();
          }
        })
        .catch(() => {
          this.$message({
            type: "info",
            message: "已取消删除",
          });
        });
    },
    save(formName) {
      this.$refs[formName].validate(async (valid) => {
        if (valid) {
          const res = await createHost(this.form);
          if (res.code == 0) {
            this.$message({
              type: "success",
              message: "添加成功",
              showClose: true,
            });
          }
          this.getTableData();
          this.closeDialog();
        }
      });
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
    async nginx_t() {
      this.hostList = [];
      for (const hostinfo of this.multipleSelection) {
        this.hostList.push(hostinfo.ip);
      }
      console.log(this.hostList);
      const res = await checkHost(this.hostList);
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: "Check通过",
          showClose: true,
        });
      }
    },
    async nginx_s() {
      this.hostList = [];
      for (const hostinfo of this.multipleSelection) {
        this.hostList.push(hostinfo.ip);
      }
      const res = await reloadHost(this.hostList);
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: "Reload成功",
          showClose: true,
        });
      }
    },
  },
  created() {
    this.getTableData();
  },
};
</script>
<style scoped lang="scss">
.button-box {
  padding: 10px 20px;
  .el-button {
    float: right;
  }
}
.el-tag--mini {
  margin-left: 5px;
}
.warning {
  color: #dc143c;
}
.my-editor {
  background: #2d2d2d;
  color: #ccc;

  font-family: Fira code, Fira Mono, Consolas, Menlo, Courier, monospace;
  font-size: 14px;
  line-height: 1.5;
  padding: 5px;
}

// optional
.prism-editor__textarea:focus {
  outline: none;
}

// not required:
.height-300 {
  height: 60%;
}
</style>
