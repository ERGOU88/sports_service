/*
 Navicat Premium Data Transfer

 Source Server         : dasheng
 Source Server Type    : MySQL
 Source Server Version : 80012
 Source Host           : localhost:3306
 Source Schema         : fpv2

 Target Server Type    : MySQL
 Target Server Version : 80012
 File Encoding         : 65001

 Date: 14/09/2020 12:24:26
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for banner
-- ----------------------------
DROP TABLE IF EXISTS `banner`;
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

-- ----------------------------
-- Table structure for circle_attention
-- ----------------------------
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='关注的圈子';

-- ----------------------------
-- Table structure for collect_post_record
-- ----------------------------
DROP TABLE IF EXISTS `collect_post_record`;
CREATE TABLE `collect_post_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `post_id` bigint(20) NOT NULL COMMENT '帖子id',
  `status` tinyint(1) NOT NULL COMMENT '1 收藏 2 取消收藏',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `post_id` (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='收藏的帖子';

-- ----------------------------
-- Table structure for collect_video_record
-- ----------------------------
DROP TABLE IF EXISTS `collect_video_record`;
CREATE TABLE `collect_video_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `status` tinyint(1) NOT NULL COMMENT '1 收藏 2 取消收藏',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `video_id` (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='收藏的视频';

-- ----------------------------
-- Table structure for complaint
-- ----------------------------
DROP TABLE IF EXISTS `complaint`;
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
) ENGINE=InnoDB AUTO_INCREMENT=116 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='举报信息';

-- ----------------------------
-- Table structure for default_avatar
-- ----------------------------
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

-- ----------------------------
-- Table structure for feedback
-- ----------------------------
DROP TABLE IF EXISTS `feedback`;
CREATE TABLE `feedback` (
  `id` bigint(20) NOT NULL COMMENT '自增主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `phone` varchar(200) DEFAULT NULL COMMENT '手机号码',
  `text` mediumtext COMMENT '反馈的内容',
  `contact` varchar(200) DEFAULT NULL COMMENT '联系方式',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 0未回复 1已回复',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `content` mediumtext COMMENT '回复内容',
  `pics` varchar(512) DEFAULT NULL COMMENT '上传的图片，多张逗号分隔',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='反馈表';

-- ----------------------------
-- Table structure for hot_circle
-- ----------------------------
DROP TABLE IF EXISTS `hot_circle`;
CREATE TABLE `hot_circle` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `hot_circle_id` varchar(128) NOT NULL COMMENT '热门圈子id 多个用逗号分隔 例如：3,6,11,21',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='热门圈子（后台直接配置）';

-- ----------------------------
-- Table structure for hot_search
-- ----------------------------
DROP TABLE IF EXISTS `hot_search`;
CREATE TABLE `hot_search` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `hot_search_content` varchar(128) NOT NULL COMMENT '热门搜索内容 多个用逗号分隔 例如：FPV,电竞,无人机',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='热门搜索（后台配置）';

-- ----------------------------
-- Table structure for post_comment
-- ----------------------------
DROP TABLE IF EXISTS `post_comment`;
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

-- ----------------------------
-- Table structure for post_comment_parent_child
-- ----------------------------
DROP TABLE IF EXISTS `post_comment_parent_child`;
CREATE TABLE `post_comment_parent_child` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `parent_id` bigint(20) NOT NULL COMMENT '评论父id',
  `child_id` bigint(20) NOT NULL COMMENT '评论子id',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `post_id` (`parent_id`),
  KEY `from_uid` (`child_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='帖子评论表中各个评论的父子关系';

-- ----------------------------
-- Table structure for post_info
-- ----------------------------
DROP TABLE IF EXISTS `post_info`;
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

-- ----------------------------
-- Table structure for post_info_examine
-- ----------------------------
DROP TABLE IF EXISTS `post_info_examine`;
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

-- ----------------------------
-- Table structure for post_label_config
-- ----------------------------
DROP TABLE IF EXISTS `post_label_config`;
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

-- ----------------------------
-- Table structure for post_labels
-- ----------------------------
DROP TABLE IF EXISTS `post_labels`;
CREATE TABLE `post_labels` (
  `post_id` bigint(20) unsigned NOT NULL COMMENT '帖子id',
  `label_id` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标签id',
  `label_name` varchar(521) DEFAULT NULL COMMENT '标签名',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`post_id`,`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='帖子拥有的标签表';

-- ----------------------------
-- Table structure for post_statistic
-- ----------------------------
DROP TABLE IF EXISTS `post_statistic`;
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

-- ----------------------------
-- Table structure for received_at
-- ----------------------------
DROP TABLE IF EXISTS `received_at`;
CREATE TABLE `received_at` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(60) NOT NULL COMMENT '被@的用户id',
  `comment_id` bigint(20) NOT NULL COMMENT '评论id',
  `topic_type` tinyint(2) NOT NULL COMMENT '1.视频 2.帖子 3.评论',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='收到的@';

-- ----------------------------
-- Table structure for search_history
-- ----------------------------
DROP TABLE IF EXISTS `search_history`;
CREATE TABLE `search_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `search_content` varchar(128) NOT NULL COMMENT '搜索的内容',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1 正常 2 已删除',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户搜索历史表';

-- ----------------------------
-- Table structure for section_config
-- ----------------------------
DROP TABLE IF EXISTS `section_config`;
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

-- ----------------------------
-- Table structure for share_record
-- ----------------------------
DROP TABLE IF EXISTS `share_record`;
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

-- ----------------------------
-- Table structure for social_account_login
-- ----------------------------
DROP TABLE IF EXISTS `social_account_login`;
CREATE TABLE `social_account_login` (
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `unionid` varchar(256) NOT NULL DEFAULT '' COMMENT '社交平台关联id',
  `social_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '区分社交软件 1 微信关联id 2 微博关联id 3 qq关联id',
  `status` tinyint(1) DEFAULT '0' COMMENT '0 正常 1 封禁',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`user_id`,`social_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='社交平台登陆表';

-- ----------------------------
-- Table structure for social_circle
-- ----------------------------
DROP TABLE IF EXISTS `social_circle`;
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

-- ----------------------------
-- Table structure for system_log
-- ----------------------------
DROP TABLE IF EXISTS `system_log`;
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

-- ----------------------------
-- Table structure for system_message
-- ----------------------------
DROP TABLE IF EXISTS `system_message`;
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

-- ----------------------------
-- Table structure for system_notice_settings
-- ----------------------------
DROP TABLE IF EXISTS `system_notice_settings`;
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

-- ----------------------------
-- Table structure for system_push
-- ----------------------------
DROP TABLE IF EXISTS `system_push`;
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

-- ----------------------------
-- Table structure for thumbs_up
-- ----------------------------
DROP TABLE IF EXISTS `thumbs_up`;
CREATE TABLE `thumbs_up` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `type_id` bigint(20) NOT NULL COMMENT '作品id （视频id/帖子id/评论id）',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `zan_type` tinyint(2) NOT NULL COMMENT '1 视频点赞 2 帖子点赞 3 评论点赞',
  `status` tinyint(1) NOT NULL COMMENT '1赞 2取消点赞',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `type_id` (`type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='点赞表（针对帖子/视频/评论）';

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
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
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户表';

-- ----------------------------
-- Table structure for user_attention
-- ----------------------------
DROP TABLE IF EXISTS `user_attention`;
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

-- ----------------------------
-- Table structure for user_browse_post
-- ----------------------------
DROP TABLE IF EXISTS `user_browse_post`;
CREATE TABLE `user_browse_post` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `post_id` bigint(20) NOT NULL COMMENT '帖子id',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户浏览过的帖子记录';

-- ----------------------------
-- Table structure for user_browse_video
-- ----------------------------
DROP TABLE IF EXISTS `user_browse_video`;
CREATE TABLE `user_browse_video` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `video_id` bigint(20) NOT NULL COMMENT '视频id',
  `create_at` int(11) NOT NULL COMMENT '创建时间',
  `update_at` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户浏览过的视频记录';

-- ----------------------------
-- Table structure for user_ycoin
-- ----------------------------
DROP TABLE IF EXISTS `user_ycoin`;
CREATE TABLE `user_ycoin` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `ycoin` int(11) NOT NULL COMMENT '游币数',
  `update_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `ycoin` (`ycoin`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户货币表';

-- ----------------------------
-- Table structure for video_barrage
-- ----------------------------
DROP TABLE IF EXISTS `video_barrage`;
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

-- ----------------------------
-- Table structure for video_comment
-- ----------------------------
DROP TABLE IF EXISTS `video_comment`;
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

-- ----------------------------
-- Table structure for video_comment_parent_child
-- ----------------------------
DROP TABLE IF EXISTS `video_comment_parent_child`;
CREATE TABLE `video_comment_parent_child` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `parent_id` bigint(20) NOT NULL COMMENT '评论父id',
  `child_id` bigint(20) NOT NULL COMMENT '评论子id',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `post_id` (`parent_id`),
  KEY `from_uid` (`child_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频评论表中各个评论的父子关系';

-- ----------------------------
-- Table structure for video_label_config
-- ----------------------------
DROP TABLE IF EXISTS `video_label_config`;
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

-- ----------------------------
-- Table structure for video_labels
-- ----------------------------
DROP TABLE IF EXISTS `video_labels`;
CREATE TABLE `video_labels` (
  `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
  `label_id` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标签id',
  `label_name` varchar(521) DEFAULT NULL COMMENT '标签名',
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`video_id`,`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='视频拥有的标签表';

-- ----------------------------
-- Table structure for video_live
-- ----------------------------
DROP TABLE IF EXISTS `video_live`;
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

-- ----------------------------
-- Table structure for video_statistic
-- ----------------------------
DROP TABLE IF EXISTS `video_statistic`;
CREATE TABLE `video_statistic` (
  `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
  `fabulous` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '点赞数',
  `browse` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '浏览数',
  `share` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '分享数',
  `reward` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '打赏的游币数',
  `create_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='视频相关参数统计';

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos` (
  `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
  `title` mediumtext COMMENT '视频标题',
  `describe` mediumtext COMMENT '视频描述',
  `cover` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '视频封面',
  `video_addr` varchar(521) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '视频地址',
  `user_id` varchar(60) NOT NULL COMMENT '用户id',
  `sortorder` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '状态（0：展示，1：隐藏 ）',
  `is_recommend` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '推荐（0：不推荐；1：推荐）',
  `is_top` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '置顶（0：不置顶；1：置顶）',
  `video_duration` int(8) NOT NULL DEFAULT '0' COMMENT '视频时长（单位：秒）',
  `rec_content` mediumtext COMMENT '推荐理由',
  `top_content` mediumtext COMMENT '置顶理由',
  `video_width` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频宽',
  `video_height` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '视频高',
  `create_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`video_id`),
  KEY `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='视频表';

-- ----------------------------
-- Table structure for videos_examine
-- ----------------------------
DROP TABLE IF EXISTS `videos_examine`;
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

-- ----------------------------
-- Table structure for ycoin_receive_record
-- ----------------------------
DROP TABLE IF EXISTS `ycoin_receive_record`;
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

-- ----------------------------
-- Table structure for ycoin_reward_record
-- ----------------------------
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='游币打赏记录表';

-- ----------------------------
-- Table structure for ycoin_task
-- ----------------------------
DROP TABLE IF EXISTS `ycoin_task`;
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

SET FOREIGN_KEY_CHECKS = 1;
