/*
 Navicat Premium Data Transfer

 Source Server         : aws ccapy
 Source Server Type    : MySQL
 Source Server Version : 50610
 Source Host           : pay.cffwkxrnb4va.ap-northeast-1.rds.amazonaws.com:3306
 Source Schema         : ccpay

 Target Server Type    : MySQL
 Target Server Version : 50610
 File Encoding         : 65001

 Date: 06/07/2018 18:47:06
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for account_info
-- ----------------------------
DROP TABLE IF EXISTS `account_info`;
CREATE TABLE `account_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sch_id` varchar(45) NOT NULL COMMENT '商户 id',
  `types` smallint(6) NOT NULL COMMENT '帐号类型  1支付宝 2微信 3招商银行 4华夏银行 ...各种银行',
  `account` varchar(45) NOT NULL COMMENT '收款帐号',
  `open` varchar(45) DEFAULT '' COMMENT '开户人',
  `status` tinyint(4) DEFAULT NULL COMMENT '状态 0停止 1正常',
  `details` varchar(2048) DEFAULT '' COMMENT '详细 `json`',
  `stop_time` int(11) DEFAULT '0' COMMENT '停止收款时间',
  `create_time` int(11) DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商户收款账户信息';

-- ----------------------------
-- Table structure for account_way
-- ----------------------------
DROP TABLE IF EXISTS `account_way`;
CREATE TABLE `account_way` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `type` int(10) NOT NULL COMMENT '账号类型  100支付宝',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  `name` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `account_UNIQUE` (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for ach_amount_detail
-- ----------------------------
DROP TABLE IF EXISTS `ach_amount_detail`;
CREATE TABLE `ach_amount_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `ach_id` varchar(45) NOT NULL COMMENT '商户号',
  `receive_id` int(11) NOT NULL COMMENT '收款账号表id',
  `amount` int(10) NOT NULL COMMENT '收款的额度',
  `pay_code` varchar(45) NOT NULL COMMENT '付款码',
  `type` int(3) NOT NULL COMMENT '1回款 2充值',
  `account` varchar(45) NOT NULL COMMENT '结算人',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=538 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for ach_daily_order
-- ----------------------------
DROP TABLE IF EXISTS `ach_daily_order`;
CREATE TABLE `ach_daily_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `order_ids` text NOT NULL COMMENT '订单id 逗号分割',
  `amount` int(10) NOT NULL COMMENT '总金额',
  `order_num` int(10) NOT NULL COMMENT '订单数',
  `ach_id` varchar(45) NOT NULL COMMENT '商户号',
  `daily_time` varchar(45) NOT NULL COMMENT '每日时间',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for ach_order_check
-- ----------------------------
DROP TABLE IF EXISTS `ach_order_check`;
CREATE TABLE `ach_order_check` (
  `id` bigint(20) unsigned NOT NULL,
  `uid` bigint(20) DEFAULT NULL COMMENT '商户id',
  `order_id` varchar(255) DEFAULT NULL COMMENT '商户订单号',
  `cert_img` varchar(500) DEFAULT NULL COMMENT '凭证地址',
  `created` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for ach_receive_account
-- ----------------------------
DROP TABLE IF EXISTS `ach_receive_account`;
CREATE TABLE `ach_receive_account` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `ach_id` varchar(45) NOT NULL COMMENT '商户号',
  `type` int(10) NOT NULL COMMENT '账号类型  100支付宝',
  `amount` bigint(20) DEFAULT NULL COMMENT '汇款的额度',
  `use_amount` bigint(20) NOT NULL DEFAULT '0',
  `open` varchar(45) NOT NULL COMMENT '开户人',
  `account` varchar(256) NOT NULL COMMENT '收款账号',
  `stop_unix` bigint(20) DEFAULT NULL COMMENT '已经收到的额度',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 0审核中 1已通过 2未通过\n',
  `password` varchar(45) NOT NULL COMMENT '账户密码',
  `remark` text NOT NULL COMMENT '备注信息\n',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=70 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for admin
-- ----------------------------
DROP TABLE IF EXISTS `admin`;
CREATE TABLE `admin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `account` varchar(45) NOT NULL COMMENT '登录账号',
  `password` varchar(45) DEFAULT '' COMMENT '明文密码',
  `encry_password` varchar(125) DEFAULT '' COMMENT '密文密码',
  `group` tinyint(4) DEFAULT '0' COMMENT '分组',
  `content` varchar(45) DEFAULT '' COMMENT '备注',
  `create_time` int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UQE_admin_account` (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8 COMMENT='管理员信息';

-- ----------------------------
-- Table structure for admin_update_water
-- ----------------------------
DROP TABLE IF EXISTS `admin_update_water`;
CREATE TABLE `admin_update_water` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `admin_id` int(11) NOT NULL,
  `url` varchar(45) NOT NULL,
  `title` varchar(45) DEFAULT '',
  `datas` varchar(1024) DEFAULT '',
  `ip` varchar(45) DEFAULT '',
  `content` varchar(125) DEFAULT '',
  `update_time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=67165 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for amount_detail
-- ----------------------------
DROP TABLE IF EXISTS `amount_detail`;
CREATE TABLE `amount_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `account` varchar(45) NOT NULL COMMENT '影响的用户',
  `amount` int(10) NOT NULL COMMENT '金额',
  `type` int(3) NOT NULL COMMENT '1 授权获取 2收取给别人 3 充值 4花费',
  `charge_type` int(3) NOT NULL COMMENT '1 自己的 2父亲的 3 爷爷的',
  `ext` varchar(2000) DEFAULT NULL COMMENT '扩展参数',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10270 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for badmin_update_water
-- ----------------------------
DROP TABLE IF EXISTS `badmin_update_water`;
CREATE TABLE `badmin_update_water` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `admin_id` int(11) NOT NULL,
  `url` varchar(45) NOT NULL,
  `title` varchar(45) DEFAULT '',
  `datas` varchar(1024) DEFAULT '',
  `ip` varchar(45) DEFAULT '',
  `content` varchar(125) DEFAULT '',
  `update_time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3413 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for black_ip
-- ----------------------------
DROP TABLE IF EXISTS `black_ip`;
CREATE TABLE `black_ip` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `ip` varchar(128) NOT NULL COMMENT 'ip',
  `expire` bigint(20) NOT NULL COMMENT '到期时间戳',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_UNIQUE` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for business_amount_operate
-- ----------------------------
DROP TABLE IF EXISTS `business_amount_operate`;
CREATE TABLE `business_amount_operate` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `business_id` int(11) NOT NULL,
  `operate_id` int(11) NOT NULL COMMENT '操作人',
  `amount` float DEFAULT '0' COMMENT '操作金额',
  `content` varchar(45) DEFAULT '' COMMENT '备注',
  `operate_time` int(11) DEFAULT '0' COMMENT '操作时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商户金额操作';

-- ----------------------------
-- Table structure for business_real_auth
-- ----------------------------
DROP TABLE IF EXISTS `business_real_auth`;
CREATE TABLE `business_real_auth` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ach_id` varchar(45) NOT NULL COMMENT '商户标识',
  `real_name` varchar(45) NOT NULL COMMENT '真实姓名',
  `phone` varchar(45) NOT NULL COMMENT '手机',
  `idcard` varchar(45) NOT NULL COMMENT '身份证',
  `address` varchar(45) DEFAULT '' COMMENT '地址',
  `front_pic` varchar(125) DEFAULT '' COMMENT '身份证正面',
  `back_pic` varchar(125) DEFAULT '' COMMENT '身份证反面',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态 1.审核中 2.通过',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COMMENT='商户实名认证';

-- ----------------------------
-- Table structure for business_recharge
-- ----------------------------
DROP TABLE IF EXISTS `business_recharge`;
CREATE TABLE `business_recharge` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ach_id` varchar(45) NOT NULL COMMENT '商户号',
  `recharge_num` int(11) DEFAULT '0' COMMENT '充值金额',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态：1.审核中 2.通过',
  `mark` varchar(50) DEFAULT '' COMMENT '操作备注',
  `recharge_time` int(11) DEFAULT '0' COMMENT '充值时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for business_user
-- ----------------------------
DROP TABLE IF EXISTS `business_user`;
CREATE TABLE `business_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `account` varchar(45) NOT NULL COMMENT '账户',
  `password` varchar(45) DEFAULT '' COMMENT '密码',
  `encry_password` varchar(45) DEFAULT '' COMMENT '加密密码',
  `name` varchar(45) DEFAULT '' COMMENT '名称',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态 0审核中 1审核通过 2封停',
  `amount` float DEFAULT '0' COMMENT '总额度',
  `odd_amount` float DEFAULT '0' COMMENT '剩余额度',
  `ratio` int(11) DEFAULT '0' COMMENT '分成比例',
  `types` tinyint(4) DEFAULT '0' COMMENT '0签约商户 1平台商户',
  `create_time` int(11) DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `account_UNIQUE` (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for cashing
-- ----------------------------
DROP TABLE IF EXISTS `cashing`;
CREATE TABLE `cashing` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `cash_num` int(11) NOT NULL,
  `arrival_num` int(11) NOT NULL,
  `operate` int(11) DEFAULT '0',
  `status` tinyint(4) DEFAULT '0',
  `create_time` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='提现';

-- ----------------------------
-- Table structure for clear_order
-- ----------------------------
DROP TABLE IF EXISTS `clear_order`;
CREATE TABLE `clear_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `ach_id` varchar(45) NOT NULL COMMENT '商户id',
  `bind_orders` varchar(2000) NOT NULL COMMENT '绑定的的订单',
  `pay_code` varchar(128) NOT NULL COMMENT '付款码 对帐用',
  `status` int(3) NOT NULL COMMENT '订单状态 1未处理 2已处理',
  `price` int(11) NOT NULL COMMENT '应收金额',
  `ticket` varchar(2000) NOT NULL COMMENT '票据',
  `pay_info` varchar(2000) NOT NULL COMMENT '收款方式',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  `account` varchar(128) NOT NULL,
  `open` varchar(128) DEFAULT NULL COMMENT '开户人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=384 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for clear_order_details
-- ----------------------------
DROP TABLE IF EXISTS `clear_order_details`;
CREATE TABLE `clear_order_details` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pay_account` int(11) NOT NULL COMMENT '支付账户',
  `pay_code` varchar(45) NOT NULL COMMENT '结账码',
  `amount` int(11) DEFAULT '0' COMMENT '付款金额',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态',
  `created` int(11) DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for common_config
-- ----------------------------
DROP TABLE IF EXISTS `common_config`;
CREATE TABLE `common_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `config_key` varchar(45) NOT NULL,
  `config_value` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `config_key_UNIQUE` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for cut_amount_apply
-- ----------------------------
DROP TABLE IF EXISTS `cut_amount_apply`;
CREATE TABLE `cut_amount_apply` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `account` varchar(45) NOT NULL COMMENT '代理用户',
  `amount` int(10) NOT NULL,
  `status` int(3) NOT NULL COMMENT ' 1未处理 2已处理',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for feedback
-- ----------------------------
DROP TABLE IF EXISTS `feedback`;
CREATE TABLE `feedback` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ach_id` varchar(45) NOT NULL COMMENT '账户标识',
  `title` varchar(125) NOT NULL COMMENT '反馈标题',
  `content` varchar(255) NOT NULL COMMENT '反馈内容',
  `back_content` varchar(255) DEFAULT '' COMMENT '反馈的反馈内容',
  `status` tinyint(4) DEFAULT '1' COMMENT '反馈状态 1未回复 2已回复',
  `create_time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=55 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for feedback_pics
-- ----------------------------
DROP TABLE IF EXISTS `feedback_pics`;
CREATE TABLE `feedback_pics` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `feedback_id` int(11) NOT NULL,
  `url` varchar(125) NOT NULL COMMENT '图片地址',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='反馈图片';

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `from_id` varchar(45) NOT NULL COMMENT '发送人id',
  `to_id` varchar(45) NOT NULL COMMENT '接受人id',
  `type` int(5) NOT NULL COMMENT '消息类型 结算消息100 额度消息 200 系统消息300',
  `alert` varchar(1024) DEFAULT '' COMMENT '通知内容',
  `body` text COMMENT '消息业务体 json',
  `msg_id` varchar(256) NOT NULL COMMENT '消息id',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=368163 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for pay_api
-- ----------------------------
DROP TABLE IF EXISTS `pay_api`;
CREATE TABLE `pay_api` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ach_id` varchar(45) NOT NULL COMMENT '商户',
  `uid` varchar(45) NOT NULL COMMENT 'uid',
  `token` varchar(45) NOT NULL COMMENT 'token',
  `white_list` varchar(125) DEFAULT '' COMMENT 'Ip白名单',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ach_id_UNIQUE` (`ach_id`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for pay_like
-- ----------------------------
DROP TABLE IF EXISTS `pay_like`;
CREATE TABLE `pay_like` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `account` varchar(45) NOT NULL COMMENT '开户人',
  `type` int(10) NOT NULL COMMENT '账号类型  100支付宝',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=442 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for platform_account
-- ----------------------------
DROP TABLE IF EXISTS `platform_account`;
CREATE TABLE `platform_account` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `account` varchar(45) NOT NULL,
  `open` varchar(45) DEFAULT '',
  `account_type` tinyint(4) DEFAULT '0',
  `bank_type` smallint(6) DEFAULT '0',
  `status` tinyint(4) DEFAULT '0',
  `create_time` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for platform_transfer
-- ----------------------------
DROP TABLE IF EXISTS `platform_transfer`;
CREATE TABLE `platform_transfer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ach_id` varchar(45) NOT NULL COMMENT '商户号',
  `amount` int(10) NOT NULL COMMENT '转账金额（单位：分）',
  `operate_id` int(10) NOT NULL COMMENT '操作人(管理员id)',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `ach_id` (`ach_id`)
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for preduce_requests
-- ----------------------------
DROP TABLE IF EXISTS `preduce_requests`;
CREATE TABLE `preduce_requests` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `proxy_id` int(11) NOT NULL,
  `reduce_num` float NOT NULL,
  `operate` int(11) DEFAULT '0',
  `status` tinyint(4) DEFAULT '0',
  `reduce_time` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='降额申请';

-- ----------------------------
-- Table structure for product
-- ----------------------------
DROP TABLE IF EXISTS `product`;
CREATE TABLE `product` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `busi_account` varchar(45) NOT NULL COMMENT '商品发布者',
  `name` varchar(45) NOT NULL COMMENT '商品名称',
  `price` smallint(6) DEFAULT '0' COMMENT '价格',
  `create_time` int(11) DEFAULT '0' COMMENT '创建时间',
  `delete_time` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商户商品';

-- ----------------------------
-- Table structure for profit_order
-- ----------------------------
DROP TABLE IF EXISTS `profit_order`;
CREATE TABLE `profit_order` (
  `id` bigint(20) NOT NULL COMMENT 'ID',
  `account` varchar(45) NOT NULL COMMENT '金字塔顶尖的人物',
  `amount` int(10) NOT NULL COMMENT '分入金额',
  `status` int(3) NOT NULL COMMENT '订单状态 1未处理 2已处理',
  `date` varchar(128) DEFAULT NULL COMMENT '2015-03-02',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for profit_sub_order
-- ----------------------------
DROP TABLE IF EXISTS `profit_sub_order`;
CREATE TABLE `profit_sub_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `profit_order_id` bigint(20) NOT NULL COMMENT '父分润单',
  `account` varchar(45) NOT NULL COMMENT '分润的用户',
  `amount` int(10) NOT NULL COMMENT '分入金额',
  `bind_profits` varchar(2000) NOT NULL COMMENT '绑定分润的的订单',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  `status` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=342 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for proxy
-- ----------------------------
DROP TABLE IF EXISTS `proxy`;
CREATE TABLE `proxy` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `account` varchar(45) NOT NULL COMMENT '账户',
  `password` varchar(45) NOT NULL COMMENT '密码(加密)',
  `name` varchar(45) DEFAULT '' COMMENT '代理名称',
  `real_name` varchar(45) DEFAULT '' COMMENT '真实姓名',
  `parent_account` varchar(45) DEFAULT NULL COMMENT '上一级代理account',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态 1 正常 2 冻结',
  `amount` int(11) DEFAULT '0' COMMENT '暂时不用,用户的余额没想好，不存在此处',
  `share_code` varchar(45) NOT NULL COMMENT '邀请码',
  `level` tinyint(1) NOT NULL COMMENT '用户等级，总代为1',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  `check_status` int(11) DEFAULT NULL COMMENT '0未审核 1 审核通过 2 审核不通过',
  `pay_app_status` int(3) NOT NULL DEFAULT '1',
  `app_env_status` int(3) DEFAULT NULL COMMENT '0',
  `version` varchar(128) DEFAULT NULL COMMENT '当前app版本号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `share_code_UNIQUE` (`share_code`),
  UNIQUE KEY `account_UNIQUE` (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=139104 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for proxy_account_manage
-- ----------------------------
DROP TABLE IF EXISTS `proxy_account_manage`;
CREATE TABLE `proxy_account_manage` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `account` varchar(45) NOT NULL COMMENT '代理用户',
  `ali_account_nick` varchar(256) NOT NULL COMMENT '支付宝账户昵称',
  `ali_account` varchar(256) NOT NULL COMMENT '支付宝账户',
  `wx_account_nick` varchar(256) NOT NULL COMMENT '微信账户昵称',
  `wx_account` varchar(256) NOT NULL COMMENT '微信账户',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  `ali_id` varchar(20) NOT NULL,
  `wx_id` varchar(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ali_id` (`ali_id`),
  UNIQUE KEY `wx_id` (`wx_id`)
) ENGINE=InnoDB AUTO_INCREMENT=227 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for proxy_amount_operate
-- ----------------------------
DROP TABLE IF EXISTS `proxy_amount_operate`;
CREATE TABLE `proxy_amount_operate` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `proxy_id` varchar(45) NOT NULL,
  `operate_id` int(11) NOT NULL COMMENT '操作人',
  `amount` float DEFAULT '0' COMMENT '操作金额',
  `content` varchar(45) DEFAULT '' COMMENT '备注',
  `operate_time` int(11) DEFAULT '0' COMMENT '操作时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=171 DEFAULT CHARSET=utf8 COMMENT='代理金额操作';

-- ----------------------------
-- Table structure for proxy_node
-- ----------------------------
DROP TABLE IF EXISTS `proxy_node`;
CREATE TABLE `proxy_node` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `node_name` varchar(45) NOT NULL COMMENT '节点名称',
  `owner_account` varchar(45) NOT NULL COMMENT '当前节点的账户',
  `father_account` varchar(45) DEFAULT NULL COMMENT '父节点账户',
  `level` int(5) NOT NULL COMMENT '节点的等级',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for proxy_order
-- ----------------------------
DROP TABLE IF EXISTS `proxy_order`;
CREATE TABLE `proxy_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `account` varchar(45) NOT NULL COMMENT '影响的用户',
  `amount` int(10) NOT NULL COMMENT '金额',
  `status` int(3) NOT NULL COMMENT '订单状态 1未支付 2已经支付',
  `order_id` varchar(128) NOT NULL COMMENT '订单id',
  `ach_order_id` varchar(128) NOT NULL DEFAULT '0' COMMENT '商户订单号',
  `pay_type` int(3) NOT NULL COMMENT '支付类型 100支付宝 200 微信',
  `ach_id` varchar(128) NOT NULL COMMENT '商户号',
  `ach_uid` varchar(128) DEFAULT NULL COMMENT '商户UID',
  `product_id` varchar(128) NOT NULL COMMENT '商品id',
  `is_clear` int(3) NOT NULL COMMENT '是否结账',
  `extend` varchar(2056) NOT NULL COMMENT '业务扩展字段',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(2000) DEFAULT NULL COMMENT '备注',
  `push_msg` varchar(2000) DEFAULT NULL COMMENT '提醒app监听收款通知',
  `notify_msg` text COMMENT '收款通知',
  `un_known_id` varchar(128) DEFAULT NULL COMMENT '绑定不知名id',
  `qr` varchar(2000) DEFAULT NULL COMMENT '二维码信息',
  `manage_account` varchar(4000) DEFAULT NULL COMMENT '支付宝微信账号信息',
  PRIMARY KEY (`id`),
  UNIQUE KEY `proxy_order__id` (`order_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5737 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for proxy_pay
-- ----------------------------
DROP TABLE IF EXISTS `proxy_pay`;
CREATE TABLE `proxy_pay` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `account` varchar(45) NOT NULL COMMENT '代理用户',
  `amount` int(10) NOT NULL,
  `pay_code` varchar(128) NOT NULL COMMENT '付款码 对帐用',
  `status` int(3) NOT NULL COMMENT '0没有上传凭证 1未处理 2已处理',
  `ticket` varchar(2000) NOT NULL COMMENT '票据',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  `pay_info` varchar(2000) DEFAULT NULL,
  `open` varchar(128) DEFAULT NULL COMMENT '开户人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=164 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for proxy_profit
-- ----------------------------
DROP TABLE IF EXISTS `proxy_profit`;
CREATE TABLE `proxy_profit` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `account` varchar(45) NOT NULL COMMENT '影响的用户',
  `rate` varchar(50) NOT NULL COMMENT '当前状态的分成比例 收款代理和平台的比例／爷爷和儿子的比例',
  `amount` int(10) NOT NULL COMMENT '分入金额',
  `order_id` varchar(128) NOT NULL COMMENT '订单id',
  `divide_type` int(3) NOT NULL COMMENT '分成类型  1 自己  2 爸爸 3爷爷',
  `team_id` varchar(50) NOT NULL COMMENT '组id',
  `status` int(3) NOT NULL COMMENT '订单状态 1未处理 2已处理',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5176 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for proxy_qr
-- ----------------------------
DROP TABLE IF EXISTS `proxy_qr`;
CREATE TABLE `proxy_qr` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `account` varchar(45) NOT NULL COMMENT '用户',
  `pay_type` int(10) NOT NULL COMMENT '100:支付宝 200微信',
  `qr_url` varchar(512) NOT NULL COMMENT '二维码连接',
  `amount` int(10) NOT NULL COMMENT '二维码金额',
  `qr_md5` varchar(128) NOT NULL COMMENT '二维码md5',
  `status` int(3) NOT NULL COMMENT '1未审核 2 审核通过 ',
  `mark` varchar(512) NOT NULL COMMENT '备注',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  `parse_url` varchar(128) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `proxy_qr_qr_md5_uindex` (`qr_md5`)
) ENGINE=InnoDB AUTO_INCREMENT=2316 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for proxy_real_auth
-- ----------------------------
DROP TABLE IF EXISTS `proxy_real_auth`;
CREATE TABLE `proxy_real_auth` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `account` varchar(45) NOT NULL COMMENT '代理标识',
  `real_name` varchar(45) NOT NULL COMMENT '真实姓名',
  `phone` varchar(45) NOT NULL COMMENT '手机',
  `idcard` varchar(45) NOT NULL COMMENT '身份证',
  `address` varchar(45) DEFAULT '' COMMENT '地址',
  `front_pic` varchar(125) DEFAULT '' COMMENT '身份证正面',
  `back_pic` varchar(125) DEFAULT '' COMMENT '身份证反面',
  `status` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=115057 DEFAULT CHARSET=utf8mb4 COMMENT='代理实名认证';

-- ----------------------------
-- Table structure for proxy_recharge
-- ----------------------------
DROP TABLE IF EXISTS `proxy_recharge`;
CREATE TABLE `proxy_recharge` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `proxy_id` int(11) NOT NULL,
  `recharge_num` float DEFAULT '0',
  `status` tinyint(4) DEFAULT '0',
  `recharge_time` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for qr
-- ----------------------------
DROP TABLE IF EXISTS `qr`;
CREATE TABLE `qr` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pid` int(11) NOT NULL,
  `url` varchar(125) NOT NULL,
  `status` tinyint(4) DEFAULT '0',
  `create_time` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='二维码';

-- ----------------------------
-- Table structure for qr_consts
-- ----------------------------
DROP TABLE IF EXISTS `qr_consts`;
CREATE TABLE `qr_consts` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `amount` int(11) NOT NULL COMMENT '金额(分)',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for reduce_requests
-- ----------------------------
DROP TABLE IF EXISTS `reduce_requests`;
CREATE TABLE `reduce_requests` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `business_id` int(11) NOT NULL,
  `reduce_num` float NOT NULL,
  `operate` int(11) DEFAULT '0',
  `status` tinyint(4) DEFAULT '0',
  `reduce_time` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='降额申请';

-- ----------------------------
-- Table structure for temp_proxy_order_0703
-- ----------------------------
DROP TABLE IF EXISTS `temp_proxy_order_0703`;
CREATE TABLE `temp_proxy_order_0703` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `pay_type` int(11) DEFAULT NULL,
  `msg_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=201807030312422718 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- ----------------------------
-- Table structure for test
-- ----------------------------
DROP TABLE IF EXISTS `test`;
CREATE TABLE `test` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `content` varchar(45) DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for unknown_order
-- ----------------------------
DROP TABLE IF EXISTS `unknown_order`;
CREATE TABLE `unknown_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `account` varchar(45) NOT NULL COMMENT '影响的用户',
  `amount` int(10) NOT NULL COMMENT '金额',
  `created` int(10) DEFAULT NULL COMMENT '创建时间',
  `updated` int(10) DEFAULT NULL COMMENT '更新时间',
  `msg` text COMMENT '收到的消息体',
  `pay_time` bigint(20) DEFAULT NULL COMMENT '支付的时间时间',
  `title` text COMMENT '收到通知的详细信息',
  `pay_type` int(10) DEFAULT NULL COMMENT '100 支付宝;200微信',
  `msg_id` varchar(128) DEFAULT NULL COMMENT '消息id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unknown_order_msg_id_uindex` (`msg_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4755 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_auth
-- ----------------------------
DROP TABLE IF EXISTS `user_auth`;
CREATE TABLE `user_auth` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `pid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '上级ID，0为顶级',
  `auth_name` varchar(64) NOT NULL DEFAULT '0' COMMENT '权限名称',
  `auth_url` varchar(255) NOT NULL DEFAULT '0' COMMENT 'URL地址',
  `sort` int(1) unsigned NOT NULL DEFAULT '999' COMMENT '排序，越小越前',
  `is_show` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否显示，0-隐藏，1-显示',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '操作者ID',
  `create_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建者ID',
  `update_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改者ID',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '状态，1-正常，0-删除',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `class` varchar(30) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=utf8mb4 COMMENT='权限因子';

-- ----------------------------
-- Table structure for user_role
-- ----------------------------
DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_name` varchar(32) NOT NULL DEFAULT '0' COMMENT '角色名称',
  `detail` varchar(255) NOT NULL DEFAULT '0' COMMENT '备注',
  `create_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建者ID',
  `update_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改这ID',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '状态1-正常，0-删除',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '添加时间',
  `update_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='角色表';

-- ----------------------------
-- Table structure for user_role_auth
-- ----------------------------
DROP TABLE IF EXISTS `user_role_auth`;
CREATE TABLE `user_role_auth` (
  `role_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '角色ID',
  `auth_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '权限ID',
  PRIMARY KEY (`role_id`,`auth_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限和角色关系表';

SET FOREIGN_KEY_CHECKS = 1;
