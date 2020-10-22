<template>
  <div class="app-container">
<!--    <div align="right" style="margin-bottom: 20px;margin-top: 20px">-->
<!--    <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="queryBefore">-->
<!--      搜索-->
<!--    </el-button>-->

<!--    </div>-->
    <div class="filter-container" style="display: flex;">
      <el-input v-model="listQuery.gameName" placeholder="游戏名称" style="width: 200px;margin-right: 10px;" class="filter-item" />
      <!-- <el-select v-model="listQuery.cateid" placeholder="类别" clearable style="width: 90px" class="filter-item">
        <el-option v-for="item in typeList" :key="item.id" :label="item.name" :value="item.id" />
      </el-select> -->
      <el-cascader v-model="listQuery.cateid" :options="typeList" :show-all-levels="false" filterable :props="{ checkStrictly: true, value: 'cateid', label: 'name', children: 'child'}" clearable></el-cascader>
      <el-select v-model="listQuery.status" placeholder="状态" clearable class="filter-item" style="width: 130px;margin: 0 10px;">
        <el-option v-for="item in stateList" :key="item.id" :label="item.name" :value="item.id" />
      </el-select>
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="queryBefore">
        搜索
      </el-button>
    </div>
    <el-table v-loading="loading" :data="list" border fit highlight-current-row style="width: 100%">
      <el-table-column align="center" label="视频ID" width="80">
        <template slot-scope="scope">
          <span>{{ scope.row.video_id }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="视频标题">
        <template slot-scope="scope">
          <span>{{ scope.row.title}}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="视频描述">
        <template slot-scope="scope">
          <span>{{ scope.row.describe }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="视频封面">
        <template slot-scope="scope">
          <img :src=scope.row.cover height="120" >
        </template>
      </el-table-column>

      <el-table-column width="500px" align="center" label="视频地址">
        <template slot-scope="scope">
          <el-button :type="'primary'" size="mini" v-if="playerOptions[scope.$index].isVideoVisible === false" @click="watchVideo(scope.$index)">点我观看视频</el-button>
          <div class="video" v-if="playerOptions[scope.$index].isVideoVisible">
            <el-button :type="'info'" size="mini" @click="handleVideoVisible(scope.$index)">点我关闭视频</el-button>
            <video-player
              class="video-player vjs-custom-skin"
              ref="videoPlayer"
              :playsinline="true"
              :options="playerOptions[scope.$index]"
              @play="onPlayerPlay($event)"
            ></video-player>
          </div>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="视频时长">
        <template slot-scope="scope">
          <span>{{ scope.row.video_duration }}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="评论内容">
        <template slot-scope="scope">
          <span>{{ scope.row.content }}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="点赞量">
        <template slot-scope="scope">
          <span>{{ scope.row.like_num }}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="回复量">
        <template slot-scope="scope">
          <span>{{ scope.row.reply_num }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="用户id">
        <template slot-scope="scope">
          <span>{{ scope.row.user_id }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="评论时间">
        <template slot-scope="scope">
          <span>{{ scope.row.create_at | formatDate}}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="操作" width="130">
        <template slot-scope="scope">
          <el-button :type="(scope.row.status === 1)?'primary':'info'" size="mini" @click="handleDelComment(scope.row, 1)" :disabled="!(scope.row.status === 1)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.size" @pagination="listData" />
  </div>
</template>
<script>
  import {
    videoCommentList,
    delVideoComment,
  } from '@/api/comment'
  import Pagination from '@/components/Pagination'
  import {formatDate} from '@/utils/format-date'
  export default {
    components: { Pagination },
    filters: {
      formatDate(time) {
        time = time * 1000
        let date = new Date(time)
        return formatDate(date, 'yyyy-MM-dd hh:mm')
      },
    },
    data() {
      return {
        playerOptions:[],
        total: 0,
        list: [],
        listQuery: {
          page: 1,
          size: 10,
        },
        loading: 1,
      }
    },

    created() {
      this.listData();
    },

    methods: {
      async listData() {
        const res = await videoCommentList(this.listQuery);
        console.log(res)
        if (res.code === 200) {
          console.log(res.data.list)
          this.list = res.data.list;
          this.total = res.data.total;
          for (let v of this.list) {
            let videoConf = {
              playbackRates: [0.7, 1.0, 1.5, 2.0], //播放速度
              autoplay: false, //如果true,浏览器准备好时开始回放。
              muted: false, // 默认情况下将会消除任何音频。
              loop: false, // 导致视频一结束就重新开始。
              preload: 'auto', // 建议浏览器在<video>加载元素后是否应该开始下载视频数据。auto浏览器选择最佳行为,立即开始加载视频（如果浏览器支持）
              language: 'zh-CN',
              aspectRatio: '16:9', // 将播放器置于流畅模式，并在计算播放器的动态大小时使用该值。值应该代表一个比例 - 用冒号分隔的两个数字（例如"16:9"或"4:3"）
              fluid: false, // 当true时，Video.js player将拥有流体大小。换句话说，它将按比例缩放以适应其容器。
              sources: [{
                type: "",
                src: v.video_addr// url地址
              }],
              poster: v.cover, // 封面地址
              width: document.documentElement.clientWidth,
              notSupportedMessage: '此视频暂无法播放，请稍后再试', //允许覆盖Video.js无法播放媒体源时显示的默认信息。
              controlBar: {
                timeDivider: true,
                durationDisplay: true,
                remainingTimeDisplay: false,
                fullscreenToggle: true  //全屏按钮
              },
              isVideoVisible: false,
            }
            this.playerOptions.push(videoConf)
          }
        } else {
          this.list = [];
          this.$message.error(res.message)
        }

        this.loading = 0;
      },

      // 重新拉取页面信息
      refreshList() {
        this.listData();
      },

      handleDelComment(row, state) {
        const ids = row.id;
        this.$confirm('是否确认删除评论id为"' + ids + '"的数据项?', "警告", {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning"
          }).then(() => {
            return this.delVideoComment(ids, state)
        }).catch(()=> {});

      },

      // 修改评论状态
      async delVideoComment(id, state) {
        const res = await delVideoComment({
          comment_id: id.toString(),
          status: state
        });
        console.log(res);
        if (res.code === 200) {
          this.refreshList()
          this.$message.success("评论id为" + id + "，删除成功")
        } else {
          this.$message.error(res.message)
        }
      },

      // 控制视频弹窗开关
      handleVideoVisible(index) {
        this.playerOptions[index].isVideoVisible = !this.playerOptions[index].isVideoVisible;
      },
      // 打开视频
      watchVideo(index) {
        this.playerOptions[index].isVideoVisible = true;
      },

      //方法
      full(element) {
        //做兼容处理
        if (element.requestFullscreen) {
          element.requestFullscreen();
        } else {
          var docHtml = document.documentElement;
          var docBody = document.body;
          var videobox = document.getElementsByClassName("video-player");
          var cssText = "width:100%;height:100%;overflow:hidden;";
          docHtml.style.cssText = cssText;
          docBody.style.cssText = cssText;
          videobox.style.cssText = cssText + ";" + "margin:0px;padding:0px;";
          document.IsFullScreen = true;
        }
      },
      onPlayerPlay(player) {
        this.full(player)
      },
    },
  }
</script>
