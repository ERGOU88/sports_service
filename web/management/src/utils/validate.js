/**
 * Created by PanJiaChen on 16/11/18.
 */

/**
 * @param {string} path
 * @returns {Boolean}
 */
export function isExternal(path) {
  return /^(https?:|mailto:|tel:)/.test(path)
}

/**
 * @param {string} str
 * @returns {Boolean}
 */
export function validUsername(str) {
  const valid_map = ['admin', 'editor']
  return valid_map.indexOf(str.trim()) >= 0
}

export function isTel(str) {
  const reg = /^1[3456789]\d{9}$/
  return reg.test(str)
}

export function isNumber(value) {
  const r = /^\+?[1-9][0-9]*$/; // 正整数
  return r.test(value)
}
