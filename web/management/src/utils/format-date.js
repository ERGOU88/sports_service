// 时间戳（秒） 转 日期
export function formatDate(value) {
    let date = new Date(value);
    let y = date.getFullYear();
    let MM = date.getMonth() + 1;
    MM = MM < 10 ? ('0' + MM) : MM;
    let d = date.getDate();
    d = d < 10 ? ('0' + d) : d;
    let h = date.getHours();
    h = h < 10 ? ('0' + h) : h;
    let m = date.getMinutes();
    m = m < 10 ? ('0' + m) : m;
    let s = date.getSeconds();
    s = s < 10 ? ('0' + s) : s;
    return y + '-' + MM + '-' + d + ' ' + h + ':' + m + ':' + s;
}

// 秒数 转 天/时/分/秒
export function secondToDate(msd) {
  let time =msd

  if (null != time && "" != time) {
    if (time > 60 && time < 60 * 60) {
      time = parseInt(time / 60.0) + "分钟" + parseInt((parseFloat(time / 60.0) -

        parseInt(time / 60.0)) * 60) + "秒";

    }

    else if (time >= 60 * 60 && time < 60 * 60 * 24) {
      time = parseInt(time / 3600.0) + "小时" + parseInt((parseFloat(time / 3600.0) -

        parseInt(time / 3600.0)) * 60) + "分钟" +

        parseInt((parseFloat((parseFloat(time / 3600.0) - parseInt(time / 3600.0)) * 60) -

          parseInt((parseFloat(time / 3600.0) - parseInt(time / 3600.0)) * 60)) * 60) + "秒";

    } else if (time >= 60 * 60 * 24) {
      time = parseInt(time / 3600.0/24) + "天" +parseInt((parseFloat(time / 3600.0/24)-

        parseInt(time / 3600.0/24))*24) + "小时" + parseInt((parseFloat(time / 3600.0) -

        parseInt(time / 3600.0)) * 60) + "分钟" +

        parseInt((parseFloat((parseFloat(time / 3600.0) - parseInt(time / 3600.0)) * 60) -

          parseInt((parseFloat(time / 3600.0) - parseInt(time / 3600.0)) * 60)) * 60) + "秒";

    }

    else {
      time = parseInt(time) + "秒";

    }

  }

  return time;

}
