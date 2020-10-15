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

      <el-table-column width="90px" align="center" label="视频大小">
        <template slot-scope="scope">
          <span>{{ scope.row.size }}</span>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" label="审核状态" width="110">
        <template slot-scope="scope">
          <el-tag type="warning" effect="dark" v-if="scope.row.status==0">待审核</el-tag>
          <el-tag type="danger" effect="dark" v-if="scope.row.status==2">审核失败</el-tag>
        </template>
      </el-table-column>

      <el-table-column align="center" label="操作" width="300">
        <template slot-scope="scope">
          <el-button :type="(scope.row.status === 0)?'primary':'info'" size="mini" @click="handleEditState(scope.row, 1)" :disabled="!(scope.row.status === 0)">通过</el-button>
          <el-button :type="(scope.row.status === 2)?'primary':'info'" size="mini" @click="handleEditState(scope.row, 3)" :disabled="!(scope.row.status === 2)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.size" @pagination="listData" />
  </div>
</template>
<script>
  import {
    videoReviewList,
    editVideoStatus,
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
        const res = await videoReviewList(this.listQuery);
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

      // 重新拉取页面信息
      refreshList() {
        this.listData();
      },

      handleEditState(row, state) {
        const ids = row.video_id;
        if (state === 3) {
          this.$confirm('是否确认删除视频id为"' + ids + '"的数据项?', "警告", {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning"
          }).then(() => {
            return this.editVideoStatus(ids, state)
          }).catch(()=> {});
        }

        if (state === 1) {
          this.$confirm('是否确认通过视频id为"' + ids + '"的数据项?', "提示", {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning"
          }).then(() => {
            return this.editVideoStatus(ids, state)
          }).catch(()=> {});
        }
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
          if (state === 3) {
            this.$message.success("视频id为" + id + "，删除成功")
          }

          if (state === 1) {
            this.$message.success("视频id为" + id + "，已通过审核")
          }

        } else {
          this.$message.error(res.message)
        }
      },

    },

  }
</script>
