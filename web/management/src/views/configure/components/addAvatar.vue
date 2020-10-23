<template>
  <el-dialog title="新增头像" :before-close="handleClose" :visible.sync="ishide" class="add">
    <el-form :model="form" ref="ruleForm" :rules="rules" label-position="right" label-width="90px">
      <el-form-item label="排序权重" prop="sortorder"><el-input v-model="form.sortorder" type="sortorder"></el-input></el-form-item>
      <el-form-item label="选择头像" prop="avatar"> <upload @uploadFile='uploadFile'></upload></el-form-item>
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
    addAvatar,
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

      const validateAvatar = (rule, value, callback) => {
        const tmpcnt = value.lastIndexOf('.')
        const exname = value.substring(tmpcnt + 1)
        const names = ['jpg', 'jpeg', 'png']
        if (names.indexOf(exname) < 0) {
          callback(new Error('不支持的图片格式!请重新选择！'))
        } else {
          callback()
        }
      }

      return {
        rules: {
          sortorder: [
            {
              validator: validateSortorder,
              trigger: 'blur'
            }
          ],
          avatar: [
            { required: true, message: '请确认是否已上传头像' },
            {
              validator: validateAvatar,
              trigger: 'blur'
            }
          ]
        },
        form: {
          sortorder: 0,
          avatar: '',
        },
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

      // 添加系统头像
      handleAddAvatar(form) {
        this.addSysAvatar(form)
      },

      uploadFile(fileUrl) {
        this.form.avatar = fileUrl
      },

      async addSysAvatar(form) {
        const res = await addAvatar({
          avatar: form.avatar,
          sortorder: parseInt(form.sortorder),
        });
        console.log(res);
        if (res.code === 200) {
          this.$message.success("头像添加成功")
          this.handleClose()
        } else {
          this.$message.error(res.message)
        }
      },
      submitForm() {
        this.$refs.ruleForm.validate(async valid => {
          if (valid) {
            console.log(this.form.name)
            this.handleAddAvatar(this.form)
          } else {
            console.log('error submit!!');
            this.$message.error("添加失败")
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
