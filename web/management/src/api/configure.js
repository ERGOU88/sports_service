import request from '@/utils/request'

// 热门搜索列表
export function hotSearchList() {
  return new Promise((resolve, reject) => {
    request({
      url: '/config/hot/search',
      method: 'get',
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 添加热搜
export function addHotSearch(params) {
  return new Promise((resolve, reject) => {
    request({
      url: '/config/add/hot/search',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 删除热搜
export function delHotSearch(params) {
  return new Promise((resolve, reject) => {
    request({
      url: '/config/del/hot/search',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 设置热搜权重
export function setSort(params) {
  return new Promise((resolve, reject) => {
    request({
      url: '/config/set/hot/sort',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 设置热搜状态
export function setStatus(params) {
  return new Promise((resolve, reject) => {
    request({
      url: '/config/set/hot/status',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}
