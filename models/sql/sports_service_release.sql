# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 192.168.5.12 (MySQL 5.7.34)
# Database: sports_service
# Generation Time: 2021-09-02 06:05:04 +0000
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

DROP TABLE IF EXISTS `admin_user`;

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



# Dump of table app_version_control
# ------------------------------------------------------------

DROP TABLE IF EXISTS `app_version_control`;

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
  `describe` varchar(500) NOT NULL DEFAULT '' COMMENT '版本说明',
  PRIMARY KEY (`id`),
  KEY `version_code` (`version_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='app版本控制';



# Dump of table banner
# ------------------------------------------------------------

DROP TABLE IF EXISTS `banner`;

CREATE TABLE `banner` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '标题',
  `cover` varchar(512) NOT NULL DEFAULT '' COMMENT 'banner封面',
  `explain` varchar(255) NOT NULL DEFAULT '' COMMENT '说明',
  `jump_url` varchar(512) NOT NULL DEFAULT '' COMMENT '跳转地址',
  `share_url` varchar(255) NOT NULL DEFAULT '' COMMENT '分享地址',
  `type` int(1) NOT NULL DEFAULT '1' COMMENT '1 首页 2 直播页 3 官网banner',
  `start_time` int(11) NOT NULL DEFAULT '0' COMMENT '上架时间',
  `end_time` int(11) NOT NULL DEFAULT '0' COMMENT '下架时间',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0待上架 1上架 2 已过期',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `jump_type` tinyint(1) NOT NULL COMMENT '跳转类型 0 站内跳转 1 站外跳转',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='banner配置表';



# Dump of table circle_attention
# ------------------------------------------------------------

DROP TABLE IF EXISTS `circle_attention`;

CREATE TABLE `circle_attention` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(60) NOT NULL COMMENT '关注圈子的用户id',
  `circle_id` int(11) NOT NULL COMMENT '圈子id',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1表示关注 2表示取消关注',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `circle_id` (`circle_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='关注的圈子';



# Dump of table collect_record
# ------------------------------------------------------------

DROP TABLE IF EXISTS `collect_record`;

CREATE TABLE `collect_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `to_user_id` varchar(60) NOT NULL COMMENT '作品发布者用户id',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='收藏的视频';



# Dump of table comment_report
# ------------------------------------------------------------

DROP TABLE IF EXISTS `comment_report`;

CREATE TABLE `comment_report` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(60) NOT NULL DEFAULT '' COMMENT '举报人用户id',
  `comment_id` bigint(20) NOT NULL COMMENT '评论id',
  `reason` varchar(300) NOT NULL COMMENT '举报理由',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `comment_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1视频评论 2帖子评论',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `comment_id` (`comment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论举报';



# Dump of table community_section
# ------------------------------------------------------------

DROP TABLE IF EXISTS `community_section`;

CREATE TABLE `community_section` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `section_name` varchar(100) NOT NULL COMMENT '板块名称',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 未操作 1 展示  2 隐藏',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='社区版块';



# Dump of table community_topic
# ------------------------------------------------------------

DROP TABLE IF EXISTS `community_topic`;

CREATE TABLE `community_topic` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `topic_name` varchar(100) NOT NULL COMMENT '话题名称',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 未操作 1 展示  2 隐藏',
  `is_hot` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否热门话题 1 热门',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  `cover` varchar(256) NOT NULL DEFAULT '' COMMENT '话题封面',
  `describe` varchar(1000) NOT NULL DEFAULT '' COMMENT '话题描述',
  `section_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属板块id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='社区话题';



# Dump of table complaint
# ------------------------------------------------------------

DROP TABLE IF EXISTS `complaint`;

CREATE TABLE `complaint` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(60) NOT NULL COMMENT '举报人',
  `to_uid` varchar(60) NOT NULL COMMENT '被举报人',
  `reason` mediumtext COMMENT '原因',
  `compose_id` bigint(20) NOT NULL COMMENT '举报的作品id（视频/帖子id）',
  `complaint_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0 举报其他 1 举报视频 2 举报帖子',
  `cover` varchar(512) NOT NULL DEFAULT '' COMMENT '图片地址  逗号分隔',
  `is_dispose` tinyint(1) DEFAULT '1' COMMENT '是否受理 1未受理 2受理',
  `content` mediumtext COMMENT '回复内容',
  `create_at` int(11) NOT NULL,
  `update_at` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='举报信息';



# Dump of table default_avatar
# ------------------------------------------------------------

DROP TABLE IF EXISTS `default_avatar`;

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

DROP TABLE IF EXISTS `feedback`;

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

DROP TABLE IF EXISTS `hot_circle`;

CREATE TABLE `hot_circle` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `hot_circle_id` varchar(128) NOT NULL COMMENT '热门圈子id 多个用逗号分隔 例如：3,6,11,21',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='热门圈子（后台直接配置）';



# Dump of table hot_search
# ------------------------------------------------------------

DROP TABLE IF EXISTS `hot_search`;

CREATE TABLE `hot_search` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `hot_search_content` varchar(128) NOT NULL COMMENT '热门搜索内容 如：FPV、电竞',
  `status` tinyint(1) DEFAULT '0' COMMENT '0 展示 1 隐藏',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='热门搜索（后台配置）';



# Dump of table posting_apply_cream
# ------------------------------------------------------------

DROP TABLE IF EXISTS `posting_apply_cream`;

CREATE TABLE `posting_apply_cream` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(60) NOT NULL DEFAULT '' COMMENT '用户id',
  `post_id` bigint(20) NOT NULL COMMENT '帖子id',
  `reason` varchar(100) NOT NULL DEFAULT '' COMMENT '理由',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 待审核 1 通过 2 拒绝',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `post_id` (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子申精记录表';



# Dump of table posting_comment
# ------------------------------------------------------------

DROP TABLE IF EXISTS `posting_comment`;

CREATE TABLE `posting_comment` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '评论id',
  `user_id` varchar(60) NOT NULL COMMENT '评论人userId',
  `post_id` bigint(20) NOT NULL COMMENT '帖子id',
  `parent_comment_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '父评论id',
  `parent_comment_user_id` varchar(60) NOT NULL DEFAULT '' COMMENT '父评论的用户id',
  `reply_comment_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '被回复的评论id',
  `reply_comment_user_id` varchar(60) DEFAULT '' COMMENT '被回复的评论用户id',
  `comment_level` tinyint(4) NOT NULL DEFAULT '1' COMMENT '评论等级[ 1 一级评论 默认 ，2 二级评论]',
  `content` varchar(1000) NOT NULL DEFAULT '' COMMENT '评论的内容',
  `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '状态 (1 有效，0 逻辑删除)',
  `is_top` tinyint(2) NOT NULL DEFAULT '0' COMMENT '置顶状态[ 1 置顶，0 不置顶 默认 ]',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `parent_comment_id` (`parent_comment_id`),
  KEY `status` (`status`),
  KEY `create_time` (`create_at`),
  KEY `post_id` (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子评论表';



# Dump of table posting_info
# ------------------------------------------------------------

DROP TABLE IF EXISTS `posting_info`;

CREATE TABLE `posting_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '帖子id',
  `section_id` int(11) NOT NULL COMMENT '版块id',
  `title` mediumtext COMMENT '帖子标题',
  `describe` mediumtext COMMENT '帖子描述',
  `content` mediumtext COMMENT '帖子内容 图片列表/json结构 例如转发的视频 完整结构',
  `video_id` bigint(20) NOT NULL COMMENT '关联的视频id',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `posting_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '帖子类型  0 纯文本 1 图文 2 视频 + 文字 ',
  `content_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '内容类型 0 发布 1 转发视频 2 转发帖子',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 审核状态（0：审核中，1：审核通过 2：审核不通过 3：逻辑删除）',
  `is_recommend` tinyint(1) NOT NULL DEFAULT '0' COMMENT '推荐（0：不推荐；1：推荐）',
  `is_top` tinyint(1) NOT NULL DEFAULT '0' COMMENT '置顶（0：不置顶；1：置顶；）',
  `is_cream` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否精华帖（0: 不是 1: 是）',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子表';



# Dump of table posting_report
# ------------------------------------------------------------

DROP TABLE IF EXISTS `posting_report`;

CREATE TABLE `posting_report` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(60) NOT NULL DEFAULT '' COMMENT '用户id',
  `post_id` bigint(20) NOT NULL COMMENT '帖子id',
  `reason` varchar(100) NOT NULL DEFAULT '' COMMENT '举报理由',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `post_id` (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子举报';



# Dump of table posting_statistic
# ------------------------------------------------------------

DROP TABLE IF EXISTS `posting_statistic`;

CREATE TABLE `posting_statistic` (
  `posting_id` bigint(20) NOT NULL COMMENT '帖子id',
  `fabulous_num` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
  `browse_num` int(11) NOT NULL DEFAULT '0' COMMENT '浏览数',
  `share_num` int(11) NOT NULL DEFAULT '0' COMMENT '分享/转发数',
  `reward_num` int(11) NOT NULL DEFAULT '0' COMMENT '打赏的游币数',
  `comment_num` int(11) NOT NULL DEFAULT '0' COMMENT '评论数',
  `collect_num` int(11) NOT NULL DEFAULT '0' COMMENT '收藏数',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '修改时间',
  `heat_num` int(11) NOT NULL DEFAULT '0' COMMENT '热度 点赞数+ 浏览数 + 评论数',
  PRIMARY KEY (`posting_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子统计';



# Dump of table posting_topic
# ------------------------------------------------------------

DROP TABLE IF EXISTS `posting_topic`;

CREATE TABLE `posting_topic` (
  `posting_id` bigint(20) NOT NULL COMMENT '帖子id',
  `topic_id` int(11) NOT NULL COMMENT '话题id',
  `topic_name` varchar(521) DEFAULT NULL COMMENT '话题名',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '修改时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '帖子审核通过 则status为1 其他情况默认为0',
  PRIMARY KEY (`posting_id`,`topic_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子所属话题';



# Dump of table received_at
# ------------------------------------------------------------

DROP TABLE IF EXISTS `received_at`;

CREATE TABLE `received_at` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `to_user_id` varchar(60) NOT NULL COMMENT '被@的用户id',
  `user_id` varchar(60) NOT NULL COMMENT '执行@的用户id',
  `compose_id` bigint(20) NOT NULL COMMENT '视频id/帖子id/视频评论id/帖子评论id',
  `topic_type` tinyint(2) NOT NULL COMMENT '1.视频评论、回复中@ 2.帖子评论、回复中@ 3.视频评论/回复 4.帖子评论/回复 5.发布帖子时候@的用户',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `comment_level` tinyint(4) NOT NULL DEFAULT '1' COMMENT '评论等级[ 1 一级评论 默认 ，2 二级评论]',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '@状态 1 正常 0 作品待审核',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `to_user_id` (`to_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收到的@';



# Dump of table search_history
# ------------------------------------------------------------

DROP TABLE IF EXISTS `search_history`;

CREATE TABLE `search_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `search_content` varchar(128) NOT NULL COMMENT '搜索的内容',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1 正常 2 已删除',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户搜索历史表';



# Dump of table share_record
# ------------------------------------------------------------

DROP TABLE IF EXISTS `share_record`;

CREATE TABLE `share_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `compose_id` bigint(20) NOT NULL COMMENT '作品id（视频/帖子id）',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `content` text NOT NULL COMMENT '分享的整体内容（json字符串）',
  `share_type` tinyint(2) NOT NULL COMMENT '分享/转发类型 1 分享视频 2 分享帖子',
  `share_platform` tinyint(2) NOT NULL COMMENT '分享平台 1 微信 2 微博 3 qq 4 app内',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0正常 1废弃',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `compose_id` (`compose_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户分享记录';



# Dump of table social_account_login
# ------------------------------------------------------------

DROP TABLE IF EXISTS `social_account_login`;

CREATE TABLE `social_account_login` (
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `unionid` varchar(256) NOT NULL DEFAULT '' COMMENT '社交平台关联id',
  `social_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '区分社交软件 1 微信关联id 2 微博关联id 3 qq关联id',
  `status` tinyint(1) DEFAULT '0' COMMENT '0 正常 1 封禁',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`user_id`,`social_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='社交平台登陆表';



# Dump of table system_log
# ------------------------------------------------------------

DROP TABLE IF EXISTS `system_log`;

CREATE TABLE `system_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `sys_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '系统账号ID',
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

DROP TABLE IF EXISTS `system_message`;

CREATE TABLE `system_message` (
  `system_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '系统通知ID',
  `send_id` varchar(60) NOT NULL DEFAULT '' COMMENT '发送者ID（后台用户）',
  `receive_id` varchar(60) NOT NULL DEFAULT '' COMMENT '接收者id',
  `cover` varchar(256) NOT NULL DEFAULT '' COMMENT '消息封面',
  `send_default` tinyint(2) NOT NULL DEFAULT '0' COMMENT '1时发送所有用户，0时则不采用',
  `system_topic` mediumtext NOT NULL COMMENT '通知标题',
  `system_content` mediumtext NOT NULL COMMENT '通知内容',
  `send_time` int(11) NOT NULL COMMENT '发送时间',
  `expire_time` int(11) NOT NULL DEFAULT '0' COMMENT '过期时间',
  `send_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0.默认为系统通知 1 活动通知 ',
  `extra` varchar(1024) NOT NULL DEFAULT ' ' COMMENT '附件内容 例如：奖励',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 未读 1 已读  默认未读',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `android_task_id` varchar(256) NOT NULL DEFAULT '' COMMENT '友盟任务id[android端]',
  `ios_task_id` varchar(256) NOT NULL DEFAULT '' COMMENT '友盟任务id[ios端]',
  `umeng_platform` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 所有 1 android 2 ios',
  `send_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 已发送 1 未发送 2 已撤回 3 已删除',
  PRIMARY KEY (`system_id`),
  KEY `receive_id` (`receive_id`),
  KEY `send_type` (`send_type`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统消息';



# Dump of table system_notice_settings
# ------------------------------------------------------------

DROP TABLE IF EXISTS `system_notice_settings`;

CREATE TABLE `system_notice_settings` (
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `comment_push_set` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态（0：接收推送；1：拒绝推送 包含评论/回复推送）',
  `thumb_up_push_set` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态（0：接收推送；1：拒绝推送 点赞推送）',
  `attention_push_set` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态（0：接收推送；1：拒绝推送 关注推送）',
  `share_push_set` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态（0：接收推送；1：拒绝推送 包含评论/回复）',
  `slot_push_set` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态（0：接收推送；1：拒绝推送 投币推送）',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户系统推送通知设置表';



# Dump of table tencent_cloud_events
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tencent_cloud_events`;

CREATE TABLE `tencent_cloud_events` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `file_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '腾讯云文件id',
  `event` mediumtext COMMENT '事件内容（json字符串）',
  `create_at` int(11) NOT NULL,
  `event_type` tinyint(3) NOT NULL DEFAULT '0' COMMENT '0 上传事件 1 视频转码事件',
  `compose_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '作品id 视频/帖子id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='腾讯事件回调';



# Dump of table thumbs_up
# ------------------------------------------------------------

DROP TABLE IF EXISTS `thumbs_up`;

CREATE TABLE `thumbs_up` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `type_id` bigint(20) NOT NULL COMMENT '作品id （视频id/帖子id/评论id）',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `to_user_id` varchar(60) NOT NULL COMMENT '被点赞的用户id',
  `zan_type` tinyint(2) NOT NULL COMMENT '1 视频点赞 2 帖子点赞 3 视频评论点赞 4 帖子评论点赞',
  `status` tinyint(1) NOT NULL COMMENT '1赞 0未点赞',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `to_user_id` (`to_user_id`),
  KEY `type_id` (`type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='点赞表（针对帖子/视频/评论）';



# Dump of table user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `nick_name` varchar(45) NOT NULL DEFAULT '' COMMENT '昵称',
  `mobile_num` bigint(20) NOT NULL COMMENT '手机号码',
  `password` varchar(128) NOT NULL COMMENT '用户密码',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `gender` tinyint(1) NOT NULL DEFAULT '0' COMMENT '性别 0人妖 1男性 2女性',
  `born` varchar(128) NOT NULL DEFAULT '' COMMENT '出生日期',
  `age` int(3) NOT NULL DEFAULT '0' COMMENT '年龄',
  `avatar` varchar(300) NOT NULL DEFAULT '' COMMENT '头像',
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
  `reg_ip` varchar(30) DEFAULT ' ' COMMENT '注册ip',
  `device_token` varchar(100) NOT NULL DEFAULT '' COMMENT '设备token',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';



# Dump of table user_alread_read_msg
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_alread_read_msg`;

CREATE TABLE `user_alread_read_msg` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `system_id` bigint(20) NOT NULL COMMENT '系统消息ID',
  `user_id` varchar(60) NOT NULL DEFAULT '' COMMENT '用户id',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `system_id` (`system_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='记录用户已读的系统消息';



# Dump of table user_attention
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_attention`;

CREATE TABLE `user_attention` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(60) NOT NULL COMMENT '被关注的用户id',
  `attention_uid` varchar(60) NOT NULL COMMENT '关注的用户id',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1表示关注 0表示取消关注',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `attention_uid` (`attention_uid`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户关注表';



# Dump of table user_browse_record
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_browse_record`;

CREATE TABLE `user_browse_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `compose_id` bigint(20) NOT NULL COMMENT '作品id',
  `compose_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0 视频 1 帖子',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户浏览过的作品记录（包含帖子、视频）';



# Dump of table user_play_duration_record
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_play_duration_record`;

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

DROP TABLE IF EXISTS `user_ycoin`;

CREATE TABLE `user_ycoin` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `ycoin` int(11) NOT NULL COMMENT '游币数',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `ycoin` (`ycoin`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户货币表';



# Dump of table video_album
# ------------------------------------------------------------

DROP TABLE IF EXISTS `video_album`;

CREATE TABLE `video_album` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '专辑id',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `album_name` varchar(60) NOT NULL COMMENT '专辑名称',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 正常 1 废弃',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频专辑表';



# Dump of table video_album_detail
# ------------------------------------------------------------

DROP TABLE IF EXISTS `video_album_detail`;

CREATE TABLE `video_album_detail` (
  `album_id` bigint(20) NOT NULL COMMENT '专辑id',
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `album_name` varchar(60) DEFAULT NULL COMMENT '专辑名',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0为正常 1为废弃',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`album_id`,`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频专辑详情 [记录专辑下的视频]';



# Dump of table video_barrage
# ------------------------------------------------------------

DROP TABLE IF EXISTS `video_barrage`;

CREATE TABLE `video_barrage` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `video_cur_duration` int(8) NOT NULL COMMENT '视频当前时长节点（单位：秒）',
  `content` varchar(512) NOT NULL DEFAULT '' COMMENT '弹幕内容',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `color` varchar(100) NOT NULL DEFAULT '' COMMENT '弹幕字体颜色',
  `font` varchar(100) NOT NULL DEFAULT '' COMMENT '弹幕字体',
  `barrage_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '预留字段',
  `location` tinyint(2) NOT NULL DEFAULT '0' COMMENT '弹幕位置',
  `send_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '弹幕发送时间',
  PRIMARY KEY (`id`),
  KEY `video_id` (`video_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频弹幕';



# Dump of table video_comment
# ------------------------------------------------------------

DROP TABLE IF EXISTS `video_comment`;

CREATE TABLE `video_comment` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '评论id',
  `user_id` varchar(60) NOT NULL COMMENT '评论人userId',
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `parent_comment_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '父评论id',
  `parent_comment_user_id` varchar(60) NOT NULL DEFAULT '' COMMENT '父评论的用户id',
  `reply_comment_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '被回复的评论id',
  `reply_comment_user_id` varchar(60) DEFAULT '' COMMENT '被回复的评论用户id',
  `comment_level` tinyint(4) NOT NULL DEFAULT '1' COMMENT '评论等级[ 1 一级评论 默认 ，2 二级评论]',
  `content` varchar(1000) NOT NULL DEFAULT '' COMMENT '评论的内容',
  `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '状态 (1 有效，0 逻辑删除)',
  `is_top` tinyint(2) NOT NULL DEFAULT '0' COMMENT '置顶状态[ 1 置顶，0 不置顶 默认 ]',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `parent_comment_id` (`parent_comment_id`),
  KEY `status` (`status`),
  KEY `create_time` (`create_at`),
  KEY `comment_index` (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户评论表';



# Dump of table video_label_config
# ------------------------------------------------------------

DROP TABLE IF EXISTS `video_label_config`;

CREATE TABLE `video_label_config` (
  `label_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '标签id',
  `pid` int(11) NOT NULL COMMENT '父类id 0为1级分类',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` int(1) NOT NULL DEFAULT '1' COMMENT '类别状态1-正常,2-已废弃',
  `label_name` varchar(64) NOT NULL COMMENT '标签名称',
  `icon` varchar(256) NOT NULL DEFAULT '' COMMENT '标签icon',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`label_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COMMENT='视频标签配置';



# Dump of table video_labels
# ------------------------------------------------------------

DROP TABLE IF EXISTS `video_labels`;

CREATE TABLE `video_labels` (
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `label_id` varchar(60) NOT NULL DEFAULT '' COMMENT '标签id',
  `label_name` varchar(521) DEFAULT NULL COMMENT '标签名',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '视频审核通过 则status为1 其他情况默认为0',
  PRIMARY KEY (`video_id`,`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频拥有的标签表';



# Dump of table video_live
# ------------------------------------------------------------

DROP TABLE IF EXISTS `video_live`;

CREATE TABLE `video_live` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `anchor_id` varchar(60) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主播id',
  `room_id` varchar(60) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '房间id',
  `cover` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '直播封面',
  `rtmp_addr` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'rtmp地址',
  `flv_addr` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'flv地址',
  `hls_addr` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'hls地址',
  `stream_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '推流url',
  `stream_key` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '推流密钥',
  `play_time` int(11) NOT NULL COMMENT '开播时间',
  `end_time` int(11) NOT NULL COMMENT '结束时间',
  `income_ycoin` int(11) DEFAULT NULL COMMENT '本次直播收益（游币）',
  `status` tinyint(1) DEFAULT '0' COMMENT '状态 0未直播 1直播中 2异常',
  `describe` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '描述',
  `tags` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '直播标签',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '记录创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '记录更新时间',
  `duration` bigint(20) DEFAULT '0' COMMENT '时长（秒）',
  `live_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '直播类型（0：管理员[sys_user]，1：用户[user]）',
  `manager` int(11) NOT NULL DEFAULT '0' COMMENT '后台操作用户',
  `sequence` varchar(50) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '推流唯一标识',
  PRIMARY KEY (`id`),
  KEY `anchor_id` (`anchor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频直播表';



# Dump of table video_report
# ------------------------------------------------------------

DROP TABLE IF EXISTS `video_report`;

CREATE TABLE `video_report` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(60) NOT NULL DEFAULT '' COMMENT '用户id',
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `reason` varchar(100) NOT NULL DEFAULT '' COMMENT '举报理由',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `video_id` (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频举报';



# Dump of table video_statistic
# ------------------------------------------------------------

DROP TABLE IF EXISTS `video_statistic`;

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



# Dump of table video_subarea
# ------------------------------------------------------------

DROP TABLE IF EXISTS `video_subarea`;

CREATE TABLE `video_subarea` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '分区id',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `subarea_name` varchar(60) DEFAULT NULL COMMENT '分区名',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0正常 1废弃',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频分区配置表';



# Dump of table videos
# ------------------------------------------------------------

DROP TABLE IF EXISTS `videos`;

CREATE TABLE `videos` (
  `video_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '视频id',
  `title` mediumtext COMMENT '视频标题',
  `describe` mediumtext COMMENT '视频描述',
  `cover` varchar(521) NOT NULL DEFAULT '' COMMENT '视频封面',
  `video_addr` varchar(521) NOT NULL DEFAULT '' COMMENT '视频地址',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `user_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '添加用户类型（0：管理员[sys_user]，1：用户[user]）',
  `sortorder` int(11) NOT NULL DEFAULT '1' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '审核状态（0：审核中，1：审核通过 2：审核不通过 3：逻辑删除）',
  `is_recommend` tinyint(1) NOT NULL DEFAULT '0' COMMENT '推荐（0：不推荐；1：推荐）',
  `is_top` tinyint(1) NOT NULL DEFAULT '0' COMMENT '置顶（0：不置顶；1：置顶；）',
  `video_duration` int(8) NOT NULL DEFAULT '0' COMMENT '视频时长（单位：毫秒）',
  `rec_content` mediumtext COMMENT '推荐理由',
  `top_content` mediumtext COMMENT '置顶理由',
  `video_width` bigint(20) NOT NULL DEFAULT '0' COMMENT '视频宽',
  `video_height` bigint(20) NOT NULL DEFAULT '0' COMMENT '视频高',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '修改时间',
  `file_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '腾讯云文件id',
  `size` bigint(20) NOT NULL DEFAULT '0' COMMENT '视频大小（字节数）',
  `play_info` mediumtext NOT NULL COMMENT '视频转码数据',
  `pub_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1 首页发布视频 2 社区发布视频',
  `subarea` int(11) NOT NULL DEFAULT '0' COMMENT '视频所属分区',
  `album` bigint(20) NOT NULL DEFAULT '0' COMMENT '视频所属专辑',
  `ai_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'AI审核状态（0：未使用AI审核，1：AI审核通过 2：AI审核不通过 3：AI建议复审',
  PRIMARY KEY (`video_id`),
  KEY `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='视频表';



# Dump of table videos_examine
# ------------------------------------------------------------

DROP TABLE IF EXISTS `videos_examine`;

CREATE TABLE `videos_examine` (
  `video_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '视频id',
  `title` mediumtext COMMENT '视频标题',
  `describe` mediumtext COMMENT '视频描述',
  `cover` varchar(521) NOT NULL DEFAULT '' COMMENT '视频封面',
  `video_addr` varchar(521) NOT NULL DEFAULT '' COMMENT '视频地址',
  `label_id` varchar(521) DEFAULT '' COMMENT '标签id',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `user_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '添加用户类型（0：管理员[sys_user]，1：用户[user]）',
  `sortorder` int(11) NOT NULL DEFAULT '1' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '审核状态（0：无操作，1：审核通过 2：审核不通过 3：删除）',
  `is_recommend` tinyint(1) NOT NULL DEFAULT '0' COMMENT '推荐（0：不推荐；1：推荐）',
  `is_top` tinyint(1) NOT NULL DEFAULT '0' COMMENT '置顶（0：不置顶；1：置顶；）',
  `video_duration` int(8) NOT NULL DEFAULT '0' COMMENT '视频时长（单位：秒）',
  `label_name` mediumtext COMMENT '标签',
  `rec_content` mediumtext COMMENT '推荐理由',
  `top_content` mediumtext COMMENT '置顶理由',
  `video_width` bigint(20) NOT NULL DEFAULT '0' COMMENT '视频宽',
  `video_height` bigint(20) NOT NULL DEFAULT '0' COMMENT '视频高',
  `manager` int(11) NOT NULL DEFAULT '0' COMMENT '后台操作用户',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`video_id`),
  KEY `user_id` (`user_id`) USING BTREE,
  FULLTEXT KEY `label_name` (`label_name`),
  FULLTEXT KEY `label_id` (`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频审核表';



# Dump of table world_map
# ------------------------------------------------------------

DROP TABLE IF EXISTS `world_map`;

CREATE TABLE `world_map` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `name` varchar(50) NOT NULL COMMENT '国家(省份/城市)名称',
  `code` char(4) NOT NULL COMMENT '国家(省份/城市)编码',
  `pid` int(11) NOT NULL DEFAULT '0' COMMENT '父级id',
  `layer` tinyint(2) NOT NULL DEFAULT '0' COMMENT '层级',
  `sortorder` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 展示 1 隐藏',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `pid` (`pid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='世界国家/省份/城市 信息';



# Dump of table ycoin_receive_record
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ycoin_receive_record`;

CREATE TABLE `ycoin_receive_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
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

DROP TABLE IF EXISTS `ycoin_reward_record`;

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

DROP TABLE IF EXISTS `ycoin_task`;

CREATE TABLE `ycoin_task` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
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
