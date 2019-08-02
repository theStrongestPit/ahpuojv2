<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`adminIndex`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    el-form(:model="form", :rules="rules", ref="form", label-width="6em")
      el-form-item(label="名称", prop="name")
        el-input(placeholder="请输入名称",v-model="form.name",:autofocus="true")
      el-form-item(label="开始时间", prop="start_time")
        el-date-picker(v-model="form.start_time", type="datetime", format="yyyy-MM-dd HH:mm:ss",value-format="yyyy-MM-dd HH:mm:ss", placeholder="选择开始时间", style="width:100%")
      el-form-item(label="结束时间", prop="end_time")
        el-date-picker(v-model="form.end_time", type="datetime", format="yyyy-MM-dd HH:mm:ss",value-format="yyyy-MM-dd HH:mm:ss", placeholder="选择结束时间", style="width:100%")
      el-form-item(label="题目编号")
        el-input(placeholder="使用 , 半角逗号 分隔题目ID列表，格式如:1000,1001,1002",v-model="form.problems",:autofocus="true")
      el-form-item(label="简介")
        tinymce-editor(v-model="form.description",:height="300")
      el-form-item(label="公开度")
        el-switch(v-model="form.private", active-text="私有", inactive-text="公开", inactive-color="#99cc33", :active-value="1", :inactive-value="0")
      el-form-item(label="团队模式")
        el-switch(v-model="form.team_mode", active-text="团队", inactive-text="个人", inactive-color="#99cc33", :active-value="1", :inactive-value="0",@change="notifyModeChange")
      el-form-item(label="语言")
        el-checkbox-group(v-model="selectedLangList", @change="calcMask")
          el-checkbox(v-for="item,index in langList",:key="item.name", :label="item.name", :disabled="!item.available" :checked="(form.langmask & (1<<index))>0")
      el-form-item(class="fl")
        el-button(type="primary",@click="submit") 提交
</template>

<script>
import TinymceEditor from "@/web-common/components/tinymce_editor.vue";
import {
  createContest,
  editContest,
  getContest
} from "@/web-admin/js/api/contest.js";
import { getLanguageList } from "@/web-user/js/api/nologin.js";
export default {
  components: {
    TinymceEditor
  },
  data() {
    let validateStartTime = (rule, value, callback) => {
      if (value === "") {
        callback(new Error("请输入开始时间"));
      } else if (this.form.end_time != "") {
        let start_time = new Date(this.form.start_time);
        let end_time = new Date(this.form.end_time);
        if (start_time.getTime() < end_time.getTime()) {
          this.$refs.form.validateField("end_time");
          callback();
        } else {
          callback();
        }
      } else {
        callback();
      }
    };
    let validateEndTime = (rule, value, callback) => {
      if (value === "") {
        callback(new Error("请输入结束时间"));
      } else if (this.form.start_time != "") {
        let start_time = new Date(this.form.start_time);
        let end_time = new Date(this.form.end_time);
        if (start_time.getTime() < end_time.getTime()) {
          callback();
        } else {
          callback(new Error("结束必须大于开始时间"));
        }
      } else {
        callback();
      }
    };
    return {
      langList: [],
      selectedLangList: [],
      form: {
        name: "",
        start_time: "",
        end_time: "",
        problems: "",
        description: "",
        langmask: 0,
        private: 0,
        team_mode: 0
      },
      rules: {
        name: [
          {
            required: true,
            message: "请输入问题名称",
            trigger: "blur"
          },
          {
            max: 20,
            message: "超出长度限制",
            trigger: "blur"
          }
        ],
        start_time: [
          {
            validator: validateStartTime,
            trigger: "blur"
          }
        ],
        end_time: [
          {
            validator: validateEndTime,
            trigger: "blur"
          }
        ]
      }
    };
  },
  async mounted() {
    let res = await getLanguageList();
    this.langList = res.data.languages;
    this.init();
  },
  methods: {
    async init() {
      const self = this;

      if (self.$route.name == "adminEditContest") {
        try {
          let id = self.$route.params.id;
          let res = await getContest(id);
          self.form = res.data.contest;
          this.$notify({
            title: "提示",
            message:
              "更改公开度与团队模式后私有比赛的参赛人员需要重新进行设置！",
            duration: 6000
          });

          // 重新计算selectlist
          for (let i in self.langList) {
            if (self.langList[i].available && self.form.langmask & (1 << i)) {
              self.selectedLangList.push(self.langList[i].name);
            }
          }
        } catch (err) {
          console.log(err);
          self.$router.replace({ name: "admin404Page" });
        }
      } else {
        Object.assign(self.form, {
          name: "",
          start_time: "",
          end_time: "",
          problems: "",
          description: "",
          langmask: 0,
          private: 0,
          team_mode: 0
        });
      }
      // 计算默认的langmask

      for (let i in this.langList) {
        if (this.langList[i].available) {
          this.form.langmask = this.form.langmask | (1 << i);
        }
      }
    },
    async submit() {
      const self = this;
      self.$refs["form"].validate(async valid => {
        if (valid) {
          try {
            let res;
            if (self.$route.name == "adminAddContest") {
              res = await createContest(self.form);
            } else {
              let id = self.$route.params.id;
              res = await editContest(id, self.form);
            }
            self.$message({
              message: res.data.message,
              type: "success"
            });
            self.$router.push({ name: "adminContestList" });
          } catch (err) {
            self.$message({
              message: err.response.data.message,
              type: "error"
            });
          }
        } else {
          self.$message({
            message: "表单必填项不能为空",
            type: "error"
          });
          return false;
        }
      });
    },
    calcMask(selectedList) {
      let langmask = 0;
      for (let i = 0; i < this.selectedLangList.length; i++) {
        for (let j = 0; j < this.langList.length; j++) {
          if (this.selectedLangList[i] == this.langList[j].name) {
            langmask = langmask | (1 << j);
          }
        }
      }
      this.form.langmask = langmask;
    },
    notifyModeChange() {
      if (this.form.team_mode == 1) {
        this.$notify({
          title: "提示",
          message: "团队模式只能为私有",
          duration: 6000
        });
      }
    }
  },
  computed: {
    private() {
      return this.form.private;
    },
    teamMode() {
      return this.form.team_mode;
    }
  },
  watch: {
    $route(to, from) {
      if (to.name == "adminAddContest" || to.name == "adminEditContest") {
        this.init();
      }
    },
    private(to, from) {
      if (to == 0) {
        this.form.team_mode = 0;
      }
    },
    teamMode(to, from) {
      if (to == 1) {
        this.form.private = 1;
      }
    }
  }
};
</script>

<style lang="scss" scoped>
</style>