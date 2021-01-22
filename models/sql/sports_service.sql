# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.0.12)
# Database: fpv2
# Generation Time: 2021-01-21 06:35:40 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table admin_user
# ------------------------------------------------------------

CREATE TABLE `admin_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `phone` char(11) DEFAULT '' COMMENT '手机号',
  `username` varchar(255) NOT NULL COMMENT '用户名',
  `password` char(32) NOT NULL COMMENT '登录密码',
  `salt` char(10) NOT NULL COMMENT 'salt',
  `sub_account` bigint(11) NOT NULL DEFAULT '0' COMMENT '子账号标识：0则为主账号，大于0则为子账号，值等于主账号ID',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  `create_at` int(11) NOT NULL COMMENT '注册时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `admin_user` WRITE;
/*!40000 ALTER TABLE `admin_user` DISABLE KEYS */;

INSERT INTO `admin_user` (`id`, `phone`, `username`, `password`, `salt`, `sub_account`, `update_at`, `create_at`)
VALUES
	(4,'','admin','de8b658e4f6a6d1e5e923e951e0f1407','XAVLsyWU',0,1603426460,1603426460);

/*!40000 ALTER TABLE `admin_user` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table app_version_control
# ------------------------------------------------------------

CREATE TABLE `app_version_control` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `version_name` varchar(128) NOT NULL DEFAULT '' COMMENT '版本名称',
  `version` varchar(60) NOT NULL DEFAULT '' COMMENT '版本号',
  `version_code` int(8) NOT NULL DEFAULT '0' COMMENT '版本code',
  `size` varchar(128) NOT NULL DEFAULT '' COMMENT '包大小',
  `is_force` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 不需要强更 1 需要强更',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 可用 1 不可用',
  `platform` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 android 1 ios',
  `upgrade_url` varchar(256) NOT NULL DEFAULT '' COMMENT '更新包地址',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `describe` varchar(500) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '版本说明',
  PRIMARY KEY (`id`),
  KEY `version_code` (`version_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='app版本控制';



# Dump of table banner
# ------------------------------------------------------------

CREATE TABLE `banner` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标题',
  `cover` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'banner封面',
  `explain` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '说明',
  `jump_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '跳转地址',
  `share_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '分享地址',
  `type` int(1) NOT NULL DEFAULT '1' COMMENT '1 首页 2 直播页 3 官网banner',
  `start_time` int(11) NOT NULL DEFAULT '0' COMMENT '上架时间',
  `end_time` int(11) NOT NULL DEFAULT '0' COMMENT '下架时间',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0待上架 1上架 2 已过期',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `jump_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0站内跳转 1站外跳转 2 无跳转',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='banner配置表';

LOCK TABLES `banner` WRITE;
/*!40000 ALTER TABLE `banner` DISABLE KEYS */;

INSERT INTO `banner` (`id`, `title`, `cover`, `explain`, `jump_url`, `share_url`, `type`, `start_time`, `end_time`, `sortorder`, `status`, `create_at`, `update_at`, `jump_type`)
VALUES
	(8,'首页','http://192.168.50.148:13002/upload/2020_10_23/539517102248169472.jpeg','说点啥','www.baidu.com','www.google.com',1,1603433790,1603865793,10,2,1603433797,1603433797,0);

/*!40000 ALTER TABLE `banner` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table browse_detail_record
# ------------------------------------------------------------

CREATE TABLE `browse_detail_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `study_duration` int(8) NOT NULL COMMENT '当前播放的时长（秒）',
  `total_duration` int(8) NOT NULL COMMENT '当前视频总时长（秒）',
  `cur_progress` decimal(5,2) NOT NULL DEFAULT '0.00' COMMENT '当前播放的进度',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `video_id` (`video_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户观看视频的详情记录表';



# Dump of table circle_attention
# ------------------------------------------------------------

CREATE TABLE `circle_attention` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(60) NOT NULL COMMENT '关注圈子的用户id',
  `circle_id` int(11) NOT NULL COMMENT '圈子id',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1表示关注 2表示取消关注',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `circle_id` (`circle_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='关注的圈子';



# Dump of table collect_record
# ------------------------------------------------------------

CREATE TABLE `collect_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户id',
  `to_user_id` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '作品发布者用户id',
  `compose_id` bigint(20) NOT NULL COMMENT '作品id',
  `compose_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0 视频 1 帖子',
  `status` tinyint(1) NOT NULL COMMENT '1 收藏 0 取消收藏',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `video_id` (`compose_id`),
  KEY `to_user_id` (`to_user_id`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='收藏的视频';



# Dump of table comment_report
# ------------------------------------------------------------

CREATE TABLE `comment_report` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(60) NOT NULL DEFAULT '' COMMENT '举报人用户id',
  `comment_id` bigint(20) NOT NULL COMMENT '评论id',
  `reason` varchar(300) NOT NULL COMMENT '举报理由',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `comment_id` (`comment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频举报';



# Dump of table complaint
# ------------------------------------------------------------

CREATE TABLE `complaint` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(60) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '举报人',
  `to_uid` varchar(60) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '被举报人',
  `reason` mediumtext COLLATE utf8mb4_unicode_ci COMMENT '原因',
  `compose_id` bigint(20) NOT NULL COMMENT '举报的作品id（视频/帖子id）',
  `complaint_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0 举报其他 1 举报视频 2 举报帖子',
  `cover` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '图片地址  逗号分隔',
  `is_dispose` tinyint(1) DEFAULT '1' COMMENT '是否受理 1未受理 2受理',
  `content` mediumtext COLLATE utf8mb4_unicode_ci COMMENT '回复内容',
  `create_at` int(11) NOT NULL,
  `update_at` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='举报信息';



# Dump of table default_avatar
# ------------------------------------------------------------

CREATE TABLE `default_avatar` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `avatar` varchar(128) NOT NULL COMMENT '头像地址',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0展示 1不展示',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='系统默认头像';

LOCK TABLES `default_avatar` WRITE;
/*!40000 ALTER TABLE `default_avatar` DISABLE KEYS */;

INSERT INTO `default_avatar` (`id`, `avatar`, `sortorder`, `create_at`, `update_at`, `status`)
VALUES
	(1,'https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=2863307160,2081500909&fm=26&gp=0.jpg?imageView2/1/w/80/h/80',100,1600426884,1600426884,0),
	(6,'http://192.168.50.148:13002/upload/2020_10_23/539519399892094976.png',2,1603434339,1603434339,0),
	(5,'http://192.168.50.148:13002/upload/2020_10_23/539519314500259840.png',1,1603434317,1603434317,0),
	(7,'http://192.168.50.148:13002/upload/2020_10_23/539519450638979072.png',3,1603434352,1603434352,0),
	(8,'http://192.168.50.148:13002/upload/2020_10_23/539519490564558848.png',4,1603434361,1603434361,0);

/*!40000 ALTER TABLE `default_avatar` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table feedback
# ------------------------------------------------------------

CREATE TABLE `feedback` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `phone` varchar(200) DEFAULT NULL COMMENT '手机号码',
  `describe` mediumtext COMMENT '描述问题内容',
  `problem` varchar(500) NOT NULL DEFAULT '' COMMENT '遇到的问题',
  `contact` varchar(200) DEFAULT NULL COMMENT '联系方式',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 0未回复 1已回复',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `content` mediumtext COMMENT '回复内容',
  `pics` varchar(512) DEFAULT NULL COMMENT '上传的图片，多张逗号分隔',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='反馈表';



# Dump of table hot_circle
# ------------------------------------------------------------

CREATE TABLE `hot_circle` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `hot_circle_id` varchar(128) NOT NULL COMMENT '热门圈子id 多个用逗号分隔 例如：3,6,11,21',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='热门圈子（后台直接配置）';



# Dump of table hot_search
# ------------------------------------------------------------

CREATE TABLE `hot_search` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `hot_search_content` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '热门搜索内容 如：FPV、电竞',
  `status` tinyint(1) DEFAULT '0' COMMENT '0 展示 1 隐藏',
  `sortorder` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '排序权重',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='热门搜索（后台配置）';

LOCK TABLES `hot_search` WRITE;
/*!40000 ALTER TABLE `hot_search` DISABLE KEYS */;

INSERT INTO `hot_search` (`id`, `hot_search_content`, `status`, `sortorder`, `create_at`, `update_at`)
VALUES
	(1,'电竞',0,100,1600000000,1603339947),
	(3,'FPV',0,1000000,1600000000,1603778803),
	(5,'拼装无人机',0,1001,1600000000,1603264178),
	(55,'电竞比赛',0,1111,1603269824,1603269824);

/*!40000 ALTER TABLE `hot_search` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table post_comment
# ------------------------------------------------------------

CREATE TABLE `post_comment` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '评论id',
  `user_id` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '评论人userId',
  `user_name` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '评论人名称',
  `avatar` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像',
  `post_id` bigint(20) NOT NULL COMMENT '帖子id',
  `parent_comment_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '父评论id',
  `parent_comment_user_id` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '父评论的用户id',
  `reply_comment_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '被回复的评论id',
  `reply_comment_user_id` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '被回复的评论用户id',
  `comment_level` tinyint(4) NOT NULL DEFAULT '1' COMMENT '评论等级[ 1 一级评论 默认 ，2 二级评论]',
  `content` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '评论的内容',
  `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '状态 (1 有效，0 逻辑删除)',
  `is_top` tinyint(2) NOT NULL DEFAULT '0' COMMENT '置顶状态[ 1 置顶，0 不置顶 默认 ]',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_post_id` (`post_id`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE,
  KEY `idx_create_time` (`create_at`),
  KEY `idx_parent_comment_id` (`parent_comment_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='帖子评论表';



# Dump of table post_info
# ------------------------------------------------------------

CREATE TABLE `post_info` (
  `post_id` bigint(20) unsigned NOT NULL COMMENT '帖子id',
  `title` mediumtext COMMENT '帖子标题',
  `describe` mediumtext COMMENT '帖子描述',
  `cover` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '封面',
  `video_addr` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频地址',
  `circle_id` int(11) NOT NULL COMMENT '圈子id',
  `section_id` int(11) NOT NULL COMMENT '板块id',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `user_type` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '添加用户类型（0：管理员[sys_user]，1：用户[user]）',
  `sortorder` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '权重',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '状态（0：展示 1：隐藏 2：删除 预留状态）',
  `is_recommend` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '推荐（0：不推荐；1：推荐）',
  `is_top` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '置顶（0：不置顶；1：置顶；）',
  `create_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `video_duration` int(8) NOT NULL DEFAULT '0' COMMENT '视频时长（单位：秒）',
  `rec_content` mediumtext COMMENT '推荐理由',
  `top_content` mediumtext COMMENT '置顶理由',
  `video_width` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频宽',
  `video_height` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频高',
  PRIMARY KEY (`post_id`),
  KEY `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='帖子表';



# Dump of table post_info_examine
# ------------------------------------------------------------

CREATE TABLE `post_info_examine` (
  `post_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '帖子id',
  `title` mediumtext COMMENT '帖子标题',
  `describe` mediumtext COMMENT '帖子描述',
  `cover` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '封面',
  `video_addr` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频地址',
  `circle_id` int(11) NOT NULL COMMENT '圈子id',
  `section_id` int(11) NOT NULL COMMENT '板块id',
  `label_id` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '标签id',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `user_type` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '添加用户类型（0：管理员[sys_user]，1：用户[user]）',
  `sortorder` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '权重',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '审核状态（0：无操作，1：审核通过，2：审核不通过 3: 删除）',
  `is_recommend` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '推荐（0：不推荐；1：推荐）',
  `is_top` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '置顶（0：不置顶；1：置顶；）',
  `manager` int(11) NOT NULL DEFAULT '0' COMMENT '后台操作用户',
  `video_duration` int(8) NOT NULL DEFAULT '0' COMMENT '视频时长（单位：秒）',
  `label_name` mediumtext COMMENT '标签',
  `rec_content` mediumtext COMMENT '推荐理由',
  `top_content` mediumtext COMMENT '置顶理由',
  `video_width` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频宽',
  `video_height` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频高',
  `create_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`post_id`),
  KEY `user_id` (`user_id`) USING BTREE,
  FULLTEXT KEY `label_name` (`label_name`),
  FULLTEXT KEY `label_id` (`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='帖子审核表';



# Dump of table post_label_config
# ------------------------------------------------------------

CREATE TABLE `post_label_config` (
  `label_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '标签id',
  `pid` int(11) NOT NULL COMMENT '父类id 0为1级分类',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` int(1) NOT NULL DEFAULT '1' COMMENT '类别状态1-正常,2-已废弃',
  `label_name` varchar(64) NOT NULL COMMENT '标签名称',
  `icon` varchar(256) NOT NULL DEFAULT '' COMMENT '标签icon',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`label_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COMMENT='帖子标签配置';



# Dump of table post_labels
# ------------------------------------------------------------

CREATE TABLE `post_labels` (
  `post_id` bigint(20) unsigned NOT NULL COMMENT '帖子id',
  `label_id` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标签id',
  `label_name` varchar(521) DEFAULT NULL COMMENT '标签名',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`post_id`,`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='帖子拥有的标签表';



# Dump of table post_statistic
# ------------------------------------------------------------

CREATE TABLE `post_statistic` (
  `post_id` bigint(20) unsigned NOT NULL COMMENT '帖子id',
  `fabulous` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '点赞数',
  `browse` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '浏览数',
  `share` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '分享数',
  `reward` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '打赏的游币数',
  `create_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='帖子相关参数统计';



# Dump of table received_at
# ------------------------------------------------------------

CREATE TABLE `received_at` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `to_user_id` varchar(60) NOT NULL COMMENT '被@的用户id',
  `user_id` varchar(60) NOT NULL COMMENT '执行@的用户id',
  `comment_id` bigint(20) NOT NULL COMMENT '评论id',
  `topic_type` tinyint(2) NOT NULL COMMENT '1.视频 2.帖子 3.评论',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `comment_level` tinyint(4) NOT NULL DEFAULT '1' COMMENT '评论等级[ 1 一级评论 默认 ，2 二级评论]',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `to_user_id` (`to_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收到的@';

LOCK TABLES `received_at` WRITE;
/*!40000 ALTER TABLE `received_at` DISABLE KEYS */;

INSERT INTO `received_at` (`id`, `to_user_id`, `user_id`, `comment_id`, `topic_type`, `create_at`, `comment_level`)
VALUES
	(1,'202009181548217779','202009181548217779',1,3,1600859121,2),
	(2,'202009181548217779','202009181548217779',43,3,1600859134,2),
	(3,'202009181548217779','202009181548217779',44,3,1600859397,2),
	(4,'202009181548217779','202009181548217779',45,3,1600859494,2),
	(5,'202009181548217779','202009181548217779',46,3,1600859516,2),
	(6,'202009181548217779','202009181548217779',47,3,1600859526,2),
	(7,'1234','202009181548217779',22,3,1600859695,2),
	(8,'202009181548217779','202009181548217779',48,3,1600859723,2),
	(9,'1234','202009181548217779',20,3,1600859877,2),
	(10,'1234','202009181548217779',21,3,1600860029,2),
	(11,'1234','202009181548217779',22,3,1600860443,2),
	(12,'202009181548217779','202009181548217779',49,3,1600860532,0),
	(13,'202009181548217779','202009181548217779',50,3,1600860630,0),
	(15,'202009181548217779','202009181548217779',51,3,1601193324,2),
	(16,'202009181548217779','202009181548217779',52,3,1601193343,2),
	(17,'202009181548217779','202009181548217779',53,3,1601193578,2),
	(18,'202009181548217779','202009181548217779',54,3,1601193594,2),
	(19,'202009181548217779','202009181548217779',55,3,1601193607,2),
	(20,'202009181548217779','202009181548217779',31,3,1601193614,2),
	(21,'202009181548217779','202009181548217779',32,3,1601193619,2),
	(22,'202009181548217779','202009181548217779',33,3,1601193625,2),
	(23,'202009181548217779','202009101933004667',57,3,1604494389,1),
	(24,'202009181548217779','202009101933004667',58,3,1604494436,1),
	(25,'202009101933004667','202009181548217779',59,3,1604494436,2),
	(26,'202009181548217779','202009101933004667',60,3,1604494610,1),
	(27,'202009101933004667','202009181548217779',61,3,1604494610,2),
	(28,'202009181548217779','202009101933004667',62,3,1605236118,1),
	(29,'202009101933004667','202009181548217779',63,3,1605236118,2);

/*!40000 ALTER TABLE `received_at` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table search_history
# ------------------------------------------------------------

CREATE TABLE `search_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `search_content` varchar(128) NOT NULL COMMENT '搜索的内容',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1 正常 2 已删除',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户搜索历史表';



# Dump of table section_config
# ------------------------------------------------------------

CREATE TABLE `section_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `circle_id` int(11) NOT NULL COMMENT '圈子id',
  `section_name` varchar(100) NOT NULL COMMENT '板块名称',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 未操作 1 展示  2 隐藏',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `circle_id` (`circle_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='所属圈子的板块配置表（后台设置）';



# Dump of table share_record
# ------------------------------------------------------------

CREATE TABLE `share_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `compose_id` bigint(20) NOT NULL COMMENT '作品id（视频/帖子id）',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分享的整体内容（json）',
  `share_type` tinyint(2) NOT NULL COMMENT '分享类型 1 分享视频 2 分享帖子',
  `share_platform` tinyint(2) NOT NULL COMMENT '分享平台 1 微信 2 微博 3 qq',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0展示 1不展示',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户分享记录';



# Dump of table social_account_login
# ------------------------------------------------------------

CREATE TABLE `social_account_login` (
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `unionid` varchar(256) NOT NULL DEFAULT '' COMMENT '社交平台关联id',
  `social_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '区分社交软件 1 微信关联id 2 微博关联id 3 qq关联id',
  `status` tinyint(1) DEFAULT '0' COMMENT '0 正常 1 封禁',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`user_id`,`social_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='社交平台登陆表';



# Dump of table social_circle
# ------------------------------------------------------------

CREATE TABLE `social_circle` (
  `circle_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '圈子id',
  `circle_name` varchar(100) NOT NULL COMMENT '圈子名称',
  `cover` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '圈子封面',
  `describe` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '描述',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 未操作 1 展示  2 隐藏',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`circle_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='圈子表';



# Dump of table system_log
# ------------------------------------------------------------

CREATE TABLE `system_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sys_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '系统账号ID',
  `sys_user` varchar(255) DEFAULT '' COMMENT '用户昵称',
  `sys_role` varchar(255) DEFAULT '' COMMENT '用户角色',
  `log_type` varchar(200) DEFAULT '' COMMENT '记录类型',
  `log_cont` mediumtext COMMENT '操作内容',
  `log_text` mediumtext COMMENT '备注',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统记录';



# Dump of table system_message
# ------------------------------------------------------------

CREATE TABLE `system_message` (
  `system_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '系统通知ID',
  `send_id` varchar(60) NOT NULL DEFAULT '' COMMENT '发送者ID（后台用户）',
  `receive_id` varchar(60) NOT NULL DEFAULT '' COMMENT '接收者id',
  `cover` varchar(256) NOT NULL DEFAULT '' COMMENT '消息封面',
  `send_default` tinyint(2) NOT NULL DEFAULT '0' COMMENT '1时发送所有用户，0时则不采用',
  `system_topic` mediumtext CHARACTER SET utf8mb4 NOT NULL COMMENT '通知标题',
  `system_content` mediumtext CHARACTER SET utf8mb4 NOT NULL COMMENT '通知内容',
  `send_time` int(11) NOT NULL COMMENT '发送时间',
  `expire_time` int(11) NOT NULL DEFAULT '0' COMMENT '过期时间',
  `send_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0.默认为系统通知 1 活动通知  2 待支付订单延时提示消息（15分钟 用户端） 3. 待咨询订单通知(开始前1天 用户端及咨询师端) 4.待咨询订单通知(开始前1小时 用户端及咨询师端) 5. 咨询师写评估报告通知（结束后1小时 咨询师端）6. 咨询师写评估报告通知（结束后24小时 咨询师端）',
  `extra` varchar(1024) NOT NULL DEFAULT ' ' COMMENT '附件内容 例如：奖励',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 未读 1 已读  默认未读',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `android_task_id` varchar(256) CHARACTER SET utf8mb4 COLLATE NOT NULL DEFAULT '' COMMENT '友盟任务id[android端]',
  `ios_task_id` varchar(256) CHARACTER SET utf8mb4 COLLATE NOT NULL DEFAULT '' COMMENT '友盟任务id[ios端]',
  `umeng_platform` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 所有 1 android 2 ios',
  `send_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 已发送 1 未发送 2 已撤回 3 已删除',
  PRIMARY KEY (`system_id`),
  KEY `receive_id` (`receive_id`),
  KEY `send_type` (`send_type`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统消息';



# Dump of table system_notice_settings
# ------------------------------------------------------------

CREATE TABLE `system_notice_settings` (
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `comment_push_set` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '状态（0：接收推送；1：拒绝推送 包含评论/回复推送）',
  `thumb_up_push_set` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '状态（0：接收推送；1：拒绝推送 点赞推送）',
  `attention_push_set` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '状态（0：接收推送；1：拒绝推送 关注推送）',
  `share_push_set` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '状态（0：接收推送；1：拒绝推送 包含评论/回复）',
  `slot_push_set` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '状态（0：接收推送；1：拒绝推送 投币推送）',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='用户系统推送通知设置表';

LOCK TABLES `system_notice_settings` WRITE;
/*!40000 ALTER TABLE `system_notice_settings` DISABLE KEYS */;

INSERT INTO `system_notice_settings` (`user_id`, `comment_push_set`, `thumb_up_push_set`, `attention_push_set`, `share_push_set`, `slot_push_set`, `create_at`, `update_at`)
VALUES
	('202009181540233466',0,0,0,0,0,1600414824,1600414824),
	('202009181543188713',0,0,0,0,0,1600414998,1600414998),
	('202009181544573421',0,0,0,0,0,1600415097,1600415097),
	('202009181548217779',1,1,0,1,1,1600415301,1605184836),
	('202010101545291936',0,0,0,0,0,1602315929,1602315929);

/*!40000 ALTER TABLE `system_notice_settings` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table system_push
# ------------------------------------------------------------

CREATE TABLE `system_push` (
  `push_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '系统推送ID',
  `send_id` int(11) NOT NULL COMMENT '发送者ID（后台用户id）',
  `receive_id` varchar(60) NOT NULL COMMENT '接收者id',
  `send_default` tinyint(2) NOT NULL COMMENT '1时发送所有用户，0时则不采用',
  `push_topic` varchar(60) NOT NULL COMMENT '推送标题',
  `push_content` varchar(255) NOT NULL COMMENT '推送内容',
  `send_time` int(11) NOT NULL COMMENT '发送时间',
  `send_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0.默认为系统推送 1.收到@推送 2.收到点赞推送 3.收到收藏推送 4.收到分享推送 5.收到评论/回复推送 6.特殊业务的系统奖励推送 7.活动推送',
  `extra` varchar(2056) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT ' ' COMMENT '透传数据 ',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 成功 1 失败',
  PRIMARY KEY (`push_id`),
  KEY `receive_id` (`receive_id`),
  KEY `send_type` (`send_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统推送';



# Dump of table tencent_cloud_events
# ------------------------------------------------------------

CREATE TABLE `tencent_cloud_events` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `file_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '腾讯云文件id',
  `event` mediumtext COMMENT '事件内容（json字符串）',
  `create_at` int(11) NOT NULL,
  `event_type` tinyint(3) NOT NULL DEFAULT '0' COMMENT '事件类型 0 视频上传事件 1 视频转码事件',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='腾讯事件回调';

LOCK TABLES `tencent_cloud_events` WRITE;
/*!40000 ALTER TABLE `tencent_cloud_events` DISABLE KEYS */;

INSERT INTO `tencent_cloud_events` (`id`, `file_id`, `event`, `create_at`, `event_type`)
VALUES
	(1,5285890809593139008,'{\"EventHandle\":\"8095179751097087535\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-eea7cd24a4357257272e4dea1bf703b0t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809593139008\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/e3caf6115285890809593139008/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e3caf6115285890809593139008/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455525,\"Height\":540,\"Width\":960,\"Size\":2972305,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"68e6714b61c2e00814053d8a5e397515\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385018,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e3caf6115285890809593139008/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e3caf6115285890809593139008/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e3caf6115285890809593139008/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605704631,1),
	(2,5285890809595400796,'{\"EventHandle\":\"8095023532436265478\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-0d674dbffc359f2c3ea2248aad22b885t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809595400796\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/8483cf785285890809595400796/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8483cf785285890809595400796/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455605,\"Height\":540,\"Width\":960,\"Size\":2972823,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"558256a06d07a10ea694ca5dc9d7b2fc\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385098,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8483cf785285890809595400796/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8483cf785285890809595400796/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8483cf785285890809595400796/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605704631,1),
	(3,5285890809593139036,'{\"EventHandle\":\"8094950583651013973\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-802f77747a1616287282e6ceeae98282t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809593139036\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/e3caf6725285890809593139036/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e3caf6725285890809593139036/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455605,\"Height\":540,\"Width\":960,\"Size\":2972825,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"8f1da0e1b19c9681e1339915223e43f6\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385098,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e3caf6725285890809593139036/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e3caf6725285890809593139036/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e3caf6725285890809593139036/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605704631,1),
	(4,5285890809595400848,'{\"EventHandle\":\"8094953596110117584\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-5b0cac60573e52ddc807b4e9e9beb269t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809595400848\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/8483d3165285890809595400848/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8483d3165285890809595400848/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455572,\"Height\":540,\"Width\":960,\"Size\":2972611,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"50bb63b805aa8407fa334fd01e0cef5b\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385065,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8483d3165285890809595400848/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8483d3165285890809595400848/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8483d3165285890809595400848/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605704690,1),
	(5,5285890809589814946,'{\"EventHandle\":\"8095104471295143232\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-9dfa918adc94415c755de322862b4775t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809589814946\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d593e2015285890809589814946/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d593e2015285890809589814946/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455581,\"Height\":540,\"Width\":960,\"Size\":2972667,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"40edcb7c3ccddb56d387a889b8881662\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385074,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d593e2015285890809589814946/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d593e2015285890809589814946/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d593e2015285890809589814946/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605751856,1),
	(6,5285890809595472555,'{\"EventHandle\":\"8095204083927343729\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-56e6744fe42a8bda1b15c01a0bb187f7t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809595472555\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/85038aba5285890809595472555/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/85038aba5285890809595472555/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455570,\"Height\":540,\"Width\":960,\"Size\":2972596,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"2ce59f943ded7f85af77251d40c1a7ff\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385063,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/85038aba5285890809595472555/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/85038aba5285890809595472555/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/85038aba5285890809595472555/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605759835,1),
	(7,5285890809589811843,'{\"EventHandle\":\"8095151089185686028\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-dd0b60697391dc24757dafa921fe2f0bt0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809589811843\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592389a5285890809589811843/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592389a5285890809589811843/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455558,\"Height\":540,\"Width\":960,\"Size\":2972519,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"9310179d0f3aa6b45a9cf58c013b174a\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385051,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592389a5285890809589811843/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716419,\"Height\":720,\"Width\":1280,\"Size\":4674637,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"3dab092ce4a93c6fa58f06a43b41ca79\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582179,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592389a5285890809589811843/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592389a5285890809589811843/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605762355,1),
	(8,5285890809595472752,'{\"EventHandle\":\"8095178859869222032\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-5969d596d7766a95d3976fcabb71c5a9t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809595472752\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850393395285890809595472752/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850393395285890809595472752/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455572,\"Height\":540,\"Width\":960,\"Size\":2972613,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"1ec2118d459c3f968ec1870f6543f478\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385065,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850393395285890809595472752/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850393395285890809595472752/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850393395285890809595472752/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605764287,1),
	(9,5285890809595472796,'{\"EventHandle\":\"8095022800235174065\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-d02426eb12502d7856ec1447c9e96401t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809595472796\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850393c15285890809595472796/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850393c15285890809595472796/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455567,\"Height\":540,\"Width\":960,\"Size\":2972579,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"47294ee4adb38470fb214df33fc34ef4\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385060,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850393c15285890809595472796/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850393c15285890809595472796/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850393c15285890809595472796/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605764288,1),
	(10,5285890809589812027,'{\"EventHandle\":\"8095039700766603119\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-800e2ce37d9be1a47a9f379d4d85e01ct0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809589812027\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592a2b55285890809589812027/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a2b55285890809589812027/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455558,\"Height\":540,\"Width\":960,\"Size\":2972519,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"9310179d0f3aa6b45a9cf58c013b174a\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385051,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a2b55285890809589812027/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a2b55285890809589812027/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a2b55285890809589812027/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605765230,1),
	(11,5285890809589812046,'{\"EventHandle\":\"8095202929856864344\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-f49b92a057275c91d9ef67e122222945t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809589812046\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592a2f65285890809589812046/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a2f65285890809589812046/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455615,\"Height\":540,\"Width\":960,\"Size\":2972889,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"6ff7aeea228fa166ad3c1011af9ae50b\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385108,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a2f65285890809589812046/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a2f65285890809589812046/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a2f65285890809589812046/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605767663,1),
	(12,5285890809593212821,'{\"EventHandle\":\"8095106650895495479\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-24934e460b99bea1211773b33ec983a5t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809593212821\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/e5f8394c5285890809593212821/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f8394c5285890809593212821/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455596,\"Height\":540,\"Width\":960,\"Size\":2972769,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"4843642d4e1861edd7cf39706e016b49\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385089,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f8394c5285890809593212821/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f8394c5285890809593212821/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f8394c5285890809593212821/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605767663,1),
	(13,5285890809589812089,'{\"EventHandle\":\"8094956019500457784\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809589812089\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T06:09:21Z\",\"UpdateTime\":\"2020-11-19T06:09:34Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a37d5285890809589812089/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592a37d5285890809589812089/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009101933004667\\\",\\\"task_id\\\":3690708}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809589812089\"},\"ProcedureTaskId\":\"1251103923-procedurev2-fa395e7e9f5856b13a631ce3a45e59a5t0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605767664,0),
	(14,5285890809589812089,'{\"EventHandle\":\"8095167327986753430\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-fa395e7e9f5856b13a631ce3a45e59a5t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809589812089\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592a37d5285890809589812089/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a37d5285890809589812089/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455554,\"Height\":540,\"Width\":960,\"Size\":2972496,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"95926d4d7d86f292f2b0e83d266a4521\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385047,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a37d5285890809589812089/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a37d5285890809589812089/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a37d5285890809589812089/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605767664,1),
	(15,5285890809589812104,'{\"EventHandle\":\"8094995908120307655\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809589812104\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T06:15:19Z\",\"UpdateTime\":\"2020-11-19T06:15:34Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a6b15285890809589812104/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592a6b15285890809589812104/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009101933004667\\\",\\\"task_id\\\":10185224}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809589812104\"},\"ProcedureTaskId\":\"1251103923-procedurev2-24f4d0e87bf57b81fc3ba857f17ea085t0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605767722,0),
	(16,5285890809589812104,'{\"EventHandle\":\"8095217546969217774\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-24f4d0e87bf57b81fc3ba857f17ea085t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809589812104\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592a6b15285890809589812104/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a6b15285890809589812104/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455615,\"Height\":540,\"Width\":960,\"Size\":2972889,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"6ff7aeea228fa166ad3c1011af9ae50b\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385108,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a6b15285890809589812104/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716419,\"Height\":720,\"Width\":1280,\"Size\":4674637,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"3dab092ce4a93c6fa58f06a43b41ca79\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582179,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a6b15285890809589812104/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a6b15285890809589812104/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605767722,1),
	(17,5285890809595473006,'{\"EventHandle\":\"8094998142378339287\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809595473006\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T06:20:55Z\",\"UpdateTime\":\"2020-11-19T06:21:07Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850401325285890809595473006/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850401325285890809595473006/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009181548217779\\\",\\\"task_id\\\":12204011}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809595473006\"},\"ProcedureTaskId\":\"1251103923-procedurev2-4f311d1f506adf63380b4503d23234act0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605767722,0),
	(18,5285890809595473006,'{\"EventHandle\":\"8094965783656280315\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-4f311d1f506adf63380b4503d23234act0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809595473006\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850401325285890809595473006/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850401325285890809595473006/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455583,\"Height\":540,\"Width\":960,\"Size\":2972682,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"703dbe3cf73a513c513701563f5d99bb\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385076,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850401325285890809595473006/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850401325285890809595473006/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850401325285890809595473006/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605767722,1),
	(19,5285890809593212934,'{\"EventHandle\":\"8095080424705163650\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809593212934\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T06:22:06Z\",\"UpdateTime\":\"2020-11-19T06:22:19Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f83db15285890809593212934/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/e5f83db15285890809593212934/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009181548217779\\\",\\\"task_id\\\":2877155}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809593212934\"},\"ProcedureTaskId\":\"1251103923-procedurev2-47566eb9766c7f775cf8c0486814d4f7t0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605767722,0),
	(20,5285890809593212934,'{\"EventHandle\":\"8095165325766925297\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-47566eb9766c7f775cf8c0486814d4f7t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809593212934\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/e5f83db15285890809593212934/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f83db15285890809593212934/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455572,\"Height\":540,\"Width\":960,\"Size\":2972611,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"50bb63b805aa8407fa334fd01e0cef5b\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385065,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f83db15285890809593212934/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f83db15285890809593212934/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f83db15285890809593212934/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605767723,1),
	(21,5285890809589815965,'{\"EventHandle\":\"8094991111662100237\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809589815965\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T06:23:03Z\",\"UpdateTime\":\"2020-11-19T06:23:16Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ea35285890809589815965/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d5946ea35285890809589815965/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009181548217779\\\",\\\"task_id\\\":15552957}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809589815965\"},\"ProcedureTaskId\":\"1251103923-procedurev2-66bda4849cc2f853b86170e84c8f6f88t0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605767723,0),
	(22,5285890809589815965,'{\"EventHandle\":\"8095192669317221806\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-66bda4849cc2f853b86170e84c8f6f88t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809589815965\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d5946ea35285890809589815965/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ea35285890809589815965/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455591,\"Height\":540,\"Width\":960,\"Size\":2972733,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"4452e1d63234d549fe1822ac8dac7258\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385084,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ea35285890809589815965/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ea35285890809589815965/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ea35285890809589815965/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605767723,1),
	(23,5285890809589815978,'{\"EventHandle\":\"8095077137852988296\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809589815978\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T06:26:40Z\",\"UpdateTime\":\"2020-11-19T06:26:51Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ec75285890809589815978/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d5946ec75285890809589815978/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009181548217779\\\",\\\"task_id\\\":2054783}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809589815978\"},\"ProcedureTaskId\":\"1251103923-procedurev2-b1a69da22739ef94e2245755b68e7691t0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605767723,0),
	(24,5285890809589815978,'{\"EventHandle\":\"8095112757710371923\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-b1a69da22739ef94e2245755b68e7691t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809589815978\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d5946ec75285890809589815978/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ec75285890809589815978/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455548,\"Height\":540,\"Width\":960,\"Size\":2972453,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"62e73fd23da4fa7def4a2fef80b9f390\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385041,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ec75285890809589815978/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716419,\"Height\":720,\"Width\":1280,\"Size\":4674637,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"3dab092ce4a93c6fa58f06a43b41ca79\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582179,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ec75285890809589815978/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ec75285890809589815978/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605767781,1),
	(25,5285890809595473070,'{\"EventHandle\":\"8095082915035285002\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809595473070\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T06:32:23Z\",\"UpdateTime\":\"2020-11-19T06:32:36Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850402135285890809595473070/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850402135285890809595473070/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009181548217779\\\",\\\"task_id\\\":12616971}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809595473070\"},\"ProcedureTaskId\":\"1251103923-procedurev2-83cf81788dcf7313649be14e9d70f23dt0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605767782,0),
	(26,5285890809595473070,'{\"EventHandle\":\"8094987410972760146\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-83cf81788dcf7313649be14e9d70f23dt0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809595473070\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850402135285890809595473070/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850402135285890809595473070/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455548,\"Height\":540,\"Width\":960,\"Size\":2972453,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"62e73fd23da4fa7def4a2fef80b9f390\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385041,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850402135285890809595473070/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850402135285890809595473070/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850402135285890809595473070/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605767782,1),
	(27,5285890809589812240,'{\"EventHandle\":\"8095170406274671919\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809589812240\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T06:37:43Z\",\"UpdateTime\":\"2020-11-19T06:37:56Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592ab725285890809589812240/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592ab725285890809589812240/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009101933004667\\\",\\\"task_id\\\":9409760}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809589812240\"},\"ProcedureTaskId\":\"1251103923-procedurev2-7fd3e5bde031cc6d027a6e0d2877f3e2t0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605767902,0),
	(28,5285890809589812240,'{\"EventHandle\":\"8094965025339617178\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-7fd3e5bde031cc6d027a6e0d2877f3e2t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809589812240\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592ab725285890809589812240/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592ab725285890809589812240/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455479,\"Height\":540,\"Width\":960,\"Size\":2972003,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"751ba5fe0ce531c6e5e8dde117b2218a\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":384972,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100030},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592ab725285890809589812240/v.f100030.mp4\",\"Definition\":100030,\"Bitrate\":716294,\"Height\":720,\"Width\":1280,\"Size\":4673821,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"94f366a804682e820eeb896613cb2193\",\"AudioStreamSet\":[{\"Bitrate\":128166,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":582054,\"Height\":720,\"Width\":1280,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":0,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592ab725285890809589812240/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592ab725285890809589812240/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605767903,1),
	(29,5285890809595473942,'{\"EventHandle\":\"8095152197780210704\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809595473942\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T09:15:23Z\",\"UpdateTime\":\"2020-11-19T09:15:41Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850427fb5285890809595473942/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850427fb5285890809595473942/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009101933004667\\\",\\\"task_id\\\":7265444}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809595473942\"},\"ProcedureTaskId\":\"1251103923-procedurev2-f0d214b03acbb6859077466e5b43d839t0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605777587,0),
	(30,5285890809595473942,'{\"EventHandle\":\"8095052670670655239\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-f0d214b03acbb6859077466e5b43d839t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809595473942\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850427fb5285890809595473942/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850427fb5285890809595473942/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455546,\"Height\":540,\"Width\":960,\"Size\":2972439,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"6baee0213fef56cbbe2418c749e68b5e\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385039,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":10,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850427fb5285890809595473942/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850427fb5285890809595473942/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605777587,1),
	(31,5285890809589813076,'{\"EventHandle\":\"8095046905271600266\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809589813076\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T09:16:30Z\",\"UpdateTime\":\"2020-11-19T09:16:43Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5932fba5285890809589813076/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d5932fba5285890809589813076/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009101933004667\\\",\\\"task_id\\\":12164274}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809589813076\"},\"ProcedureTaskId\":\"1251103923-procedurev2-f33b41c19630d63532c0d9246ab4ce47t0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605777587,0),
	(32,5285890809589813076,'{\"EventHandle\":\"8095189680469559358\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-f33b41c19630d63532c0d9246ab4ce47t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809589813076\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d5932fba5285890809589813076/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5932fba5285890809589813076/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455548,\"Height\":540,\"Width\":960,\"Size\":2972453,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"62e73fd23da4fa7def4a2fef80b9f390\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385041,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":10,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5932fba5285890809589813076/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5932fba5285890809589813076/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605777587,1),
	(33,5285890809595473951,'{\"EventHandle\":\"8095019511585215158\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809595473951\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T09:18:59Z\",\"UpdateTime\":\"2020-11-19T09:19:09Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8504281b5285890809595473951/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/8504281b5285890809595473951/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009101933004667\\\",\\\"task_id\\\":2322355}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809595473951\"},\"ProcedureTaskId\":\"1251103923-procedurev2-ff29f4db70b1fdc1f313fbdc4a83bdbbt0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605777587,0),
	(34,5285890809595473951,'{\"EventHandle\":\"8094983300575154392\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-ff29f4db70b1fdc1f313fbdc4a83bdbbt0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809595473951\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/8504281b5285890809595473951/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8504281b5285890809595473951/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455558,\"Height\":540,\"Width\":960,\"Size\":2972519,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"9310179d0f3aa6b45a9cf58c013b174a\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385051,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":10,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8504281b5285890809595473951/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8504281b5285890809595473951/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605777588,1),
	(35,5285890809595473991,'{\"EventHandle\":\"8095019843479704400\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809595473991\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T09:29:38Z\",\"UpdateTime\":\"2020-11-19T09:29:42Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8504289f5285890809595473991/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/8504289f5285890809595473991/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009181548217779\\\",\\\"task_id\\\":9685651}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809595473991\"},\"ProcedureTaskId\":\"1251103923-procedurev2-315938bbfe6133a5b89f622a5c80af66t0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605778186,0),
	(36,5285890809589813158,'{\"EventHandle\":\"8094973712970301287\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809589813158\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T09:33:28Z\",\"UpdateTime\":\"2020-11-19T09:33:39Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59333bb5285890809589813158/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d59333bb5285890809589813158/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009101933004667\\\",\\\"task_id\\\":3065354}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809589813158\"},\"ProcedureTaskId\":\"1251103923-procedurev2-ee13fc57351211c678027163413b7898t0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605778426,0),
	(37,5285890809589813158,'{\"EventHandle\":\"8095038253242692813\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-ee13fc57351211c678027163413b7898t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809589813158\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d59333bb5285890809589813158/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59333bb5285890809589813158/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455565,\"Height\":540,\"Width\":960,\"Size\":2972562,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"de97e333f5bb06e84d4b805ba0e4dcfd\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385058,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":10,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59333bb5285890809589813158/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59333bb5285890809589813158/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605778426,1),
	(38,5285890809589813231,'{\"EventHandle\":\"8095179979748914841\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809589813231\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T09:51:09Z\",\"UpdateTime\":\"2020-11-19T09:51:18Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59337b35285890809589813231/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d59337b35285890809589813231/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009181548217779\\\",\\\"task_id\\\":15693014}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809589813231\"},\"ProcedureTaskId\":\"1251103923-procedurev2-8f92e67f9a82ab9d0d26fe20541c8172t0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605779700,0),
	(39,5285890809595474105,'{\"EventHandle\":\"8095080469108031895\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809595474105\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-19T09:53:27Z\",\"UpdateTime\":\"2020-11-19T09:53:37Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850491d35285890809595474105/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850491d35285890809595474105/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009181548217779\\\",\\\"task_id\\\":15628249}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809595474105\"},\"ProcedureTaskId\":\"1251103923-procedurev2-855679ad25063110e008f2477b300e0et0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605779700,0),
	(40,5285890809589819733,'{\"EventHandle\":\"8095116173047972489\",\"EventType\":\"NewFileUpload\",\"FileUploadEvent\":{\"FileId\":\"5285890809589819733\",\"MediaBasicInfo\":{\"Name\":\"test\",\"Description\":\"\",\"CreateTime\":\"2020-11-20T02:14:54Z\",\"UpdateTime\":\"2020-11-20T02:15:05Z\",\"ExpireTime\":\"9999-12-31T23:59:59Z\",\"ClassId\":0,\"ClassName\":\"其他\",\"ClassPath\":\"-1\",\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59697405285890809589819733/coverBySnapshot/coverBySnapshot_10_0.jpg\",\"Type\":\"mp4\",\"MediaUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d59697405285890809589819733/f0.mp4\",\"SourceInfo\":{\"SourceType\":\"Upload\",\"SourceContext\":\"{\\\"user_id\\\":\\\"202009101933004667\\\",\\\"task_id\\\":10488807}\"},\"StorageRegion\":\"ap-shanghai\",\"Vid\":\"5285890809589819733\"},\"ProcedureTaskId\":\"1251103923-procedurev2-db4f506a01e33c2fbc6e0b6ffca19f54t0\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293}}}',1605838517,0),
	(41,5285890809589819733,'{\"EventHandle\":\"8095031698200635617\",\"EventType\":\"ProcedureStateChanged\",\"ProcedureStateChangeEvent\":{\"TaskId\":\"1251103923-procedurev2-db4f506a01e33c2fbc6e0b6ffca19f54t0\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"FileId\":\"5285890809589819733\",\"FileName\":\"test\",\"FileUrl\":\"https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d59697405285890809589819733/f0.mp4\",\"MetaData\":{\"Size\":4372373,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Bitrate\":664569,\"Height\":480,\"Width\":854,\"Duration\":52.20833206176758,\"Rotate\":0,\"VideoStreamSet\":[{\"Bitrate\":537875,\"Height\":480,\"Width\":854,\"Codec\":\"h264\",\"Fps\":24}],\"AudioStreamSet\":[{\"Bitrate\":126694,\"SamplingRate\":48000,\"Codec\":\"aac\"}],\"VideoDuration\":52.20833206176758,\"AudioDuration\":51.9466667175293},\"MediaProcessResultSet\":[{\"Type\":\"Transcode\",\"TranscodeTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":100020},\"Output\":{\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59697405285890809589819733/v.f100020.mp4\",\"Definition\":100020,\"Bitrate\":455605,\"Height\":540,\"Width\":960,\"Size\":2972823,\"Duration\":52.2,\"Container\":\"mov,mp4,m4a,3gp,3g2,mj2\",\"Md5\":\"558256a06d07a10ea694ca5dc9d7b2fc\",\"AudioStreamSet\":[{\"Bitrate\":64083,\"SamplingRate\":44100,\"Codec\":\"aac\"}],\"VideoStreamSet\":[{\"Bitrate\":385098,\"Height\":540,\"Width\":960,\"Codec\":\"h264\",\"Fps\":25}]}}},{\"Type\":\"SnapshotByTimeOffset\",\"SnapshotByTimeOffsetTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"TimeOffsetSet\":[0]},\"Output\":{\"Definition\":10,\"PicInfoSet\":[{\"TimeOffset\":0,\"Url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59697405285890809589819733/snapshotByTimeOffset/snapshotByTimeOffset_10_0.jpg\"}]}}},{\"Type\":\"CoverBySnapshot\",\"CoverBySnapshotTask\":{\"Status\":\"SUCCESS\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"Input\":{\"Definition\":10,\"PositionType\":\"Time\",\"PositionValue\":2},\"Output\":{\"CoverUrl\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59697405285890809589819733/coverBySnapshot/coverBySnapshot_10_0.jpg\"}}}],\"TasksPriority\":0,\"TasksNotifyMode\":\"\",\"SessionContext\":\"\",\"SessionId\":\"\"}}',1605838518,1);

/*!40000 ALTER TABLE `tencent_cloud_events` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table thumbs_up
# ------------------------------------------------------------

CREATE TABLE `thumbs_up` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `type_id` bigint(20) NOT NULL COMMENT '作品id （视频id/帖子id/评论id）',
  `user_id` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户id',
  `to_user_id` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '被点赞的用户id',
  `zan_type` tinyint(2) NOT NULL COMMENT '1 视频点赞 2 帖子点赞 3 评论点赞',
  `status` tinyint(1) NOT NULL COMMENT '1赞 2取消点赞',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `to_user_id` (`to_user_id`),
  KEY `type_id` (`type_id`),
  KEY `thumbs_up_id` (`type_id`,`zan_type`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='点赞表（针对帖子/视频/评论）';

LOCK TABLES `thumbs_up` WRITE;
/*!40000 ALTER TABLE `thumbs_up` DISABLE KEYS */;

INSERT INTO `thumbs_up` (`id`, `type_id`, `user_id`, `to_user_id`, `zan_type`, `status`, `create_at`)
VALUES
	(10,59,'202009181548217779','202010101545291936',1,1,1603765530),
	(11,43,'202009101933004667','202010101545291936',3,1,1603765530),
	(12,59,'202010101545291936','202010101545291936',1,1,1603765530),
	(13,43,'202010101545291936','202010101545291936',3,1,1603765530),
	(14,59,'202009101933004667','202010101545291936',1,1,1603765530),
	(15,59,'202010101545291936','202010101545291936',1,1,1603765530),
	(16,61,'202010101545291936','202010101545291936',1,1,1603765530);

/*!40000 ALTER TABLE `thumbs_up` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user
# ------------------------------------------------------------

CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `nick_name` varchar(45) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称',
  `mobile_num` bigint(20) NOT NULL COMMENT '手机号码',
  `password` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户密码',
  `user_id` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户id',
  `gender` tinyint(1) NOT NULL DEFAULT '0' COMMENT '性别 0人妖 1男性 2女性',
  `born` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '出生日期',
  `age` int(3) NOT NULL DEFAULT '0' COMMENT '年龄',
  `avatar` varchar(300) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '头像',
  `status` tinyint(1) DEFAULT '0' COMMENT '0 正常 1 封禁',
  `last_login_time` int(11) DEFAULT NULL COMMENT '最后登录时间',
  `signature` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '签名',
  `device_type` tinyint(2) DEFAULT NULL COMMENT '设备类型 0 android 1 iOS 2 web',
  `city` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '城市',
  `is_anchor` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0不是主播 1为主播',
  `channel_id` int(11) NOT NULL DEFAULT '0' COMMENT '渠道id',
  `background_img` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '背景图',
  `title` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '称号/特殊身份',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `user_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '用户类型 0 手机号 1 微信 2 QQ 3 微博',
  `country` int(3) NOT NULL DEFAULT '0' COMMENT '国家',
  `reg_ip` varchar(30) COLLATE utf8mb4_general_ci DEFAULT ' ' COMMENT '注册ip',
  `device_token` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '设备token',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='用户表';

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;

INSERT INTO `user` (`id`, `nick_name`, `mobile_num`, `password`, `user_id`, `gender`, `born`, `age`, `avatar`, `status`, `last_login_time`, `signature`, `device_type`, `city`, `is_anchor`, `channel_id`, `background_img`, `title`, `create_at`, `update_at`, `user_type`, `country`, `reg_ip`, `device_token`)
VALUES
	(1,'陈二go',13177656222,'','202009101933004667',1,'1993-06-20',27,'https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=2863307160,2081500909&fm=26&gp=0.jpg?imageView2/1/w/80/h/80',0,1603683458,'菩提本无树，明镜亦非台',NULL,'',0,0,'','',1603705779,1603776629,0,1,' ',''),
	(2,'大狗',13689896669,'','202009181548217779',0,'',0,'http://192.168.50.148:13002/upload/2020_10_23/539519450638979072.png',0,1603673405,'',NULL,'',0,0,'','',1603701234,1603701234,0,0,' ',''),
	(3,'小狗',18779487623,'','202010101545291936',0,'',0,'http://192.168.50.148:13002/upload/2020_10_23/539519490564558848.png',0,1603701039,'',NULL,'',0,0,'','',1603702345,1603701234,0,0,' ','');

/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user_attention
# ------------------------------------------------------------

CREATE TABLE `user_attention` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '被关注的用户id',
  `attention_uid` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '关注的用户id',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1表示关注 0表示取消关注',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `attention_uid` (`attention_uid`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户关注表';

LOCK TABLES `user_attention` WRITE;
/*!40000 ALTER TABLE `user_attention` DISABLE KEYS */;

INSERT INTO `user_attention` (`id`, `user_id`, `attention_uid`, `status`, `create_at`)
VALUES
	(4,'202010101545291936','202009101933004667',1,1603763768);

/*!40000 ALTER TABLE `user_attention` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user_browse_record
# ------------------------------------------------------------

CREATE TABLE `user_browse_record` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `user_id` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户id',
  `compose_id` bigint(20) NOT NULL COMMENT '作品id',
  `compose_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0 视频 1 帖子',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户浏览过的作品记录（包含帖子、视频）';

LOCK TABLES `user_browse_record` WRITE;
/*!40000 ALTER TABLE `user_browse_record` DISABLE KEYS */;

INSERT INTO `user_browse_record` (`id`, `user_id`, `compose_id`, `compose_type`, `create_at`, `update_at`)
VALUES
	(1,'1',61,0,0,0),
	(2,'1',62,0,0,0);

/*!40000 ALTER TABLE `user_browse_record` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user_play_duration_record
# ------------------------------------------------------------

CREATE TABLE `user_play_duration_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `play_duration` int(8) NOT NULL COMMENT '当前播放的时长（秒）',
  `total_duration` int(8) NOT NULL COMMENT '当前视频总时长（秒）',
  `cur_progress` decimal(5,2) NOT NULL DEFAULT '0.00' COMMENT '当前播放的进度',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `video_id` (`video_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户播放视频时长记录';



# Dump of table user_ycoin
# ------------------------------------------------------------

CREATE TABLE `user_ycoin` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `ycoin` int(11) NOT NULL COMMENT '游币数',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `ycoin` (`ycoin`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户货币表';



# Dump of table video_barrage
# ------------------------------------------------------------

CREATE TABLE `video_barrage` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
  `video_cur_duration` int(8) NOT NULL COMMENT '视频当前时长节点（单位：秒）',
  `content` varchar(512) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '弹幕内容',
  `user_id` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户id',
  `color` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '弹幕字体颜色',
  `font` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '弹幕字体',
  `barrage_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '预留字段',
  `location` tinyint(2) NOT NULL DEFAULT '0' COMMENT '弹幕位置',
  `send_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '弹幕发送时间',
  PRIMARY KEY (`id`),
  KEY `video_id` (`video_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='视频弹幕';



# Dump of table video_comment
# ------------------------------------------------------------

CREATE TABLE `video_comment` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '评论id',
  `user_id` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论人userId',
  `user_name` varchar(45) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '评论人名称',
  `avatar` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '头像',
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `parent_comment_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '父评论id',
  `parent_comment_user_id` varchar(60) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '父评论的用户id',
  `reply_comment_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '被回复的评论id',
  `reply_comment_user_id` varchar(60) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '被回复的评论用户id',
  `comment_level` tinyint(4) NOT NULL DEFAULT '1' COMMENT '评论等级[ 1 一级评论 默认 ，2 二级评论]',
  `content` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '评论的内容',
  `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '状态 (1 有效，0 逻辑删除)',
  `is_top` tinyint(2) NOT NULL DEFAULT '0' COMMENT '置顶状态[ 1 置顶，0 不置顶 默认 ]',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_video_id` (`video_id`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE,
  KEY `idx_create_time` (`create_at`),
  KEY `idx_parent_comment_id` (`parent_comment_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='视频评论表';

LOCK TABLES `video_comment` WRITE;
/*!40000 ALTER TABLE `video_comment` DISABLE KEYS */;

INSERT INTO `video_comment` (`id`, `user_id`, `user_name`, `avatar`, `video_id`, `parent_comment_id`, `parent_comment_user_id`, `reply_comment_id`, `reply_comment_user_id`, `comment_level`, `content`, `status`, `is_top`, `create_at`)
VALUES
	(1,'202009181548217779','','',59,0,'',0,'',1,'测试评论1',1,0,1600426890),
	(43,'202010101545291936','','',61,0,'',0,'',1,'测试评论2',1,0,1600426802),
	(44,'202009101933004667','','',61,0,'',0,'',1,'测试评论9',1,0,1600426802),
	(45,'202010101545291936','','',61,0,'',0,'',1,'测试评论8',1,0,1600426802),
	(46,'202009181548217779','','',61,0,'',0,'',1,'测试评论10',1,0,1600426802),
	(47,'202009101933004667','','',61,0,'',0,'',1,'测试评论7',1,0,1600426802),
	(48,'202010101545291936','','',61,0,'',0,'',1,'测试评论6',1,0,1600426802),
	(49,'202010101545291936','','https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=2863307160,2081500909&fm=26&gp=0.jpg?imageView2/1/w/80/h/80',61,0,'',0,'',1,'测试评论5',1,0,1600426802),
	(54,'202009181548217779','','https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=2863307160,2081500909&fm=26&gp=0.jpg?imageView2/1/w/80/h/80',61,0,'',0,'',1,'测试评论15',1,0,1600426802),
	(55,'202009181548217779','','https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=2863307160,2081500909&fm=26&gp=0.jpg?imageView2/1/w/80/h/80',61,0,'',0,'',1,'测试评论3',1,0,1600426802),
	(56,'202010101545291936','','',62,0,'',0,'',1,'测试评论13',1,0,1600426802),
	(57,'202009101933004667','陈二go','https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=2863307160,2081500909&fm=26&gp=0.jpg?imageView2/1/w/80/h/80',97,0,'',0,'',1,'我是1级评论',1,0,1604494389),
	(58,'202009101933004667','陈二go','https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=2863307160,2081500909&fm=26&gp=0.jpg?imageView2/1/w/80/h/80',97,0,'',0,'',1,'我是1级评论',1,0,1604494436),
	(59,'202009181548217779','大狗','http://192.168.50.148:13002/upload/2020_10_23/539519450638979072.png',97,57,'202009101933004667',57,'202009101933004667',2,'评论回复no.1',1,0,1604494436),
	(60,'202009101933004667','陈二go','https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=2863307160,2081500909&fm=26&gp=0.jpg?imageView2/1/w/80/h/80',97,0,'',0,'',1,'我是1级评论',1,0,1604494610),
	(61,'202009181548217779','大狗','http://192.168.50.148:13002/upload/2020_10_23/539519450638979072.png',97,57,'202009101933004667',57,'202009101933004667',2,'评论回复no.1',1,0,1604494610),
	(62,'202009101933004667','陈二go','https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=2863307160,2081500909&fm=26&gp=0.jpg?imageView2/1/w/80/h/80',97,0,'',0,'',1,'我是1级评论',1,0,1605236118),
	(63,'202009181548217779','大狗','http://192.168.50.148:13002/upload/2020_10_23/539519450638979072.png',97,57,'202009101933004667',57,'202009101933004667',2,'评论回复no.1',1,0,1605236118);

/*!40000 ALTER TABLE `video_comment` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table video_label_config
# ------------------------------------------------------------

CREATE TABLE `video_label_config` (
  `label_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '标签id',
  `pid` int(11) NOT NULL COMMENT '父类id 0为1级分类',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` int(1) NOT NULL DEFAULT '1' COMMENT '类别状态1-正常,2-已废弃',
  `label_name` varchar(64) NOT NULL COMMENT '标签名称',
  `icon` varchar(256) NOT NULL DEFAULT '' COMMENT '标签icon',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`label_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COMMENT='视频标签配置';

LOCK TABLES `video_label_config` WRITE;
/*!40000 ALTER TABLE `video_label_config` DISABLE KEYS */;

INSERT INTO `video_label_config` (`label_id`, `pid`, `sortorder`, `status`, `label_name`, `icon`, `create_at`, `update_at`)
VALUES
	(1,0,0,1,'无人机竞技','',1600426884,0),
	(5,0,123,1,'电竞','',1603424035,1603424035);

/*!40000 ALTER TABLE `video_label_config` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table video_labels
# ------------------------------------------------------------

CREATE TABLE `video_labels` (
  `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
  `label_id` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标签id',
  `label_name` varchar(521) DEFAULT NULL COMMENT '标签名',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '视频审核通过 则status为1 其他情况默认为0',
  PRIMARY KEY (`video_id`,`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='视频拥有的标签表';

LOCK TABLES `video_labels` WRITE;
/*!40000 ALTER TABLE `video_labels` DISABLE KEYS */;

INSERT INTO `video_labels` (`video_id`, `label_id`, `label_name`, `create_at`, `status`)
VALUES
	(10,'1','todo:通过标签id获取名称',1600426884,1),
	(10,'3','todo:通过标签id获取名称',1600426884,1),
	(11,'1','todo:通过标签id获取名称',1600426884,1),
	(11,'3','todo:通过标签id获取名称',1600426884,1),
	(12,'1','todo:通过标签id获取名称',1600426885,1),
	(12,'3','todo:通过标签id获取名称',1600426885,1),
	(13,'1','todo:通过标签id获取名称',1600426886,1),
	(13,'3','todo:通过标签id获取名称',1600426886,1),
	(14,'1','todo:通过标签id获取名称',1600426886,1),
	(14,'3','todo:通过标签id获取名称',1600426886,1),
	(15,'1','todo:通过标签id获取名称',1600426887,1),
	(15,'3','todo:通过标签id获取名称',1600426887,1),
	(16,'1','todo:通过标签id获取名称',1600426887,1),
	(16,'3','todo:通过标签id获取名称',1600426887,1),
	(17,'1','todo:通过标签id获取名称',1600426888,1),
	(17,'3','todo:通过标签id获取名称',1600426888,1),
	(18,'1','todo:通过标签id获取名称',1600426888,1),
	(18,'3','todo:通过标签id获取名称',1600426888,1),
	(19,'1','todo:通过标签id获取名称',1600426889,1),
	(19,'3','todo:通过标签id获取名称',1600426889,1),
	(20,'1','todo:通过标签id获取名称',1600426889,1),
	(20,'3','todo:通过标签id获取名称',1600426889,1),
	(21,'1','todo:通过标签id获取名称',1600426890,1),
	(21,'3','todo:通过标签id获取名称',1600426890,1),
	(22,'1','todo:通过标签id获取名称',1600426890,1),
	(22,'3','todo:通过标签id获取名称',1600426890,1),
	(23,'1','todo:通过标签id获取名称',1600426891,1),
	(23,'3','todo:通过标签id获取名称',1600426891,1),
	(24,'1','todo:通过标签id获取名称',1600426891,1),
	(24,'3','todo:通过标签id获取名称',1600426891,1),
	(25,'1','todo:通过标签id获取名称',1600426892,1),
	(25,'3','todo:通过标签id获取名称',1600426892,1),
	(26,'1','todo:通过标签id获取名称',1600426892,1),
	(26,'3','todo:通过标签id获取名称',1600426892,1),
	(27,'1','todo:通过标签id获取名称',1600426892,1),
	(27,'3','todo:通过标签id获取名称',1600426892,1),
	(28,'1','todo:通过标签id获取名称',1600426893,1),
	(28,'3','todo:通过标签id获取名称',1600426893,1),
	(29,'1','todo:通过标签id获取名称',1600426893,1),
	(29,'3','todo:通过标签id获取名称',1600426893,1),
	(30,'1','todo:通过标签id获取名称',1600426894,1),
	(30,'3','todo:通过标签id获取名称',1600426894,1),
	(31,'1','todo:通过标签id获取名称',1600426894,1),
	(31,'3','todo:通过标签id获取名称',1600426894,1),
	(32,'1','todo:通过标签id获取名称',1600426895,1),
	(32,'3','todo:通过标签id获取名称',1600426895,1),
	(33,'1','todo:通过标签id获取名称',1600426895,1),
	(33,'3','todo:通过标签id获取名称',1600426895,1),
	(34,'1','todo:通过标签id获取名称',1600426895,1),
	(34,'3','todo:通过标签id获取名称',1600426895,1),
	(35,'1','todo:通过标签id获取名称',1600426896,1),
	(35,'3','todo:通过标签id获取名称',1600426896,1),
	(36,'1','todo:通过标签id获取名称',1600426896,1),
	(36,'3','todo:通过标签id获取名称',1600426896,1),
	(37,'1','todo:通过标签id获取名称',1600426897,1),
	(37,'3','todo:通过标签id获取名称',1600426897,1),
	(38,'1','todo:通过标签id获取名称',1600426897,1),
	(38,'3','todo:通过标签id获取名称',1600426897,1),
	(39,'1','todo:通过标签id获取名称',1600426897,1),
	(39,'3','todo:通过标签id获取名称',1600426897,1),
	(40,'1','todo:通过标签id获取名称',1600426898,1),
	(40,'3','todo:通过标签id获取名称',1600426898,1),
	(41,'1','todo:通过标签id获取名称',1600426898,1),
	(41,'3','todo:通过标签id获取名称',1600426898,1),
	(42,'1','todo:通过标签id获取名称',1600426899,1),
	(42,'3','todo:通过标签id获取名称',1600426899,1),
	(43,'1','todo:通过标签id获取名称',1600426899,1),
	(43,'3','todo:通过标签id获取名称',1600426899,1),
	(44,'1','todo:通过标签id获取名称',1600426899,1),
	(44,'3','todo:通过标签id获取名称',1600426899,1),
	(45,'1','todo:通过标签id获取名称',1600426900,1),
	(45,'3','todo:通过标签id获取名称',1600426900,1),
	(46,'1','todo:通过标签id获取名称',1600426900,1),
	(46,'3','todo:通过标签id获取名称',1600426900,1),
	(47,'1','todo:通过标签id获取名称',1600426901,1),
	(47,'3','todo:通过标签id获取名称',1600426901,1),
	(48,'1','todo:通过标签id获取名称',1600426901,1),
	(48,'3','todo:通过标签id获取名称',1600426901,1),
	(49,'1','todo:通过标签id获取名称',1600426901,1),
	(49,'3','todo:通过标签id获取名称',1600426901,1),
	(50,'1','todo:通过标签id获取名称',1600426902,1),
	(50,'3','todo:通过标签id获取名称',1600426902,1),
	(51,'1','todo:通过标签id获取名称',1600426902,1),
	(51,'3','todo:通过标签id获取名称',1600426902,1),
	(52,'1','todo:通过标签id获取名称',1600426902,1),
	(52,'3','todo:通过标签id获取名称',1600426902,1),
	(53,'1','todo:通过标签id获取名称',1600426903,1),
	(53,'3','todo:通过标签id获取名称',1600426903,1),
	(54,'1','todo:通过标签id获取名称',1600426903,1),
	(54,'3','todo:通过标签id获取名称',1600426903,1),
	(55,'1','todo:通过标签id获取名称',1600426903,1),
	(55,'3','todo:通过标签id获取名称',1600426903,1),
	(56,'1','todo:通过标签id获取名称',1600426904,1),
	(56,'3','todo:通过标签id获取名称',1600426904,1),
	(57,'1','todo:通过标签id获取名称',1600426904,1),
	(57,'3','todo:通过标签id获取名称',1600426904,1),
	(58,'1','todo:通过标签id获取名称',1600426904,1),
	(58,'3','todo:通过标签id获取名称',1600426904,1),
	(59,'1','todo:通过标签id获取名称',1600426905,1),
	(59,'3','todo:通过标签id获取名称',1600426905,1),
	(60,'1','todo:通过标签id获取名称',1600426905,1),
	(60,'3','todo:通过标签id获取名称',1600426905,1),
	(61,'1','todo:通过标签id获取名称',1600426906,1),
	(61,'3','todo:通过标签id获取名称',1600426906,1),
	(62,'1','todo:通过标签id获取名称',1600426907,1),
	(62,'3','todo:通过标签id获取名称',1600426907,1),
	(63,'1','todo:通过标签id获取名称',1600426907,1),
	(63,'3','todo:通过标签id获取名称',1600426907,1),
	(64,'1','todo:通过标签id获取名称',1600426908,1),
	(64,'3','todo:通过标签id获取名称',1600426908,1),
	(65,'1','todo:通过标签id获取名称',1600426908,1),
	(65,'3','todo:通过标签id获取名称',1600426908,1);

/*!40000 ALTER TABLE `video_labels` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table video_live
# ------------------------------------------------------------

CREATE TABLE `video_live` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `anchor_id` varchar(60) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主播id',
  `room_id` varchar(60) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '房间id',
  `cover` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '直播封面',
  `rtmp_addr` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'rtmp地址',
  `flv_addr` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'flv地址',
  `hls_addr` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'hls地址',
  `stream_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '推流url',
  `stream_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '推流密钥',
  `play_time` int(11) NOT NULL COMMENT '开播时间',
  `end_time` int(11) NOT NULL COMMENT '结束时间',
  `income_ycoin` int(11) DEFAULT NULL COMMENT '本次直播收益（游币）',
  `status` tinyint(1) DEFAULT '0' COMMENT '状态 0未直播 1直播中 2异常',
  `describe` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '描述',
  `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '直播标签',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '记录创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '记录更新时间',
  `duration` bigint(20) DEFAULT '0' COMMENT '时长（秒）',
  `live_type` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '直播类型（0：管理员[sys_user]，1：用户[user]）',
  `manager` int(11) NOT NULL DEFAULT '0' COMMENT '后台操作用户',
  `sequence` varchar(50) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '推流唯一标识',
  PRIMARY KEY (`id`),
  KEY `anchor_id` (`anchor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='视频直播表';



# Dump of table video_report
# ------------------------------------------------------------

CREATE TABLE `video_report` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(60) NOT NULL DEFAULT '' COMMENT '用户id',
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `video_id` (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频举报';



# Dump of table video_statistic
# ------------------------------------------------------------

CREATE TABLE `video_statistic` (
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `fabulous_num` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
  `browse_num` int(11) NOT NULL DEFAULT '0' COMMENT '浏览数',
  `share_num` int(11) NOT NULL DEFAULT '0' COMMENT '分享数',
  `reward_num` int(11) NOT NULL DEFAULT '0' COMMENT '打赏的游币数',
  `barrage_num` int(11) NOT NULL DEFAULT '0' COMMENT '弹幕数',
  `comment_num` int(11) NOT NULL DEFAULT '0' COMMENT '评论数',
  `collect_num` int(11) NOT NULL DEFAULT '0' COMMENT '收藏数',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频相关参数统计';

LOCK TABLES `video_statistic` WRITE;
/*!40000 ALTER TABLE `video_statistic` DISABLE KEYS */;

INSERT INTO `video_statistic` (`video_id`, `fabulous_num`, `browse_num`, `share_num`, `reward_num`, `barrage_num`, `comment_num`, `collect_num`, `create_at`, `update_at`)
VALUES
	(10,1,0,0,0,0,0,0,1600426884,1600427096),
	(11,2,0,0,0,0,0,0,1600426884,1600427161),
	(12,2,0,0,0,0,0,0,1600426885,1600427242),
	(13,3,0,0,0,0,0,0,1600426886,1600427402),
	(14,0,0,0,0,0,0,0,1600426886,1600426886),
	(15,0,0,0,0,0,0,0,1600426887,1600426887),
	(16,0,0,0,0,0,0,0,1600426887,1600426887),
	(17,0,0,0,0,0,0,0,1600426888,1600426888),
	(18,0,0,0,0,0,0,0,1600426888,1600426888),
	(19,0,0,0,0,0,0,0,1600426889,1600426889),
	(20,0,0,0,0,0,0,0,1600426889,1600426889),
	(21,0,0,0,0,0,0,0,1600426890,1600426890),
	(22,0,0,0,0,0,0,0,1600426890,1600426890),
	(23,0,0,0,0,0,0,0,1600426891,1600426891),
	(24,0,0,0,0,0,0,0,1600426891,1600426891),
	(25,0,0,0,0,0,0,0,1600426892,1600426892),
	(26,0,0,0,0,0,0,0,1600426892,1600426892),
	(27,0,0,0,0,0,0,0,1600426892,1600426892),
	(28,0,0,0,0,0,0,0,1600426893,1600426893),
	(29,0,0,0,0,0,0,0,1600426893,1600426893),
	(30,0,0,0,0,0,0,0,1600426894,1600426894),
	(31,0,0,0,0,0,0,0,1600426894,1600426894),
	(32,0,0,0,0,0,0,0,1600426895,1600426895),
	(33,0,0,0,0,0,0,0,1600426895,1600426895),
	(34,0,0,0,0,0,0,0,1600426895,1600426895),
	(35,0,0,0,0,0,0,0,1600426896,1600426896),
	(36,0,0,0,0,0,0,0,1600426896,1600426896),
	(37,0,0,0,0,0,0,0,1600426897,1600426897),
	(38,0,0,0,0,0,0,0,1600426897,1600426897),
	(39,0,0,0,0,0,0,0,1600426897,1600426897),
	(40,8767,0,0,0,0,0,0,1600426898,1600426898),
	(41,667,0,0,0,0,0,0,1600426898,1600426898),
	(42,5435,0,0,0,0,0,0,1600426899,1600426899),
	(43,123,0,0,0,0,0,0,1600426899,1600426899),
	(44,32,0,0,0,0,0,0,1600426899,1600426899),
	(45,345,0,0,0,0,0,0,1600426900,1600426900),
	(46,52,0,0,0,0,0,0,1600426900,1600426900),
	(47,54,0,0,0,0,0,0,1600426901,1600426901),
	(48,42,0,0,0,0,0,0,1600426901,1600426901),
	(49,23,0,0,0,0,0,0,1600426901,1600426901),
	(50,43,0,0,0,0,0,0,1600426902,1600426902),
	(51,564,0,0,0,0,0,0,1600426902,1600426902),
	(52,4533,0,0,0,0,0,0,1600426902,1600426902),
	(53,14,0,0,0,0,0,0,1600426903,1600426903),
	(54,166,0,0,0,0,0,0,1600426903,1600426903),
	(55,13,0,0,0,0,0,0,1600426903,1600426903),
	(56,12,0,0,0,0,0,0,1600426904,1600426904),
	(57,11,0,0,0,0,0,0,1600426904,1600426904),
	(58,10,0,0,0,0,0,0,1600426904,1600426904),
	(59,1000,1,2,3,4,11,6,1600426905,1603765530),
	(60,9,0,0,0,0,0,0,1600426905,1600426905),
	(61,8,0,0,0,0,0,0,1600426906,1600426906),
	(62,7,0,0,0,0,0,0,1600426907,1600426907),
	(63,6,0,0,0,0,0,0,1600426907,1600426907),
	(64,5,0,0,0,0,0,0,1600426908,1600426908),
	(65,4,0,0,0,0,0,0,1600426908,1600426908),
	(98,0,0,0,0,0,0,0,1605703330,1605703330),
	(99,0,0,0,0,0,0,0,1605703575,1605703575),
	(100,0,0,0,0,0,0,0,1605703883,1605703883),
	(101,0,0,0,0,0,0,0,1605704513,1605704513),
	(102,0,0,0,0,0,0,0,1605751801,1605751801),
	(103,0,0,0,0,0,0,0,1605759787,1605759787),
	(104,0,0,0,0,0,0,0,1605762300,1605762300),
	(105,0,0,0,0,0,0,0,1605763500,1605763500),
	(106,0,0,0,0,0,0,0,1605764183,1605764183),
	(107,0,0,0,0,0,0,0,1605765115,1605765115),
	(108,0,0,0,0,0,0,0,1605765496,1605765496),
	(109,0,0,0,0,0,0,0,1605765656,1605765656),
	(110,0,0,0,0,0,0,0,1605767664,1605767664),
	(111,0,0,0,0,0,0,0,1605767722,1605767722),
	(112,0,0,0,0,0,0,0,1605767722,1605767722),
	(113,0,0,0,0,0,0,0,1605767722,1605767722),
	(114,0,0,0,0,0,0,0,1605767723,1605767723),
	(115,0,0,0,0,0,0,0,1605767723,1605767723),
	(116,0,0,0,0,0,0,0,1605767782,1605767782),
	(117,0,0,0,0,0,0,0,1605767902,1605767902),
	(118,0,0,0,0,0,0,0,1605777587,1605777587),
	(119,0,0,0,0,0,0,0,1605777587,1605777587),
	(120,0,0,0,0,0,0,0,1605777587,1605777587),
	(121,0,0,0,0,0,0,0,1605778186,1605778186),
	(122,0,0,0,0,0,0,0,1605778426,1605778426),
	(123,0,0,0,0,0,0,0,1605779700,1605779700),
	(124,0,0,0,0,0,0,0,1605779700,1605779700),
	(125,0,0,0,0,0,0,0,1605838517,1605838517);

/*!40000 ALTER TABLE `video_statistic` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table videos
# ------------------------------------------------------------

CREATE TABLE `videos` (
  `video_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '视频id',
  `title` mediumtext COLLATE utf8mb4_general_ci COMMENT '视频标题',
  `describe` mediumtext COLLATE utf8mb4_general_ci COMMENT '视频描述',
  `cover` varchar(521) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '视频封面',
  `video_addr` varchar(521) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '视频地址',
  `user_id` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户id',
  `user_type` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '添加用户类型（0：管理员[sys_user]，1：用户[user]）',
  `sortorder` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '排序',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '审核状态（0：审核中，1：审核通过 2：审核不通过 3：逻辑删除）',
  `is_recommend` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '推荐（0：不推荐；1：推荐）',
  `is_top` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '置顶（0：不置顶；1：置顶；）',
  `video_duration` int(8) NOT NULL DEFAULT '0' COMMENT '视频时长（单位：毫秒）',
  `rec_content` mediumtext COLLATE utf8mb4_general_ci COMMENT '推荐理由',
  `top_content` mediumtext COLLATE utf8mb4_general_ci COMMENT '置顶理由',
  `video_width` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频宽',
  `video_height` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频高',
  `create_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `file_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '腾讯云文件id',
  `size` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频大小（字节数）',
  `play_info` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '视频转码信息',
  PRIMARY KEY (`video_id`),
  KEY `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='视频表';

LOCK TABLES `videos` WRITE;
/*!40000 ALTER TABLE `videos` DISABLE KEYS */;

INSERT INTO `videos` (`video_id`, `title`, `describe`, `cover`, `video_addr`, `user_id`, `user_type`, `sortorder`, `status`, `is_recommend`, `is_top`, `video_duration`, `rec_content`, `top_content`, `video_width`, `video_height`, `create_at`, `update_at`, `file_id`, `size`, `play_info`)
VALUES
	(59,'fpv电竞无人机','了解一哈','https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=2863307160,2081500909&fm=26&gp=0.jpg?imageView2/1/w/80/h/80','http://vjs.zencdn.net/v/oceans.mp4','202010101545291936',0,1,1,0,0,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(61,'【油管搬运】<em class=\"keyword\">FPV穿越机</em> 直观感受PID各项调整带来的反应','ig vs rng','https://photo.mac69.com/180205/18020526/a9yPQozt0g.jpg','https://www.runoob.com/try/demo_source/movie.mp4','202009101933004667',0,1,1,0,1,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(62,'【油管搬运】<em class=\"keyword\">FPV穿越机</em> 直观感受PID各项调整带来的反应','ig vs rng','https://photo.mac69.com/180205/18020526/a9yPQozt0g.jpg','https://www.runoob.com/try/demo_source/movie.mp4','202010101545291936',0,1,0,0,0,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(63,'电竞无人机2020年度小组赛','ig vs rng','https://photo.mac69.com/180205/18020526/a9yPQozt0g.jpg','https://www.runoob.com/try/demo_source/movie.mp4','202009181548217779',0,1,0,0,0,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(89,'电竞无人机2020年度小组赛','ig vs rng','https://photo.mac69.com/180205/18020526/a9yPQozt0g.jpg','https://www.runoob.com/try/demo_source/movie.mp4','202009181548217779',0,1,0,0,0,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(90,'电竞无人机2020年度小组赛','ig vs rng','https://photo.mac69.com/180205/18020526/a9yPQozt0g.jpg','https://www.runoob.com/try/demo_source/movie.mp4','202009181548217779',0,1,0,0,0,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(91,'电竞无人机2020年度小组赛','ig vs rng','https://photo.mac69.com/180205/18020526/a9yPQozt0g.jpg','https://www.runoob.com/try/demo_source/movie.mp4','202009181548217779',0,1,0,0,0,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(92,'电竞无人机2020年度小组赛','ig vs rng','https://photo.mac69.com/180205/18020526/a9yPQozt0g.jpg','https://www.runoob.com/try/demo_source/movie.mp4','202009181548217779',0,1,0,0,0,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(93,'电竞无人机2020年度小组赛','ig vs rng','https://photo.mac69.com/180205/18020526/a9yPQozt0g.jpg','https://www.runoob.com/try/demo_source/movie.mp4','202009181548217779',0,1,0,0,0,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(94,'电竞无人机2020年度小组赛','ig vs rng','https://photo.mac69.com/180205/18020526/a9yPQozt0g.jpg','https://www.runoob.com/try/demo_source/movie.mp4','202009181548217779',0,1,0,0,0,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(95,'电竞无人机2020年度小组赛','ig vs rng','https://photo.mac69.com/180205/18020526/a9yPQozt0g.jpg','https://www.runoob.com/try/demo_source/movie.mp4','202009181548217779',0,1,0,0,0,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(96,'电竞无人机2020年度小组赛','ig vs rng','https://photo.mac69.com/180205/18020526/a9yPQozt0g.jpg','https://www.runoob.com/try/demo_source/movie.mp4','202009181548217779',0,1,0,0,0,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(97,'电竞无人机2020年度小组赛','ig vs rng','https://photo.mac69.com/180205/18020526/a9yPQozt0g.jpg','https://www.runoob.com/try/demo_source/movie.mp4','202009181548217779',0,1,1,0,0,0,NULL,NULL,0,0,0,0,0,0,'[{\"type\":\"1\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100010.mp4\",\"size\":436472,\"duration\":7480},{\"type\":\"3\",\"url\":\"http://1251103923.vod2.myqcloud.com/d2a664e5vodtranscq1251103923/085694cc5285890809406194027/v.f100030.mp4\",\"size\":1456466,\"duration\":7480}]'),
	(98,'test','test','','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/e3caf6115285890809593139008/f0.mp4','202009101933004667',1,0,0,0,0,0,'','',0,0,1605703330,1605703330,5285890809593139008,0,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e3caf6115285890809593139008/v.f100020.mp4\",\"size\":2972305,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e3caf6115285890809593139008/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(99,'test','test','','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/8483cf785285890809595400796/f0.mp4','202009101933004667',1,0,0,0,0,0,'','',0,0,1605703575,1605703575,5285890809595400796,0,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8483cf785285890809595400796/v.f100020.mp4\",\"size\":2972823,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8483cf785285890809595400796/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(100,'test','test','','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/e3caf6725285890809593139036/f0.mp4','202009101933004667',1,0,0,0,0,0,'','',0,0,1605703883,1605703883,5285890809593139036,0,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e3caf6725285890809593139036/v.f100020.mp4\",\"size\":2972825,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e3caf6725285890809593139036/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(101,'test','test','','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/8483d3165285890809595400848/f0.mp4','202009101933004667',1,0,0,0,0,0,'','',0,0,1605704513,1605704513,5285890809595400848,0,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8483d3165285890809595400848/v.f100020.mp4\",\"size\":2972611,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8483d3165285890809595400848/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(102,'test','test','','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d593e2015285890809589814946/f0.mp4','202009101933004667',1,0,0,0,0,0,'','',0,0,1605751801,1605751801,5285890809589814946,0,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d593e2015285890809589814946/v.f100020.mp4\",\"size\":2972667,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d593e2015285890809589814946/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(103,'【<em class=\"keyword\">fpv穿越机</em>】杭州.城东公园 函道树林穿梭↖(^ω^)↗','大野生通taycan gopro session5 加rs','http://i0.hdslb.com/bfs/archive/bd94eb25dd9cd631a8a35c8c6e75ec7be93f4716.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/85038aba5285890809595472555/f0.mp4','202009101933004667',1,0,0,0,0,0,'','',0,0,1605759787,1605759787,5285890809595472555,0,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/85038aba5285890809595472555/v.f100020.mp4\",\"size\":2972596,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/85038aba5285890809595472555/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(104,'【独家】2018香港<em class=\"keyword\">FPV穿越机</em>大赛（唯一女飞手现身）','电子游戏+虚拟现实+机械工程的合体，此坑深不可测！','http://i0.hdslb.com/bfs/archive/8cb3c085b2b54796918845e3e36064a2eb463716.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592389a5285890809589811843/f0.mp4','202010101545291936',1,0,0,0,0,0,'','',0,0,1605762300,1605762300,5285890809589811843,0,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592389a5285890809589811843/v.f100020.mp4\",\"size\":2972519,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592389a5285890809589811843/v.f100030.mp4\",\"size\":4674637,\"duration\":52200}]'),
	(105,'【油管搬运】<em class=\"keyword\">FPV穿越机</em>花式打杆教程第二弹','https://www.youtube.com/watch?v=_a2MsgEWzRU','http://i2.hdslb.com/bfs/archive/f3caee6d3aa32faeb6de1c59f39194691fd4d3b2.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850393395285890809595472752/f0.mp4','202010101545291936',1,0,0,0,0,0,'','',0,0,1605763500,1605763500,5285890809595472752,0,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850393395285890809595472752/v.f100020.mp4\",\"size\":2972613,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850393395285890809595472752/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(106,'【炫技】<em class=\"keyword\">FPV穿越机</em>冠军 勇闯烟花阵（1:08惊艳倒放）','海浪FPV\n电子游戏+虚拟现实+机械工程的合体，此坑深不可测！','http://i2.hdslb.com/bfs/archive/2821451d626938adee86847f97e45f7207bdd75e.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850393c15285890809595472796/f0.mp4','202010101545291936',1,0,0,0,0,0,'','',0,0,1605764183,1605764183,5285890809595472796,0,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850393c15285890809595472796/v.f100020.mp4\",\"size\":2972579,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850393c15285890809595472796/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(107,'油管大神<em class=\"keyword\">FPV穿越机</em>牛B打杆第四集，建议收藏学习','转自youtube, 关注、点赞、转发交流，给你带来最新最热门的镜外视频','http://i2.hdslb.com/bfs/archive/ffd8fd80d83c3355f63ab620a7c229f21c00553e.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592a2b55285890809589812027/f0.mp4','202009181548217779',1,0,0,0,0,0,'','',0,0,1605765115,1605765115,5285890809589812027,0,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a2b55285890809589812027/v.f100020.mp4\",\"size\":2972519,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a2b55285890809589812027/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(108,'【贴地飞行】【长板DH/FR】当<em class=\"keyword\">FPV穿越机</em>遇上长板速降-挪威风光大片','https://youtu.be/x6S6_6jDv8Q','http://i1.hdslb.com/bfs/archive/416d5dbf8b85dbdd63d6587f59bc17f7e4aed815.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592a2f65285890809589812046/f0.mp4','202009181548217779',1,0,0,0,0,0,'','',0,0,1605765496,1605765496,5285890809589812046,0,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a2f65285890809589812046/v.f100020.mp4\",\"size\":2972889,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a2f65285890809589812046/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(109,'<em class=\"keyword\">FPV穿越机</em>4K航拍-广汉连山麦田日出','这个航线速度可以飞9分钟，应该算是高效率了吧？','http://i2.hdslb.com/bfs/archive/bd4585f829fdbb0c28aa015cacd55297a4c65e92.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/e5f8394c5285890809593212821/f0.mp4','202009101933004667',1,0,0,0,0,0,'','',0,0,1605765656,1605765656,5285890809593212821,0,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f8394c5285890809593212821/v.f100020.mp4\",\"size\":2972769,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f8394c5285890809593212821/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(110,'【<em class=\"keyword\">FPV穿越机</em>】穿越小树林，树枝没装到一根，树叶没割到一片，确撞到了一根电线，极速翻滚炸鸡????????','-','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a37d5285890809589812089/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592a37d5285890809589812089/f0.mp4','202009101933004667',1,0,0,0,0,52208,'','',854,480,1605767664,1605767664,5285890809589812089,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a37d5285890809589812089/v.f100020.mp4\",\"size\":2972496,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a37d5285890809589812089/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(111,'<em class=\"keyword\">fpv穿越机</em>，做核酸检测的时候偷偷飞了一块电，时间紧张，不喜勿喷','曼巴149 3寸涵道拍摄\n新疆小区已经封了20天了\n这是核酸检测时候偷偷飞的','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a6b15285890809589812104/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592a6b15285890809589812104/f0.mp4','202009101933004667',1,0,0,0,0,52208,'','',854,480,1605767722,1605767722,5285890809589812104,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a6b15285890809589812104/v.f100020.mp4\",\"size\":2972889,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592a6b15285890809589812104/v.f100030.mp4\",\"size\":4674637,\"duration\":52200}]'),
	(112,'体验第一视角飞行，<em class=\"keyword\">FPV穿越机</em>','第一人称飞行视角','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850401325285890809595473006/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850401325285890809595473006/f0.mp4','202009181548217779',1,0,0,0,0,52208,'','',854,480,1605767722,1605767722,5285890809595473006,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850401325285890809595473006/v.f100020.mp4\",\"size\":2972682,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850401325285890809595473006/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(113,'用<em class=\"keyword\">FPV穿越机</em>陪孩们玩泵道，Beetle Hom芯动基地首飞','入坑穿越机，练了三个月的模拟器终于敢飞真机了，想飞好还是很难的，在芯动基地陪孩子们练练车。创世泰克 Beetle Hom 续航5分多，很给力。','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f83db15285890809593212934/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/e5f83db15285890809593212934/f0.mp4','202009181548217779',1,0,1,0,0,52208,'','',854,480,1605767722,1605767722,5285890809593212934,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f83db15285890809593212934/v.f100020.mp4\",\"size\":2972611,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/e5f83db15285890809593212934/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(114,'<em class=\"keyword\">FPV穿越机</em>  废弃工厂一日游，无炸不欢','垃圾佬在废弃工厂飞穿越机，摄像头太低了 速度跟不上了\n战损一颗致盈EX2306 2750kv\nfoxxer 弹弓镜头一个\n飞机老了......','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ea35285890809589815965/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d5946ea35285890809589815965/f0.mp4','202009181548217779',1,0,1,0,0,52208,'','',854,480,1605767723,1605767723,5285890809589815965,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ea35285890809589815965/v.f100020.mp4\",\"size\":2972733,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ea35285890809589815965/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(115,'大疆高清图传天空端无人机<em class=\"keyword\">FPV穿越机</em>GPS，INAV固件返航悬停失控返航一键返航','可切手动模式，可翻滚，可实现一键返航，失控返航，远航好机。','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ec75285890809589815978/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d5946ec75285890809589815978/f0.mp4','202009181548217779',1,0,1,0,0,52208,'','',854,480,1605767723,1605767723,5285890809589815978,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ec75285890809589815978/v.f100020.mp4\",\"size\":2972453,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5946ec75285890809589815978/v.f100030.mp4\",\"size\":4674637,\"duration\":52200}]'),
	(116,'<em class=\"keyword\">Fpv穿越机</em>mini pix飞控悬停一键返航定高定点','mini pix飞控+gps+osd\n一键返航，定高定点\n飞控还是不错的。','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850402135285890809595473070/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850402135285890809595473070/f0.mp4','202009181548217779',1,0,1,0,0,52208,'','',854,480,1605767782,1605767782,5285890809595473070,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850402135285890809595473070/v.f100020.mp4\",\"size\":2972453,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850402135285890809595473070/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(117,'<em class=\"keyword\">FPV穿越机</em>200元带接收机带OSD 测评 新手入门推荐，还要什么模拟器，这款就是地球OL模拟器， 安全便宜的小飞机，空心杯<em class=\"keyword\">穿越机FPV</em>','-','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592ab725285890809589812240/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d592ab725285890809589812240/f0.mp4','202009101933004667',1,0,1,0,0,52208,'','',854,480,1605767902,1605767902,5285890809589812240,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592ab725285890809589812240/v.f100020.mp4\",\"size\":2972003,\"duration\":52200},{\"type\":\"3\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d592ab725285890809589812240/v.f100030.mp4\",\"size\":4673821,\"duration\":52200}]'),
	(118,'【<em class=\"keyword\">FPV</em> <em class=\"keyword\">穿越机</em>】CINEMATIC 电影','CONFIG  6S\nFrame: Iflight XL6\nFC: Iflight Succex-E F4\nESC: 4in1 Iflight SucceX-E 45A\nMotor: Iflight Xing\nVTX: Iflight Force\nRX: Crossfire Nanox RX\n\nFrame: iflight XL5 V4\nFC: SucceX F4 TwingG\nESC: 4in1 Iflight SucceX 50A\nMotor: Iflight Xing 2306 1700KV\nVTX','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850427fb5285890809595473942/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850427fb5285890809595473942/f0.mp4','202009101933004667',1,0,0,0,0,52208,'','',854,480,1605777587,1605777587,5285890809595473942,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850427fb5285890809595473942/v.f100020.mp4\",\"size\":2972439,\"duration\":52200}]'),
	(119,'#<em class=\"keyword\">FPV</em>#<em class=\"keyword\">穿越机</em>#威海#那香海#视频随拍#零散记录','-','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5932fba5285890809589813076/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d5932fba5285890809589813076/f0.mp4','202009101933004667',1,0,0,0,0,52208,'','',854,480,1605777587,1605777587,5285890809589813076,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d5932fba5285890809589813076/v.f100020.mp4\",\"size\":2972453,\"duration\":52200}]'),
	(120,'LDARC <em class=\"keyword\">FPV穿越机</em>绿洲无线模拟器遥控器设置视频教程','LDARC 绿洲无线航模模拟器支持DCL，DRL，FPV Air2 , Freerider, LIFTOFF , Rotor Rush , NEXT-CGM等。','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8504281b5285890809595473951/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/8504281b5285890809595473951/f0.mp4','202009101933004667',1,0,0,0,0,52208,'','',854,480,1605777587,1605777587,5285890809595473951,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8504281b5285890809595473951/v.f100020.mp4\",\"size\":2972519,\"duration\":52200}]'),
	(121,'咯咯咯','哈哈哈','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/8504289f5285890809595473991/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/8504289f5285890809595473991/f0.mp4','202009181548217779',1,0,0,0,0,52208,'','',854,480,1605778186,1605778186,5285890809595473991,4372373,''),
	(122,'咯咯咯','哈哈哈','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59333bb5285890809589813158/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d59333bb5285890809589813158/f0.mp4','202009101933004667',1,0,0,0,0,52208,'','',854,480,1605778426,1605778426,5285890809589813158,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59333bb5285890809589813158/v.f100020.mp4\",\"size\":2972562,\"duration\":52200}]'),
	(123,'FPV穿越机\n4K航拍-广汉连山麦田日出','FPV穿越机\n4K航拍-广汉连山麦田日出','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59337b35285890809589813231/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d59337b35285890809589813231/f0.mp4','202009181548217779',1,0,0,0,0,52208,'','',854,480,1605779700,1605779700,5285890809589813231,4372373,''),
	(124,'FPV穿越机\n4K航拍-广汉连山麦田日出','FPV穿越机\n4K航拍-广汉连山麦田日出','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/850491d35285890809595474105/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/850491d35285890809595474105/f0.mp4','202009181548217779',1,0,0,0,0,52208,'','',854,480,1605779700,1605779700,5285890809595474105,4372373,''),
	(125,'FPV穿越机\n4K航拍-广汉连山麦田日出','FPV穿越机\n4K航拍-广汉连山麦田日出','https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59697405285890809589819733/coverBySnapshot/coverBySnapshot_10_0.jpg','https://fpv-vedios.youzu.com/5aafee04vodsh1251103923/d59697405285890809589819733/f0.mp4','202009101933004667',1,0,0,0,0,52208,'','',854,480,1605838517,1605838517,5285890809589819733,4372373,'[{\"type\":\"2\",\"url\":\"https://fpv-vedios.youzu.com/3199dbacvodtranssh1251103923/d59697405285890809589819733/v.f100020.mp4\",\"size\":2972823,\"duration\":52200}]');

/*!40000 ALTER TABLE `videos` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table videos_examine
# ------------------------------------------------------------

CREATE TABLE `videos_examine` (
  `video_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '视频id',
  `title` mediumtext COMMENT '视频标题',
  `describe` mediumtext COMMENT '视频描述',
  `cover` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频封面',
  `video_addr` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频地址',
  `label_id` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '标签id',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `user_type` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '添加用户类型（0：管理员[sys_user]，1：用户[user]）',
  `sortorder` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '排序',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '审核状态（0：无操作，1：审核通过 2：审核不通过 3：删除）',
  `is_recommend` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '推荐（0：不推荐；1：推荐）',
  `is_top` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '置顶（0：不置顶；1：置顶；）',
  `video_duration` int(8) NOT NULL DEFAULT '0' COMMENT '视频时长（单位：秒）',
  `label_name` mediumtext COMMENT '标签',
  `rec_content` mediumtext COMMENT '推荐理由',
  `top_content` mediumtext COMMENT '置顶理由',
  `video_width` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频宽',
  `video_height` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频高',
  `manager` int(11) NOT NULL DEFAULT '0' COMMENT '后台操作用户',
  `create_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`video_id`),
  KEY `user_id` (`user_id`) USING BTREE,
  FULLTEXT KEY `label_name` (`label_name`),
  FULLTEXT KEY `label_id` (`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='视频审核表';



# Dump of table world_map
# ------------------------------------------------------------

CREATE TABLE `world_map` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `name` varchar(50) NOT NULL COMMENT '国家(省份/城市)名称',
  `code` char(4) NOT NULL COMMENT '国家(省份/城市)编码',
  `pid` int(11) NOT NULL DEFAULT '0' COMMENT '父级id',
  `layer` tinyint(2) NOT NULL DEFAULT '0' COMMENT '层级',
  `sortorder` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 展示 1 隐藏',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `pid` (`pid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='世界国家/省份/城市 信息';

LOCK TABLES `world_map` WRITE;
/*!40000 ALTER TABLE `world_map` DISABLE KEYS */;

INSERT INTO `world_map` (`id`, `name`, `code`, `pid`, `layer`, `sortorder`, `status`, `create_at`, `update_at`)
VALUES
	(1,'中国','CHN',0,1,4,0,0,0),
	(2,'阿尔巴尼亚','ALB',0,1,0,0,0,0),
	(3,'阿尔及利亚','DZA',0,1,0,0,0,0),
	(4,'阿富汗','AFG',0,1,0,0,0,0),
	(5,'阿根廷','ARG',0,1,0,0,0,0),
	(6,'阿联酋','ARE',0,1,0,0,0,0),
	(7,'阿鲁巴','ABW',0,1,0,0,0,0),
	(8,'阿曼','OMN',0,1,0,0,0,0),
	(9,'阿塞拜疆','AZE',0,1,0,0,0,0),
	(10,'阿森松岛','ASC',0,1,0,0,0,0),
	(11,'埃及','EGY',0,1,0,0,0,0),
	(12,'埃塞俄比亚','ETH',0,1,0,0,0,0),
	(13,'爱尔兰','IRL',0,1,0,0,0,0),
	(14,'爱沙尼亚','EST',0,1,0,0,0,0),
	(15,'安道尔','AND',0,1,0,0,0,0),
	(16,'安哥拉','AGO',0,1,0,0,0,0),
	(17,'安圭拉','AIA',0,1,0,0,0,0),
	(18,'安提瓜岛和巴布达','ATG',0,1,0,0,0,0),
	(19,'奥地利','AUT',0,1,0,0,0,0),
	(20,'奥兰群岛','ALA',0,1,0,0,0,0),
	(21,'澳大利亚','AUS',0,1,0,0,0,0),
	(22,'巴巴多斯岛','BRB',0,1,0,0,0,0),
	(23,'巴布亚新几内亚','PNG',0,1,0,0,0,0),
	(24,'巴哈马','BHS',0,1,0,0,0,0),
	(25,'巴基斯坦','PAK',0,1,0,0,0,0),
	(26,'巴拉圭','PRY',0,1,0,0,0,0),
	(27,'巴勒斯坦','PSE',0,1,0,0,0,0),
	(28,'巴林','BHR',0,1,0,0,0,0),
	(29,'巴拿马','PAN',0,1,0,0,0,0),
	(30,'巴西','BRA',0,1,0,0,0,0),
	(31,'白俄罗斯','BLR',0,1,0,0,0,0),
	(32,'百慕大','BMU',0,1,0,0,0,0),
	(33,'保加利亚','BGR',0,1,0,0,0,0),
	(34,'北马里亚纳群岛','MNP',0,1,0,0,0,0),
	(35,'贝宁','BEN',0,1,0,0,0,0),
	(36,'比利时','BEL',0,1,0,0,0,0),
	(37,'冰岛','ISL',0,1,0,0,0,0),
	(38,'波多黎各','PRI',0,1,0,0,0,0),
	(39,'波兰','POL',0,1,0,0,0,0),
	(40,'波斯尼亚和黑塞哥维那','BIH',0,1,0,0,0,0),
	(41,'玻利维亚','BOL',0,1,0,0,0,0),
	(42,'伯利兹','BLZ',0,1,0,0,0,0),
	(43,'博茨瓦纳','BWA',0,1,0,0,0,0),
	(44,'不丹','BTN',0,1,0,0,0,0),
	(45,'布基纳法索','BFA',0,1,0,0,0,0),
	(46,'布隆迪','BDI',0,1,0,0,0,0),
	(47,'布韦岛','BVT',0,1,0,0,0,0),
	(48,'朝鲜','PRK',0,1,0,0,0,0),
	(49,'丹麦','DNK',0,1,0,0,0,0),
	(50,'德国','DEU',0,1,0,0,0,0),
	(51,'东帝汶','TLS',0,1,0,0,0,0),
	(52,'多哥','TGO',0,1,0,0,0,0),
	(53,'多米尼加','DMA',0,1,0,0,0,0),
	(54,'多米尼加共和国','DOM',0,1,0,0,0,0),
	(55,'俄罗斯','RUS',0,1,0,0,0,0),
	(56,'厄瓜多尔','ECU',0,1,0,0,0,0),
	(57,'厄立特里亚','ERI',0,1,0,0,0,0),
	(58,'法国','FRA',0,1,0,0,0,0),
	(59,'法罗群岛','FRO',0,1,0,0,0,0),
	(60,'法属波利尼西亚','PYF',0,1,0,0,0,0),
	(61,'法属圭亚那','GUF',0,1,0,0,0,0),
	(62,'法属南部领地','ATF',0,1,0,0,0,0),
	(63,'梵蒂冈','VAT',0,1,0,0,0,0),
	(64,'菲律宾','PHL',0,1,0,0,0,0),
	(65,'斐济','FJI',0,1,0,0,0,0),
	(66,'芬兰','FIN',0,1,0,0,0,0),
	(67,'佛得角','CPV',0,1,0,0,0,0),
	(68,'弗兰克群岛','FLK',0,1,0,0,0,0),
	(69,'冈比亚','GMB',0,1,0,0,0,0),
	(70,'刚果共和国','COG',0,1,0,0,0,0),
	(71,'刚果民主共和国','COD',0,1,0,0,0,0),
	(72,'哥伦比亚','COL',0,1,0,0,0,0),
	(73,'哥斯达黎加','CRI',0,1,0,0,0,0),
	(74,'格恩西岛','GGY',0,1,0,0,0,0),
	(75,'格林纳达','GRD',0,1,0,0,0,0),
	(76,'格陵兰','GRL',0,1,0,0,0,0),
	(77,'格鲁吉亚','GEG',0,1,0,0,0,0),
	(78,'古巴','CUB',0,1,0,0,0,0),
	(79,'瓜德罗普','GLP',0,1,0,0,0,0),
	(80,'关岛','GUM',0,1,0,0,0,0),
	(81,'圭亚那','GUY',0,1,0,0,0,0),
	(82,'哈萨克斯坦','KAZ',0,1,0,0,0,0),
	(83,'海地','HTI',0,1,0,0,0,0),
	(84,'韩国','KOR',0,1,0,0,0,0),
	(85,'荷兰','NLD',0,1,0,0,0,0),
	(86,'荷属安地列斯','ANT',0,1,0,0,0,0),
	(87,'赫德和麦克唐纳群岛','HMD',0,1,0,0,0,0),
	(88,'黑山','MEG',0,1,0,0,0,0),
	(89,'洪都拉斯','HND',0,1,0,0,0,0),
	(90,'基里巴斯','KIR',0,1,0,0,0,0),
	(91,'吉布提','DJI',0,1,0,0,0,0),
	(92,'吉尔吉斯斯坦','KGZ',0,1,0,0,0,0),
	(93,'几内亚','GIN',0,1,0,0,0,0),
	(94,'几内亚比绍','GNB',0,1,0,0,0,0),
	(95,'加拿大','CAN',0,1,0,0,0,0),
	(96,'加纳','GHA',0,1,0,0,0,0),
	(97,'加蓬','GAB',0,1,0,0,0,0),
	(98,'柬埔寨','KHM',0,1,0,0,0,0),
	(99,'捷克','CZE',0,1,0,0,0,0),
	(100,'津巴布韦','ZWE',0,1,0,0,0,0),
	(101,'喀麦隆','CMR',0,1,0,0,0,0),
	(102,'卡塔尔','QAT',0,1,0,0,0,0),
	(103,'开曼群岛','CYM',0,1,0,0,0,0),
	(104,'科科斯群岛','CCK',0,1,0,0,0,0),
	(105,'科摩罗','COM',0,1,0,0,0,0),
	(106,'科特迪瓦','CIV',0,1,0,0,0,0),
	(107,'科威特','KWT',0,1,0,0,0,0),
	(108,'克罗地亚','HRV',0,1,0,0,0,0),
	(109,'肯尼亚','KEN',0,1,0,0,0,0),
	(110,'库克群岛','COK',0,1,0,0,0,0),
	(111,'拉脱维亚','LVA',0,1,0,0,0,0),
	(112,'莱索托','LSO',0,1,0,0,0,0),
	(113,'老挝','LAO',0,1,0,0,0,0),
	(114,'黎巴嫩','LBN',0,1,0,0,0,0),
	(115,'立陶宛','LTU',0,1,0,0,0,0),
	(116,'利比里亚','LBR',0,1,0,0,0,0),
	(117,'利比亚','LBY',0,1,0,0,0,0),
	(118,'列支敦士登','LIE',0,1,0,0,0,0),
	(119,'留尼汪岛','REU',0,1,0,0,0,0),
	(120,'卢森堡','LUX',0,1,0,0,0,0),
	(121,'卢旺达','RWA',0,1,0,0,0,0),
	(122,'罗马尼亚','ROU',0,1,0,0,0,0),
	(123,'马达加斯加','MDG',0,1,0,0,0,0),
	(124,'马尔代夫','MDV',0,1,0,0,0,0),
	(125,'马耳他','MLT',0,1,0,0,0,0),
	(126,'马拉维','MWI',0,1,0,0,0,0),
	(127,'马来西亚','MYS',0,1,0,0,0,0),
	(128,'马里','MLI',0,1,0,0,0,0),
	(129,'马其顿','MKD',0,1,0,0,0,0),
	(130,'马绍尔群岛','MHL',0,1,0,0,0,0),
	(131,'马提尼克','MTQ',0,1,0,0,0,0),
	(132,'马约特岛','MYT',0,1,0,0,0,0),
	(133,'曼岛','IMN',0,1,0,0,0,0),
	(134,'毛里求斯','MUS',0,1,0,0,0,0),
	(135,'毛里塔尼亚','MRT',0,1,0,0,0,0),
	(136,'美国','USA',0,1,0,0,0,0),
	(137,'美属萨摩亚','ASM',0,1,0,0,0,0),
	(138,'美属外岛','UMI',0,1,0,0,0,0),
	(139,'蒙古','MNG',0,1,0,0,0,0),
	(140,'蒙特塞拉特','MSR',0,1,0,0,0,0),
	(141,'孟加拉国','BGD',0,1,0,0,0,0),
	(142,'秘鲁','PER',0,1,0,0,0,0),
	(143,'密克罗尼西亚','FSM',0,1,0,0,0,0),
	(144,'缅甸','MMR',0,1,0,0,0,0),
	(145,'摩尔多瓦','MDA',0,1,0,0,0,0),
	(146,'摩洛哥','MAR',0,1,0,0,0,0),
	(147,'摩纳哥','MCO',0,1,0,0,0,0),
	(148,'莫桑比克','MOZ',0,1,0,0,0,0),
	(149,'墨西哥','MEX',0,1,0,0,0,0),
	(150,'纳米比亚','NAM',0,1,0,0,0,0),
	(151,'南非','ZAF',0,1,0,0,0,0),
	(152,'南极洲','ATA',0,1,0,0,0,0),
	(153,'南乔治亚和南桑德威奇群岛','SGS',0,1,0,0,0,0),
	(154,'瑙鲁','NRU',0,1,0,0,0,0),
	(155,'尼泊尔','NPL',0,1,0,0,0,0),
	(156,'尼加拉瓜','NIC',0,1,0,0,0,0),
	(157,'尼日尔','NER',0,1,0,0,0,0),
	(158,'尼日利亚','NGA',0,1,0,0,0,0),
	(159,'纽埃','NIU',0,1,0,0,0,0),
	(160,'挪威','NOR',0,1,0,0,0,0),
	(161,'诺福克','NFK',0,1,0,0,0,0),
	(162,'帕劳群岛','PLW',0,1,0,0,0,0),
	(163,'皮特凯恩','PCN',0,1,0,0,0,0),
	(164,'葡萄牙','PRT',0,1,0,0,0,0),
	(165,'乔治亚','GEO',0,1,0,0,0,0),
	(166,'日本','JPN',0,1,0,0,0,0),
	(167,'瑞典','SWE',0,1,0,0,0,0),
	(168,'瑞士','CHE',0,1,0,0,0,0),
	(169,'萨尔瓦多','SLV',0,1,0,0,0,0),
	(170,'萨摩亚','WSM',0,1,0,0,0,0),
	(171,'塞尔维亚','SCG',0,1,0,0,0,0),
	(172,'塞拉利昂','SLE',0,1,0,0,0,0),
	(173,'塞内加尔','SEN',0,1,0,0,0,0),
	(174,'塞浦路斯','CYP',0,1,0,0,0,0),
	(175,'塞舌尔','SYC',0,1,0,0,0,0),
	(176,'沙特阿拉伯','SAU',0,1,0,0,0,0),
	(177,'圣诞岛','CXR',0,1,0,0,0,0),
	(178,'圣多美和普林西比','STP',0,1,0,0,0,0),
	(179,'圣赫勒拿','SHN',0,1,0,0,0,0),
	(180,'圣基茨和尼维斯','KNA',0,1,0,0,0,0),
	(181,'圣卢西亚','LCA',0,1,0,0,0,0),
	(182,'圣马力诺','SMR',0,1,0,0,0,0),
	(183,'圣皮埃尔和米克隆群岛','SPM',0,1,0,0,0,0),
	(184,'圣文森特和格林纳丁斯','VCT',0,1,0,0,0,0),
	(185,'斯里兰卡','LKA',0,1,0,0,0,0),
	(186,'斯洛伐克','SVK',0,1,0,0,0,0),
	(187,'斯洛文尼亚','SVN',0,1,0,0,0,0),
	(188,'斯瓦尔巴和扬马廷','SJM',0,1,0,0,0,0),
	(189,'斯威士兰','SWZ',0,1,0,0,0,0),
	(190,'苏丹','SDN',0,1,0,0,0,0),
	(191,'苏里南','SUR',0,1,0,0,0,0),
	(192,'所罗门群岛','SLB',0,1,0,0,0,0),
	(193,'索马里','SOM',0,1,0,0,0,0),
	(194,'塔吉克斯坦','TJK',0,1,0,0,0,0),
	(195,'泰国','THA',0,1,0,0,0,0),
	(196,'坦桑尼亚','TZA',0,1,0,0,0,0),
	(197,'汤加','TON',0,1,0,0,0,0),
	(198,'特克斯和凯克特斯群岛','TCA',0,1,0,0,0,0),
	(199,'特里斯坦达昆哈','TAA',0,1,0,0,0,0),
	(200,'特立尼达和多巴哥','TTO',0,1,0,0,0,0),
	(201,'突尼斯','TUN',0,1,0,0,0,0),
	(202,'图瓦卢','TUV',0,1,0,0,0,0),
	(203,'土耳其','TUR',0,1,0,0,0,0),
	(204,'土库曼斯坦','TKM',0,1,0,0,0,0),
	(205,'托克劳','TKL',0,1,0,0,0,0),
	(206,'瓦利斯和福图纳','WLF',0,1,0,0,0,0),
	(207,'瓦努阿图','VUT',0,1,0,0,0,0),
	(208,'危地马拉','GTM',0,1,0,0,0,0),
	(209,'美属维尔京群岛','VIR',0,1,0,0,0,0),
	(210,'英属维尔京群岛','VGB',0,1,0,0,0,0),
	(211,'委内瑞拉','VEN',0,1,0,0,0,0),
	(212,'文莱','BRN',0,1,0,0,0,0),
	(213,'乌干达','UGA',0,1,0,0,0,0),
	(214,'乌克兰','UKR',0,1,0,0,0,0),
	(215,'乌拉圭','URY',0,1,0,0,0,0),
	(216,'乌兹别克斯坦','UZB',0,1,0,0,0,0),
	(217,'西班牙','ESP',0,1,0,0,0,0),
	(218,'希腊','GRC',0,1,0,0,0,0),
	(219,'新加坡','SGP',0,1,0,0,0,0),
	(220,'新喀里多尼亚','NCL',0,1,0,0,0,0),
	(221,'新西兰','NZL',0,1,0,0,0,0),
	(222,'匈牙利','HUN',0,1,0,0,0,0),
	(223,'叙利亚','SYR',0,1,0,0,0,0),
	(224,'牙买加','JAM',0,1,0,0,0,0),
	(225,'亚美尼亚','ARM',0,1,0,0,0,0),
	(226,'也门','YEM',0,1,0,0,0,0),
	(227,'伊拉克','IRQ',0,1,0,0,0,0),
	(228,'伊朗','IRN',0,1,0,0,0,0),
	(229,'以色列','ISR',0,1,0,0,0,0),
	(230,'意大利','ITA',0,1,0,0,0,0),
	(231,'印度','IND',0,1,0,0,0,0),
	(232,'印度尼西亚','IDN',0,1,0,0,0,0),
	(233,'英国','GBR',0,1,0,0,0,0),
	(234,'英属印度洋领地','IOT',0,1,0,0,0,0),
	(235,'约旦','JOR',0,1,0,0,0,0),
	(236,'越南','VNM',0,1,0,0,0,0),
	(237,'赞比亚','ZMB',0,1,0,0,0,0),
	(238,'泽西岛','JEY',0,1,0,0,0,0),
	(239,'乍得','TCD',0,1,0,0,0,0),
	(240,'直布罗陀','GIB',0,1,0,0,0,0),
	(241,'智利','CHL',0,1,0,0,0,0),
	(242,'中非共和国','CAF',0,1,0,0,0,0),
	(243,'香港','CHN',0,1,3,0,0,0),
	(244,'澳门','CHN',0,1,2,0,0,0),
	(245,'台湾','CHN',0,1,1,0,0,0);

/*!40000 ALTER TABLE `world_map` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table ycoin_receive_record
# ------------------------------------------------------------

CREATE TABLE `ycoin_receive_record` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(128) NOT NULL COMMENT '任务名称',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `explain` varchar(256) NOT NULL DEFAULT '' COMMENT '任务描述',
  `task_type` tinyint(1) NOT NULL COMMENT '任务类型 1.登陆应用',
  `ycoin` int(11) NOT NULL COMMENT '获取游币数',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态，0为展示 1为不展示',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户获取游币记录表';



# Dump of table ycoin_reward_record
# ------------------------------------------------------------

CREATE TABLE `ycoin_reward_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `giver_id` varchar(60) NOT NULL COMMENT '赠予人uid',
  `donne_id` varchar(60) NOT NULL COMMENT '受赠人uid',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1-正常展示, 2-不展示',
  `reward_ycoin` int(11) NOT NULL COMMENT '打赏的游币数量',
  `record_id` varchar(100) NOT NULL DEFAULT ' ' COMMENT '打赏的视频记录id（暂时只有视频可以打赏）',
  `reward_type` tinyint(1) NOT NULL COMMENT '打赏类型 1 打赏视频',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `giver_id` (`giver_id`),
  KEY `donne_id` (`donne_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='游币打赏记录表';



# Dump of table ycoin_task
# ------------------------------------------------------------

CREATE TABLE `ycoin_task` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(128) NOT NULL COMMENT '任务名称',
  `explain` varchar(256) NOT NULL DEFAULT '' COMMENT '说明',
  `task_type` tinyint(1) NOT NULL COMMENT '任务类型 1.登陆应用',
  `ycoin` int(11) NOT NULL COMMENT '每次获取游币数',
  `count` int(6) NOT NULL COMMENT '限制次数 0没有限制',
  `period_limit` tinyint(1) NOT NULL COMMENT '限购周期 0永久 1.1天 2.1周 3.1月 4.1年 ',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态，0为关闭任务 1为开启任务',
  `task_icon` varchar(256) NOT NULL DEFAULT '' COMMENT '任务图标',
  `describe` varchar(256) NOT NULL DEFAULT '' COMMENT '描述',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COMMENT='游币获取任务表';




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
