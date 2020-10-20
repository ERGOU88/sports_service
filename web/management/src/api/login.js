import request from "@/utils/request";

export function adminLogin(params) {
  return new Promise((resolve, reject) => {
    request({
      url: '/admin/login',
      method: 'post',
      data: params,
    }).then(res => {
      resolve(res)
    }).catch(err => {
      resolve(err)
    })
  })
}
