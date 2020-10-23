<template>
  <el-dialog title="新增banner" :before-close="handleClose" :visible.sync="ishide" class="add">
    <el-form :model="form" ref="ruleForm" :rules="rules" label-position="right" label-width="90px">
      <el-form-item label="banner名称" prop="title">
        <el-input v-model="form.title"></el-input>
      </el-form-item>
      <el-form-item label="排序权重" prop="sortorder"><el-input v-model="form.sortorder" type="sortorder"></el-input></el-form-item>

      <el-form-item label="跳转地址" prop="jump_url"><el-input v-model="form.jump_url" type="jump_url"></el-input></el-form-item>

      <el-form-item label="分享地址" prop="share_url"><el-input v-model="form.share_url" type="share_url"></el-input></el-form-item>

      <el-form-item label="说明" prop="explain"><el-input v-model="form.explain" type="explain"></el-input></el-form-item>

      <el-form-item label="banner类型" prop="explain">
        <el-select v-model="form.type" placeholder="banner类型" clearable class="filter-item" style="width: 130px;margin: 0 0px;">
          <el-option v-for="item in types" :key="item.id" :label="item.name" :value="item.id" />
        </el-select>
      </el-form-item>

      <el-form-item label="banner图" prop="cover"> <upload @uploadFile='uploadFile'></upload></el-form-item>
    <el-form-item label="开始时间" prop="start_time">
      <template slot-scope="scope">
        <el-row :gutter="0">
          <el-col :span="7">
            <el-date-picker v-model="form.start_time" type="datetime" placeholder="选择日期时间" default-time="12:00:00" format="yyyy-MM-dd HH:mm:ss"> </el-date-picker>
          </el-col>
        </el-row>
      </template>
    </el-form-item>

    <el-form-item label="结束时间" prop="end_time">
      <template slot-scope="scope">
        <el-row :gutter="0">
          <el-col :span="7">
            <el-date-picker v-model="form.end_time" type="datetime" placeholder="选择日期时间" default-time="12:00:00" format="yyyy-MM-dd HH:mm:ss"> </el-date-picker>
          </el-col>
        </el-row>
      </template>
    </el-form-item>

    </el-form>

    <div slot="footer" class="dialog-footer">
      <el-button type="primary" :loading="loading" @click.native.prevent="submitForm">确定</el-button>
    </div>
  </el-dialog>
</template>

<script type="text/javascript">
  import {isNumber} from '@/utils/validate';
  import upload from "./upload.vue"
  import {
    addBanner,
  } from '@/api/configure'
  export default {
    data() {
      const validateSortorder = (rule, value, callback) => {
        if (isNumber(value)) {
          callback();
        } else {
          callback(new Error('请输入正确的权重值（需大于0）'));
        }
      };

      const validateCover = (rule, value, callback) => {
        const tmpcnt = value.lastIndexOf('.')
        const exname = value.substring(tmpcnt + 1)
        const names = ['jpg', 'jpeg', 'png']
        if (names.indexOf(exname) < 0) {
          callback(new Error('不支持的图片格式!请重新选择图片！'))
        } else {
          callback()
        }
      }

      return {
        rules: {
          title: [
            { required: true, message: '请输入标题名称' },
            { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
          ],
          sortorder: [
            {
              validator: validateSortorder,
              trigger: 'blur'
            }
          ],

          start_time: [
            { required: true, message: '请选择开始日期' }
          ],
          end_time: [
            { required: true, message: '请选择结束日期' }
          ],
          cover: [
            { required: true, message: '请确认是否已上传banner图' },
            {
              validator: validateCover,
              trigger: 'blur'
            }
          ]
        },
        form: {
          name: '',
          sortorder: 0,
          cover: '',
          explain: '',
          jump_url: '',
          share_url: '',
          start_time: 0,
          end_time: 0,
          // 1 首页 2 直播页 3 官网banner
          type: 1,
          title: ''
        },
        types: [
          {
            id: 1,
            name: '首页'
          },
          {
            id: 2,
            name: '直播页'
          },
          {
            id: 3,
            name: '官网'
          }
        ],
        loading: false,
        formLabelWidth: '120px',
        ishide: this.show,
        fileUrl: ''
      };
    },
    components: {
      upload,
    },
    props: {
      show: {
        type: Boolean,
        default: true
      }
    },
    created() {},
    methods: {
      handleClose() {
        this.$emit('handleClose');
      },

      // 添加banner
      handleAddBanner(form) {
        this.addBannerInfo(form)
      },

      uploadFile(fileUrl) {
        this.form.cover = fileUrl
      },

      async addBannerInfo(form) {
        const startTime = new Date(form.start_time)
        const endTime = new Date(form.end_time)

        form.start_time = Date.parse(startTime);
        form.start_time = form.start_time / 1000;

        form.end_time = Date.parse(endTime);
        form.end_time = form.end_time / 1000;

        if (form.end_time <= form.start_time || form.end_time - form.start_time < 1800) {
          this.$message.error("Banner自动下架时间 必须大于 上架时间 (且上架时长必须大于30分钟)")
          return
        }

        const timestamp = (new Date()).valueOf() / 1000
        if (timestamp >= form.end_time) {
          this.$message.error("下架时间不能小于当前时间")
          return
        }

        const res = await addBanner({
            cover: form.cover,
            explain: form.explain,
            jump_url: form.jump_url,
            share_url: form.share_url,
            sortorder: parseInt(form.sortorder),
            start_time: form.start_time,
            end_time: form.end_time,
            // 1 首页 2 直播页 3 官网banner
            type: form.type,
            title: form.title
          });
          console.log(res);
          if (res.code === 200) {
            this.$message.success("banner添加成功")
            this.handleClose()
          } else {
            this.$message.error(res.message)
          }
      },
      submitForm() {
        this.$refs.ruleForm.validate(async valid => {
          if (valid) {
            console.log(this.form.name)
            this.handleAddBanner(this.form)
          } else {
            console.log('error submit!!');
            return false;
          }
        });
      }
    },
  }
</script>

<style lang="scss" scoped>
  .add {
    line-height: initial;
    .el-dialog {
      width: 60%;
      padding: 0 40px;
      max-width: 800px;

      .el-dialog__header {
        border-bottom: 1px solid #ebebeb;
        font-size: 24px;
        font-weight: 600;
      }

      .el-dialog__body {
        border-bottom: 1px solid #ebebeb;
        // padding-bottom: 0;
      }

      .el-form {
        // width: 500px;
        max-width: 500px;
        margin: 0 auto;

        .el-form-item {
          // margin-bottom: 35px;
        }

        .input-msg {
          position: absolute;
          font-size: 12px;
          padding-top: 4px;
          line-height: 1;
          top: 100%;
          left: 0;
          color: #424242;
        }

        label {
          font-weight: initial;
          padding-right: 20px;
        }

        .input-code {
          max-width: 200px;
          display: inline-block;
        }

        .input-code-btn {
          padding-left: 10px;
          display: inline-block;

          button {
            width: 130px;
            padding: 10px 0;
          }
        }
      }

      .dialog-footer {
        text-align: center;

        .diaglog-check {
          margin: 0 0 20px 0;

          span {
            font-size: 14px;
            color: #1f1f1f;

            a {
              color: #409eff;
              text-decoration: underline;
            }
          }
        }
      }
    }
  }
</style>
