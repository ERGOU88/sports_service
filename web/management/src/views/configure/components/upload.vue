<!-- 文件上传组件 -->
<template>
  <div class="uplod">
    <el-dialog :visible.sync="dialogVisible">
      <img width="100%" :src="imageUrl" alt="">
    </el-dialog>
    <el-upload
      :http-request="Upload"
      action=""
      list-type="picture-card"
      ref="upload"
      :limit="1"
      :before-upload="beforeUpload"
      :on-preview="handlePictureCardPreview"
      :on-remove="handleRemove"
      :auto-upload="true"
    >
      <i class="el-icon-plus"></i>
    </el-upload>
<!--    <el-button style="margin-left: 10px;" size="small" type="success" @click="submitUpload">上传到服务器</el-button>-->
    <!--    <div class="width30">-->
    <!--      <el-progress :percentage="percentage"></el-progress>-->
    <!--    </div>-->
  </div>
</template>

<script>
  import {fileUpload} from '@/api/user'
  export default{
    data() {
      return {
        imageUrl:'', // 记录已上传图片
        percentage: 0, // 记录进度条
        dialogVisible: false,
        dialogImageUrl: '',
        isPass: false,
      }
    },
    created() {
    },
    methods: {
      handleRemove(file, fileList) {
        console.log(file)
      },
      handlePictureCardPreview(file) { // 点击预览图片
        this.dialogImageUrl = file.url;
        this.dialogVisible = true;
      },
      handleAvatarSuccess(res, file) {
        this.imageUrl = URL.createObjectURL(file.raw);
      },
      progressBar(p, _checkpoint) { // 返回上传进度 参数 p 上传进度, _checkpoint 返回的断点信息
        console.log(p); // Object的上传进度。
        this.percentage = parseInt(p * 100);
      },
      beforeUpload(file) {
        this.fileUpload(file)
      },

      async fileUpload(file) {
        const res = await fileUpload(file);
        console.log(res);
        if (res.code === 200) {
          this.$message.success("banner图已选中")
          console.log(res.data.file_url)
          this.$emit('uploadFile', res.data.file_url)
        } else {
          this.$message.error(res.message)
        }
      },
      submitUpload() {
        this.$refs.upload.submit();
      },
      Upload(file) {
        // debugger;
        // const that = this
        // // 判断扩展名
        // const tmpcnt = file.file.name.lastIndexOf('.')
        // const exname = file.file.name.substring(tmpcnt + 1)
        // const names = ['jpg', 'jpeg', 'webp', 'png', 'bmp', 'apk', 'rar']
        // that.name = file.file.name;
        // if (names.indexOf(exname) < 0) {
        //   this.$message.error('不支持的格式!')
        //   return
        // }
        //
        // this.isPass = true

      }
    }
  }
</script>

<style scoped lang="scss">
  .width30{max-width: 300px;}
</style>
