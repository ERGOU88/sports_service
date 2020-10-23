<template>
  <div class="app-container">
    <el-table v-loading="loading" :data="list" border fit highlight-current-row style="width: 100%">
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="scope">
          <span>{{ scope.row.id }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="用户id">
        <template slot-scope="scope">
          <span>{{ scope.row.user_id}}</span>
        </template>
      </el-table-column>


      <el-table-column width="120px" align="center" label="用户头像">
        <template slot-scope="scope">
          <img :src=scope.row.avatar height="120" >
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="手机号">
        <template slot-scope="scope">
          <span>{{ scope.row.mobile_num }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="最后登陆时间">
        <template slot-scope="scope">
          <span>{{ scope.row.last_login_tm | formatDate}}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="性别">
        <template slot-scope="scope">
          <span v-if="scope.row.gender===1">男</span>
          <span v-if="scope.row.gender===2">女</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="个性签名">
        <template slot-scope="scope">
          <span>{{ scope.row.signature }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="生日">
        <template slot-scope="scope">
          <span>{{ scope.row.born }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="年龄">
        <template slot-scope="scope">
          <span>{{ scope.row.age }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="国家">
        <template slot-scope="scope">
          <span>{{ scope.row.country_cn }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="注册ip">
        <template slot-scope="scope">
          <span>{{ scope.row.reg_ip }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="登陆方式">
        <template slot-scope="scope">
          <span v-if="scope.row.user_type===0">手机号码登陆</span>
          <span v-if="scope.row.user_type===1">微信授权登陆</span>
          <span v-if="scope.row.user_type===2">QQ授权登陆</span>
          <span v-if="scope.row.user_type===3">微博授权登陆</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="点赞数">
        <template slot-scope="scope">
          <span>{{ scope.row.total_likes }}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="粉丝数">
        <template slot-scope="scope">
          <span>{{ scope.row.total_fans }}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="收藏数">
        <template slot-scope="scope">
          <span>{{ scope.row.total_collect }}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="发布数">
        <template slot-scope="scope">
          <span>{{ scope.row.total_publish }}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="浏览量">
        <template slot-scope="scope">
          <span>{{ scope.row.total_browse }}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="评论数">
        <template slot-scope="scope">
          <span>{{ scope.row.total_comment}}</span>
        </template>
      </el-table-column>

      <el-table-column width="90px" align="center" label="弹幕数">
        <template slot-scope="scope">
          <span>{{ scope.row.total_barrage}}</span>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" label="状态" width="110" align="center">
        <template slot-scope="scope">
          <el-tag type="success" effect="dark" v-if="scope.row.status===0">正常</el-tag>
          <el-tag type="warning" effect="dark" v-if="scope.row.status===1">封禁</el-tag>
        </template>
      </el-table-column>

      <el-table-column align="center" label="操作" width="200">
        <template slot-scope="scope">
          <el-button :type="(scope.row.status === 0)?'primary':'info'" size="mini" @click="handleForbid(scope.row.id)" :disabled="!(scope.row.status === 0)">封禁</el-button>
          <el-button :type="(scope.row.status === 1)?'primary':'info'" size="mini" @click="handleUnForbid(scope.row.id)" :disabled="!(scope.row.status === 1)">解封</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.size" @pagination="listData" />
  </div>
</template>
<script>
  import {
    userList,
    forbidUser,
    unForbidUser,
  } from '@/api/user'
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
        const res = await userList(this.listQuery);
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

      // 封禁用户
      handleForbid(id) {
        this.$confirm('是否确认封禁id为"' + id + '"的用户?', "警告", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        }).then(() => {
          return this.forbidUser(id)
        }).catch(()=> {});
      },

      // 解除封禁
      handleUnForbid(id) {
        this.$confirm('是否确认解封id为"' + id + '"的用户?', "警告", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        }).then(() => {
          return this.unForbidUser(id)
        }).catch(()=> {});
      },

      async unForbidUser(id) {
        const res = await unForbidUser({
           id: id.toString(),
        });
        console.log(res);
        if (res.code === 200) {
          this.$message.success("解除封禁成功")
          this.refreshList()
        } else {
          this.$message.error(res.message)
        }
      },

      async forbidUser(id) {
        const res = await forbidUser({
           id: id.toString(),
        });
        console.log(res);
        if (res.code === 200) {
          this.$message.success("封禁成功")
          this.refreshList()
        } else {
          this.$message.error(res.message)
        }
      },
    },

  }
</script>
