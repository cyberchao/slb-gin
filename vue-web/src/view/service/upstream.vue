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
        <el-form-item label="关键字">
          <el-input placeholder="请输入" v-model="searchInfo.name"></el-input>
        </el-form-item>

        <el-form-item>
          <el-button @click="onSubmit" type="primary">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="addupstream()" type="primary"
            >新增upstream</el-button
          >
        </el-form-item>
      </el-form>
    </div>

    <el-table :data="tableData" border stripe>
      <el-table-column type="expand">
        <template slot-scope="props">
          <el-table
            ref="filterTable"
            :data="props.row.serverlist"
            style="width: 100%"
          >
            <el-table-column prop="ip" label="IP"> </el-table-column>
            <el-table-column prop="port" label="http端口"> </el-table-column>
            <el-table-column prop="weight" label="权重"> </el-table-column>
            <el-table-column
              prop="status"
              label="状态"
              :filters="[
                { text: 'up', value: 'up' },
                { text: 'down', value: 'down' },
              ]"
              :filter-method="filterstatus"
              filter-placement="bottom-end"
            >
              <template slot-scope="scope">
                <el-tag
                  :type="scope.row.status === 'up' ? 'success' : 'warning'"
                  disable-transitions
                  >{{ scope.row.status }}</el-tag
                >
              </template>
            </el-table-column>
          </el-table>
        </template>
      </el-table-column>
      <el-table-column
        label="Name"
        min-width="50"
        prop="name"
      ></el-table-column>
      <el-table-column
        label="集群"
        min-width="50"
        prop="cluster"
      ></el-table-column>
      <el-table-column label="环境" min-width="50" prop="env"></el-table-column>
      <el-table-column
        label="转发策略"
        min-width="150"
        prop="forward"
      ></el-table-column>

      <el-table-column fixed="right" label="操作" width="300">
        <template slot-scope="scope">
          <el-button
            @click="editUpstream(scope.row)"
            size="small"
            type="primary"
            icon="el-icon-bank-card"
            >修改</el-button
          >
          <el-button
            @click="deleteUpstream(scope.row)"
            size="small"
            type="danger"
            icon="el-icon-delete"
            >删除</el-button
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
    <el-dialog :title="dialogaddTitle" :visible.sync="dialogaddFormVisible">
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
        <el-form-item label="name" prop="name">
          <el-input placeholder="name" v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item v-for="server in serverList" :key="server.key">
          <span :inline="true" class="demo-form-inline">
            <el-input placeholder="ip" prop="server_ip" v-model="server.ip">
            </el-input>
            <el-input
              placeholder="port"
              prop="server_port"
              v-model="server.port"
            >
            </el-input>
            <el-input
              placeholder="权重"
              prop="server_weight"
              v-model="server.weight"
            >
            </el-input>
            <el-select v-model="server.status" placeholder="请选择">
              <el-option
                v-for="item in status"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </span>
          <el-button @click.prevent="removeServer(server)"
            >删除server</el-button
          >
        </el-form-item>

        <el-form-item>
          <el-button @click="addServer">添加server</el-button>
        </el-form-item>
        <el-form-item label="转发策略" prop="forward">
          <el-select v-model="form.forward" placeholder="选择">
            <el-option
              v-for="item in forwards"
              :key="item.value"
              :label="item.value"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div class="dialog-footer" slot="footer">
        <el-button @click="closeaddDialog">取 消</el-button>
        <el-button @click="save('form')" type="primary">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="dialogeditTitle" :visible.sync="dialogeditFormVisible">
      <el-form
        :inline="true"
        ref="form"
        :rules="rules"
        :model="form"
        class="demo-form-inline"
      >
        <el-form-item v-for="server in servereditList" :key="server.key">
          <span :inline="true" class="demo-form-inline">
            <el-input placeholder="ip" prop="server_ip" v-model="server.ip">
            </el-input>
            <el-input
              placeholder="port"
              prop="server_port"
              v-model="server.port"
            >
            </el-input>
            <el-input
              placeholder="权重"
              prop="server_weight"
              v-model="server.weight"
            >
            </el-input>
            <el-select v-model="server.status" placeholder="请选择">
              <el-option
                v-for="item in status"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </span>
          <el-button @click.prevent="removeServer(server)"
            >删除server</el-button
          >
        </el-form-item>

        <el-form-item>
          <el-button @click="addupdateServer">添加server</el-button>
        </el-form-item>
        <el-form-item label="转发策略" prop="forward">
          <el-select v-model="form.forward" placeholder="选择">
            <el-option
              v-for="item in forwards"
              :key="item.value"
              :label="item.value"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div class="dialog-footer" slot="footer">
        <el-button @click="closeeditDialog">取 消</el-button>
        <el-button @click="update()" type="primary">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
// 获取列表内容封装在mixins内部  getTableData方法 初始化已封装完成 条件搜索时候 请把条件安好后台定制的结构体字段 放到 this.searchInfo 中即可实现条件搜索

import {
  getUpstreamList,
  createUpstream,
  updateUpstream,
  deleteUpstream,
} from "@/api/upstream";
import infoList from "@/mixins/infoList";
const forwards = [{ value: "chash" }, { value: "roundbin" }];
const status = [{ value: "up" }, { value: "down" }];

export default {
  name: "Upstream",
  mixins: [infoList],
  data() {
    return {
      page: 1,
      total: 10,
      pageSize: 10,
      tableData: [],
      searchInfo: {},
      lineNumbers: true,
      dialogaddFormVisible: false,
      dialogaddTitle: "新增Upstream",
      dialogeditFormVisible: false,
      dialogeditTitle: "修改Upstream",
      listApi: getUpstreamList,
      serverList: [],
      servereditList: [],
      form: {
        id: "",
        env: "",
        cluster: "",
        name: "",
        forward: "",
      },
      updateform: {},
      nginx_clusters: this.GLOBAL.nginx_clusters,
      envs: this.GLOBAL.envs,
      status: status,
      forwards: forwards,
      rules: {
        cluster: [
          { required: true, message: "请选择nginx集群", trigger: "change" },
        ],
        env: [{ required: true, message: "请选择环境", trigger: "change" }],
        name: [{ required: true, message: "请填写name", trigger: "change" }],

        serverList: [
          { required: true, message: "请填写server列表", trigger: "change" },
        ],

        server_ip: [
          { required: true, message: "请填写server ip", trigger: "change" },
        ],

        server_weight: [
          { required: true, message: "请填写server权重", trigger: "change" },
        ],

        server_status: [
          { required: true, message: "请选择状态", trigger: "change" },
        ],

        forward: [
          { required: true, message: "请填写转发策略", trigger: "change" },
        ],
      },
    };
  },
  methods: {
    addServer() {
      this.serverList.push({
        ip: "",
        weight: "",
        status: "",
      });
    },
    addupdateServer() {
      this.servereditList.push({
        ip: "",
        weight: "",
        status: "",
      });
    },
    removeServer(item) {
      var index = this.serverList.indexOf(item);
      if (index !== -1) {
        this.serverList.splice(index, 1);
      }
    },
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
    closeaddDialog() {
      this.dialogaddFormVisible = false;
    },
    addupstream() {
      this.dialogaddFormVisible = true;
    },
    openaddDialog(row) {
      this.dialogaddTitlethis = "新增Upstream";
      this.code = row.src;
      this.dialogaddFormVisible = true;
    },
    closeeditDialog() {
      this.dialogeditFormVisible = false;
    },
    editUpstream(row) {
      this.servereditList = [];
      this.updateform.id = row._id;
      console.log("edit upstream:", row);
      for (const server of row.serverlist) {
        this.servereditList.push(server);
      }
      this.dialogeditFormVisible = true;
    },
    openeditDialog() {
      this.dialogeditFormVisible = true;
    },
    async deleteUpstream(row) {
      console.log("row:", row);
      this.$confirm("确认是否删除该server", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(async () => {
          const res = await deleteUpstream({
            id: row._id,
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
    save(formName) {
      this.$refs[formName].validate(async (valid) => {
        if (valid) {
          this.form["serverList"] = this.serverList;
          var postdata = this.form;
          console.log("add data:", postdata);
          const res = await createUpstream(postdata);
          if (res.code == 0) {
            this.$message({
              type: "success",
              message: "添加成功",
              showClose: true,
            });
          }
          this.getTableData();
          this.closeaddDialog();
        }
      });
    },
    async update() {
      this.updateform.form = this.servereditList;
      console.log(this.updateform);
      const res = await updateUpstream(this.updateform);
      console.log("res", res);
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: "修改成功",
          showClose: true,
        });
      }
      this.getTableData();
      this.closeeditDialog();
    },
    filterstatus(value, row) {
      console.log(row, value);
      return row.status === value;
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
