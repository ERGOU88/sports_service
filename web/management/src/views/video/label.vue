<template>
  <div class="app-container">
    <div align="right" style="width: 50%;margin-bottom: 20px;margin-top: 20px">
      <el-button :type="'primary'" size="mini" @click="addShow=true">新增视频标签</el-button>
    </div>
    <el-table v-loading="loading" :data="list" border fit highlight-current-row style="width: 50%">
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="scope">
          <span>{{ scope.row.label_id }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="排序权重">
        <template slot-scope="scope">
          <span>{{ scope.row.sortorder}}</span>
        </template>
      </el-table-column>


      <el-table-column width="120px" align="center" label="标签名称">
        <template slot-scope="scope">
          <span>{{ scope.row.label_name }}</span>
        </template>
      </el-table-column>

<!--      <el-table-column width="120px" align="center" label="icon">-->
<!--        <template slot-scope="scope">-->
<!--          <img :src=scope.row.icon height="120" >-->
<!--        </template>-->
<!--      </el-table-column>-->

      <el-table-column width="120px" align="center" label="创建时间">
        <template slot-scope="scope">
          <span>{{ scope.row.create_at | formatDate}}</span>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" label="当前状态" width="110" align="center">
        <template slot-scope="scope">
          <el-tag type="success" effect="dark" v-if="scope.row.status===1">正常</el-tag>
          <el-tag type="warning" effect="dark" v-if="scope.row.status===2">屏蔽</el-tag>
        </template>
      </el-table-column>


      <el-table-column align="center" label="操作" width="130">
        <template slot-scope="scope">
          <el-button :type="'primary'" size="mini" @click="handleDel(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <add v-if="addShow" :show='addShow' @handleClose='handleClose'></add>
  </div>
</template>
<script>
  import {
    videoLabelList,
    delLabel,
  } from '@/api/video'
  import add from './components/addVideoLabel.vue'
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
        const res = await videoLabelList();
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
        const id = row.label_id;
        this.$confirm('是否确认删除id为"' + id + '"的标签?', "警告", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        }).then(() => {
          return this.delLabel(id)
        }).catch(()=> {});

      },

      // 删除标签
      async delLabel(id) {
        const res = await delLabel({
          label_id: id.toString(),
        });
        console.log(res);
        if (res.code === 200) {
          this.refreshList()
          this.$message.success("标签id为" + id + "，删除成功")
        } else {
          this.$message.error(res.message)
        }
      },
    }
  }
</script>
