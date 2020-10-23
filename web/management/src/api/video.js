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

// 待审核的视频列表
export function videoReviewList(query) {
  return new Promise((resolve, reject) => {
    request ({
      url: '/video/review/list',
      method: 'get',
      params: query,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 视频标签列表
export function videoLabelList(query) {
  return new Promise((resolve, reject) => {
    request({
      url: '/video/label/list',
      method: 'get',
      params: query,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 删除标签
export function delLabel(params) {
  return new Promise((resolve, reject) => {
    request ({
      url: '/video/del/label',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 添加视频标签
export function addVideoLabel(params) {
  return new Promise((resolve, reject) => {
    request({
      url: '/video/add/label',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}
