import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/vue-admin-template/user/login',
    method: 'post',
    data
  })
}

export function getInfo(token) {
  return request({
    url: '/vue-admin-template/user/info',
    method: 'get',
    params: { token }
  })
}

export function logout() {
  return request({
    url: '/vue-admin-template/user/logout',
    method: 'post'
  })
}

// 获取用户列表
export function userList(query) {
  return new Promise((resolve, reject) => {
    request({
      url: '/user/list',
      method: 'get',
      params: query,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 封禁用户
export function forbidUser(params) {
  return new Promise((resolve, reject) => {
    request({
      url: '/user/forbid',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

// 解除封禁
export function unForbidUser(params) {
  return new Promise((resolve, reject) => {
    request({
      url: '/user/unforbid',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}

