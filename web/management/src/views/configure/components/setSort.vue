<template>
  <el-dialog title="修改热搜词权重" :before-close="handleSetClose" :visible.sync="ishide" class="add">
    <el-form :model="form" ref="ruleForm" :rules="rules" label-position="right" label-width="90px">
      <el-form-item label="排序权重" prop="sortorder"><el-input v-model="form.sortorder" type="sortorder"></el-input></el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button type="primary" :loading="loading" @click.native.prevent="submitForm">确定</el-button>
    </div>
  </el-dialog>
</template>

<script type="text/javascript">
  import {isNumber} from '@/utils/validate';
  import {
    setSort,
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

      return {
        rules: {
          sortorder: [
            {
              validator: validateSortorder,
              trigger: 'blur'
            }
          ]
        },
        form: {
          sortorder: this.sortorder,
        },
        loading: false,
        formLabelWidth: '120px',
        ishide: this.show
      };
    },
    props: {
      show: {
        type: Boolean,
        default: true
      },
      id: {
        type: Number,
      },
      sortorder: {
        type: Number,
      }
    },
    created() {},
    methods: {
      handleSetClose() {
        this.$emit('handleSetClose');
      },

      // 设置热搜权重
      handleSetSort(id, sort) {
        this.setSort(id, sort)
      },

      async setSort(id, sort) {
        const res = await setSort({
          id: id,
          sortorder: sort,
        });
        console.log(res);
        if (res.code === 200) {
          this.$message.success("权重设置成功")
          this.handleSetClose()
        } else {
          this.$message.error(res.message)
        }
      },
      // 生命周期钩子函数
      cread(){
        console.log(this.id)
      },
      submitForm() {
        this.$refs.ruleForm.validate(async valid => {
          if (valid) {
            this.cread()
            this.handleSetSort(this.id, parseInt(this.form.sortorder))
          } else {
            console.log('error submit!!');
            return false;
          }
        });
      }
    },
  };
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
