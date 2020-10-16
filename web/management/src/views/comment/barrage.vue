<template>
  <div class="app-container">
    <el-table v-loading="loading" :data="list" border fit highlight-current-row style="width: 100%">
      <el-table-column align="center" label="视频ID" width="130">
        <template slot-scope="scope">
          <span>{{ scope.row.video_id }}</span>
        </template>
      </el-table-column>

      <el-table-column width="130px" align="center" label="弹幕ID">
        <template slot-scope="scope">
          <span>{{ scope.row.id}}</span>
        </template>
      </el-table-column>

      <el-table-column width="160px" align="center" label="视频标题">
        <template slot-scope="scope">
          <span>{{ scope.row.title}}</span>
        </template>
      </el-table-column>

      <el-table-column width="160px" align="center" label="视频链接">
        <template slot-scope="scope">
          <span>{{ scope.row.video_addr }}</span>
        </template>
      </el-table-column>

      <el-table-column width="160px" align="center" label="弹幕内容">
        <template slot-scope="scope">
          <span>{{ scope.row.content }}</span>
        </template>
      </el-table-column>

      <el-table-column width="160px" align="center" label="用户id">
        <template slot-scope="scope">
          <span>{{ scope.row.user_id }}</span>
        </template>
      </el-table-column>

      <el-table-column width="160px" align="center" label="时间节点">
        <template slot-scope="scope">
          <span>{{ scope.row.video_cur_duration | secondToDate }}</span>
        </template>
      </el-table-column>

      <el-table-column width="160px" align="center" label="发布时间">
        <template slot-scope="scope">
          <span>{{ scope.row.send_time | formatDate}}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="操作" width="130">
        <template slot-scope="scope">
          <el-button :type="'primary'" size="mini" @click="handleDelBarrage(scope.row)" >删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.size" @pagination="listData" />
  </div>
</template>
<script>
  import {
    videoBarrageList,
    delVideoBarrage,
  } from '@/api/comment'
  import Pagination from '@/components/Pagination'
  import {
    formatDate,
    secondToDate,
  } from '@/utils/format-date'
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
      }
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
        const res = await videoBarrageList(this.listQuery);
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

      handleDelBarrage(row) {
        const ids = row.id;
        this.$confirm('是否确认删除id为"' + ids + '"的弹幕?', "警告", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        }).then(() => {
          return this.delVideoBarrage(ids)
        }).catch(()=> {});

      },

      // 删除弹幕
      async delVideoBarrage(id) {
        const res = await delVideoBarrage({
          id: id.toString(),
        });
        console.log(res);
        if (res.code === 200) {
          this.refreshList()
          this.$message.success("弹幕id为" + id + "，删除成功")
        } else {
          this.$message.error(res.message)
        }
      },
    },

  }
</script>
