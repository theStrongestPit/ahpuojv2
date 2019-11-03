<template lang="pug">
  .content
    title {{`${series?series.name:''} - AHPUOJ`}}
    .content__main
      .siderbar
        ul.siderbar__item__list
          li
            .header 系列赛信息
          li 模式：
              template(v-if="series")          
                oj-tag(:type="series.team_mode == 0 ? 'success':'primary'")  {{ series.team_mode == 0?"个人系列赛":"团体系列赛" }}

      .main
        h1.content__panel__title {{series?series.name:''}}
        .main__section(style="min-height:200px;") 
          h3 系列赛简介
          div(v-if="series",v-html="series.description")
        .main__section
          h3 竞赛列表
          el-table.dataTable(v-if="series",:data="series.contestinfos", style="width: 100%")
            el-table-column(width="90")
              template(slot-scope="scope")
                oj-tag(v-if="scope.row.status==1", type="success") 未开始
                oj-tag(v-if="scope.row.status==2",type="primary") 进行中
                oj-tag(v-if="scope.row.status==3",type="danger") 已结束
            el-table-column(label="名称", min-width="180")
              template(slot-scope="scope")
                router-link(:to="{name:'contest',params:{id:scope.row.id}}") {{scope.row.name}}
        .main__section
          h3 参赛人员信息
          el-table.dataTable(v-if="series",:data="userRankList", style="width: 100%")
            el-table-column(label="用户名",min-width="160")
              template(slot-scope="scope")
                  router-link(:to="{name:'userinfo',params:{id:scope.row.user.id}}")  {{scope.row.user.username}}
            el-table-column(label="昵称",min-width="160")
              template(slot-scope="scope")
                  router-link(:to="{name:'userinfo',params:{id:scope.row.user.id}}")  {{scope.row.user.nick}}
            template(v-for="item in series.contestinfos")
              //- el-table-column(:label="`${item.name}(排名)`",min-width="160")
              el-table-column(:label="`${item.name}(通过)`",min-width="160")
                template(slot-scope="scope")
                  span {{scope.row.solved[""+item.id]?scope.row.solved[""+item.id]:"--"   }}
              
</template>

<script>
import OjTag from '@/web-common/components/ojtag';
import {getSeries} from '@/web-user/js/api/nologin.js';
import {getLanguageList} from '@/web-user/js/api/nologin.js';
import {EventBus} from '@/web-common/eventbus';
import {submitJudgeCode} from '@/web-user/js/api/user.js';
export default {
  components: {
    OjTag
  },
  data() {
    return {
      seeable: false,
      reason: '',
      series: null,
      userRankList: [],
      langList: []
    };
  },
  mounted() {
    this.init();
  },
  methods: {
    async init() {
      console.log('initing');
      const self = this;
      let res = await getLanguageList();
      this.langList = res.data.languages;
      let id = self.$route.params.id;
      try {
        let res = await getSeries(id);
        console.log(res);
        let data = res.data;
        self.series = data.series;
        self.userRankList = data.userranklist;
      } catch (err) {
        console.log(err);
        self.$router.replace({name: '404Page'});
      }
    }
  },
  beforeRouteUpdate(to, from, next) {
    console.log('beforeRouteUpdate!!');
    this.init();
    next();
  }
};
</script>
<style lang="scss" scoped>
</style>