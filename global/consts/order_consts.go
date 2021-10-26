package consts


const (
    // 1001 场馆预约 2101 临时卡 2201 次卡 2311 购买月卡 2321 购买季卡 2331 半年卡 2341 购买年卡 3001 私教（教练）订单 3002 课程订单
    // 4001 充值订单 5101 实体商品 5102 线下结算
    ORDER_TYPE_APPOINTMENT_VENUE  = 1001
    ORDER_TYPE_INTERIM_CARD       = 2101
    ORDER_TYPE_EXPERIENCE_CARD    = 2201
    ORDER_TYPE_MONTH_CARD         = 2311
    ORDER_TYPE_SEANSON_CARD       = 2321
    ORDER_TYPE_HALF_YEAR_CARD     = 2331
    ORDER_TYPE_YEAR_CARD          = 2341
    ORDER_TYPE_APPOINTMENT_COACH  = 3001
    ORDER_TYPE_APPOINTMENT_COURSE = 3002
    ORDER_TYPE_PHYSICAL_GOODS     = 5101
    ORDER_TYPE_SETTLEMENT         = 5102
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
    ORDER_TYPE_WAIT = iota
    ORDER_TYPE_UNPAID
    ORDER_TYPE_PAID
    ORDER_TYPE_COMPLETED
    ORDER_TYPE_REFUND_WAIT
    ORDER_TYPE_REFUND_SUCCESS
    ORDER_TYPE_EXPIRE
)


// 1001表示android
// 1002表示ios
const (
    PLT_TYPE_ANDROID = 1001
    PLT_TYPE_IOS     = 1002
)

// 支付方式  支付宝 alipay 2 微信 weixin
const (
    ALIPAY  = "alipay"
    WEICHAT = "weixin"
)


// 执行类型
// 1 执行退款流程
// 2 查询退款金额、手续费
const (
    EXECUTE_TYPE_REFUND = 1
    EXECUTE_TYPE_QUERY  = 2
)

// 1000 预约场馆
// 2000 卡类
// 3000 预约私交/课程
// 5000 实物类
const (
    PRODUCT_CATEGORY_APPOINTMENT_VENUE  = 1000
    PRODUCT_CATEGORY_CARD               = 2000
    PRODUCT_CATEGORY_APPOINTMENT_COURSE = 3000
    PRODUCT_CATEGORY_MATERIAL_OBJECT    = 5000
)

const (
    EVALUATE_DEFAULT_AVATAR = "https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/default_avatar.png"
)
