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

      <el-table-column width="120px" align="center" label="视频封面">
        <template slot-scope="scope">
          <span>{{ scope.row.cover }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="视频链接">
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

      <el-table-column align="center" label="操作" width="300">
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

    },

  }
</script>
