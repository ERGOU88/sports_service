package models

type Orders struct {
	Id                int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	OrderId           string `json:"order_id" xorm:"not null default '' comment('订单编号') VARCHAR(50)"`
	Extra             string `json:"extra" xorm:"comment('记录订单相关扩展数据') TEXT"`
	OrderTypeName     string `json:"order_type_name" xorm:"not null default '' comment('订单类型名称') VARCHAR(50)"`
	Transaction       string `json:"transaction" xorm:"not null default '' comment('第三方订单号') VARCHAR(100)"`
	DeliveryCode      string `json:"delivery_code" xorm:"not null default '' comment('提货编码') VARCHAR(50)"`
	PayStatus         int    `json:"pay_status" xorm:"not null default 0 comment('支付状态 0 待付款 1 已超时/已取消 2 已付款') index TINYINT(2)"`
	DeliveryStatus    int    `json:"delivery_status" xorm:"not null default 0 comment('配送状态 0 未配送 1 已配送 2 已签收') TINYINT(2)"`
	RefundStatus      int    `json:"refund_status" xorm:"not null default 0 comment('退款状态 0 未退款 1 退款中 2 已退款') TINYINT(2)"`
	PayType           int    `json:"pay_type" xorm:"not null default 0 comment('支付方式 2001 支付宝 3001微信') INT(8)"`
	DeliveryType      int    `json:"delivery_type" xorm:"not null default 0 comment('配送方式') INT(10)"`
	DeliveryTypeName  string `json:"delivery_type_name" xorm:"not null default '' comment('配送方式名称') VARCHAR(50)"`
	UserId            string `json:"user_id" xorm:"not null default '' comment('买家uid') VARCHAR(60)"`
	ProductAmount     int    `json:"product_amount" xorm:"not null default 0 comment('商品总金额 (分) ') INT(10)"`
	DeliveryAmount    int    `json:"delivery_amount" xorm:"not null default 0 comment('配送费用（分）') INT(10)"`
	OrderAmount       int    `json:"order_amount" xorm:"not null default 0 comment('订单合计金额 (分)') INT(10)"`
	DiscountAmount    int    `json:"discount_amount" xorm:"not null default 0 comment('订单优惠金额 (分)') INT(10)"`
	PayAmount         int    `json:"pay_amount" xorm:"not null default 0 comment('应付金额') INT(10)"`
	PayTime           int    `json:"pay_time" xorm:"not null default 0 comment('订单支付时间') INT(11)"`
	DeliveryTime      int    `json:"delivery_time" xorm:"not null default 0 comment('订单配送时间') INT(11)"`
	SignTime          int    `json:"sign_time" xorm:"not null default 0 comment('订单签收时间') INT(11)"`
	FinishTime        int    `json:"finish_time" xorm:"not null default 0 comment('订单完成时间') INT(11)"`
	CloseTime         int    `json:"close_time" xorm:"not null default 0 comment('订单关闭时间') INT(11)"`
	IsEvaluate        int    `json:"is_evaluate" xorm:"not null default 0 comment('是否允许订单评价 0 允许 1不允许') INT(11)"`
	IsDelete          int    `json:"is_delete" xorm:"not null default 0 comment('是否删除 0 未删 1 已删') INT(11)"`
	IsEnableRefund    int    `json:"is_enable_refund" xorm:"not null default 0 comment('是否允许退款 0 允许 1 不允许') INT(11)"`
	Remark            string `json:"remark" xorm:"not null default '' comment('卖家留言') VARCHAR(255)"`
	EvaluateStatus    int    `json:"evaluate_status" xorm:"not null default 0 comment('评价状态，0：未评价，1：已评价，2：已追评') TINYINT(2)"`
	RefundAmount      int    `json:"refund_amount" xorm:"not null default 0 comment('订单退款金额（分）') INT(10)"`
	OrderType         int    `json:"order_type" xorm:"not null comment('下单方式：1001 APP下单，1002 小程序下单，1003第三方推广渠道购买') INT(8)"`
	CreateAt          int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt          int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	ActionType        int    `json:"action_type" xorm:"not null default 0 comment('1 商品详情页下单 2 购物车下单') TINYINT(2)"`
	DeliveryTelephone string `json:"delivery_telephone" xorm:"not null default '0' comment('承运人电话') VARCHAR(30)"`
}
