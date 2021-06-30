<template>
  <div id="app">
    <el-button-group>
      <el-button type="primary" icon="el-icon-view">检查</el-button>
      <el-button type="primary" icon="el-icon-paperclip">保存</el-button>
      <el-button type="primary" icon="el-icon-upload">发布</el-button>
    </el-button-group>
    <prism-editor
      class="my-editor height-300"
      v-model="code"
      :highlight="highlighter"
      :line-numbers="lineNumbers"
    ></prism-editor>
  </div>
</template>

<script>
import { PrismEditor } from "vue-prism-editor";
import "vue-prism-editor/dist/prismeditor.min.css"; // import the styles somewhere

// import highlighting library (you can use any library you want just return html string)
import { highlight, languages } from "prismjs/components/prism-core";
import "prismjs/components/prism-clike";
import "prismjs/components/prism-javascript";
import "prismjs/themes/prism-coy.css"; // import syntax highlighting styles

export default {
  components: {
    PrismEditor,
  },
  data: () => ({
    code: "",
    lineNumbers: true,
  }),
  methods: {
    highlighter(code) {
      return highlight(code, languages.js); //returns html
    },
  },
};
</script>

<style lang="scss">
// required class
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
