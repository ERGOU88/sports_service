import request from '@/utils/request'

// 视频列表
export function videoList(query) {
  return new Promise((resolve, reject) => {
    request({
      url: '/video/list',
      method: 'get',
      params: query,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 编辑视频状态
export function editVideoStatus(params) {
  return new Promise((resolve, reject) => {
    request({
      url: '/video/edit/status',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 编辑视频置顶/取消置顶
export function editTopStatus(params) {
  return new Promise((resolve, reject) => {
    request ({
      url: '/video/edit/top',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}


// export function videoList(query) {
//   return request({
//     url: '/video/list',
//     method: 'get',
//     params: query
//   })
// }
