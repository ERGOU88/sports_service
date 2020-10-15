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
          <span>{{ scope.row.cover }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="视频地址">
        <template slot-scope="scope">
          <span>{{ scope.row.video_addr }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="视频标签">
        <template slot-scope="scope">
          <span v-for="item in scope.row.labels" :label="item.label_id.toString()" :key="item.label_id"> {{item.label_name}} </span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="视频时长">
        <template slot-scope="scope">
          <span>{{ scope.row.video_duration }}</span>
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
          <el-tag type="success" effect="dark" v-if="scope.row.is_top==1">已置顶</el-tag>
          <el-tag type="warning" effect="dark" v-if="scope.row.is_top==0">未置顶</el-tag>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" label="审核状态" width="110" align="center">
        <template slot-scope="scope">
          <el-tag type="success" effect="dark" v-if="scope.row.status==1">审核成功</el-tag>
          <el-tag type="warning" effect="dark" v-if="scope.row.status==0">待审核</el-tag>
          <el-tag type="danger" effect="dark" v-if="scope.row.status==2">审核失败</el-tag>
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
          this.list = res.data.list;
          this.total = res.data.total;
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

    },

  }
</script>
