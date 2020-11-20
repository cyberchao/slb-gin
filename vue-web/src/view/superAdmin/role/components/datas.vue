<template>
  <div>
    <div class="clearflex" style="margin:18px">
      <el-button
        @click="authDataEnter"
        class="fl-right"
        size="small"
        type="primary"
        >确 定</el-button
      >
      <el-button @click="all" class="fl-left" size="small" type="primary"
        >全选</el-button
      >
      <el-button @click="self" class="fl-left" size="small" type="primary"
        >本角色</el-button
      >
      <el-button
        @click="selfAndChildren"
        class="fl-left"
        size="small"
        type="primary"
        >本角色及子角色</el-button
      >
    </div>
    <el-checkbox-group v-model="dataRoleId" @change="selectRole">
      <el-checkbox v-for="(item, key) in roles" :label="item" :key="key">{{
        item.roleName
      }}</el-checkbox>
    </el-checkbox-group>
  </div>
</template>
<script>
import { setDataRole } from "@/api/role";
export default {
  name: "Datas",
  data() {
    return {
      roles: [],
      dataRoleId: [],
      needConfirm: false,
    };
  },
  props: {
    row: {
      default: function() {
        return {};
      },
      type: Object,
    },
    role: {
      default: function() {
        return {};
      },
      type: Array,
    },
  },
  methods: {
    // 暴露给外层使用的切换拦截统一方法
    enterAndNext() {
      this.authDataEnter();
    },
    all() {
      this.dataRoleId = [...this.roles];
      this.row.dataRoleId = this.dataRoleId;
      this.needConfirm = true;
    },
    self() {
      this.dataRoleId = this.roles.filter(
        (item) => item.roleId === this.row.roleId
      );
      this.row.dataRoleId = this.dataRoleId;
      this.needConfirm = true;
    },
    selfAndChildren() {
      const arrBox = [];
      this.getChildrenId(this.row, arrBox);
      this.dataRoleId = this.roles.filter(
        (item) => arrBox.indexOf(item.roleId) > -1
      );
      this.row.dataRoleId = this.dataRoleId;
      this.needConfirm = true;
    },
    getChildrenId(row, arrBox) {
      arrBox.push(row.roleId);
      row.children &&
        row.children.map((item) => {
          this.getChildrenId(item, arrBox);
        });
    },
    // 提交
    async authDataEnter() {
      const res = await setDataRole(this.row);
      if (res.code == 0) {
        this.$message({ type: "success", message: "资源设置成功" });
      }
    },
    //   平铺角色
    roundRole(roles) {
      roles &&
        roles.map((item) => {
          const obj = {};
          obj.roleId = item.roleId;
          obj.roleName = item.roleName;
          this.roles.push(obj);
          if (item.children && item.children.length) {
            this.roundRole(item.children);
          }
        });
    },
    //   选择
    selectRole() {
      this.row.dataRoleId = this.dataRoleId;
      this.needConfirm = true;
    },
  },
  created() {
    this.roles = [];
    this.dataRoleId = [];
    this.roundRole(this.role);
    this.row.dataRoleId &&
      this.row.dataRoleId.map((item) => {
        const obj =
          this.roles &&
          this.roles.filter((au) => au.roleId === item.roleId) &&
          this.roles.filter((au) => au.roleId === item.roleId)[0];
        this.dataRoleId.push(obj);
      });
  },
};
</script>
<style lang="less"></style>
