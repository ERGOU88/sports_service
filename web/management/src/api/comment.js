import request from '@/utils/request'

// 视频评论/回复列表
export function videoCommentList(query) {
  return new Promise((resolve, reject) => {
    request({
      url: '/comment/list',
      method: 'get',
      params: query,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}
