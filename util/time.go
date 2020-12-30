package util

import (
  "bytes"
  "math"
  "strconv"
  "time"
)

type TimeS struct {
  timeStr string
  option  string
  num     int64
}

//时间戳与时间的互相转换
//@param            interface   int64与string格式的数据  前者转换格式化时间 后者转换时间戳
//@param            option      时间戳转换的格式
//@return           args        格式化时间
func (ts *TimeS) GetTimeStrOrStamp(param interface{}, option string) interface{} {
  //根据值的类型不同赋予不同的字段
  ts.option = option
  var result interface{}
  switch param.(type) {
  case int64:
    p := param.(int64)
    ts.num = p
    result = ts.GetStampToFormat()
  case string:
    p := param.(string)
    ts.timeStr = p
    result = ts.GetFormatToStamp()
  default:
    p := param.(int64)
    ts.num = p
    result = ts.GetStampToFormat()
  }

  return result
}

//时间转换时间戳
//@param            num         时间
//@param            option         时间戳转换的格式
//@return            args        格式化时间
func (ts *TimeS) GetFormatToStamp() int64 {
  var timeStr string
  timeStr = ts.timeStr
  option := ts.option
  loc, _ := time.LoadLocation("Asia/Shanghai") //设置时区
  switch option {
  case "YmdHis":
    tt, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
    return tt.Unix()
  case "YmdHi":
    tt, _ := time.ParseInLocation("2006-01-02 15:04", timeStr, loc)
    return tt.Unix()
  case "Ymd":
    tt, _ := time.ParseInLocation("2006-01-02", timeStr, loc)
    return tt.Unix()
  case "ANSIC":
    tt, _ := time.ParseInLocation("Mon Jan _2 15:04:05 2006", timeStr, loc)
    return tt.Unix()
  case "UnixDate":
    tt, _ := time.ParseInLocation("Mon Jan _2 15:04:05 MST 2006", timeStr, loc)
    return tt.Unix()
  case "RFC822Z":
    tt, _ := time.ParseInLocation("02 Jan 06 15:04 -0700", timeStr, loc)
    return tt.Unix()
  case "RFC850":
    tt, _ := time.ParseInLocation("Monday, 02-Jan-06 15:04:05 MST", timeStr, loc)
    return tt.Unix()
  case "RFC1123":
    tt, _ := time.ParseInLocation("Mon, 02 Jan 2006 15:04:05 MST", timeStr, loc)
    return tt.Unix()
  case "RFC1123Z":
    tt, _ := time.ParseInLocation("Mon, 02 Jan 2006 15:04:05 -0700", timeStr, loc)
    return tt.Unix()
  case "RFC3339":
    tt, _ := time.ParseInLocation("2006-01-02T15:04:05Z07:00", timeStr, loc)
    return tt.Unix()
  case "RFC3339Nano":
    tt, _ := time.ParseInLocation("2006-01-02T15:04:05.999999999Z07:00", timeStr, loc)
    return tt.Unix()
  case "Kitchen":
    tt, _ := time.ParseInLocation("3:04PM", timeStr, loc)
    return tt.Unix()
  case "Stamp":
    tt, _ := time.ParseInLocation("Jan _2 15:04:05", timeStr, loc)
    return tt.Unix()
  case "StampMilli":
    tt, _ := time.ParseInLocation("Jan _2 15:04:05.000", timeStr, loc)
    return tt.Unix()
  case "StampMicro":
    tt, _ := time.ParseInLocation("Jan _2 15:04:05.000000", timeStr, loc)
    return tt.Unix()
  case "StampNano":
    tt, _ := time.ParseInLocation("Jan _2 15:04:05.000000000", timeStr, loc)
    return tt.Unix()
  default:
    tt, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
    return tt.Unix()
  }

}

func (ts *TimeS) GetStampToFormat() string {
  num := ts.num
  option := ts.option
  format := make(map[string]interface{})
  format["YmdHis"] = "2006-01-02 15:04:05"
  switch option {
  case "YmdHis":
    return time.Unix(num, 0).Format("2006-01-02 15:04:05")
  case "YmdHi":
    return time.Unix(num, 0).Format("2006-01-02 15:04")
  case "Ymd":
    return time.Unix(num, 0).Format("2006-01-02")
  case "ANSIC":
    return time.Unix(num, 0).Format("Mon Jan _2 15:04:05 2006")
  case "UnixDate":
    return time.Unix(num, 0).Format("Mon Jan _2 15:04:05 MST 2006")
  case "RFC822Z":
    return time.Unix(num, 0).Format("02 Jan 06 15:04 -0700")
  case "RFC850":
    return time.Unix(num, 0).Format("Monday, 02-Jan-06 15:04:05 MST")
  case "RFC1123":
    return time.Unix(num, 0).Format("Mon, 02 Jan 2006 15:04:05 MST")
  case "RFC1123Z":
    return time.Unix(num, 0).Format("Mon, 02 Jan 2006 15:04:05 -0700")
  case "RFC3339":
    return time.Unix(num, 0).Format("2006-01-02T15:04:05Z07:00")
  case "RFC3339Nano":
    return time.Unix(num, 0).Format("2006-01-02T15:04:05.999999999Z07:00")
  case "Kitchen":
    return time.Unix(num, 0).Format("3:04PM")
  case "Stamp":
    return time.Unix(num, 0).Format("Jan _2 15:04:05")
  case "StampMilli":
    return time.Unix(num, 0).Format("Jan _2 15:04:05.000")
  case "StampMicro":
    return time.Unix(num, 0).Format("Jan _2 15:04:05.000000")
  case "StampNano":
    return time.Unix(num, 0).Format("Jan _2 15:04:05.000000000")
  default:
    return time.Unix(num, 0).Format("2006-01-02 15:04:05")
  }

}

//@des 时间转换函数
//@param atime string 要转换的时间戳（秒）
//@return string
func StrTime (atime int64) string{
  var byTime = []int64{365*24*60*60,24*60*60, 60*60, 60, 1}
  var unit = []string{"年前", "天前", "小时前", "分钟前", "秒钟前"}
  now := time.Now().Unix()
  ct := now - atime
  if ct < 0 {
    return "刚刚"
  }

  var res string
  for i := 0; i < len(byTime); i++ {
    if ct < byTime[i] {
      continue
    }

    var temp = math.Floor(float64(ct / byTime[i]))
    ct = ct % byTime[i]
    if temp > 0 {
      var tempStr string
      tempStr = strconv.FormatFloat(temp,'f',-1,64)
      // 字符串拼接
      res = MergeString(tempStr,unit[i])
    }
    // 精确到最大单位，即："2天前"这种形式，如果想要"2天12小时36分钟48秒前"这种形式，把此处break去掉，然后把字符串拼接调整下即可
    break
  }

  return res
}


//@des 拼接字符串
//@param args ...string 要被拼接的字符串序列
//@return string
func MergeString (args ...string) string {
  buffer := bytes.Buffer{}
  for i := 0; i < len(args); i++ {
    buffer.WriteString(args[i])
  }

  return buffer.String()
}
