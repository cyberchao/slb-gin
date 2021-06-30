<template>
  <div>
    <el-row :gutter="12">
      <el-col :span="12">
        <el-card shadow="hover">
          <div slot="header" class="clearfix">
            <span>填写信息</span>
          </div>
          <el-col :offset="2" :span="16">
            <el-form
              ref="form"
              :rules="rules"
              :model="form"
              label-width="120px"
            >
              <el-form-item label="负载类型">
                <el-radio-group v-model="radio">
                  <el-radio :label="7">七层负载</el-radio>
                  <el-radio :label="4">四层负载</el-radio>
                </el-radio-group>
              </el-form-item>
              <el-form-item label="集群" prop="cluster">
                <el-select placeholder="请选择" v-model="form.cluster">
                  <el-option
                    :key="item.value"
                    :label="item.value"
                    :value="item.value"
                    v-for="item in nginx_clusters"
                  ></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="环境" prop="env">
                <el-select placeholder="请选择" v-model="form.env">
                  <el-option
                    :key="item.value"
                    :label="item.value"
                    :value="item.value"
                    v-for="item in envs"
                  ></el-option>
                </el-select>
              </el-form-item>

              <el-form-item label="域名" prop="server_name">
                <el-input v-model="form.server_name"></el-input>
                <el-tooltip
                  class="item"
                  effect="dark"
                  content="Left Center 提示文字"
                  placement="left"
                >
                  <el-button>左边</el-button>
                </el-tooltip>
              </el-form-item>
              <el-form-item label="端口" prop="listen">
                <el-input v-model="form.listen"></el-input>
              </el-form-item>
              <el-form-item label="ssl">
                <el-switch v-model="form.ssl"></el-switch>
              </el-form-item>
              <!-- <el-form-item label="include">
                <el-input style="width: 80%; margin-right: 10px;"></el-input
                ><el-button @click.prevent="removeDomain(domain)"
                  >删除</el-button
                >
              </el-form-item> -->

              <!-- <el-button icon="el-icon-circle-plus-outline" circle></el-button> -->

              <el-form-item
                v-for="(loc, index) in locs"
                :label="'location-' + index"
                :key="loc.key"
              >
                <el-card shadow="never">
                  <div class="demo-input-suffix">
                    <el-input placeholder="匹配策略" v-model="loc.rule">
                    </el-input>
                    <el-input placeholder="转发规则" v-model="loc.path">
                    </el-input>
                    <el-input
                      type="textarea"
                      :autosize="{ minRows: 2 }"
                      placeholder="请输入内容"
                      v-model="loc.other"
                    >
                    </el-input>
                  </div>
                </el-card>
                <el-button @click.prevent="removeLoc(loc)">删除</el-button>
              </el-form-item>

              <el-form-item>
                <el-button @click="addLoc">添加策略</el-button>
              </el-form-item>
              <el-form-item label="其它">
                <el-input
                  type="textarea"
                  :autosize="{ minRows: 2 }"
                  placeholder="请输入内容"
                  v-model="form.other"
                >
                </el-input>
              </el-form-item>
              <el-form-item label="描述">
                <el-input v-model="description"></el-input>
              </el-form-item>
            </el-form>
          </el-col>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <div slot="header" class="clearfix">
            <span>数据生成</span>
          </div>
          <div id="app">
            <el-button-group>
              <el-button type="primary" icon="el-icon-view">检查</el-button>
              <el-button
                type="primary"
                icon="el-icon-paperclip"
                @click="save('form')"
                >保存</el-button
              >
              <el-button type="primary" icon="el-icon-upload">发布</el-button>
            </el-button-group>
            <prism-editor
              class="my-editor height-300"
              v-model="text"
              :highlight="highlighter"
              :line-numbers="lineNumbers"
            ></prism-editor>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-tab-pane label="vServer管理" name="first">
      <el-row> </el-row>
    </el-tab-pane>
  </div>
</template>

<script>
import { PrismEditor } from "vue-prism-editor";
import "vue-prism-editor/dist/prismeditor.min.css"; // import the styles somewhere

// import highlighting library (you can use any library you want just return html string)
import { highlight, languages } from "prismjs/components/prism-core";
import "prismjs/components/prism-clike";
import "prismjs/components/prism-javascript";
import "prismjs/themes/prism-funky.css"; // import syntax highlighting styles
import { createServer } from "@/api/server";

const nginx_clusters = [
  {
    value: "nginx-ssl",
    type: "success",
  },
  {
    value: "nginx-bc",
    type: "",
  },
  {
    value: "nginx-cproxy",
    type: "warning",
  },
  {
    value: "pafm-fs",
    type: "danger",
  },
];
const envs = [
  {
    value: "uat",
  },
  {
    value: "hd",
  },
  {
    value: "prd_bx",
  },
  {
    value: "prd_wgq",
  },
];
export default {
  name: "Servercreate",
  components: {
    PrismEditor,
  },
  data() {
    return {
      radio: 7,
      formother: "",
      lineNumbers: true,
      locs: [
        {
          rule: "",
          path: "",
          other: "",
        },
      ],
      form: {
        env: "",
        cluster: "",
        description: "",
        server_name: "",
        listen: "",
        ssl: false,
        other: "",
      },
      nginx_clusters: nginx_clusters,
      envs: envs,
      description: "",
      version: 1,
      text: "",
      // locstring:"",
      rules: {
        cluster: [
          { required: true, message: "请选择nginx集群", trigger: "change" },
        ],
        env: [{ required: true, message: "请选择环境", trigger: "change" }],
        server_name: [
          { required: true, message: "请填写域名", trigger: "change" },
        ],
        listen: [{ required: true, message: "请填写端口", trigger: "change" }],
      },
    };
  },
  watch: {
    locstring: function() {
      var locstr = "";
      for (const loc of this.locs) {
        locstr += `
        location ${loc.rule}  {
            ${loc.path} ${loc.other ? "\n            " + loc.other : ""} ;
        }\n`;
      }
      return locstr;
    },
    code: {
      get: function() {
        let tmp = `server {
        listen      ${this.form.listen} ${this.form.ssl ? "ssl" : ""};
        server_name  ${this.form.server_name};
${this.form.other}
${this.locstring}
        access_log /wls/applogs/nginx/${this.form.server_name}.access.log main;
        error_log /wls/applogs/nginx/${this.form.server_name}.error.log warn;
}`;
        this.text = tmp;
        return tmp;
      },

      set: function(v) {
        console.log(v);
        this.coder = v;
        return v;
      },
    },
  },
  methods: {
    highlighter(code) {
      return highlight(code, languages.js); //returns html
    },
    removeLoc(item) {
      var index = this.locs.indexOf(item);
      if (index !== -1) {
        this.locs.splice(index, 1);
      }
    },
    addLoc() {
      this.locs.push({
        rule: "",
        path: "",
        other: "",
      });
    },
    save(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          var filename = this.form.server_name.split(",")[0] + ".conf";
          var postdata = {
            env: this.form.env,
            cluster: this.form.cluster,
            code: this.code,
            description: this.description,
            verison: this.verison,
            filename: filename,
          };
          console.log(postdata);
          createServer(postdata).then((res) => {
            if (res.code == 0) {
              this.$message({
                type: "success",
                message: "保存成功!",
              });
            }
          });
        } else {
          alert(this.form.env);
          console.log("error submit!!");
          return false;
        }
      });
    },
    // save() {
    //   var postdata = {
    //     env: "",
    //     cluster: "",
    //     code: this.code,
    //     description: this.description,
    //     verison: this.verison,
    //   };
    //   console.log(postdata);
    //   createServer(postdata).then((res) => {
    //     if (res.code == 0) {
    //       this.$message({
    //         type: "success",
    //         message: "保存成功!",
    //       });
    //     }
    //   });
    // },
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
