<template>

  <div class="app-container">
    <el-table v-loading="loading" :data="list" border fit highlight-current-row style="width: 100%">
      <el-table-column align="center" label="ID" width="80">
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

      <el-table-column width="120px" align="center" label="发布者">
        <template slot-scope="scope">
          <span>{{ scope.row.user_id }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="发布时间">
        <template slot-scope="scope">
          <span>{{ scope.row.create_at | formatDate}}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="视频封面">
        <template slot-scope="scope">
          <img :src=scope.row.cover height="120" >
        </template>
      </el-table-column>

      <el-table-column width="500px" align="center" label="视频链接">
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

      <el-table-column width="120px" align="center" label="视频标签">
        <template slot-scope="scope">
          <span v-for="item in scope.row.labels" :label="item.label_id.toString()" :key="item.label_id"> {{item.label_name}} </span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="视频时长">
        <template slot-scope="scope">
          <span>{{ scope.row.video_duration | secondToDate}}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="弹幕量">
        <template slot-scope="scope">
          <span>{{ scope.row.barrage_num }}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="点赞量">
        <template slot-scope="scope">
          <span>{{ scope.row.fabulous_num }}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="评论量">
        <template slot-scope="scope">
          <span>{{ scope.row.comment_num }}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="分享量">
        <template slot-scope="scope">
          <span>{{ scope.row.share_num }}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="播放量">
        <template slot-scope="scope">
          <span>{{ scope.row.browse_num}}</span>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" label="是否置顶" width="110" align="center">
        <template slot-scope="scope">
          <el-tag type="success" effect="dark" v-if="scope.row.is_top===1">已置顶</el-tag>
          <el-tag type="warning" effect="dark" v-if="scope.row.is_top===0">未置顶</el-tag>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" label="审核状态" width="110" align="center">
        <template slot-scope="scope">
          <el-tag type="success" effect="dark" v-if="scope.row.status===1">已通过</el-tag>
          <el-tag type="warning" effect="dark" v-if="scope.row.status===0">待审核</el-tag>
          <el-tag type="danger" effect="dark" v-if="scope.row.status===2">未通过</el-tag>
        </template>
      </el-table-column>

      <el-table-column align="center" label="操作" width="300">
        <template slot-scope="scope">
          <el-button :type="(scope.row.is_top === 0)?'primary':'info'" size="mini" @click="handleEditTop(scope.row, 1)" :disabled="!(scope.row.is_top === 0)">置顶</el-button>
          <el-button :type="(scope.row.is_top === 1)?'primary':'info'" size="mini" @click="handleEditTop(scope.row, 0)" :disabled="!(scope.row.is_top === 1)">取消置顶</el-button>
          <el-button :type="'primary'" size="mini" @click="handleDel(scope.row, 3)">删除</el-button>
        </template>
      </el-table-column>

<!--      <el-table-column min-width="300px" label="Title">-->
<!--        <template slot-scope="{row}">-->
<!--          <router-link :to="'/video/edit/'+row.id" class="link-type">-->
<!--            <span>{{ row.title }}</span>-->
<!--          </router-link>-->
<!--        </template>-->
<!--      </el-table-column>-->

<!--      <el-table-column align="center" label="Actions" width="120">-->
<!--        <template slot-scope="scope">-->
<!--          <router-link :to="'/video/edit/'+scope.row.id">-->
<!--            <el-button type="primary" size="small" icon="el-icon-edit">-->
<!--              修改-->
<!--            </el-button>-->
<!--          </router-link>-->
<!--        </template>-->
<!--      </el-table-column>-->
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.size" @pagination="listData" />
  </div>
</template>
<script>
  import {
    videoList,
    editVideoStatus,
    editTopStatus
  } from '@/api/video'
  import Pagination from '@/components/Pagination'
  import {formatDate, formatFileSize, secondToDate} from '@/utils/format-date'
  export default {
    components: { Pagination },
    filters: {
      formatDate(time) {
        time = time * 1000
        let date = new Date(time)
        return formatDate(date, 'yyyy-MM-dd hh:mm')
      },
      secondToDate(time) {
        return secondToDate(time)
      },
      formatFileSize(size) {
        return formatFileSize(size)
      }
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
        const res = await videoList(this.listQuery);
        console.log(res)
        if (res.code === 200) {
          //this.playerOptions['sources'][0]['src'] = 'http://vjs.zencdn.net/v/oceans.mp4';
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
      // 编辑后重新拉取页面信息
      refreshList() {
        this.listData();
      },

      // 编辑置顶/取消置顶
      handleEditTop(row, state) {
        this.editTopStatus(row.video_id, state)
      },

      async editTopStatus(id, state) {
        const res = await editTopStatus({
          video_id: id.toString(),
          status: state
        });
        console.log(res);
        if (res.code === 200) {
          if (state === 1) {
            this.$message.success("id为" + id + "的视频，置顶成功")
          } else {
            this.$message.success("id为" + id + "的视频，已取消置顶")
          }
          this.refreshList()
        } else {
          this.$message.error(res.message)
        }
      },

      handleDel(row, state) {
        const ids = row.video_id;
        this.$confirm('是否确认删除视频id为"' + ids + '"的数据项?', "警告", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        }).then(() => {
          return this.editVideoStatus(ids, state)
        }).catch(()=> {});

      },

      // 修改视频状态
      async editVideoStatus(id, state) {
        const res = await editVideoStatus({
          video_id: id.toString(),
          status: state
        });
        console.log(res);
        if (res.code === 200) {
          this.refreshList()
          this.$message.success("视频id为" + id + "，删除成功")
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
        //this.full(player)
      },
    },

  }
</script>
