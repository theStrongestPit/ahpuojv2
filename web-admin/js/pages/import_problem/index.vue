<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    el-card.box-card.problem-import-card
      el-upload.upload-demo(ref="upload", action="/api/admin/problemset", 
      :headers="{Authorization:$store.getters.token}", 
      accept=".xml", :auto-upload="false", :on-success="handleUploadSuccess", :on-error="handleUploadError", :limit="1")
        el-button(slot="trigger", size="small", type="primary") 选取文件
        el-button(style="margin-left: 10px;",size="small",type="success",@click="submitUpload") 导入问题
        .el-upload__tip(slot="tip") 选择fps格式的xml文件导入
</template>

<script>
export default {
  data() {
    return {};
  },
  methods: {
    submitUpload() {
      this.$refs.upload.submit();
    },
    handleUploadSuccess(response, file, fileList) {
      console.log(response);
      this.$notify({
        title: '成功',
        message: response.message,
        type: 'success'
      });
      let message = '';
      for (let index in response.info) {
        message += response.info[index] + '</br>';
      }
      this.$alert(message, {
        dangerouslyUseHTMLString: true
      });
    },
    handleUploadError(err, file, fileList) {
      console.log('error', response);
    }
  }
};
</script>

<style lang="scss" scoped>
.problem-import-card {
  width: 400px;
  height: 300px;
}
</style>