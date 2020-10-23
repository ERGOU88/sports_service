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

// 删除视频评论
export function delVideoComment(params) {
  return new Promise((resolve, reject) => {
    request({
      url: '/comment/delete',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 视频弹幕列表
export function videoBarrageList(query) {
  return new Promise((resolve, reject) => {
    request({
      url: '/comment/barrage',
      method: 'get',
      params: query,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 删除视频弹幕
export function delVideoBarrage(params) {
  return new Promise((resolve, reject) => {
    request({
      url: '/comment/barrage/delete',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

