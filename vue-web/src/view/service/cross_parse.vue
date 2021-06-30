<template>
  <div>
    <el-input
      type="textarea"
      :rows="15"
      placeholder="请输入内容"
      v-model="textarea"
    >
    </el-input>
    <el-button @click="postdata">提交</el-button>
    <el-input type="textarea" :rows="15" placeholder="" v-model="parsedata">
    </el-input>
  </div>
</template>

<script>
import { cross_parse } from "@/api/crossplane";
export default {
  data() {
    return {
      code: "",
      textarea: "",
      parsedata: "",
    };
  },
  methods: {
    postdata() {
      cross_parse(this.textarea).then((res) => {
        this.code = res.code;
        this.parsedata = JSON.stringify(res.data, null, 4);
      });
    },
  },
};
</script>
