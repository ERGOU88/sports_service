<template>
  <div class="app-container">
    <div align="right" style="width: 100%;margin-bottom: 20px;margin-top: 20px">
      <el-button :type="'primary'" size="mini" @click="addShow=true">新增banner</el-button>
    </div>
    <el-table v-loading="loading" :data="list" border fit highlight-current-row style="width: 100%">
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="scope">
          <span>{{ scope.row.id }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="排序权重">
        <template slot-scope="scope">
          <span>{{ scope.row.sortorder}}</span>
        </template>
      </el-table-column>


      <el-table-column width="200px" align="center" label="封面">
        <template slot-scope="scope">
          <img :src=scope.row.cover height="150" >
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="跳转链接">
        <template slot-scope="scope">
          <span>{{ scope.row.jump_url }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="分享链接">
        <template slot-scope="scope">
          <span>{{ scope.row.share_url }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="标题">
        <template slot-scope="scope">
          <span>{{ scope.row.title }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="说明">
        <template slot-scope="scope">
          <span>{{ scope.row.explain }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="创建时间">
        <template slot-scope="scope">
          <span>{{ scope.row.create_at | formatDate}}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="开始时间">
        <template slot-scope="scope">
          <span>{{ scope.row.start_time | formatDate}}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="结束时间">
        <template slot-scope="scope">
          <span>{{ scope.row.end_time | formatDate}}</span>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" label="banner位置" width="110" align="center">
        <template slot-scope="scope">
          <el-tag type="success" effect="dark" v-if="scope.row.type===1">首页banner</el-tag>
          <el-tag type="warning" effect="dark" v-if="scope.row.type===2">直播页banner</el-tag>
          <el-tag type="danger" effect="dark" v-if="scope.row.type===3">官网banner</el-tag>
        </template>
      </el-table-column>

      <el-table-column align="center" label="操作" width="100">
        <template slot-scope="scope">
          <el-button :type="'primary'" size="mini" @click="handleDel(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <add v-if="addShow" :show='addShow' @handleClose='handleClose'></add>
    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.size" @pagination="listData" />
   </div>
</template>
<script>
  import {
    bannerList,
    delBanner
  } from '@/api/configure'
  import Pagination from '@/components/Pagination'
  import add from './components/addBanner.vue'
  import {formatDate} from '@/utils/format-date'
  export default {
    filters: {
      formatDate(time) {
        time = time * 1000
        let date = new Date(time)
        return formatDate(date, 'yyyy-MM-dd hh:mm')
      },
    },
    data() {
      return {
        addShow: false,
        setSortShow: false,
        list: [],
        listQuery: {
          page: 1,
          size: 10,
        },
        total: 0,
        loading: 1,
        id: 0,
        sortorder: 0,
      }
    },
    components: {
      Pagination,
      add,
    },

    created() {
      this.listData();
    },

    methods: {
      handleClose() {
        this.addShow = false
        this.refreshList()
      },
      async listData() {
        const res = await bannerList(this.listQuery);
        console.log(res)
        if (res.code === 200) {
          this.list = res.data.list;
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

      handleDel(row) {
        const id = row.id;
        this.$confirm('是否确认删除banner ID为"' + id + '"的数据项?', "警告", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        }).then(() => {
          return this.delBanner(id)
        }).catch(()=> {});

      },

      // 删除banner
      async delBanner(id) {
        const res = await delBanner({
          id: id.toString(),
        });
        console.log(res);
        if (res.code === 200) {
          this.refreshList()
          this.$message.success("banner ID为" + id + "，删除成功")
        } else {
          this.$message.error(res.message)
        }
      },
    }
  }
</script>
