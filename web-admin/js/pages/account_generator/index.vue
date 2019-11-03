<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    el-row
      el-col(:xs="24" :md="12" class="content__main__item")
        el-card
          h2(class="content__generator__title") 比赛用户账号生成器
          el-form(:model="leftForm", :rules="leftRules", ref="leftForm", label-width="60px")
            el-form-item(label="前缀",prop="prefix")
              el-input(v-model="leftForm.prefix")
            el-form-item(label="数量",prop="number")
              el-input(v-model.number="leftForm.number")
            el-form-item
              el-button(type="success",@click="submitCompeteAccount") 生成
      el-col(:xs="24" :md="12" class="content__main__item")
        el-card
          h2(class="content__generator__title") 用户账号生成器
          el-form(:model="rightForm" :rules="rightRules",ref="rightForm",label-width="60px")
            el-form-item
              el-input(v-model="rightForm.userList",type="textarea",:rows="12",resize="none")
            el-form-item
              el-button(type="success",@click="submitUserAccount") 生成
            el-alert(title="可以将学生学号从Excel整列复制过来，批量生成用户账号，初始密码默认为123456。" ,type="info", :closable="false")
  el-dialog(title="账号列表",:visible.sync="dialogOperatorInfoVisible",:close-on-click-modal="false",width="600px")
    el-row
      el-col(el-col :xs="24" :md="8" class="infolist")
        p(v-for="(item,index) in info" :key="index" class="dialog__info") {{item}}
      el-col(el-col :xs="24" :md="16")
        el-table(:data="tableData")
          el-table-column(property="username",label="用户名",width="150")
          el-table-column(property="password",label="密码")
    .dialog-footer(slot="footer")
      el-button(@click="dialogOperatorInfoVisible = false") 取消
      el-button(type="primary",@click="dialogOperatorInfoVisible = false") 确定
  </template>
<script>
import {
  generateCompeteAccount,
  generateUserAccount
} from '@/web-admin/js/api/generator.js';

export default {
  data() {
    return {
      dialogOperatorInfoVisible: false,
      info: '',
      leftForm: {
        prefix: '',
        number: ''
      },
      leftRules: {
        prefix: [
          {
            required: true,
            message: '请输入前缀',
            trigger: 'blur'
          }
        ],
        number: [
          {
            required: true,
            message: '请输入数量',
            trigger: 'blur'
          },
          {
            type: 'integer',
            min: 0,
            max: 100,
            message: '请输入1-100之间的整数',
            trigger: 'blur'
          }
        ]
      },
      rightForm: {
        userList: ''
      },
      rightRules: {
        userList: [
          {
            required: true,
            message: '请输入用户名列表',
            trigger: 'blur'
          }
        ]
      },
      tableData: []
    };
  },
  mounted() {},
  methods: {
    submitCompeteAccount() {
      const self = this;
      self.info = '';
      self.$refs['leftForm'].validate(async valid => {
        if (valid) {
          try {
            let res = await generateCompeteAccount(self.leftForm);
            self.dialogOperatorInfoVisible = true;
            self.tableData = res.data.users;
            self.info = res.data.info;
            console.log(self.info);
            self.$message({
              message: res.data.message,
              type: 'success'
            });
          } catch (err) {
            console.log(err);
            self.$message({
              message: err.response.data.message,
              type: 'error'
            });
          }
        } else {
          return false;
        }
      });
    },
    submitUserAccount() {
      const self = this;
      self.info = '';
      self.$refs['rightForm'].validate(async valid => {
        console.log(self.rightForm.userList);
        if (valid) {
          try {
            let res = await generateUserAccount(self.rightForm);
            self.dialogOperatorInfoVisible = true;
            self.tableData = res.data.users;
            self.info = res.data.info;
            console.log(self.info);
            self.$message({
              message: res.data.message,
              type: 'success'
            });
          } catch (err) {
            console.log(err);
            self.$message({
              message: err.response.data.message,
              type: 'error'
            });
          }
        } else {
          return false;
        }
      });
    }
  }
};
</script>

<style lang="scss" scoped>
.content__generator__title {
  margin-bottom: 1em;
  color: $--color-level3;
}
.content_generator_card {
  height: 400px;
}
.content__main__item {
  height: 40em;
  padding: 2em;
  .el-card {
    height: 100%;
  }
}
.dialog__info {
  font-size: 16px;
  text-align: left;
  padding: 0.5em;
}
</style>