# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.0.12)
# Database: fpv2
# Generation Time: 2020-09-21 01:59:48 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


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
  `start_time` int(11) NOT NULL DEFAULT '0',
  `end_time` int(11) NOT NULL DEFAULT '0',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0待上架 1上架 2下架',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='banner配置表';



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
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `to_user_id` varchar(60) NOT NULL COMMENT '作品发布者用户id',
  `compose_id` bigint(20) NOT NULL COMMENT '作品id',
  `compose_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0 视频 1 帖子',
  `status` tinyint(1) NOT NULL COMMENT '1 收藏 2 取消收藏',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `video_id` (`compose_id`),
  KEY `to_user_id` (`to_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='收藏的视频';



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



# Dump of table feedback
# ------------------------------------------------------------

CREATE TABLE `feedback` (
  `id` bigint(20) NOT NULL COMMENT '自增主键',
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
  `hot_search_content` varchar(128) NOT NULL COMMENT '热门搜索内容 多个用逗号分隔 例如：FPV,电竞,无人机',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='热门搜索（后台配置）';



# Dump of table post_comment
# ------------------------------------------------------------

CREATE TABLE `post_comment` (
  `comment_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '评论id',
  `post_id` bigint(20) NOT NULL COMMENT '帖子id',
  `content` varchar(1024) NOT NULL COMMENT '内容',
  `from_uid` varchar(60) NOT NULL COMMENT '评论内容的用户id',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0展示 1不展示',
  PRIMARY KEY (`comment_id`),
  KEY `post_id` (`post_id`),
  KEY `from_uid` (`from_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='帖子评论表';



# Dump of table post_comment_parent_child
# ------------------------------------------------------------

CREATE TABLE `post_comment_parent_child` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `parent_id` bigint(20) NOT NULL COMMENT '评论父id',
  `child_id` bigint(20) NOT NULL COMMENT '评论子id',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `post_id` (`parent_id`),
  KEY `from_uid` (`child_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='帖子评论表中各个评论的父子关系';



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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='帖子表';



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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='帖子审核表';



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
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='帖子标签配置';



# Dump of table post_labels
# ------------------------------------------------------------

CREATE TABLE `post_labels` (
  `post_id` bigint(20) unsigned NOT NULL COMMENT '帖子id',
  `label_id` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标签id',
  `label_name` varchar(521) DEFAULT NULL COMMENT '标签名',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`post_id`,`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='帖子拥有的标签表';



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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='帖子相关参数统计';



# Dump of table received_at
# ------------------------------------------------------------

CREATE TABLE `received_at` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `to_user_id` varchar(60) NOT NULL COMMENT '被@的用户id',
  `user_id` varchar(60) NOT NULL COMMENT '执行@的用户id',
  `comment_id` bigint(20) NOT NULL COMMENT '评论id',
  `topic_type` tinyint(2) NOT NULL COMMENT '1.视频 2.帖子 3.评论',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `to_user_id` (`to_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='收到的@';



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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户搜索历史表';



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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='所属圈子的板块配置表（后台设置）';



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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户分享记录';



# Dump of table social_account_login
# ------------------------------------------------------------

CREATE TABLE `social_account_login` (
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `unionid` varchar(256) NOT NULL DEFAULT '' COMMENT '社交平台关联id',
  `social_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '区分社交软件 1 微信关联id 2 微博关联id 3 qq关联id',
  `status` tinyint(1) DEFAULT '0' COMMENT '0 正常 1 封禁',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`user_id`,`social_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='社交平台登陆表';



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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='圈子表';



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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统记录';



# Dump of table system_message
# ------------------------------------------------------------

CREATE TABLE `system_message` (
  `system_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '系统通知ID',
  `send_id` int(11) NOT NULL COMMENT '发送者ID（后台用户id）',
  `receive_id` varchar(60) NOT NULL COMMENT '接收者id',
  `send_default` tinyint(2) NOT NULL COMMENT '1时发送所有用户，0时则不采用',
  `system_topic` varchar(60) NOT NULL COMMENT '通知标题',
  `system_content` varchar(255) NOT NULL COMMENT '通知内容',
  `send_time` int(11) NOT NULL COMMENT '发送时间',
  `expire_time` int(11) NOT NULL COMMENT '过期时间',
  `send_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0.默认为系统通知 1.收到@通知 2.收到点赞通知 3.收到收藏通知 4.收到分享通知 5.收到评论/回复通知 6.特殊业务的系统奖励通知 7.活动通知',
  `extra` varchar(2056) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT ' ' COMMENT '附件内容 例如：奖励',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 未读 1 已读  默认未读',
  PRIMARY KEY (`system_id`),
  KEY `receive_id` (`receive_id`),
  KEY `send_type` (`send_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统消息';



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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户系统推送通知设置表';

LOCK TABLES `system_notice_settings` WRITE;
/*!40000 ALTER TABLE `system_notice_settings` DISABLE KEYS */;

INSERT INTO `system_notice_settings` (`user_id`, `comment_push_set`, `thumb_up_push_set`, `attention_push_set`, `share_push_set`, `slot_push_set`, `create_at`, `update_at`)
VALUES
	('202009181540233466',0,0,0,0,0,1600414824,1600414824),
	('202009181543188713',0,0,0,0,0,1600414998,1600414998),
	('202009181544573421',0,0,0,0,0,1600415097,1600415097),
	('202009181548217779',0,0,0,0,0,1600415301,1600415301);

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统推送';



# Dump of table thumbs_up
# ------------------------------------------------------------

CREATE TABLE `thumbs_up` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `type_id` bigint(20) NOT NULL COMMENT '作品id （视频id/帖子id/评论id）',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `to_user_id` varchar(60) NOT NULL COMMENT '被点赞的用户id',
  `zan_type` tinyint(2) NOT NULL COMMENT '1 视频点赞 2 帖子点赞 3 评论点赞',
  `status` tinyint(1) NOT NULL COMMENT '1赞 2取消点赞',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `to_user_id` (`to_user_id`),
  KEY `type_id` (`type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='点赞表（针对帖子/视频/评论）';

LOCK TABLES `thumbs_up` WRITE;
/*!40000 ALTER TABLE `thumbs_up` DISABLE KEYS */;

INSERT INTO `thumbs_up` (`id`, `type_id`, `user_id`, `to_user_id`, `zan_type`, `status`, `create_at`)
VALUES
	(1,2,'202009181548217779','202009101933004667',1,2,1600426777),
	(2,10,'202009181548217779','202009101933004667',1,1,1600427096),
	(3,11,'202009181548217779','202009101933004667',1,1,1600427161),
	(4,12,'202009181548217779','202009101933004667',1,2,1600427250),
	(5,13,'202009181548217779','202009101933004667',1,1,1600427402);

/*!40000 ALTER TABLE `thumbs_up` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user
# ------------------------------------------------------------

CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `nick_name` varchar(45) NOT NULL DEFAULT '' COMMENT '昵称',
  `mobile_num` bigint(20) NOT NULL COMMENT '手机号码',
  `password` varchar(128) NOT NULL COMMENT '用户密码',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `gender` tinyint(1) NOT NULL DEFAULT '0' COMMENT '性别 0男性 1女性',
  `born` varchar(128) NOT NULL DEFAULT '' COMMENT '出生日期',
  `age` int(3) NOT NULL DEFAULT '0' COMMENT '年龄',
  `avatar` varchar(100) NOT NULL DEFAULT '' COMMENT '头像',
  `status` tinyint(1) DEFAULT '0' COMMENT '0 正常 1 封禁',
  `last_login_time` int(11) DEFAULT NULL COMMENT '最后登录时间',
  `signature` varchar(200) NOT NULL DEFAULT '' COMMENT '签名',
  `device_type` tinyint(2) DEFAULT NULL COMMENT '设备类型 0 android 1 iOS 2 web',
  `city` varchar(64) NOT NULL DEFAULT '' COMMENT '城市',
  `is_anchor` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0不是主播 1为主播',
  `channel_id` int(11) NOT NULL DEFAULT '0' COMMENT '渠道id',
  `background_img` varchar(255) NOT NULL DEFAULT '' COMMENT '背景图',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '称号/特殊身份',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `user_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '用户类型 0 手机号 1 微信 2 QQ 3 微博',
  `country` int(3) NOT NULL DEFAULT '0' COMMENT '国家',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户表';

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;

INSERT INTO `user` (`id`, `nick_name`, `mobile_num`, `password`, `user_id`, `gender`, `born`, `age`, `avatar`, `status`, `last_login_time`, `signature`, `device_type`, `city`, `is_anchor`, `channel_id`, `background_img`, `title`, `create_at`, `update_at`, `user_type`, `country`)
VALUES
	(2,'用户6222',13177656222,'','202009101933004667',0,'',0,'',0,1599737580,'',0,'',0,0,'','',1599737580,1599737580,0,0),
	(7,'用户6666',13177666666,'','202009181548217779',0,'',0,'',0,1600415301,'',0,'',0,0,'','',1600415301,1600415301,0,0);

/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user_attention
# ------------------------------------------------------------

CREATE TABLE `user_attention` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(60) NOT NULL COMMENT '被关注的用户id',
  `attention_uid` varchar(60) NOT NULL COMMENT '关注的用户id',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1表示关注 2表示取消关注',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `attention_uid` (`attention_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户关注表';

LOCK TABLES `user_attention` WRITE;
/*!40000 ALTER TABLE `user_attention` DISABLE KEYS */;

INSERT INTO `user_attention` (`id`, `user_id`, `attention_uid`, `status`, `create_at`)
VALUES
	(2,'202009101933004667','202009181548217779',1,1600419087);

/*!40000 ALTER TABLE `user_attention` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user_browse_record
# ------------------------------------------------------------

CREATE TABLE `user_browse_record` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `compose_id` bigint(20) NOT NULL COMMENT '作品id',
  `compose_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0 视频 1 帖子',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户浏览过的作品记录（包含帖子、视频）';



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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户货币表';



# Dump of table video_barrage
# ------------------------------------------------------------

CREATE TABLE `video_barrage` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
  `video_cur_duration` int(8) NOT NULL COMMENT '视频当前时长节点（单位：秒）',
  `content` varchar(512) NOT NULL DEFAULT '' COMMENT '弹幕内容',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `color` varchar(100) NOT NULL DEFAULT '' COMMENT '弹幕字体颜色',
  `font` varchar(100) NOT NULL DEFAULT '' COMMENT '弹幕字体',
  `barrage_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '预留字段',
  `location` tinyint(2) NOT NULL DEFAULT '0' COMMENT '弹幕位置',
  `send_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '弹幕发送时间',
  PRIMARY KEY (`id`),
  KEY `video_id` (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频弹幕';



# Dump of table video_comment
# ------------------------------------------------------------

CREATE TABLE `video_comment` (
  `comment_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '评论id',
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `content` varchar(1024) NOT NULL COMMENT '内容',
  `from_uid` varchar(60) NOT NULL COMMENT '评论内容的用户id',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0展示 1不展示',
  PRIMARY KEY (`comment_id`),
  KEY `video_id` (`video_id`),
  KEY `from_uid` (`from_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频评论表';



# Dump of table video_comment_parent_child
# ------------------------------------------------------------

CREATE TABLE `video_comment_parent_child` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `parent_id` bigint(20) NOT NULL COMMENT '评论父id',
  `child_id` bigint(20) NOT NULL COMMENT '评论子id',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `post_id` (`parent_id`),
  KEY `from_uid` (`child_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频评论表中各个评论的父子关系';



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
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频标签配置';



# Dump of table video_labels
# ------------------------------------------------------------

CREATE TABLE `video_labels` (
  `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
  `label_id` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标签id',
  `label_name` varchar(521) DEFAULT NULL COMMENT '标签名',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`video_id`,`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='视频拥有的标签表';

LOCK TABLES `video_labels` WRITE;
/*!40000 ALTER TABLE `video_labels` DISABLE KEYS */;

INSERT INTO `video_labels` (`video_id`, `label_id`, `label_name`, `create_at`)
VALUES
	(0,'1','todo:通过标签id获取名称',1600422665),
	(0,'2','todo:通过标签id获取名称',1600422665),
	(0,'3','todo:通过标签id获取名称',1600422665),
	(10,'1','todo:通过标签id获取名称',1600426884),
	(10,'2','todo:通过标签id获取名称',1600426884),
	(10,'3','todo:通过标签id获取名称',1600426884),
	(11,'1','todo:通过标签id获取名称',1600426884),
	(11,'2','todo:通过标签id获取名称',1600426884),
	(11,'3','todo:通过标签id获取名称',1600426884),
	(12,'1','todo:通过标签id获取名称',1600426885),
	(12,'2','todo:通过标签id获取名称',1600426885),
	(12,'3','todo:通过标签id获取名称',1600426885),
	(13,'1','todo:通过标签id获取名称',1600426886),
	(13,'2','todo:通过标签id获取名称',1600426886),
	(13,'3','todo:通过标签id获取名称',1600426886),
	(14,'1','todo:通过标签id获取名称',1600426886),
	(14,'2','todo:通过标签id获取名称',1600426886),
	(14,'3','todo:通过标签id获取名称',1600426886),
	(15,'1','todo:通过标签id获取名称',1600426887),
	(15,'2','todo:通过标签id获取名称',1600426887),
	(15,'3','todo:通过标签id获取名称',1600426887),
	(16,'1','todo:通过标签id获取名称',1600426887),
	(16,'2','todo:通过标签id获取名称',1600426887),
	(16,'3','todo:通过标签id获取名称',1600426887),
	(17,'1','todo:通过标签id获取名称',1600426888),
	(17,'2','todo:通过标签id获取名称',1600426888),
	(17,'3','todo:通过标签id获取名称',1600426888),
	(18,'1','todo:通过标签id获取名称',1600426888),
	(18,'2','todo:通过标签id获取名称',1600426888),
	(18,'3','todo:通过标签id获取名称',1600426888),
	(19,'1','todo:通过标签id获取名称',1600426889),
	(19,'2','todo:通过标签id获取名称',1600426889),
	(19,'3','todo:通过标签id获取名称',1600426889),
	(20,'1','todo:通过标签id获取名称',1600426889),
	(20,'2','todo:通过标签id获取名称',1600426889),
	(20,'3','todo:通过标签id获取名称',1600426889),
	(21,'1','todo:通过标签id获取名称',1600426890),
	(21,'2','todo:通过标签id获取名称',1600426890),
	(21,'3','todo:通过标签id获取名称',1600426890),
	(22,'1','todo:通过标签id获取名称',1600426890),
	(22,'2','todo:通过标签id获取名称',1600426890),
	(22,'3','todo:通过标签id获取名称',1600426890),
	(23,'1','todo:通过标签id获取名称',1600426891),
	(23,'2','todo:通过标签id获取名称',1600426891),
	(23,'3','todo:通过标签id获取名称',1600426891),
	(24,'1','todo:通过标签id获取名称',1600426891),
	(24,'2','todo:通过标签id获取名称',1600426891),
	(24,'3','todo:通过标签id获取名称',1600426891),
	(25,'1','todo:通过标签id获取名称',1600426892),
	(25,'2','todo:通过标签id获取名称',1600426892),
	(25,'3','todo:通过标签id获取名称',1600426892),
	(26,'1','todo:通过标签id获取名称',1600426892),
	(26,'2','todo:通过标签id获取名称',1600426892),
	(26,'3','todo:通过标签id获取名称',1600426892),
	(27,'1','todo:通过标签id获取名称',1600426892),
	(27,'2','todo:通过标签id获取名称',1600426892),
	(27,'3','todo:通过标签id获取名称',1600426892),
	(28,'1','todo:通过标签id获取名称',1600426893),
	(28,'2','todo:通过标签id获取名称',1600426893),
	(28,'3','todo:通过标签id获取名称',1600426893),
	(29,'1','todo:通过标签id获取名称',1600426893),
	(29,'2','todo:通过标签id获取名称',1600426893),
	(29,'3','todo:通过标签id获取名称',1600426893),
	(30,'1','todo:通过标签id获取名称',1600426894),
	(30,'2','todo:通过标签id获取名称',1600426894),
	(30,'3','todo:通过标签id获取名称',1600426894),
	(31,'1','todo:通过标签id获取名称',1600426894),
	(31,'2','todo:通过标签id获取名称',1600426894),
	(31,'3','todo:通过标签id获取名称',1600426894),
	(32,'1','todo:通过标签id获取名称',1600426895),
	(32,'2','todo:通过标签id获取名称',1600426895),
	(32,'3','todo:通过标签id获取名称',1600426895),
	(33,'1','todo:通过标签id获取名称',1600426895),
	(33,'2','todo:通过标签id获取名称',1600426895),
	(33,'3','todo:通过标签id获取名称',1600426895),
	(34,'1','todo:通过标签id获取名称',1600426895),
	(34,'2','todo:通过标签id获取名称',1600426895),
	(34,'3','todo:通过标签id获取名称',1600426895),
	(35,'1','todo:通过标签id获取名称',1600426896),
	(35,'2','todo:通过标签id获取名称',1600426896),
	(35,'3','todo:通过标签id获取名称',1600426896),
	(36,'1','todo:通过标签id获取名称',1600426896),
	(36,'2','todo:通过标签id获取名称',1600426896),
	(36,'3','todo:通过标签id获取名称',1600426896),
	(37,'1','todo:通过标签id获取名称',1600426897),
	(37,'2','todo:通过标签id获取名称',1600426897),
	(37,'3','todo:通过标签id获取名称',1600426897),
	(38,'1','todo:通过标签id获取名称',1600426897),
	(38,'2','todo:通过标签id获取名称',1600426897),
	(38,'3','todo:通过标签id获取名称',1600426897),
	(39,'1','todo:通过标签id获取名称',1600426897),
	(39,'2','todo:通过标签id获取名称',1600426897),
	(39,'3','todo:通过标签id获取名称',1600426897),
	(40,'1','todo:通过标签id获取名称',1600426898),
	(40,'2','todo:通过标签id获取名称',1600426898),
	(40,'3','todo:通过标签id获取名称',1600426898),
	(41,'1','todo:通过标签id获取名称',1600426898),
	(41,'2','todo:通过标签id获取名称',1600426898),
	(41,'3','todo:通过标签id获取名称',1600426898),
	(42,'1','todo:通过标签id获取名称',1600426899),
	(42,'2','todo:通过标签id获取名称',1600426899),
	(42,'3','todo:通过标签id获取名称',1600426899),
	(43,'1','todo:通过标签id获取名称',1600426899),
	(43,'2','todo:通过标签id获取名称',1600426899),
	(43,'3','todo:通过标签id获取名称',1600426899),
	(44,'1','todo:通过标签id获取名称',1600426899),
	(44,'2','todo:通过标签id获取名称',1600426899),
	(44,'3','todo:通过标签id获取名称',1600426899),
	(45,'1','todo:通过标签id获取名称',1600426900),
	(45,'2','todo:通过标签id获取名称',1600426900),
	(45,'3','todo:通过标签id获取名称',1600426900),
	(46,'1','todo:通过标签id获取名称',1600426900),
	(46,'2','todo:通过标签id获取名称',1600426900),
	(46,'3','todo:通过标签id获取名称',1600426900),
	(47,'1','todo:通过标签id获取名称',1600426901),
	(47,'2','todo:通过标签id获取名称',1600426901),
	(47,'3','todo:通过标签id获取名称',1600426901),
	(48,'1','todo:通过标签id获取名称',1600426901),
	(48,'2','todo:通过标签id获取名称',1600426901),
	(48,'3','todo:通过标签id获取名称',1600426901),
	(49,'1','todo:通过标签id获取名称',1600426901),
	(49,'2','todo:通过标签id获取名称',1600426901),
	(49,'3','todo:通过标签id获取名称',1600426901),
	(50,'1','todo:通过标签id获取名称',1600426902),
	(50,'2','todo:通过标签id获取名称',1600426902),
	(50,'3','todo:通过标签id获取名称',1600426902),
	(51,'1','todo:通过标签id获取名称',1600426902),
	(51,'2','todo:通过标签id获取名称',1600426902),
	(51,'3','todo:通过标签id获取名称',1600426902),
	(52,'1','todo:通过标签id获取名称',1600426902),
	(52,'2','todo:通过标签id获取名称',1600426902),
	(52,'3','todo:通过标签id获取名称',1600426902),
	(53,'1','todo:通过标签id获取名称',1600426903),
	(53,'2','todo:通过标签id获取名称',1600426903),
	(53,'3','todo:通过标签id获取名称',1600426903),
	(54,'1','todo:通过标签id获取名称',1600426903),
	(54,'2','todo:通过标签id获取名称',1600426903),
	(54,'3','todo:通过标签id获取名称',1600426903),
	(55,'1','todo:通过标签id获取名称',1600426903),
	(55,'2','todo:通过标签id获取名称',1600426903),
	(55,'3','todo:通过标签id获取名称',1600426903),
	(56,'1','todo:通过标签id获取名称',1600426904),
	(56,'2','todo:通过标签id获取名称',1600426904),
	(56,'3','todo:通过标签id获取名称',1600426904),
	(57,'1','todo:通过标签id获取名称',1600426904),
	(57,'2','todo:通过标签id获取名称',1600426904),
	(57,'3','todo:通过标签id获取名称',1600426904),
	(58,'1','todo:通过标签id获取名称',1600426904),
	(58,'2','todo:通过标签id获取名称',1600426904),
	(58,'3','todo:通过标签id获取名称',1600426904),
	(59,'1','todo:通过标签id获取名称',1600426905),
	(59,'2','todo:通过标签id获取名称',1600426905),
	(59,'3','todo:通过标签id获取名称',1600426905),
	(60,'1','todo:通过标签id获取名称',1600426905),
	(60,'2','todo:通过标签id获取名称',1600426905),
	(60,'3','todo:通过标签id获取名称',1600426905),
	(61,'1','todo:通过标签id获取名称',1600426906),
	(61,'2','todo:通过标签id获取名称',1600426906),
	(61,'3','todo:通过标签id获取名称',1600426906),
	(62,'1','todo:通过标签id获取名称',1600426907),
	(62,'2','todo:通过标签id获取名称',1600426907),
	(62,'3','todo:通过标签id获取名称',1600426907),
	(63,'1','todo:通过标签id获取名称',1600426907),
	(63,'2','todo:通过标签id获取名称',1600426907),
	(63,'3','todo:通过标签id获取名称',1600426907),
	(64,'1','todo:通过标签id获取名称',1600426908),
	(64,'2','todo:通过标签id获取名称',1600426908),
	(64,'3','todo:通过标签id获取名称',1600426908),
	(65,'1','todo:通过标签id获取名称',1600426908),
	(65,'2','todo:通过标签id获取名称',1600426908),
	(65,'3','todo:通过标签id获取名称',1600426908);

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
  `sequence` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '推流唯一标识',
  PRIMARY KEY (`id`),
  KEY `anchor_id` (`anchor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='视频直播表';



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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频相关参数统计';

LOCK TABLES `video_statistic` WRITE;
/*!40000 ALTER TABLE `video_statistic` DISABLE KEYS */;

INSERT INTO `video_statistic` (`video_id`, `fabulous_num`, `browse_num`, `share_num`, `reward_num`, `barrage_num`, `comment_num`, `collect_num`, `create_at`, `update_at`)
VALUES
	(10,0,0,0,0,0,0,0,1600426884,1600427096),
	(11,2,0,0,0,0,0,0,1600426884,1600427161),
	(12,1,0,0,0,0,0,0,1600426885,1600427242),
	(13,1,0,0,0,0,0,0,1600426886,1600427402),
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
	(40,0,0,0,0,0,0,0,1600426898,1600426898),
	(41,0,0,0,0,0,0,0,1600426898,1600426898),
	(42,0,0,0,0,0,0,0,1600426899,1600426899),
	(43,0,0,0,0,0,0,0,1600426899,1600426899),
	(44,0,0,0,0,0,0,0,1600426899,1600426899),
	(45,0,0,0,0,0,0,0,1600426900,1600426900),
	(46,0,0,0,0,0,0,0,1600426900,1600426900),
	(47,0,0,0,0,0,0,0,1600426901,1600426901),
	(48,0,0,0,0,0,0,0,1600426901,1600426901),
	(49,0,0,0,0,0,0,0,1600426901,1600426901),
	(50,0,0,0,0,0,0,0,1600426902,1600426902),
	(51,0,0,0,0,0,0,0,1600426902,1600426902),
	(52,0,0,0,0,0,0,0,1600426902,1600426902),
	(53,0,0,0,0,0,0,0,1600426903,1600426903),
	(54,0,0,0,0,0,0,0,1600426903,1600426903),
	(55,0,0,0,0,0,0,0,1600426903,1600426903),
	(56,0,0,0,0,0,0,0,1600426904,1600426904),
	(57,0,0,0,0,0,0,0,1600426904,1600426904),
	(58,0,0,0,0,0,0,0,1600426904,1600426904),
	(59,0,0,0,0,0,0,0,1600426905,1600426905),
	(60,0,0,0,0,0,0,0,1600426905,1600426905),
	(61,0,0,0,0,0,0,0,1600426906,1600426906),
	(62,0,0,0,0,0,0,0,1600426907,1600426907),
	(63,0,0,0,0,0,0,0,1600426907,1600426907),
	(64,0,0,0,0,0,0,0,1600426908,1600426908),
	(65,0,0,0,0,0,0,0,1600426908,1600426908);

/*!40000 ALTER TABLE `video_statistic` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table videos
# ------------------------------------------------------------

CREATE TABLE `videos` (
  `video_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '视频id',
  `title` mediumtext COMMENT '视频标题',
  `describe` mediumtext COMMENT '视频描述',
  `cover` varchar(521) NOT NULL DEFAULT '' COMMENT '视频封面',
  `video_addr` varchar(521) NOT NULL DEFAULT '' COMMENT '视频地址',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `user_type` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '添加用户类型（0：管理员[sys_user]，1：用户[user]）',
  `sortorder` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '排序',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '审核状态（0：审核中，1：审核通过 2：审核不通过 3：逻辑删除）',
  `is_recommend` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '推荐（0：不推荐；1：推荐）',
  `is_top` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '置顶（0：不置顶；1：置顶；）',
  `video_duration` int(8) NOT NULL DEFAULT '0' COMMENT '视频时长（单位：秒）',
  `rec_content` mediumtext COMMENT '推荐理由',
  `top_content` mediumtext COMMENT '置顶理由',
  `video_width` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频宽',
  `video_height` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频高',
  `create_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`video_id`),
  KEY `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='视频表';

LOCK TABLES `videos` WRITE;
/*!40000 ALTER TABLE `videos` DISABLE KEYS */;

INSERT INTO `videos` (`video_id`, `title`, `describe`, `cover`, `video_addr`, `user_id`, `user_type`, `sortorder`, `status`, `is_recommend`, `is_top`, `video_duration`, `rec_content`, `top_content`, `video_width`, `video_height`, `create_at`, `update_at`)
VALUES
	(10,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426884,1600426884),
	(11,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426884,1600426884),
	(12,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426885,1600426885),
	(13,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426886,1600426886),
	(14,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426886,1600426886),
	(15,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426887,1600426887),
	(16,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426887,1600426887),
	(17,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426888,1600426888),
	(18,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426888,1600426888),
	(19,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426889,1600426889),
	(20,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426889,1600426889),
	(21,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426890,1600426890),
	(22,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426890,1600426890),
	(23,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426891,1600426891),
	(24,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426891,1600426891),
	(25,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426892,1600426892),
	(26,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426892,1600426892),
	(27,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426892,1600426892),
	(28,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426893,1600426893),
	(29,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426893,1600426893),
	(30,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426894,1600426894),
	(31,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426894,1600426894),
	(32,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426895,1600426895),
	(33,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426895,1600426895),
	(34,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426895,1600426895),
	(35,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426896,1600426896),
	(36,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426896,1600426896),
	(37,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426897,1600426897),
	(38,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426897,1600426897),
	(39,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426897,1600426897),
	(40,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426898,1600426898),
	(41,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426898,1600426898),
	(42,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426899,1600426899),
	(43,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426899,1600426899),
	(44,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426899,1600426899),
	(45,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426900,1600426900),
	(46,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426900,1600426900),
	(47,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426901,1600426901),
	(48,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426901,1600426901),
	(49,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426901,1600426901),
	(50,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426902,1600426902),
	(51,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426902,1600426902),
	(52,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426902,1600426902),
	(53,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426903,1600426903),
	(54,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426903,1600426903),
	(55,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426903,1600426903),
	(56,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426904,1600426904),
	(57,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426904,1600426904),
	(58,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426904,1600426904),
	(59,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426905,1600426905),
	(60,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426905,1600426905),
	(61,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426906,1600426906),
	(62,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426907,1600426907),
	(63,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426907,1600426907),
	(64,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426908,1600426908),
	(65,'FPV','无人机竞技','无ssss','http://111111111','202009181548217779',1,0,1,0,0,1888,'','',0,0,1600426908,1600426908);

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='视频审核表';



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
	(1,'中国','CHN',0,1,0,0,0,0),
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
	(242,'中非共和国','CAF',0,1,0,0,0,0);

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户获取游币记录表';



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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='游币打赏记录表';



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
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='游币获取任务表';




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
