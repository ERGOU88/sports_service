<template>
  <div class="app-container">
    <div align="right" style="width: 79%;margin-bottom: 20px;margin-top: 20px">
      <el-button :type="'primary'" size="mini" @click="addShow=true">新增热搜</el-button>
    </div>
    <el-table v-loading="loading" :data="list" border fit highlight-current-row style="width: 79%">
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


      <el-table-column width="120px" align="center" label="热搜词">
        <template slot-scope="scope">
          <span>{{ scope.row.hot_search_content }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="搜索次数">
        <template slot-scope="scope">
          <span>100</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="创建时间">
        <template slot-scope="scope">
          <span>{{ scope.row.create_at | formatDate}}</span>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" label="当前状态" width="110" align="center">
        <template slot-scope="scope">
          <el-tag type="success" effect="dark" v-if="scope.row.status===0">正常</el-tag>
          <el-tag type="warning" effect="dark" v-if="scope.row.status===1">屏蔽</el-tag>
        </template>
      </el-table-column>


      <el-table-column align="center" label="操作" width="400">
        <template slot-scope="scope">
          <el-button :type="(scope.row.status === 0)?'primary':'info'" size="mini" @click="handleSetStatus(scope.row, 1)" :disabled="!(scope.row.status === 0)">设为屏蔽</el-button>
          <el-button :type="(scope.row.status === 1)?'primary':'info'" size="mini" @click="handleSetStatus(scope.row, 0)" :disabled="!(scope.row.status === 1)">设为正常</el-button>
          <el-button :type="'primary'" size="mini" @click="handleSet(scope.row.id, scope.row.sortorder)">设置权重</el-button>
          <el-button :type="'primary'" size="mini" @click="handleDel(scope.row, 3)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <add v-if="addShow" :show='addShow' @handleClose='handleClose'></add>
    <set v-if="setSortShow" :show='setSortShow' @handleSetClose='handleSetClose' :id="id" :sortorder='sortorder'></set>
  </div>
</template>
<script>
  import {
    hotSearchList,
    delHotSearch,
    setStatus
  } from '@/api/configure'
  import add from './components/addHotSearch.vue'
  import set from './components/setSort.vue'
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
        loading: 1,
        id: 0,
        sortorder: 0,
      }
    },
    components: {
      add,
      set,
    },

    created() {
      this.listData();
    },

    methods: {
      handleClose() {
        this.addShow = false
        this.refreshList()
      },

      handleSet(id, sortorder) {
        this.setSortShow = true
        this.id = id
        this.sortorder = sortorder
      },
      handleSetClose() {
        this.setSortShow = false
        this.refreshList()
      },

      async listData() {
        const res = await hotSearchList();
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

      // 设置热搜内容状态
      handleSetStatus(row, state) {
        this.setStatus(row.id, state)
      },

      // 设置热搜内容状态
      async setStatus(id, state) {
        const res = await setStatus({
          id: id,
          status: state
        });
        console.log(res);
        if (res.code === 200) {
          if (state === 1) {
            this.$message.success("屏蔽成功")
          } else {
            this.$message.success("已取消屏蔽")
          }
          this.refreshList()
        } else {
          this.$message.error(res.message)
        }
      },

      handleDel(row, state) {
        const id = row.id;
        this.$confirm('是否确认删除热搜id为"' + id + '"的数据项?', "警告", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        }).then(() => {
          return this.delHotSearch(id, state)
        }).catch(()=> {});

      },

      // 删除热搜
      async delHotSearch(id, state) {
        const res = await delHotSearch({
          id: id,
          status: state
        });
        console.log(res);
        if (res.code === 200) {
          this.refreshList()
          this.$message.success("热搜id为" + id + "，删除成功")
        } else {
          this.$message.error(res.message)
        }
      },
    }
  }
</script>
