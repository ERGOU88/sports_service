# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 192.168.5.12 (MySQL 5.7.34)
# Database: venue
# Generation Time: 2021-08-30 04:33:50 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table venue_administrator
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_administrator`;

CREATE TABLE `venue_administrator` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` varchar(255) DEFAULT '',
  `job_number` varchar(255) DEFAULT '' COMMENT '工号',
  `mobile` int(11) unsigned DEFAULT '0' COMMENT '手机号',
  `name` varchar(255) DEFAULT '' COMMENT '用户名称',
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '账号',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `status` tinyint(1) unsigned DEFAULT '0' COMMENT '状态',
  `roles` varchar(255) DEFAULT NULL COMMENT '角色',
  `create_at` int(11) unsigned DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) unsigned DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';

LOCK TABLES `venue_administrator` WRITE;
/*!40000 ALTER TABLE `venue_administrator` DISABLE KEYS */;

INSERT INTO `venue_administrator` (`id`, `user_id`, `job_number`, `mobile`, `name`, `username`, `password`, `status`, `roles`, `create_at`, `update_at`)
VALUES
	(1,'','',0,'admin','admin','$2a$10$nXQM7kynFWrueT9PP62Dp.8Z8AfGCTD.RzeqclzsW5TtF09TNhZvy',0,'ROLE_ADMIN',0,0);

/*!40000 ALTER TABLE `venue_administrator` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table venue_appointment_info
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_appointment_info`;

CREATE TABLE `venue_appointment_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `time_node` varchar(128) NOT NULL DEFAULT '' COMMENT '时间节点 例如 10:00-12:00',
  `duration` int(8) NOT NULL DEFAULT '0' COMMENT '总时长（秒）',
  `real_amount` int(11) NOT NULL COMMENT '真实价格（单位：分）',
  `cur_amount` int(11) NOT NULL COMMENT '当前价格 (包含真实价格、 折扣价格（单位：分）',
  `discount_rate` int(11) NOT NULL DEFAULT '0' COMMENT '折扣率',
  `discount_amount` int(11) NOT NULL DEFAULT '0' COMMENT '优惠的金额',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 正常 1 不可用',
  `quota_num` int(10) NOT NULL DEFAULT '0' COMMENT '配额：可容纳人数/可预约人数 -1表示没有限制',
  `related_id` bigint(20) NOT NULL COMMENT '场馆id/私教课程id/大课id',
  `recommend_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '推荐id 1 低价推荐',
  `appointment_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 场馆预约 1 私教课程预约 2 大课预约',
  `week_num` tinyint(1) NOT NULL COMMENT '1 周一 2 周二 3 周三 4 周四 5 周五 6 周六 0 周日',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重 倒序',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `coach_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '当前时间点 大课关联的老师id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='场馆/私教/课程 预约配置信息表';

LOCK TABLES `venue_appointment_info` WRITE;
/*!40000 ALTER TABLE `venue_appointment_info` DISABLE KEYS */;

INSERT INTO `venue_appointment_info` (`id`, `time_node`, `duration`, `real_amount`, `cur_amount`, `discount_rate`, `discount_amount`, `status`, `quota_num`, `related_id`, `recommend_type`, `appointment_type`, `week_num`, `sortorder`, `create_at`, `update_at`, `coach_id`)
VALUES
	(1,'9:00-10:00',3600,1000,500,50,500,0,5,1,1,0,5,0,0,0,0),
	(16,'10:00-11:00',3600,1000,500,50,500,0,5,1,1,0,5,0,0,0,0),
	(17,'11:00-12:00',3600,1400,700,50,500,0,5,1,1,0,5,0,0,0,0),
	(18,'12:00-13:00',3600,1200,600,50,500,0,5,1,1,0,5,0,0,0,0),
	(19,'13:00-14:00',3600,1200,600,50,500,0,5,1,1,0,5,0,0,0,0),
	(20,'14:00-15:00',3600,1200,600,50,500,0,5,1,1,0,5,0,0,0,0),
	(21,'15:00-16:00',3600,1200,600,50,500,0,5,1,1,0,5,0,0,0,0),
	(22,'16:00-17:00',3600,1200,600,50,500,0,5,1,1,0,5,0,0,0,0),
	(23,'17:00-18:00',3600,1200,600,50,500,0,5,1,1,0,5,0,0,0,0),
	(24,'17:00-18:00',3600,1200,600,50,500,0,5,2,1,1,5,0,0,0,0),
	(25,'18:00-19:00',3600,1200,600,50,500,0,5,2,1,1,5,0,0,0,0),
	(26,'17:00-18:00',3600,1200,600,50,500,0,11,1,1,2,5,0,0,0,3),
	(27,'18:00-19:00',3600,2500,2000,80,500,0,10,1,1,2,5,0,0,0,4),
	(28,'17:00-18:00',3600,1200,600,50,500,0,10,3,1,1,5,0,0,0,0),
	(29,'18:00-19:00',3600,1200,600,50,500,0,10,3,1,1,5,0,0,0,0),
	(30,'18:00-19:00',3600,1200,600,50,500,0,11,1,1,2,5,0,0,0,3),
	(31,'17:00-18:00',3600,1200,600,50,500,0,11,1,1,2,5,0,0,0,4),
	(32,'16:00-17:00',3600,1200,600,50,500,0,5,1,1,0,6,0,0,0,0),
	(33,'17:00-18:00',3600,1200,600,50,500,0,5,1,1,0,6,0,0,0,0);

/*!40000 ALTER TABLE `venue_appointment_info` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table venue_appointment_record
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_appointment_record`;

CREATE TABLE `venue_appointment_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `related_id` bigint(20) NOT NULL COMMENT '关联id 私教/场馆/课程',
  `appointment_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 场馆预约 1 私教预约',
  `time_node` varchar(128) NOT NULL DEFAULT '' COMMENT '时间节点 例如 10:00-12:00',
  `date` varchar(30) NOT NULL DEFAULT ' ' COMMENT '预约日期',
  `pay_order_id` varchar(150) NOT NULL COMMENT '订单号',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `purchased_num` int(10) NOT NULL COMMENT '购买的数量',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '-1 软删除 0 待支付 1 订单超时/未支付/已取消 2 已支付 3 已完成  4 退款中 5 已退款 6 已过期 ',
  `seat_info` varchar(1000) NOT NULL DEFAULT '0' COMMENT '预约的座位信息 ',
  `deduction_tm` bigint(20) NOT NULL DEFAULT '0' COMMENT '抵扣会员时长',
  PRIMARY KEY (`id`),
  KEY `related_id` (`related_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='预约流水记录（私教/场馆/课程）';



# Dump of table venue_appointment_stock
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_appointment_stock`;

CREATE TABLE `venue_appointment_stock` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `date` varchar(128) NOT NULL DEFAULT '' COMMENT '日期 例如 2021-10-01',
  `time_node` varchar(128) NOT NULL COMMENT '时间节点 例如 10:00-12:00',
  `quota_num` int(10) NOT NULL DEFAULT '0' COMMENT '配额：可容纳人数/可预约人数 -1表示没有限制',
  `purchased_num` int(10) NOT NULL DEFAULT '0' COMMENT '已购买数量 [冻结库存]',
  `appointment_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 场馆预约 1 私教预约 2 课程预约',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `related_id` bigint(20) NOT NULL COMMENT '场馆id/私教id/课程id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `date` (`date`,`time_node`,`related_id`,`appointment_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='场馆/私教/课程 预约库存表';



# Dump of table venue_coach_detail
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_coach_detail`;

CREATE TABLE `venue_coach_detail` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `title` mediumtext COMMENT '抬头',
  `name` varchar(60) NOT NULL DEFAULT '' COMMENT '私教名称',
  `address` varchar(128) NOT NULL COMMENT '私教地点',
  `designation` varchar(60) NOT NULL DEFAULT '' COMMENT '认证称号',
  `describe` mediumtext COMMENT '描述',
  `areas_of_expertise` varchar(512) NOT NULL DEFAULT '' COMMENT '擅长领域',
  `cover` varchar(256) NOT NULL DEFAULT '' COMMENT '封面 ',
  `avatar` varchar(256) NOT NULL DEFAULT '' COMMENT '头像',
  `price` int(11) NOT NULL DEFAULT '0' COMMENT '私教价格（分）',
  `event_price` int(11) NOT NULL DEFAULT '0' COMMENT '活动价格 (分)',
  `event_start_time` int(11) NOT NULL DEFAULT '0' COMMENT '活动开始时间',
  `event_end_time` int(11) NOT NULL DEFAULT '0' COMMENT '活动结束时间',
  `sortorder` int(11) NOT NULL DEFAULT '1' COMMENT '排序权重',
  `is_recommend` tinyint(1) NOT NULL DEFAULT '0' COMMENT '推荐（0：不推荐；1：推荐）',
  `is_top` tinyint(1) NOT NULL DEFAULT '0' COMMENT '置顶（0：不置顶；1：置顶；）',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '修改时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 正常 1 废弃',
  `coach_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1 私教课老师 2 大课老师',
  `course_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '课程id [大课老师才会关联]',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='私教表';

LOCK TABLES `venue_coach_detail` WRITE;
/*!40000 ALTER TABLE `venue_coach_detail` DISABLE KEYS */;

INSERT INTO `venue_coach_detail` (`id`, `title`, `name`, `address`, `designation`, `describe`, `areas_of_expertise`, `cover`, `avatar`, `price`, `event_price`, `event_start_time`, `event_end_time`, `sortorder`, `is_recommend`, `is_top`, `create_at`, `update_at`, `status`, `coach_type`, `course_id`)
VALUES
	(1,'鉴黄大师[1对1教学]','1对1教学 [燕过拔毛]','王者峡谷','国服第一鉴黄师','专业 专注 专一','擅长拍黄片','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png',1,0,0,0,1,1,0,0,0,0,1,0),
	(2,'鉴黄大师[大课]','叶秋2[大课老师]','王者峡谷','国服第一鉴黄师','专业 专注 专一','擅长拍黄片','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png',1,0,0,0,1,1,0,0,0,0,2,1),
	(3,'鉴黄大师2[大课]','叶秋3[大课老师]','王者峡谷','国服第一鉴黄师','专业 专注 专一','擅长拍黄片','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png',1,0,0,0,1,1,0,0,0,0,2,1),
	(4,'鉴黄大师3[大课]','叶秋4[大课老师]','王者峡谷','国服第一鉴黄师','专业 专注 专一','擅长拍黄片','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png',1,0,0,0,1,1,0,0,0,0,2,1),
	(5,'鉴黄大师[1对1教学]','叶秋叫兽[1对1]','王者峡谷','国服第一鉴黄师','专业 专注 专一','擅长拍黄片','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png',1,0,0,0,1,1,0,0,0,0,1,0);

/*!40000 ALTER TABLE `venue_coach_detail` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table venue_coach_label_config
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_coach_label_config`;

CREATE TABLE `venue_coach_label_config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `label_name` varchar(30) NOT NULL COMMENT '标签名称',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 有效 1 废弃',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='私教标签配置';

LOCK TABLES `venue_coach_label_config` WRITE;
/*!40000 ALTER TABLE `venue_coach_label_config` DISABLE KEYS */;

INSERT INTO `venue_coach_label_config` (`id`, `label_name`, `status`, `create_at`, `update_at`)
VALUES
	(1,'可可爱爱',0,0,0),
	(2,'没有脑袋',0,0,0);

/*!40000 ALTER TABLE `venue_coach_label_config` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table venue_coach_score
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_coach_score`;

CREATE TABLE `venue_coach_score` (
  `coach_id` bigint(20) NOT NULL COMMENT '私教id',
  `total_score` int(11) NOT NULL DEFAULT '0' COMMENT '1颗星1分',
  `total_num` int(11) NOT NULL DEFAULT '0' COMMENT '评分总人数',
  `total_5_star` int(11) NOT NULL DEFAULT '0' COMMENT '五星评价总人数',
  `total_4_star` int(11) NOT NULL DEFAULT '0' COMMENT '四星评价总人数',
  `total_3_star` int(11) NOT NULL DEFAULT '0' COMMENT '三星评价总人数',
  `total_2_star` int(11) NOT NULL DEFAULT '0' COMMENT '二星评价总人数',
  `total_1_star` int(11) NOT NULL DEFAULT '0' COMMENT '一星评价总人数',
  `create_at` int(11) NOT NULL COMMENT '记录创建时间',
  `update_at` int(11) NOT NULL COMMENT '记录更新时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 正常 1 废弃',
  PRIMARY KEY (`coach_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table venue_course_detail
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_course_detail`;

CREATE TABLE `venue_course_detail` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '课程id',
  `coach_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '关联的教兽id [私教课才会关联]',
  `class_period` int(11) NOT NULL COMMENT '单课程时长（秒）',
  `title` mediumtext COMMENT '课程标题',
  `describe` mediumtext COMMENT '课程描述',
  `price` int(11) NOT NULL COMMENT '课程价格（分/每课时）',
  `event_price` int(11) NOT NULL DEFAULT '0' COMMENT '活动价格',
  `event_start_time` int(11) NOT NULL DEFAULT '0' COMMENT '活动开始时间',
  `event_end_time` int(11) NOT NULL DEFAULT '0' COMMENT '活动结束时间',
  `promotion_pic` varchar(521) NOT NULL DEFAULT '' COMMENT '宣传图',
  `icon` varchar(256) NOT NULL DEFAULT '' COMMENT '图标',
  `sortorder` int(11) NOT NULL DEFAULT '1' COMMENT '排序',
  `is_recommend` tinyint(1) NOT NULL DEFAULT '0' COMMENT '推荐（0：不推荐；1：推荐）',
  `is_top` tinyint(1) NOT NULL DEFAULT '0' COMMENT '置顶（0：不置顶；1：置顶；）',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '修改时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 正常 1 废弃',
  `course_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1 私教课 2 大课',
  `period_num` int(6) NOT NULL COMMENT '总课时数',
  `name` varchar(521) NOT NULL DEFAULT '' COMMENT '课程名称',
  PRIMARY KEY (`id`),
  KEY `coach_id` (`coach_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='私教课程详情表';

LOCK TABLES `venue_course_detail` WRITE;
/*!40000 ALTER TABLE `venue_course_detail` DISABLE KEYS */;

INSERT INTO `venue_course_detail` (`id`, `coach_id`, `class_period`, `title`, `describe`, `price`, `event_price`, `event_start_time`, `event_end_time`, `promotion_pic`, `icon`, `sortorder`, `is_recommend`, `is_top`, `create_at`, `update_at`, `status`, `course_type`, `period_num`, `name`)
VALUES
	(1,0,7200,'鉴黄系列课程','带你成为一名合格的鉴黄师',0,0,0,0,'https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png',1,0,0,0,0,0,2,0,'鉴黄课程'),
	(2,1,3600,'鉴黄私教课程1','私人1对1教兽鉴黄本领基础篇',0,0,0,0,'https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png',1,0,0,0,0,0,1,0,'鉴黄课程'),
	(3,1,3600,'鉴黄私教课程2','私人1对1教兽鉴黄本领进阶篇',0,0,0,0,'https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png',1,0,0,0,0,0,1,0,'鉴黄课程'),
	(4,0,7200,'鉴黄系列课程2','带你成为一名合格的鉴黄师',0,0,0,0,'https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png',1,0,0,0,0,0,2,0,'鉴黄课程');

/*!40000 ALTER TABLE `venue_course_detail` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table venue_entry_or_exit_records
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_entry_or_exit_records`;

CREATE TABLE `venue_entry_or_exit_records` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `action_type` tinyint(1) NOT NULL COMMENT '动作类型 1 进场 2 出场',
  `venue_name` varchar(60) NOT NULL COMMENT '场馆名称',
  `venue_id` bigint(20) NOT NULL COMMENT '场馆ID',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 正常 1 废弃',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='场馆进出场记录表';



# Dump of table venue_info
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_info`;

CREATE TABLE `venue_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '场馆ID',
  `venue_name` varchar(60) NOT NULL COMMENT '场馆名称',
  `address` varchar(100) NOT NULL COMMENT '场馆地址',
  `describe` varchar(1000) NOT NULL COMMENT '场馆介绍',
  `telephone` varchar(60) NOT NULL COMMENT '联系电话',
  `venue_images` text COMMENT '场馆图片 多张逗号分隔',
  `business_hours` varchar(100) NOT NULL COMMENT '营业时间',
  `services` varchar(300) NOT NULL COMMENT '设施及服务',
  `longitude` varchar(30) NOT NULL DEFAULT '' COMMENT '经度',
  `latitude` varchar(30) NOT NULL DEFAULT '' COMMENT '纬度',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 正常 1 废弃',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `service_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 正常营业 1 暂停营业',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='场馆信息表';

LOCK TABLES `venue_info` WRITE;
/*!40000 ALTER TABLE `venue_info` DISABLE KEYS */;

INSERT INTO `venue_info` (`id`, `venue_name`, `address`, `describe`, `telephone`, `venue_images`, `business_hours`, `services`, `longitude`, `latitude`, `status`, `create_at`, `update_at`, `service_status`)
VALUES
	(1,'X-FLY闵行店','上海市闵行区吴中路','爱fly','021-8888888','[\"https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png\",\"https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2d2166e9117a326823b3f38bc225f0088cd3659c115e1d34e3d467f5cc740442.png\",\"https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/1627963949391701.png\",\"https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/diy_ad_img.png\"]','9:00-19:00','停车 茶饮  休息室','121.448815','31.187451',0,0,0,0);

/*!40000 ALTER TABLE `venue_info` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table venue_order_product_info
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_order_product_info`;

CREATE TABLE `venue_order_product_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `pay_order_id` varchar(150) NOT NULL COMMENT '订单号',
  `product_id` bigint(20) NOT NULL COMMENT '商品id',
  `product_type` int(8) NOT NULL COMMENT '1001 场馆预约 2001 购买月卡 2002 购买季卡 2003 购买年卡 2004 体验券 3001 私教（教练）订单 3002 课程订单 4001 充值订单',
  `count` int(11) NOT NULL COMMENT '购买数量',
  `real_amount` int(11) NOT NULL COMMENT '[单个商品]真实价格（单位：分）',
  `cur_amount` int(11) NOT NULL COMMENT '[单个商品]当前价格 (包含真实价格、 折扣价格（单位：分）',
  `discount_rate` int(11) NOT NULL DEFAULT '0' COMMENT '折扣率',
  `discount_amount` int(11) NOT NULL DEFAULT '0' COMMENT '[单个商品]优惠的金额',
  `amount` int(11) NOT NULL COMMENT '商品总价',
  `receive_amount` int(11) NOT NULL DEFAULT '0' COMMENT '充值金额（钱包）',
  `duration` int(11) NOT NULL DEFAULT '0' COMMENT '购买相关服务总时长',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0 待支付 1 订单超时/未支付/已取消 2 已支付 3 已完成  4 退款中 5 已退款 6 已过期',
  `related_id` bigint(20) NOT NULL DEFAULT '1' COMMENT '场馆id/私教课id/大课老师id',
  `deduction_tm` bigint(20) NOT NULL DEFAULT '0' COMMENT '抵扣会员时长',
  `deduction_amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '抵扣的价格',
  `deduction_num` bigint(20) NOT NULL DEFAULT '0' COMMENT '抵扣的数量',
  `expire_duration` int(11) NOT NULL DEFAULT '0' COMMENT '过期时长[单个商品]',
  PRIMARY KEY (`id`),
  KEY `pay_order_id` (`pay_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单商品表';



# Dump of table venue_pay_notify
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_pay_notify`;

CREATE TABLE `venue_pay_notify` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `notify_info` varchar(1000) NOT NULL DEFAULT '' COMMENT '回调通知信息',
  `transaction` varchar(200) NOT NULL DEFAULT '' COMMENT '第三方订单号',
  `pay_type` tinyint(1) NOT NULL COMMENT '1 支付宝 2 微信 ',
  `pay_order_id` varchar(150) NOT NULL COMMENT '订单号',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `notify_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1 付款回调 2 退款回调',
  PRIMARY KEY (`id`),
  KEY `pay_order_id` (`pay_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单回调信息';



# Dump of table venue_pay_orders
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_pay_orders`;

CREATE TABLE `venue_pay_orders` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `amount` int(11) NOT NULL COMMENT '商品总价[优惠后的金额]（分）',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT ' 0 待支付 1 订单超时/未支付/已取消 2 已支付 3 已完成  4 退款中 5 已退款 6 已过期',
  `extra` varchar(1000) NOT NULL DEFAULT '' COMMENT '记录订单相关扩展数据',
  `transaction` varchar(200) NOT NULL DEFAULT '' COMMENT '第三方订单号',
  `pay_type` tinyint(1) NOT NULL COMMENT '1 支付宝 2 微信 3 钱包 4 苹果内购',
  `product_type` int(8) NOT NULL COMMENT '1001 场馆预约 2000 体验券 2001 购买月卡 2002 购买季卡 2003 购买年卡  3001 私教（教练）订单 3002 课程订单 4001 充值订单',
  `error_code` varchar(20) NOT NULL DEFAULT '' COMMENT '错误码',
  `pay_order_id` varchar(150) NOT NULL COMMENT '订单号',
  `order_type` int(8) NOT NULL COMMENT '下单方式：1001 APP下单，1002 前台购买，1003第三方推广渠道购买',
  `pay_time` int(11) NOT NULL DEFAULT '0' COMMENT '用户支付时间',
  `channel_id` int(10) unsigned DEFAULT NULL COMMENT '购买渠道，1001 android ; 1002 ios',
  `is_callback` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否接收到第三方回调 0 未接收到回调 1 已接收回调',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `subject` varchar(150) NOT NULL DEFAULT '' COMMENT '商品名称',
  `write_off_code` varchar(200) NOT NULL DEFAULT '' COMMENT '核销码',
  `refund_amount` int(11) NOT NULL DEFAULT '0' COMMENT '退款金额（分）',
  `is_delete` tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否删除0正常 1删除',
  `original_amount` int(11) NOT NULL COMMENT '订单原始金额（分）',
  `admin_id` int(11) unsigned DEFAULT '0' COMMENT '操作原',
  `refund_fee` int(11) NOT NULL DEFAULT '0' COMMENT '退款手续费（分）',
  PRIMARY KEY (`id`),
  KEY `pay_order_id` (`pay_order_id`),
  KEY `user_id` (`user_id`,`status`,`order_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单流水';



# Dump of table venue_payment_channel
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_payment_channel`;

CREATE TABLE `venue_payment_channel` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL COMMENT '名称',
  `icon` varchar(255) DEFAULT NULL COMMENT 'icon',
  `background` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL COMMENT '状态 0 可用，1禁用',
  `sort` tinyint(4) DEFAULT NULL COMMENT '排序',
  `desc` varchar(255) DEFAULT NULL,
  `rate` int(255) DEFAULT NULL,
  `app_id` varchar(255) DEFAULT NULL,
  `app_key` varchar(255) DEFAULT NULL,
  `app_secret` varchar(255) DEFAULT NULL,
  `public_key` text,
  `private_key` text,
  `create_at` int(11) DEFAULT NULL,
  `update_at` int(11) DEFAULT NULL,
  `pay_type` tinyint(2) DEFAULT NULL COMMENT '1 支付宝 2 微信',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `venue_payment_channel` WRITE;
/*!40000 ALTER TABLE `venue_payment_channel` DISABLE KEYS */;

INSERT INTO `venue_payment_channel` (`id`, `title`, `icon`, `background`, `status`, `sort`, `desc`, `rate`, `app_id`, `app_key`, `app_secret`, `public_key`, `private_key`, `create_at`, `update_at`, `pay_type`)
VALUES
	(1,'支付宝',NULL,NULL,0,NULL,NULL,NULL,'2021002164657197',NULL,NULL,'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAiSQ+YgGoDditE1XP4+RYUFwXySO3O3neZO9JQq/p+GfOMKFP7OYApj9P4w6BvjtEBnRd93PHVwEV2jG+U1CA6/W7cOFZV8sXgNFlUbMqZHjTdrvziPBHAWsUzfMxG6a8WxQ9pl0DOizyokUnbYc3v9HCZq3ZPcJvB9Aiq8fFq02A6DYrX5ZQW97vxCyzumiLIZFu61eDPvkvAwvqrrsn2+LckfKKN2DfcnpmoNnME5VzR3bZv9D5C4IaMeYoz8IyonVPgtnRA1Hth91bjnHqG94EsxQgB/cApgV4F7BXbZ5UbpcrAWCg2PxTy5lIK14iL5Y69StRhk02pxc0J+4mQwIDAQAB','MIIEpAIBAAKCAQEA29Ozn6SdfeqlCmsnNQ/rk6gvZr1jjkhDzOrYJXfOq/MGZOla3ymRSN8+6P0UV0o6QQDINNU/5ZtrJPtJzMOg7bQt8Z1ATidk/OhUjC9sTyBgwXxwjaikB73n9o/fgcgesz1ofkjDpZJeV76Cn5j4JENzaQP3xjvade3B307ReJOzEKVHjyyqfJ25DcGHqYEtdnbuYt9ijfaCfB0oEXJuUxUsIQQALzbfgTCVjdXKHjqCrRFKs1et8k12d0m2xoPTWf2YX72oFbRWqEZD13VhzL/Q2hU32/ENNXbHHrTZtTb3yFHw5Uj67pUmPzNhZFiHek0BkPVvKD89uErvIFpyWQIDAQABAoIBACoRo6iDmlhElX0e8IvpFg5V+2xQBkNudPs8Xk0dVoH1ql2Zgvh+Pf2SK7nu5Puniupxud7SiL3qNmEHbiIvthaHittYWrwaMetskvGZCcNC0QF2TRvvECUjJMc81WtC3w0yTVMNndOL5V4paVodri9ScT3BsqNPRQmYjKetr8zBLII6cpjRbT339RD7Z1FrNSKKQQWPGYH6JAd2sJR+hLHHylSvtpD2IVD3RJgVw11Ge9CKiDmCD6BkHIL5zPe3YErEWAscWsBnjF9sBdTLbkzlF83AIsDipS89dNF/WCCKejCj6A7Rl5kddu5jPSLfwCF6z0YJp7DCAzVZJXg1pgECgYEA9en/Krg9iHU5kP2d58OooHMz73srmtLhJiMINw7098r+1or309MqeSmlB4HGFUdl1D/zd5eQm3VJlixXSmlo5Ew+0xlalL2n+9fNJnbQjEieUFHKdPaBMRW7FQ7pD0UmaTWot1HxZeAt3rkNA0e6oxYPPtOjpswEP/nnvkS17mECgYEA5NfM6xdvB5aITe1uHsIn8KAoBVFuUZu9GgTZgKQgR/JW1QZ5Ba/DCRSK/gmQWbOQHhi6hyLqcUUNaQbuA/1T7ByNCCdk9Mr/zmehyh+BdYj12SQ8i3Af2H45l4xX3JDAXa9l+ba+HuaJ5W4FJcYiynGnG+uNgxNTsI+OckTxVvkCgYEAt0jWZDK5uhEU/NnqbSlJb30twlpdH6H5KYGGx/Kf5mgoFCOznu+Ogovlcnjo+EckwFOB1SrkHtoGJKWb0dxKz418bb5B4waQQ4aOYxK/US92v4qWiSKJG9qEe6eHUVhKzrOtsiSi9TlnNs9ZwY4erxrr9fmryc/Zgw1yCkAQEUECgYEA464hXzUtbmtCqeW0Tj315t4xczkVfXRprF1u2SJyS6K86a1K83FvprUdpKp3SAfzNz57NsByaMe/E+OlI6sDuEKfvqETPMpLwFwzCBpYf0wI7kWzRzgDNy4+tp0XPYd3HL7Jwq0iczQDtpTD4lVDgA+bp5ewb9zmwx/RJbeaNmECgYAaooellJxLAXT0HFsKWmIMn/ckQu+wXcg+sqFOkhKCOMXNnGiADWT5LcyEVcejWQAru3QsvTzPdundFbQR+0hSWDFiPITUopINLATDW0zvSUcbL27Z30HZUXd+LASmU/EC0I6rNbsTZ/nyImpfOul7v2T/BS+eVRjevDrFMGtoeQ==',NULL,NULL,1),
	(2,'微信',NULL,NULL,0,NULL,NULL,NULL,'wxd693805bd4a2a39e','1610121931','UyPTFdVsJPPxMYzfXBztTKwsUusxSIFw','MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAiSQ+YgGoDditE1XP4+RYUFwXySO3O3neZO9JQq/p+GfOMKFP7OYApj9P4w6BvjtEBnRd93PHVwEV2jG+U1CA6/W7cOFZV8sXgNFlUbMqZHjTdrvziPBHAWsUzfMxG6a8WxQ9pl0DOizyokUnbYc3v9HCZq3ZPcJvB9Aiq8fFq02A6DYrX5ZQW97vxCyzumiLIZFu61eDPvkvAwvqrrsn2+LckfKKN2DfcnpmoNnME5VzR3bZv9D5C4IaMeYoz8IyonVPgtnRA1Hth91bjnHqG94EsxQgB/cApgV4F7BXbZ5UbpcrAWCg2PxTy5lIK14iL5Y69StRhk02pxc0J+4mQwIDAQAB','MIIEpAIBAAKCAQEA29Ozn6SdfeqlCmsnNQ/rk6gvZr1jjkhDzOrYJXfOq/MGZOla3ymRSN8+6P0UV0o6QQDINNU/5ZtrJPtJzMOg7bQt8Z1ATidk/OhUjC9sTyBgwXxwjaikB73n9o/fgcgesz1ofkjDpZJeV76Cn5j4JENzaQP3xjvade3B307ReJOzEKVHjyyqfJ25DcGHqYEtdnbuYt9ijfaCfB0oEXJuUxUsIQQALzbfgTCVjdXKHjqCrRFKs1et8k12d0m2xoPTWf2YX72oFbRWqEZD13VhzL/Q2hU32/ENNXbHHrTZtTb3yFHw5Uj67pUmPzNhZFiHek0BkPVvKD89uErvIFpyWQIDAQABAoIBACoRo6iDmlhElX0e8IvpFg5V+2xQBkNudPs8Xk0dVoH1ql2Zgvh+Pf2SK7nu5Puniupxud7SiL3qNmEHbiIvthaHittYWrwaMetskvGZCcNC0QF2TRvvECUjJMc81WtC3w0yTVMNndOL5V4paVodri9ScT3BsqNPRQmYjKetr8zBLII6cpjRbT339RD7Z1FrNSKKQQWPGYH6JAd2sJR+hLHHylSvtpD2IVD3RJgVw11Ge9CKiDmCD6BkHIL5zPe3YErEWAscWsBnjF9sBdTLbkzlF83AIsDipS89dNF/WCCKejCj6A7Rl5kddu5jPSLfwCF6z0YJp7DCAzVZJXg1pgECgYEA9en/Krg9iHU5kP2d58OooHMz73srmtLhJiMINw7098r+1or309MqeSmlB4HGFUdl1D/zd5eQm3VJlixXSmlo5Ew+0xlalL2n+9fNJnbQjEieUFHKdPaBMRW7FQ7pD0UmaTWot1HxZeAt3rkNA0e6oxYPPtOjpswEP/nnvkS17mECgYEA5NfM6xdvB5aITe1uHsIn8KAoBVFuUZu9GgTZgKQgR/JW1QZ5Ba/DCRSK/gmQWbOQHhi6hyLqcUUNaQbuA/1T7ByNCCdk9Mr/zmehyh+BdYj12SQ8i3Af2H45l4xX3JDAXa9l+ba+HuaJ5W4FJcYiynGnG+uNgxNTsI+OckTxVvkCgYEAt0jWZDK5uhEU/NnqbSlJb30twlpdH6H5KYGGx/Kf5mgoFCOznu+Ogovlcnjo+EckwFOB1SrkHtoGJKWb0dxKz418bb5B4waQQ4aOYxK/US92v4qWiSKJG9qEe6eHUVhKzrOtsiSi9TlnNs9ZwY4erxrr9fmryc/Zgw1yCkAQEUECgYEA464hXzUtbmtCqeW0Tj315t4xczkVfXRprF1u2SJyS6K86a1K83FvprUdpKp3SAfzNz57NsByaMe/E+OlI6sDuEKfvqETPMpLwFwzCBpYf0wI7kWzRzgDNy4+tp0XPYd3HL7Jwq0iczQDtpTD4lVDgA+bp5ewb9zmwx/RJbeaNmECgYAaooellJxLAXT0HFsKWmIMn/ckQu+wXcg+sqFOkhKCOMXNnGiADWT5LcyEVcejWQAru3QsvTzPdundFbQR+0hSWDFiPITUopINLATDW0zvSUcbL27Z30HZUXd+LASmU/EC0I6rNbsTZ/nyImpfOul7v2T/BS+eVRjevDrFMGtoeQ==',NULL,NULL,2);

/*!40000 ALTER TABLE `venue_payment_channel` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table venue_personal_label_conf
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_personal_label_conf`;

CREATE TABLE `venue_personal_label_conf` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `label_name` varchar(60) NOT NULL COMMENT '标签名',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0可用 1不可用',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户个人标签配置表';

LOCK TABLES `venue_personal_label_conf` WRITE;
/*!40000 ALTER TABLE `venue_personal_label_conf` DISABLE KEYS */;

INSERT INTO `venue_personal_label_conf` (`id`, `label_name`, `status`, `create_at`, `update_at`)
VALUES
	(1,'可可爱爱',0,0,0),
	(2,'没有脑袋',0,0,0),
	(3,'帅就完事儿了',0,0,0),
	(4,'风流地毯',0,0,0),
	(5,'高高手',0,0,0);

/*!40000 ALTER TABLE `venue_personal_label_conf` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table venue_product_info
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_product_info`;

CREATE TABLE `venue_product_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `product_name` varchar(128) NOT NULL COMMENT '商品名称',
  `product_type` int(8) NOT NULL COMMENT '类型 2000体验券 2001 月卡 2002 季卡 2003 年卡 3001 储值卡',
  `product_code` varchar(255) DEFAULT NULL COMMENT '商品码',
  `real_amount` int(11) NOT NULL COMMENT '真实价格（单位：分）',
  `cur_amount` int(11) NOT NULL COMMENT '当前价格 (包含真实价格、 折扣价格（单位：分）',
  `discount_rate` int(11) NOT NULL DEFAULT '0' COMMENT '折扣率',
  `discount_amount` int(11) NOT NULL DEFAULT '0' COMMENT '优惠的金额',
  `venue_id` bigint(20) NOT NULL COMMENT '场馆id',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  `effective_duration` int(11) NOT NULL DEFAULT '0' COMMENT '可用时长（秒）',
  `expire_duration` int(11) NOT NULL DEFAULT '0' COMMENT '过期时长（秒）',
  `icon` varchar(1000) NOT NULL DEFAULT '' COMMENT '商品icon',
  `image` varchar(1000) NOT NULL DEFAULT '' COMMENT '商品图片',
  `describe` varchar(1000) NOT NULL DEFAULT '' COMMENT '商品介绍',
  `title` varchar(300) NOT NULL DEFAULT '' COMMENT '简介',
  `instance_type` tinyint(4) unsigned DEFAULT '1' COMMENT '实例类型，1: 体验卡；2: 线下食品',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='场馆商品配置表';

LOCK TABLES `venue_product_info` WRITE;
/*!40000 ALTER TABLE `venue_product_info` DISABLE KEYS */;

INSERT INTO `venue_product_info` (`id`, `product_name`, `product_type`, `product_code`, `real_amount`, `cur_amount`, `discount_rate`, `discount_amount`, `venue_id`, `create_at`, `update_at`, `effective_duration`, `expire_duration`, `icon`, `image`, `describe`, `title`, `instance_type`)
VALUES
	(1,'x-fly月卡',2001,NULL,10000,9800,98,200,1,0,0,3600,4000,'https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png','','月卡',1),
	(2,'x-fly年卡',2003,NULL,50000,39800,0,0,1,0,0,3600,4000,'https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png','','年卡',1),
	(3,'x-fly季卡',2002,NULL,20000,19800,0,0,1,0,0,3600,4000,'https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png','','季卡',1),
	(4,'x-fly体验卡',2000,NULL,20000,19800,0,0,1,0,0,3600,4000,'https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png','','体验卡',1),
	(5,'x-fly储值卡',4001,NULL,20000,19800,0,0,1,0,0,3600,4000,'https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png','https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/2b806baf6ff19d086bed3b72c444194ee3ab69e9777e9694c3688a6414cbb7d6.png','','体验卡',1);

/*!40000 ALTER TABLE `venue_product_info` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table venue_recharge_config
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_recharge_config`;

CREATE TABLE `venue_recharge_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `product_name` varchar(128) NOT NULL COMMENT '商品名称',
  `price` int(11) NOT NULL COMMENT '商品价格（软妹币：分）',
  `receive_amount` int(11) NOT NULL COMMENT '充值可得金额(分)',
  `platform_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 android端 1 iOS端',
  `is_recommend` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否推荐 0 不推荐 1 推荐',
  `ios_product_name` varchar(100) NOT NULL DEFAULT '' COMMENT 'iOS商品名称',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='充值配置表';



# Dump of table venue_recommend_conf
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_recommend_conf`;

CREATE TABLE `venue_recommend_conf` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `name` varchar(60) NOT NULL COMMENT '推荐名称',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 有效 1 废弃',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='场馆推荐角标配置表';

LOCK TABLES `venue_recommend_conf` WRITE;
/*!40000 ALTER TABLE `venue_recommend_conf` DISABLE KEYS */;

INSERT INTO `venue_recommend_conf` (`id`, `name`, `status`, `create_at`, `update_at`)
VALUES
	(1,'低价推荐',0,0,0),
	(2,'高手推荐',0,0,0),
	(3,'美女推荐',0,0,0);

/*!40000 ALTER TABLE `venue_recommend_conf` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table venue_refund_rules
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_refund_rules`;

CREATE TABLE `venue_refund_rules` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `rule_order` tinyint(2) NOT NULL DEFAULT '0' COMMENT '规则校验顺序',
  `rule_name` varchar(255) NOT NULL COMMENT '规则名称',
  `rule_min_duration` int(11) NOT NULL DEFAULT '0' COMMENT '规则校验最小时长',
  `rule_max_duration` int(11) NOT NULL DEFAULT '0' COMMENT '规则校验最大时长',
  `fee_rate` int(5) NOT NULL DEFAULT '0' COMMENT '手续费比例 例如：1.55% 则 值为155 [乘以100存储]',
  `minimum_charge` int(10) NOT NULL DEFAULT '0' COMMENT '最低收取手续费金额 单位[分]',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 正常 1 废弃',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `fee_rate_cn` varchar(255) NOT NULL DEFAULT '' COMMENT '手续费比例 中文 例如：票价5%',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='场馆退款规则[手续费]表';

LOCK TABLES `venue_refund_rules` WRITE;
/*!40000 ALTER TABLE `venue_refund_rules` DISABLE KEYS */;

INSERT INTO `venue_refund_rules` (`id`, `rule_order`, `rule_name`, `rule_min_duration`, `rule_max_duration`, `fee_rate`, `minimum_charge`, `status`, `create_at`, `update_at`, `fee_rate_cn`)
VALUES
	(1,1,'小于24小时',0,86400,2000,100,0,0,0,'票价20%'),
	(2,2,'24-48小时（不含）',86400,172800,1000,100,0,0,0,'票价10%'),
	(3,3,'2天-5天（不含）',172800,432000,500,100,0,0,0,'票价5%'),
	(4,4,'5天以上',432000,0,0,100,0,0,0,'无手续费');

/*!40000 ALTER TABLE `venue_refund_rules` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table venue_user_cardbag
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_user_cardbag`;

CREATE TABLE `venue_user_cardbag` (
  `id` int(11) DEFAULT NULL,
  `user_id` varchar(255) DEFAULT NULL COMMENT '用户ID',
  `qr_code` varchar(255) DEFAULT NULL COMMENT 'ID卡号',
  `card_type` tinyint(4) DEFAULT NULL COMMENT '卡类型:1、年卡、2半年卡、3季卡、4月卡、5、此卡',
  `admin_id` int(11) unsigned DEFAULT '0' COMMENT '管理员ID',
  `status` tinyint(4) unsigned DEFAULT '0' COMMENT '可用状态0可用、1禁用',
  `create_at` int(11) unsigned DEFAULT '0' COMMENT '创建时间',
  `duration` int(11) unsigned DEFAULT '0' COMMENT '可用时长',
  `update_at` int(11) unsigned DEFAULT '0' COMMENT '修改时间',
  `dt_start` int(11) unsigned DEFAULT '0' COMMENT '启用时间',
  `dt_end` int(11) unsigned DEFAULT '0' COMMENT '过期时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table venue_user_evaluate_record
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_user_evaluate_record`;

CREATE TABLE `venue_user_evaluate_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `coach_id` bigint(20) NOT NULL COMMENT '教练id',
  `star` tinyint(2) NOT NULL COMMENT '评价几颗星 1~5星',
  `order_id` varchar(256) NOT NULL COMMENT '订单id',
  `order_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '订单类型 默认 1 私教课程',
  `content` varchar(1000) DEFAULT '' COMMENT '评价描述',
  `label_info` varchar(1000) NOT NULL DEFAULT '' COMMENT '用户选取的评价标签信息',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 有效 1 废弃',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `coach_id` (`coach_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户对私教评价记录';



# Dump of table venue_user_label
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_user_label`;

CREATE TABLE `venue_user_label` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `date` varchar(128) NOT NULL DEFAULT '' COMMENT '日期 例如 2021-10-01',
  `time_node` varchar(128) NOT NULL DEFAULT '' COMMENT '时间节点 例如 10:00-12:00',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `label_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0为用户添加标签 1为系统添加标签',
  `label_id` bigint(20) NOT NULL COMMENT '标签id',
  `label_name` varchar(60) NOT NULL COMMENT '标签名',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 有效 1 废弃',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `venue_id` bigint(20) NOT NULL COMMENT '场馆id',
  `pay_order_id` varchar(150) NOT NULL DEFAULT '' COMMENT '订单id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='场馆用户标签表';



# Dump of table venue_vip_info
# ------------------------------------------------------------

DROP TABLE IF EXISTS `venue_vip_info`;

CREATE TABLE `venue_vip_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `level` int(2) NOT NULL DEFAULT '0' COMMENT '会员等级 预留字段',
  `venue_id` bigint(20) NOT NULL COMMENT '场馆id',
  `start_tm` bigint(20) NOT NULL COMMENT '会员开始时间戳',
  `end_tm` bigint(20) NOT NULL COMMENT '会员结束时间戳',
  `create_at` int(11) NOT NULL DEFAULT '0',
  `update_at` int(11) NOT NULL DEFAULT '0',
  `duration` bigint(20) NOT NULL COMMENT '会员在场馆内可用时长',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
