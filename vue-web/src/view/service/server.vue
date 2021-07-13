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
        <el-form-item label="域名">
          <el-input
            placeholder="server_name"
            v-model="searchInfo.server_name"
          ></el-input>
        </el-form-item>

        <el-form-item>
          <el-button @click="onSubmit" type="primary">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="addvhost()" type="primary">新增</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table :data="serverList" border stripe>
      <span prop="id"></span>
      <el-table-column
        label="域名"
        min-width="100"
        prop="server_name"
      ></el-table-column>
      <el-table-column
        label="集群"
        min-width="50"
        prop="cluster"
      ></el-table-column>
      <el-table-column label="环境" min-width="50" prop="env"></el-table-column>
      <el-table-column
        label="端口"
        min-width="50"
        prop="listen"
      ></el-table-column>
      <el-table-column
        label="上次修改"
        min-width="80"
        prop="time"
      ></el-table-column>
      <el-table-column
        label="备注"
        min-width="50"
        prop="description"
      ></el-table-column>
      <el-table-column label="状态" min-width="50" prop="status">
        <template slot-scope="scope">
          <el-tag
            :type="scope.row.status === '已发布' ? 'success' : 'info'"
            disable-transitions
            >{{ scope.row.status }}</el-tag
          >
        </template></el-table-column
      >

      <el-table-column fixed="right" label="操作" width="300">
        <template slot-scope="scope">
          <el-button
            @click="editServer(scope.row)"
            size="small"
            type="primary"
            icon="el-icon-edit"
            >操作</el-button
          >
          <el-button
            @click="deleteServer(scope.row)"
            size="small"
            type="danger"
            icon="el-icon-delete"
            >删除</el-button
          >
          <el-button
            @click="logServer(scope.row)"
            size="small"
            type="info"
            icon="el-icon-edit"
            plain
            >日志</el-button
          >
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
    <el-dialog
      :before-close="closeDialog"
      :title="dialogTitle"
      :visible.sync="dialogFormVisible"
    >
      <div id="app">
        <el-button-group>
          <el-button type="primary" icon="el-icon-view">检查</el-button>
          <el-button @click="beautify()" type="primary" icon="el-icon-view"
            >格式化</el-button
          >
          <el-button
            type="primary"
            icon="el-icon-paperclip"
            @click="updateServer()"
            >保存</el-button
          >
          <el-button
            @click="publishServer()"
            type="primary"
            icon="el-icon-upload"
            >发布</el-button
          >
        </el-button-group>
        <prism-editor
          class="my-editor height-300"
          v-model="servermodify.newcode"
          :highlight="highlighter"
          :line-numbers="lineNumbers"
        ></prism-editor>
      </div>
      <div class="warning">新增Api需要在角色管理内配置权限才可使用</div>
      <div class="dialog-footer" slot="footer">
        <el-button @click="closeDialog">取 消</el-button>
        <el-button @click="enterDialog" type="primary">确 定</el-button>
      </div>
    </el-dialog>
    <el-drawer
      :visible.sync="drawer"
      :with-header="false"
      size="50%"
      title="日志查看"
      v-if="drawer"
    >
      <el-tabs type="border-card">
        <el-tab-pane label="菜单">
          <span>{{ activeRow }}</span>
        </el-tab-pane>
      </el-tabs>
    </el-drawer>
  </div>
</template>

<script>
// 获取列表内容封装在mixins内部  getTableData方法 初始化已封装完成 条件搜索时候 请把条件安好后台定制的结构体字段 放到 this.searchInfo 中即可实现条件搜索

import {
  getServerList,
  createServer,
  updateServer,
  deleteServer,
  publishServer,
} from "@/api/server";
import router from "@/router/index";
import { PrismEditor } from "vue-prism-editor";
import "vue-prism-editor/dist/prismeditor.min.css"; // import the styles somewhere

// import highlighting library (you can use any library you want just return html string)
import { highlight, languages } from "prismjs/components/prism-core";
import "prismjs/components/prism-clike";
import "prismjs/components/prism-javascript";
import "prismjs/themes/prism-funky.css"; // import syntax highlighting styles

import Beautify from "@/api/nginx";
const instance = new Beautify({ tabs: 1 });

export default {
  name: "Server",
  components: {
    PrismEditor,
  },

  data() {
    return {
      page: 1,
      total: 10,
      pageSize: 10,
      tableData: [],
      searchInfo: {},
      serverList: [],
      lineNumbers: true,
      dialogFormVisible: false,
      dialogTitle: "新增Server",
      drawer: false,
      activeRow: {},
      form: {
        env: "",
        cluster: "",
        description: "",
      },
      servermodify: {},
      publishData: {},
      nginx_clusters: this.GLOBAL.nginx_clusters,
      envs: this.GLOBAL.envs,
      type: "",
    };
  },
  methods: {
    beautify() {
      this.servermodify.newcode = instance.parse(this.servermodify.newcode);
    },
    handleSizeChange(val) {
      this.pageSize = val;
      this.getTableData();
    },
    handleCurrentChange(val) {
      this.page = val;
      this.getTableData();
    },
    async getTableData(page = this.page, pageSize = this.pageSize) {
      this.serverList = [];
      const table = await getServerList({ page, pageSize, ...this.searchInfo });
      if (table.code == 0) {
        console.log("response.data", table.data.list);
        this.tableData = table.data.list;
        this.total = table.data.total;
        this.page = table.data.page;
        this.pageSize = table.data.pageSize;
        for (var server of this.tableData) {
          var serverRaw = {
            id: "",
            env: "",
            cluster: "",
            server_name: "",
            listen: "",
            access_log: "",
            error_log: "",
            src: "",
            description: "",
            status: "",
            time: "",
            filepath: "",
          };
          serverRaw.id = server._id;
          serverRaw.env = server.env;
          serverRaw.cluster = server.cluster;
          serverRaw.description = server.description;
          if (server.status) {
            serverRaw.status = "已发布";
          } else {
            serverRaw.status = "未发布";
          }

          serverRaw.time = server.time.slice(0, 19);
          serverRaw.src = server.src;
          serverRaw.filepath = server.filepath;
          for (var block of server.ngx.block) {
            if (block.directive == "server_name") {
              serverRaw.server_name = block.args[0];
            }
            if (block.directive == "listen") {
              serverRaw.listen = block.args[0];
            }
            if (block.directive == "access_log") {
              serverRaw.access_log = block.args[0];
            }
            if (block.directive == "error_log") {
              serverRaw.error_log = block.args[0];
            }
          }
          console.log("server", server);
          console.log("serverraw:", serverRaw);
          this.serverList.push(serverRaw);
        }
      }
    },
    highlighter(code) {
      return highlight(code, languages.js); //returns html
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
    addvhost() {
      router.push({ path: "server_create" });
      // window.location.href = "/server_create"
    },
    openDialog(row, type) {
      switch (type) {
        case "add":
          this.dialogTitlethis = "新增Server";
          break;
        case "edit":
          this.dialogTitlethis = "编辑Server";
          break;
        default:
          break;
      }
      this.type = type;
      this.code = row.src;
      this.servermodify["id"] = row.id;
      this.servermodify["newcode"] = row.src;
      this.servermodify["env"] = row.env;
      this.servermodify["cluster"] = row.cluster;
      this.servermodify["description"] = row.description;
      this.servermodify["filename"] = row.filepath;
      this.servermodify["version"] = row.version;
      this.publishData["id"] = row.id;
      this.publishData["filepath"] = row.filepath;
      this.publishData["env"] = row.env;
      this.publishData["cluster"] = row.cluster;
      this.dialogFormVisible = true;
    },
    logServer(row) {
      this.drawer = true;
      this.activeRow = row;
    },
    async editServer(row) {
      this.openDialog(row, "edit");
    },
    async deleteServer(row) {
      this.$confirm("确认是否删除该server", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(async () => {
          console.log("row:", row);
          const res = await deleteServer({
            id: row.id,
            filepath: row.filepath,
          });
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
    async updateServer() {
      console.log("更新数据：", this.servermodify);
      const res = await updateServer(this.servermodify);
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: "修改成功!",
        });
      }
    },
    async publishServer() {
      const res = await publishServer(this.publishData);
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: res.msg,
        });
      }
    },

    async enterDialog() {
      this.$refs.apiForm.validate(async (valid) => {
        if (valid) {
          switch (this.type) {
            case "add":
              {
                const res = await createServer(this.form);
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

              break;
            case "edit":
              {
                const res = await updateServer(this.form);
                if (res.code == 0) {
                  this.$message({
                    type: "success",
                    message: "编辑成功",
                    showClose: true,
                  });
                }
                this.getTableData();
                this.closeDialog();
              }
              break;
            default:
              {
                this.$message({
                  type: "error",
                  message: "未知操作",
                  showClose: true,
                });
              }
              break;
          }
        }
      });
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
