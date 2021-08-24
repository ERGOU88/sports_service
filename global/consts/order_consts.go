package consts


const (
    // 1001 场馆预约 2001 购买月卡 2002 购买季卡 2003 购买年卡 2004 体验券 3001 私教（教练）订单 3002 课程订单 4001 充值订单
    ORDER_TYPE_APPOINTMENT_VENUE  = 1001
    ORDER_TYPE_MONTH_CARD         = 2001
    ORDER_TYPE_SEANSON_CARD       = 2002
    ORDER_TYPE_YEAR_CARD          = 2003
    ORDER_TYPE_EXPERIENCE_CARD    = 2004
    ORDER_TYPE_APPOINTMENT_COACH  = 3001
    ORDER_TYPE_APPOINTMENT_COURSE = 3002
)

const (
    // 可支付时长
    PAYMENT_DURATION = 15 * 60
)


// 0 预约场馆
// 1 预约私教课
// 2 预约大课
const (
    ORDER_APPOINTMENT_VENUE_MSG = iota
    ORDER_APPOINTMENT_COACH_MSG
    ORDER_APPOINTMENT_COURSE_MSG
)

// 0 待支付
// 1 订单超时/未支付
// 2 已支付
// 3 已完成
// 4 退款中
// 5 已退款
// 6 已过期
const (
    PAY_TYPE_WAIT = iota
    PAY_TYPE_UNPAID
    PAY_TYPE_PAID
    PAY_TYPE_COMPLETED
)


// 1001表示android
// 1002表示ios
const (
    PLT_TYPE_ANDROID = 1001
    PLT_TYPE_IOS     = 1002
)

const (
    ALIPAY  = 1
    WEICHAT = 2
)
