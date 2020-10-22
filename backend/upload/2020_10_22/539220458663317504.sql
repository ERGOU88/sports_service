/*
 Navicat Premium Data Transfer

 Source Server         : 223.203.221.49
 Source Server Type    : MySQL
 Source Server Version : 50634
 Source Host           : 223.203.221.49:3306
 Source Schema         : guiquan

 Target Server Type    : MySQL
 Target Server Version : 50634
 File Encoding         : 65001

 Date: 04/08/2018 11:22:34
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for _mysql_session_store
-- ----------------------------
DROP TABLE IF EXISTS `_mysql_session_store`;
CREATE TABLE `_mysql_session_store` (
  `id` varchar(255) NOT NULL,
  `expires` bigint(20) DEFAULT NULL,
  `data` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for active
-- ----------------------------
DROP TABLE IF EXISTS `active`;
CREATE TABLE `active` (
  `id` bigint(20) NOT NULL,
  `img` varchar(200) DEFAULT NULL,
  `name` varchar(256) DEFAULT NULL,
  `title` varchar(256) DEFAULT NULL,
  `url` varchar(200) DEFAULT NULL,
  `begin_time` bigint(20) DEFAULT NULL,
  `end_time` bigint(20) DEFAULT NULL,
  `limit_city` tinyint(2) DEFAULT NULL,
  `city` varchar(30) DEFAULT NULL,
  `sort_num` int(11) DEFAULT NULL,
  `active_type` tinyint(2) DEFAULT NULL,
  `status` tinyint(2) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `del` tinyint(2) DEFAULT '0',
  `children` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for admin_group
-- ----------------------------
DROP TABLE IF EXISTS `admin_group`;
CREATE TABLE `admin_group` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL COMMENT '分组名称',
  `status` int(1) DEFAULT '1' COMMENT '状态 0 不可用 1可用',
  `rules` text COMMENT '对应规则id 多个逗号分割',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COMMENT='管理员分组';

-- ----------------------------
-- Table structure for admin_group_access
-- ----------------------------
DROP TABLE IF EXISTS `admin_group_access`;
CREATE TABLE `admin_group_access` (
  `uid` bigint(20) DEFAULT NULL,
  `group_id` int(11) unsigned NOT NULL,
  UNIQUE KEY `uid_group_id` (`uid`,`group_id`),
  KEY `uid` (`uid`),
  KEY `group_id` (`group_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='用户所属分组表';

-- ----------------------------
-- Table structure for admin_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_menu`;
CREATE TABLE `admin_menu` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT '' COMMENT '菜单名称',
  `url` varchar(255) DEFAULT NULL COMMENT '请求的URL',
  `p_id` int(11) DEFAULT '0' COMMENT '上一级的ID',
  `sort` int(3) DEFAULT '1' COMMENT '排序',
  `lv` int(11) DEFAULT NULL COMMENT '等级ID',
  `path` varchar(100) DEFAULT '0' COMMENT '级别路径',
  `icon` varchar(100) DEFAULT NULL COMMENT '菜单图标',
  `is_show` tinyint(3) NOT NULL DEFAULT '1' COMMENT '是否显示在左边菜单，1--显示',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=109 DEFAULT CHARSET=utf8 COMMENT='后台菜单';

-- ----------------------------
-- Table structure for admin_rule
-- ----------------------------
DROP TABLE IF EXISTS `admin_rule`;
CREATE TABLE `admin_rule` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` char(80) DEFAULT '' COMMENT '权限唯一标识',
  `title` varchar(50) DEFAULT '' COMMENT '权限显示名称（标题）相当于权限的描述',
  `type` tinyint(1) DEFAULT '1' COMMENT '权限类型',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '权限状态 0失效 1有效',
  `condition` char(100) DEFAULT '' COMMENT '权限条件（为空存在就验证，不为空表示按条件验证）',
  `p_id` int(11) unsigned DEFAULT '0' COMMENT '父级ID',
  `lv` int(11) DEFAULT '1' COMMENT '级别',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=MyISAM AUTO_INCREMENT=310 DEFAULT CHARSET=utf8 COMMENT='权限表';

-- ----------------------------
-- Table structure for anchor
-- ----------------------------
DROP TABLE IF EXISTS `anchor`;
CREATE TABLE `anchor` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) DEFAULT NULL,
  `room_id` int(20) DEFAULT NULL COMMENT '房间号 xingxingid',
  `channel_resp` varchar(2000) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '网易拉留信息',
  `create_chat_room_resp` varchar(2000) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '聊天室相关信息',
  `operator_id` bigint(20) DEFAULT NULL,
  `recommend` tinyint(1) DEFAULT '1' COMMENT '1 不推荐 2推荐',
  `status` tinyint(1) DEFAULT '1' COMMENT '直播状态，1未直播 2直播中 ',
  `seal` int(20) DEFAULT '1' COMMENT '封房间 解封的时间戳',
  `is_seal` tinyint(1) DEFAULT '0' COMMENT '封禁状态  0未封  1已封',
  `index` bigint(20) DEFAULT NULL,
  `cover` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '房间封面',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=120 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for anchor_detail
-- ----------------------------
DROP TABLE IF EXISTS `anchor_detail`;
CREATE TABLE `anchor_detail` (
  `id` bigint(20) NOT NULL,
  `sort` int(10) DEFAULT '1',
  `anchor_id` int(11) DEFAULT NULL COMMENT '主播列表ID',
  `cover_path` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '房间封面',
  `desc` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '直播间描述',
  `city` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '城市',
  `number` int(11) DEFAULT '0' COMMENT '直播间人数',
  `private` tinyint(1) DEFAULT '0' COMMENT '是否私有 ',
  `theme_id` bigint(20) DEFAULT NULL COMMENT '主题id',
  `user_id` bigint(20) DEFAULT NULL,
  `start_time` bigint(20) DEFAULT NULL,
  `end_time` bigint(20) DEFAULT NULL,
  `live_duration` bigint(20) DEFAULT '0' COMMENT '直播时长',
  `income_gold` bigint(20) DEFAULT NULL COMMENT '直播收益木头',
  `channel_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '绑定直播回放的id',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for anchor_detail_cover
-- ----------------------------
DROP TABLE IF EXISTS `anchor_detail_cover`;
CREATE TABLE `anchor_detail_cover` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `room_id` int(11) DEFAULT NULL COMMENT '房间详情ID',
  `cover_path` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `is_check` tinyint(1) DEFAULT NULL COMMENT '是否审核 1未审核 2审核通过 3审核不通过',
  `operator_id` bigint(20) DEFAULT NULL,
  `remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '备注',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for apple_in_app_pay_record
-- ----------------------------
DROP TABLE IF EXISTS `apple_in_app_pay_record`;
CREATE TABLE `apple_in_app_pay_record` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `md5` text,
  `transaction_id` varchar(200) DEFAULT NULL,
  `goods_name` varchar(4096) DEFAULT NULL,
  `begin_time` bigint(20) DEFAULT NULL,
  `end_time` bigint(20) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `fee` int(11) unsigned DEFAULT '0',
  `diamond` int(11) unsigned DEFAULT '0',
  `status` int(11) DEFAULT NULL,
  `err_msg` varchar(200) DEFAULT NULL,
  `0` int(11) DEFAULT NULL,
  `type` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UQE_apple_in_app_pay_record_transaction_id` (`transaction_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for call_chat
-- ----------------------------
DROP TABLE IF EXISTS `call_chat`;
CREATE TABLE `call_chat` (
  `id` bigint(20) NOT NULL,
  `call_id` varchar(256) DEFAULT NULL,
  `a_user` varchar(256) DEFAULT NULL,
  `b_user` varchar(256) DEFAULT NULL,
  `a_user_star` int(4) DEFAULT '-1',
  `b_user_star` int(4) DEFAULT '-1',
  `pay_success` int(4) DEFAULT '0',
  `end_user` varchar(256) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `voice_tip` int(11) DEFAULT NULL,
  `caller_user` varchar(256) DEFAULT NULL,
  `start_unix` int(11) DEFAULT NULL,
  `end_unix` int(11) DEFAULT NULL,
  `cause` varchar(512) DEFAULT NULL,
  `call_type` tinyint(1) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `tip` int(11) unsigned DEFAULT '0',
  `wy_cause` varchar(512) DEFAULT NULL,
  `play_type` tinyint(1) DEFAULT NULL,
  `start_sys_unix` int(11) DEFAULT NULL,
  `gold` text,
  `children` text,
  `a_name` varchar(255) DEFAULT NULL,
  `b_name` varchar(255) DEFAULT NULL,
  `caller_name` varchar(255) DEFAULT NULL,
  `task_id` varchar(512) DEFAULT NULL,
  `end_sys_unix` int(11) DEFAULT NULL,
  `user_label_id` bigint(20) NOT NULL,
  `one_word` varchar(48) DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for can_withdraw
-- ----------------------------
DROP TABLE IF EXISTS `can_withdraw`;
CREATE TABLE `can_withdraw` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `channel_id` varchar(20) NOT NULL,
  `code` varchar(255) NOT NULL,
  `amount` varchar(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) DEFAULT NULL,
  `key` varchar(64) DEFAULT NULL,
  `del` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for category_child
-- ----------------------------
DROP TABLE IF EXISTS `category_child`;
CREATE TABLE `category_child` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) DEFAULT NULL,
  `parent_id` bigint(20) DEFAULT NULL,
  `flag` tinyint(1) DEFAULT NULL,
  `del` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for channel
-- ----------------------------
DROP TABLE IF EXISTS `channel`;
CREATE TABLE `channel` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `code` varchar(256) DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `name` varchar(256) DEFAULT NULL,
  `desc` varchar(256) DEFAULT NULL,
  `act_num` int(11) DEFAULT NULL,
  `full_name` varchar(256) DEFAULT NULL,
  `charge_name` varchar(256) DEFAULT NULL,
  `join_time` bigint(20) DEFAULT NULL,
  `city` varchar(256) DEFAULT NULL,
  `coo_mode` int(11) DEFAULT NULL,
  `contract_code` varchar(256) DEFAULT NULL,
  `charge_mode` int(11) DEFAULT NULL,
  `charge_cycle` int(11) DEFAULT NULL,
  `charge_manager` varchar(256) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `del` tinyint(2) DEFAULT '0',
  `children` text,
  `freeze` tinyint(2) DEFAULT '2',
  `pay_total` bigint(20) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=92695733584531457 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for channel_active
-- ----------------------------
DROP TABLE IF EXISTS `channel_active`;
CREATE TABLE `channel_active` (
  `id` bigint(20) NOT NULL,
  `active_id` bigint(20) DEFAULT NULL,
  `channel_id` bigint(20) DEFAULT NULL,
  `reward_code` varchar(256) DEFAULT NULL,
  `enter_total` bigint(20) DEFAULT NULL,
  `download_total` bigint(20) DEFAULT NULL,
  `register_total` bigint(20) DEFAULT NULL,
  `pay_total` bigint(20) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `code` varchar(256) DEFAULT NULL,
  `qr_url` varchar(256) DEFAULT NULL,
  `exchange_total` bigint(20) DEFAULT '0',
  `del` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for channel_record
-- ----------------------------
DROP TABLE IF EXISTS `channel_record`;
CREATE TABLE `channel_record` (
  `id` bigint(20) NOT NULL,
  `code` varchar(256) DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `ip` varchar(256) DEFAULT NULL,
  `user_agent` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for channel_withdraw
-- ----------------------------
DROP TABLE IF EXISTS `channel_withdraw`;
CREATE TABLE `channel_withdraw` (
  `i_d` bigint(20) NOT NULL,
  `channel_i_d` bigint(20) NOT NULL,
  `channel_code` varchar(255) DEFAULT NULL,
  `amount` varchar(20) DEFAULT NULL,
  `start_time` bigint(20) DEFAULT NULL,
  `end_time` bigint(20) DEFAULT NULL,
  `created` bigint(20) DEFAULT NULL,
  `updated` bigint(20) DEFAULT NULL,
  `out_biz_no` varchar(255) DEFAULT NULL,
  `order_i_d` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`i_d`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for chat_red_package
-- ----------------------------
DROP TABLE IF EXISTS `chat_red_package`;
CREATE TABLE `chat_red_package` (
  `id` bigint(20) NOT NULL,
  `call_id` varchar(256) DEFAULT NULL,
  `from_user` varchar(256) DEFAULT NULL,
  `to_user` varchar(256) DEFAULT NULL,
  `gold` int(11) unsigned DEFAULT '0',
  `red_package_unix` int(11) DEFAULT NULL,
  `chat_type` int(11) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for collection_invitation
-- ----------------------------
DROP TABLE IF EXISTS `collection_invitation`;
CREATE TABLE `collection_invitation` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `topic_id` bigint(20) NOT NULL,
  `invitation_id` bigint(20) DEFAULT NULL,
  `sequence` int(11) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for collection_invitation_v2
-- ----------------------------
DROP TABLE IF EXISTS `collection_invitation_v2`;
CREATE TABLE `collection_invitation_v2` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `topic_id` bigint(20) NOT NULL,
  `invitation_id` bigint(20) DEFAULT NULL,
  `sequence` int(11) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for college
-- ----------------------------
DROP TABLE IF EXISTS `college`;
CREATE TABLE `college` (
  `coid` smallint(5) unsigned NOT NULL,
  `name` char(100) NOT NULL,
  `provinceID` mediumint(8) unsigned NOT NULL,
  `badge` varchar(255) DEFAULT NULL COMMENT '校徽',
  PRIMARY KEY (`coid`),
  KEY `provinceID` (`provinceID`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for college_school
-- ----------------------------
DROP TABLE IF EXISTS `college_school`;
CREATE TABLE `college_school` (
  `scid` int(8) unsigned NOT NULL,
  `name` char(100) NOT NULL,
  `collegeID` smallint(5) unsigned NOT NULL,
  PRIMARY KEY (`scid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `task_id` bigint(20) DEFAULT NULL,
  `content` varchar(2056) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `is_check` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for comment_invitation
-- ----------------------------
DROP TABLE IF EXISTS `comment_invitation`;
CREATE TABLE `comment_invitation` (
  `id` bigint(20) NOT NULL,
  `content` varchar(2056) DEFAULT NULL,
  `user_id` bigint(20) NOT NULL,
  `is_check` tinyint(1) DEFAULT '1',
  `is_del` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `star_times` int(11) DEFAULT '0',
  `longitude` varchar(20) DEFAULT '',
  `latitude` varchar(20) DEFAULT '',
  `to_user` bigint(20) DEFAULT NULL,
  `to_user_nick_name` varchar(45) DEFAULT NULL,
  `type` tinyint(1) DEFAULT NULL,
  `voice_time` int(11) DEFAULT '0',
  `invitation_id` bigint(20) NOT NULL DEFAULT '0',
  `user_label_id` bigint(20) NOT NULL DEFAULT '0',
  `topic_id` bigint(20) NOT NULL DEFAULT '0',
  `check_comment` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for comment_star
-- ----------------------------
DROP TABLE IF EXISTS `comment_star`;
CREATE TABLE `comment_star` (
  `id` bigint(20) NOT NULL,
  `comment_id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for complains
-- ----------------------------
DROP TABLE IF EXISTS `complains`;
CREATE TABLE `complains` (
  `id` bigint(20) NOT NULL,
  `informant_id` bigint(20) DEFAULT NULL,
  `complains_id` bigint(20) DEFAULT NULL,
  `info` varchar(255) DEFAULT NULL,
  `status` int(4) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `report_channel` int(11) DEFAULT NULL COMMENT '举报渠道',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for cv_apply
-- ----------------------------
DROP TABLE IF EXISTS `cv_apply`;
CREATE TABLE `cv_apply` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `avator` varchar(200) DEFAULT NULL,
  `voices` varchar(2048) DEFAULT NULL,
  `chan_num` int(11) DEFAULT NULL,
  `checked` tinyint(2) DEFAULT NULL,
  `code` varchar(200) DEFAULT NULL,
  `pretty_photo` varchar(200) DEFAULT NULL,
  `channel_name` varchar(255) DEFAULT NULL,
  `channel_code` varchar(255) DEFAULT NULL,
  `created` bigint(20) DEFAULT NULL,
  `disapproved_reason` varchar(255) DEFAULT NULL,
  `disapprove_reason` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for device
-- ----------------------------
DROP TABLE IF EXISTS `device`;
CREATE TABLE `device` (
  `device` varchar(60) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `device_token` varchar(200) DEFAULT NULL,
  `system` varchar(45) DEFAULT NULL,
  `model` varchar(45) DEFAULT NULL,
  `version` varchar(45) DEFAULT NULL,
  `carrier` varchar(45) DEFAULT NULL,
  `platform` varchar(45) DEFAULT NULL,
  `platform_id` varchar(45) DEFAULT NULL,
  `bundle_id` varchar(45) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`device`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for draft
-- ----------------------------
DROP TABLE IF EXISTS `draft`;
CREATE TABLE `draft` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `del` tinyint(1) DEFAULT '0',
  `title` varchar(255) DEFAULT NULL,
  `content` text,
  `topic_id` bigint(20) DEFAULT NULL,
  `label_id` bigint(20) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `topic_id` (`topic_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for draft_box
-- ----------------------------
DROP TABLE IF EXISTS `draft_box`;
CREATE TABLE `draft_box` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `del` tinyint(1) DEFAULT '0',
  `title` varchar(255) DEFAULT NULL,
  `content` text,
  `topic_id` bigint(20) DEFAULT NULL,
  `label_id` bigint(20) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for drafts
-- ----------------------------
DROP TABLE IF EXISTS `drafts`;
CREATE TABLE `drafts` (
  `id` bigint(20) NOT NULL,
  `topic_id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `label_id` bigint(20) DEFAULT NULL,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `invitation_id` bigint(20) DEFAULT NULL,
  `content` text COLLATE utf8mb4_unicode_ci,
  `expire` bigint(20) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for draw_red_package
-- ----------------------------
DROP TABLE IF EXISTS `draw_red_package`;
CREATE TABLE `draw_red_package` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `order_id` varchar(255) DEFAULT NULL,
  `before_income_gold` int(11) DEFAULT NULL,
  `after_income_gold` int(11) DEFAULT NULL,
  `error_code` varchar(255) DEFAULT NULL,
  `gold` int(11) DEFAULT NULL,
  `receive_status` tinyint(2) DEFAULT NULL,
  `status` tinyint(2) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for exchange
-- ----------------------------
DROP TABLE IF EXISTS `exchange`;
CREATE TABLE `exchange` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `price` int(11) unsigned DEFAULT '0',
  `product_num` int(11) unsigned DEFAULT '0',
  `product_type` tinyint(1) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `call_id` varchar(256) DEFAULT NULL,
  `to_user` varchar(256) DEFAULT NULL,
  `video_id` bigint(20) DEFAULT NULL,
  `topic_id` bigint(20) DEFAULT NULL,
  `party_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for feature_topic
-- ----------------------------
DROP TABLE IF EXISTS `feature_topic`;
CREATE TABLE `feature_topic` (
  `id` bigint(20) NOT NULL,
  `expire` int(11) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for feedback
-- ----------------------------
DROP TABLE IF EXISTS `feedback`;
CREATE TABLE `feedback` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `phone` bigint(20) DEFAULT NULL,
  `text` varchar(250) DEFAULT NULL,
  `contact` varchar(200) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for focus
-- ----------------------------
DROP TABLE IF EXISTS `focus`;
CREATE TABLE `focus` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `focus_id` bigint(20) NOT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for history
-- ----------------------------
DROP TABLE IF EXISTS `history`;
CREATE TABLE `history` (
  `id` bigint(20) NOT NULL,
  `attach` text,
  `body` text,
  `conv_type` varchar(20) DEFAULT '',
  `event_type` varchar(5) DEFAULT '',
  `from_account` varchar(45) DEFAULT '',
  `from_client_type` varchar(10) DEFAULT '',
  `from_device_id` varchar(60) DEFAULT '',
  `from_nick` varchar(15) DEFAULT '',
  `msg_timestamp` varchar(20) DEFAULT '',
  `msg_type` varchar(20) DEFAULT '',
  `msgid_client` varchar(60) DEFAULT '',
  `msgid_server` varchar(25) DEFAULT '',
  `to` varchar(45) DEFAULT '',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `ext` varchar(255) DEFAULT NULL,
  `longtitude` varchar(20) DEFAULT '',
  `latitude` varchar(20) DEFAULT '',
  `longitude` varchar(20) DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for interest_group
-- ----------------------------
DROP TABLE IF EXISTS `interest_group`;
CREATE TABLE `interest_group` (
  `id` int(11) NOT NULL,
  `interest_title` varchar(255) NOT NULL COMMENT '兴趣组标题',
  `interest_class` int(2) DEFAULT NULL COMMENT '兴趣分类',
  `interest_label_id` int(2) DEFAULT NULL COMMENT '兴趣标签',
  `interest_type` int(2) DEFAULT NULL COMMENT '兴趣分组类型',
  `created` varchar(20) DEFAULT NULL COMMENT '创建时间',
  `updated` varchar(20) DEFAULT NULL COMMENT '更新时间',
  `status` int(2) DEFAULT NULL COMMENT '房间状态',
  `user_id` varchar(20) DEFAULT NULL COMMENT '用户id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for interest_label
-- ----------------------------
DROP TABLE IF EXISTS `interest_label`;
CREATE TABLE `interest_label` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键自增',
  `tag_name` varchar(20) NOT NULL COMMENT '标签名称',
  `tag_type_id` int(11) NOT NULL COMMENT '1--运动，2--音乐，3--书籍动漫，4--美食，5--影视综艺，6--地方，7--游戏',
  `create` varchar(20) DEFAULT NULL COMMENT '创建时间',
  `update` varchar(20) DEFAULT NULL COMMENT '更新时间',
  `shelf_time` varchar(20) DEFAULT NULL COMMENT '上架时间',
  `user_id` varchar(20) DEFAULT NULL COMMENT '用户id',
  `type` varchar(2) DEFAULT NULL COMMENT '1--移动端，2--后台',
  `status` varchar(2) DEFAULT NULL COMMENT '1--上架，2--下架状态',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for invitation
-- ----------------------------
DROP TABLE IF EXISTS `invitation`;
CREATE TABLE `invitation` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `topic_id` bigint(20) NOT NULL,
  `sequence` int(11) DEFAULT NULL,
  `del` tinyint(1) DEFAULT '0',
  `is_check` tinyint(1) DEFAULT '0',
  `content` text,
  `label_id` bigint(20) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `recommend` tinyint(1) DEFAULT '0',
  `soft_text` tinyint(1) DEFAULT '0',
  `play_times` bigint(20) NOT NULL DEFAULT '0',
  `star_times` bigint(20) NOT NULL DEFAULT '0',
  `collect_times` bigint(20) NOT NULL DEFAULT '0',
  `recommend_time` bigint(20) DEFAULT NULL,
  `featured` tinyint(1) DEFAULT '0',
  `featured_time` bigint(20) DEFAULT NULL,
  `operate_time` bigint(20) DEFAULT NULL,
  `comment_times` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for invitation_comment_star
-- ----------------------------
DROP TABLE IF EXISTS `invitation_comment_star`;
CREATE TABLE `invitation_comment_star` (
  `id` bigint(20) NOT NULL,
  `comment_id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `del` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for invitation_label_recommend
-- ----------------------------
DROP TABLE IF EXISTS `invitation_label_recommend`;
CREATE TABLE `invitation_label_recommend` (
  `id` bigint(20) NOT NULL,
  `invitation_id` bigint(20) DEFAULT NULL,
  `topic_label_id` bigint(20) NOT NULL,
  `del` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`,`topic_label_id`),
  KEY `invitation_id` (`invitation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for invitation_star
-- ----------------------------
DROP TABLE IF EXISTS `invitation_star`;
CREATE TABLE `invitation_star` (
  `id` bigint(20) NOT NULL,
  `invitation_id` bigint(20) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `topic_id` bigint(20) DEFAULT NULL,
  `del` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for label
-- ----------------------------
DROP TABLE IF EXISTS `label`;
CREATE TABLE `label` (
  `id` bigint(20) NOT NULL,
  `name` varchar(60) DEFAULT NULL,
  `info` varchar(200) DEFAULT NULL,
  `label_type` int(11) DEFAULT NULL,
  `role` tinyint(1) DEFAULT '1',
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `color_array` varchar(200) DEFAULT NULL,
  `label_cover_path` varchar(256) DEFAULT NULL,
  `label_out_cover_path` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for label_class
-- ----------------------------
DROP TABLE IF EXISTS `label_class`;
CREATE TABLE `label_class` (
  `id` int(11) NOT NULL,
  `class_name` varchar(255) DEFAULT NULL COMMENT '分类名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for like
-- ----------------------------
DROP TABLE IF EXISTS `like`;
CREATE TABLE `like` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `like_id` bigint(20) NOT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for live_gift
-- ----------------------------
DROP TABLE IF EXISTS `live_gift`;
CREATE TABLE `live_gift` (
  `id` bigint(20) NOT NULL,
  `cover_path` varchar(255) DEFAULT NULL COMMENT '主题封面',
  `gift_name` varchar(128) DEFAULT NULL COMMENT '主题标题',
  `frame_num` int(11) DEFAULT '0' COMMENT '帧数',
  `banana_num` int(11) DEFAULT '0' COMMENT '香蕉数',
  `source_time` int(11) DEFAULT '0' COMMENT '大礼物资源时长',
  `source_name` varchar(128) DEFAULT NULL COMMENT '大礼物资源的文件名',
  `on_shelf` int(1) DEFAULT '1' COMMENT '1：上架 2下架',
  `order` int(11) DEFAULT '1' COMMENT '排序',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for live_theme
-- ----------------------------
DROP TABLE IF EXISTS `live_theme`;
CREATE TABLE `live_theme` (
  `id` bigint(20) unsigned NOT NULL,
  `cover_path` varchar(128) DEFAULT NULL COMMENT '主题封面',
  `theme_name` varchar(50) DEFAULT NULL COMMENT '主题名称',
  `theme_title` varchar(128) DEFAULT NULL COMMENT '主题标题',
  `order` int(11) DEFAULT '1' COMMENT '排序',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `del` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '删除(0：未删除；1：删除)',
  `on_self` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '1上架，0下架',
  `custom_time` int(30) unsigned NOT NULL DEFAULT '0' COMMENT '自定义时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for meet_info
-- ----------------------------
DROP TABLE IF EXISTS `meet_info`;
CREATE TABLE `meet_info` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `text` varchar(256) DEFAULT NULL,
  `gender` tinyint(1) DEFAULT NULL,
  `del` tinyint(1) DEFAULT '1',
  `created` int(11) DEFAULT '1',
  `updated` int(11) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '0',
  `avator` varchar(200) DEFAULT NULL,
  `nick_name` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for member_free_talk
-- ----------------------------
DROP TABLE IF EXISTS `member_free_talk`;
CREATE TABLE `member_free_talk` (
  `id` bigint(20) NOT NULL,
  `uid` bigint(20) DEFAULT NULL,
  `created` bigint(20) DEFAULT NULL,
  `talk_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for member_job
-- ----------------------------
DROP TABLE IF EXISTS `member_job`;
CREATE TABLE `member_job` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` char(255) DEFAULT '' COMMENT '名称',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0失效 1有效',
  `p_id` int(11) unsigned DEFAULT '0' COMMENT '父级ID',
  `path` char(255) DEFAULT '0' COMMENT '路径',
  `lv` int(11) DEFAULT '1' COMMENT '等级',
  `create_time` int(11) DEFAULT '0' COMMENT '添加时间',
  `update_time` int(11) DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='会员职业表';

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `id` bigint(20) NOT NULL,
  `avator` varchar(200) DEFAULT NULL,
  `goto_url` varchar(60) DEFAULT NULL,
  `msg_type` int(11) DEFAULT NULL,
  `content` varchar(512) DEFAULT NULL,
  `alert` varchar(512) DEFAULT NULL,
  `from_id` bigint(20) DEFAULT NULL,
  `to_id` bigint(20) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `is_read` bigint(20) DEFAULT NULL,
  `turn_text` varchar(512) DEFAULT NULL,
  `sub_msg_type` int(11) DEFAULT NULL,
  `order_id` bigint(20) DEFAULT NULL,
  `message_type` tinyint(1) DEFAULT NULL,
  `from_avator_url` varchar(256) DEFAULT NULL,
  `from_nick_name` varchar(128) DEFAULT NULL,
  `from_task_id` bigint(20) DEFAULT NULL,
  `chat_type` tinyint(1) DEFAULT NULL,
  `from_chat_url` varchar(256) DEFAULT NULL,
  `tip` int(11) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  `expiration` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for mission
-- ----------------------------
DROP TABLE IF EXISTS `mission`;
CREATE TABLE `mission` (
  `id` bigint(20) NOT NULL COMMENT 'ID',
  `code` int(6) DEFAULT '0' COMMENT '???????',
  `title` varchar(256) NOT NULL DEFAULT '' COMMENT '???????',
  `cover` varchar(256) DEFAULT '' COMMENT '????ͼƬ',
  `description` varchar(512) DEFAULT '' COMMENT '????????',
  `type` char(50) DEFAULT '' COMMENT '????????, once or everyday',
  `score` int(11) NOT NULL DEFAULT '0' COMMENT '????????',
  `times` int(11) NOT NULL DEFAULT '0' COMMENT '???ۼ????ɴ????? Ĭ?ϣ?1',
  `category` tinyint(1) NOT NULL DEFAULT '1' COMMENT '???࣬1-ÿ????????2-???????',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '???',
  `enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '?Ƿ?????????,0-?رգ?1-???',
  `disabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '?Ƿ???ʾ????,0-?رգ?1-???',
  `uri` varchar(512) DEFAULT '' COMMENT '֧??????, json [??share??, ...]???????',
  `extend` varchar(512) DEFAULT '' COMMENT '??չ??json?ַ????????',
  `created` int(11) DEFAULT '0' COMMENT '????ʱ?',
  `updated` int(11) DEFAULT '0' COMMENT '?޸?ʱ?',
  `deleted` tinyint(3) NOT NULL DEFAULT '0' COMMENT 'ɾ???ֶ',
  `property` tinyint(3) NOT NULL DEFAULT '0' COMMENT '区分类型',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='?????????';

-- ----------------------------
-- Table structure for module
-- ----------------------------
DROP TABLE IF EXISTS `module`;
CREATE TABLE `module` (
  `id` bigint(20) NOT NULL,
  `name` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for new_comment
-- ----------------------------
DROP TABLE IF EXISTS `new_comment`;
CREATE TABLE `new_comment` (
  `id` bigint(20) NOT NULL,
  `user_label_id` bigint(20) NOT NULL DEFAULT '0',
  `content` varchar(2056) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `user_id` bigint(20) NOT NULL,
  `is_check` tinyint(1) DEFAULT '1',
  `age` bigint(20) DEFAULT NULL,
  `gender` tinyint(1) DEFAULT NULL,
  `star_times` int(11) DEFAULT '0',
  `type` tinyint(1) DEFAULT '1',
  `voice_time` int(11) DEFAULT '0',
  `avtoar` varchar(200) DEFAULT NULL,
  `topic_id` bigint(20) NOT NULL DEFAULT '0',
  `nick_name` varchar(45) DEFAULT NULL,
  `is_del` tinyint(1) DEFAULT '0',
  `to_user` bigint(20) DEFAULT '0',
  `to_user_nick_name` varchar(45) DEFAULT NULL,
  `longitude` varchar(20) DEFAULT '',
  `latitude` varchar(20) DEFAULT '',
  `check_comment` varchar(256) DEFAULT '' COMMENT '审核备注，json字符串，{涉黄,涉暴}',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for open_screen
-- ----------------------------
DROP TABLE IF EXISTS `open_screen`;
CREATE TABLE `open_screen` (
  `id` int(70) unsigned NOT NULL AUTO_INCREMENT,
  `type` smallint(1) unsigned NOT NULL DEFAULT '1' COMMENT '开屏类型（1：图片，2：视频）',
  `screen_url` varchar(255) DEFAULT '-1' COMMENT '封面地址',
  `advertising_link` varchar(255) DEFAULT '-1' COMMENT '跳转地址',
  `skip_time` smallint(4) unsigned NOT NULL DEFAULT '3' COMMENT '跳过时间',
  `created` int(30) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `start_time` int(30) unsigned NOT NULL DEFAULT '0' COMMENT '开始时间',
  `end_time` int(30) unsigned NOT NULL DEFAULT '0' COMMENT '结束时间',
  `user_id` bigint(50) DEFAULT '-1' COMMENT '操作员',
  `remark` mediumtext COMMENT '备注',
  `updated` int(30) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `delete` smallint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除（0：未删除，1：已删除）',
  `up_down` smallint(1) unsigned NOT NULL DEFAULT '1' COMMENT '1：上架，0：下架',
  `name` varchar(50) DEFAULT NULL,
  `share_url` varchar(255) DEFAULT '-1' COMMENT '分享活动链接地址',
  `display_times` int(10) unsigned DEFAULT '0' COMMENT '显示次数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=55 DEFAULT CHARSET=utf8mb4 COMMENT='开屏表';

-- ----------------------------
-- Table structure for open_screen_click
-- ----------------------------
DROP TABLE IF EXISTS `open_screen_click`;
CREATE TABLE `open_screen_click` (
  `click_id` bigint(70) NOT NULL,
  `click_time` int(30) unsigned NOT NULL DEFAULT '0' COMMENT '点击时间',
  `open_id` bigint(70) DEFAULT '0',
  PRIMARY KEY (`click_id`),
  UNIQUE KEY `id` (`click_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='开屏点击记录表';

-- ----------------------------
-- Table structure for operate_info
-- ----------------------------
DROP TABLE IF EXISTS `operate_info`;
CREATE TABLE `operate_info` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `nick_name` varchar(50) DEFAULT NULL,
  `operate_type` tinyint(1) DEFAULT NULL,
  `operate_content` varchar(50) DEFAULT NULL,
  `details` varchar(512) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for order
-- ----------------------------
DROP TABLE IF EXISTS `order`;
CREATE TABLE `order` (
  `id` bigint(20) NOT NULL,
  `tip` int(11) DEFAULT '0',
  `status` tinyint(1) DEFAULT NULL,
  `send_type` tinyint(1) DEFAULT NULL,
  `initiator_id` bigint(20) DEFAULT NULL,
  `recipient_id` bigint(20) DEFAULT NULL,
  `call_id` varchar(256) DEFAULT NULL,
  `task_id` bigint(20) DEFAULT NULL,
  `expiration` int(11) DEFAULT NULL,
  `initiator_num` int(11) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for part_history
-- ----------------------------
DROP TABLE IF EXISTS `part_history`;
CREATE TABLE `part_history` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) DEFAULT NULL COMMENT '主播id',
  `cover` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '房间封面',
  `channel_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '绑定直播回放的id',
  `label` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '标签 王牌主播, 新人',
  `topic` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '主题id',
  `desc` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '直播间描述',
  `pull_url` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `numbers` int(11) DEFAULT '0' COMMENT '参与人数',
  `begin_time` int(11) DEFAULT NULL COMMENT '直播开始时间',
  `end_time` int(11) DEFAULT NULL COMMENT '直播结束时间',
  `city` varchar(40) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '城市',
  `labeltype` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '上下架(0：下架，1：上架)',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `pull_url` (`pull_url`)
) ENGINE=InnoDB AUTO_INCREMENT=4844 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for party_banner
-- ----------------------------
DROP TABLE IF EXISTS `party_banner`;
CREATE TABLE `party_banner` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cover_path` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'banner图片',
  `skip_type` tinyint(1) DEFAULT '1' COMMENT '跳转类型（1：url，2：直播间，3：回放，4：聚合，5：其他）',
  `user_id` int(11) unsigned DEFAULT '0' COMMENT '主播id',
  `room_id` int(11) unsigned DEFAULT NULL COMMENT '房间id',
  `skip_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '跳转链接',
  `channel_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '绑定直播回放的id',
  `start_time` int(50) unsigned NOT NULL DEFAULT '0' COMMENT '开始时间',
  `end_time` int(50) unsigned NOT NULL DEFAULT '0' COMMENT '结束时间',
  `sort` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `on_shelf` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '上下架(0：下架，1：上架)',
  `type` smallint(1) DEFAULT '1' COMMENT '1普通banner 2明星banner',
  `title` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '-1' COMMENT 'banner标题',
  `remark` mediumtext COLLATE utf8mb4_unicode_ci COMMENT '备注',
  `created` int(50) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int(50) unsigned DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=169 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for pay
-- ----------------------------
DROP TABLE IF EXISTS `pay`;
CREATE TABLE `pay` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `gold` int(11) DEFAULT NULL,
  `msg` varchar(200) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `transaction` varchar(200) DEFAULT NULL,
  `pay_type` int(11) DEFAULT NULL,
  `pay_way` int(11) DEFAULT NULL,
  `error_code` varchar(150) DEFAULT NULL,
  `pay_order_id` varchar(150) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `price` int(11) unsigned DEFAULT '0',
  `product_num` int(11) unsigned DEFAULT '0',
  `type` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for public_click
-- ----------------------------
DROP TABLE IF EXISTS `public_click`;
CREATE TABLE `public_click` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `other_id` int(11) unsigned DEFAULT NULL COMMENT '其他ID',
  `type` int(2) unsigned DEFAULT '1' COMMENT '其他ID的类型 1:派对bannerID ，',
  `click` int(20) unsigned DEFAULT '0' COMMENT '点击数',
  `created` int(50) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=167 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for push_member_activity_log
-- ----------------------------
DROP TABLE IF EXISTS `push_member_activity_log`;
CREATE TABLE `push_member_activity_log` (
  `id` bigint(20) NOT NULL,
  `push_message_id` bigint(20) NOT NULL COMMENT '??????ϢID',
  `to_user_id` bigint(20) NOT NULL COMMENT '???Ͷ??',
  `title` varchar(255) DEFAULT '' COMMENT '??Ϣ???',
  `res_value` varchar(255) DEFAULT '' COMMENT '???ͷ???ֵ?????л?',
  `created` bigint(20) DEFAULT '0' COMMENT '????ʱ?',
  `status` int(11) DEFAULT '0' COMMENT '״̬ 0-δ???ͣ?1-???',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='???ͻ?û???־?';

-- ----------------------------
-- Table structure for push_member_system_log
-- ----------------------------
DROP TABLE IF EXISTS `push_member_system_log`;
CREATE TABLE `push_member_system_log` (
  `id` bigint(20) NOT NULL,
  `push_message_id` bigint(20) NOT NULL COMMENT '??????ϢID',
  `to_user_id` bigint(20) NOT NULL COMMENT '???Ͷ??',
  `title` varchar(255) DEFAULT '' COMMENT '??Ϣ???',
  `res_value` varchar(255) DEFAULT '' COMMENT '???ͷ???ֵ?????л?',
  `created` bigint(20) DEFAULT '0' COMMENT '????ʱ?',
  `status` int(11) DEFAULT '0' COMMENT '״̬ 0-δ???ͣ?1-???',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='????ϵͳ?û???־?';

-- ----------------------------
-- Table structure for push_system_message
-- ----------------------------
DROP TABLE IF EXISTS `push_system_message`;
CREATE TABLE `push_system_message` (
  `id` bigint(20) NOT NULL,
  `title` varchar(255) DEFAULT '' COMMENT '???ͱ??',
  `created` bigint(20) DEFAULT '0' COMMENT '????ʱ?',
  `updated` bigint(20) DEFAULT '0' COMMENT '????ʱ?',
  `status` int(11) DEFAULT '2' COMMENT '״̬ 1 ?ϼܣ?2 ?¼',
  `content` varchar(1024) DEFAULT '' COMMENT '??Ϣ???',
  `target_type` tinyint(3) DEFAULT '1' COMMENT 'Ŀ?????ͣ?1-ȫ????2-??ͨ?û???3-?ٷ????ţ?4-???????ţ?5-??ǩ???ţ?6-?ٷ??˺ţ?7-???',
  `target_population` varchar(1024) DEFAULT '' COMMENT 'Ŀ????Ⱥ????target_type Ϊ 7 ʱ??Ч',
  `putaway_time` int(11) DEFAULT '0' COMMENT '?Զ?????Чʱ?',
  `push_status` int(11) DEFAULT '0' COMMENT '????״̬??״̬ 0-δ???ͣ?1-???',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='????ϵͳ??Ϣ?';

-- ----------------------------
-- Table structure for refund
-- ----------------------------
DROP TABLE IF EXISTS `refund`;
CREATE TABLE `refund` (
  `id` bigint(20) NOT NULL,
  `order_id` bigint(20) DEFAULT NULL,
  `transaction` varchar(128) DEFAULT NULL,
  `refund_type` tinyint(1) DEFAULT NULL,
  `total_fee` int(11) DEFAULT NULL,
  `refund_fee` int(11) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  `msg` varchar(128) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for region
-- ----------------------------
DROP TABLE IF EXISTS `region`;
CREATE TABLE `region` (
  `REGION_ID` double NOT NULL COMMENT '区域ID',
  `REGION_CODE` varchar(100) NOT NULL COMMENT '地区代码',
  `REGION_NAME` varchar(100) NOT NULL COMMENT '区域名称',
  `PARENT_ID` double NOT NULL COMMENT '上级ID',
  `REGION_LEVEL` double NOT NULL COMMENT '地区级',
  `REGION_ORDER` double NOT NULL COMMENT '地区秩序',
  `REGION_NAME_EN` varchar(100) NOT NULL COMMENT '区域英文名称',
  `REGION_SHORTNAME_EN` varchar(10) NOT NULL COMMENT '英文地区名',
  PRIMARY KEY (`REGION_ID`),
  UNIQUE KEY `REGION_ID` (`REGION_ID`) USING BTREE,
  KEY `PARENT_ID` (`PARENT_ID`) USING BTREE,
  KEY `REGION_NAME` (`REGION_NAME`) USING BTREE,
  KEY `REGION_NAME_EN` (`REGION_NAME_EN`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='全国省市区联动表';

-- ----------------------------
-- Table structure for reward_code
-- ----------------------------
DROP TABLE IF EXISTS `reward_code`;
CREATE TABLE `reward_code` (
  `id` bigint(20) NOT NULL,
  `type` int(11) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for school
-- ----------------------------
DROP TABLE IF EXISTS `school`;
CREATE TABLE `school` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` char(16) NOT NULL,
  `place` char(8) NOT NULL,
  `type` char(8) NOT NULL,
  `properties` char(8) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2577 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for score_move
-- ----------------------------
DROP TABLE IF EXISTS `score_move`;
CREATE TABLE `score_move` (
  `id` bigint(20) NOT NULL COMMENT 'ID',
  `user_id` bigint(20) NOT NULL COMMENT 'user???û?ID',
  `source` varchar(256) DEFAULT '' COMMENT '??Դ????',
  `source_id` bigint(20) NOT NULL COMMENT '??Դ??ID',
  `source_data` varchar(1024) DEFAULT '' COMMENT '????ժҪ?? ???ڼ?¼?????????ѻ??ֵĹؼ???Ϣ, ???һ???ӰƱ: {"name": "cdkey", "value":"abcdefg"}',
  `total` int(11) NOT NULL DEFAULT '0' COMMENT '????, ???????ӣ? ?????????ּ??٣? ????',
  `created` int(11) DEFAULT '0' COMMENT '????ʱ?',
  `updated` int(11) DEFAULT '0' COMMENT '?޸?ʱ?',
  `date` varchar(20) DEFAULT NULL COMMENT '???????',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='??????ˮ';

-- ----------------------------
-- Table structure for self_tip
-- ----------------------------
DROP TABLE IF EXISTS `self_tip`;
CREATE TABLE `self_tip` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT 'user???û?ID',
  `title` varchar(128) DEFAULT '' COMMENT '???',
  `comment` varchar(256) DEFAULT '' COMMENT '???',
  `created` int(11) DEFAULT '0' COMMENT '????ʱ?',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=10 DEFAULT CHARSET=utf8 COMMENT='??????ˮ';

-- ----------------------------
-- Table structure for share_record
-- ----------------------------
DROP TABLE IF EXISTS `share_record`;
CREATE TABLE `share_record` (
  `id` bigint(20) NOT NULL,
  `type` int(11) DEFAULT NULL,
  `reward_code` varchar(256) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `enter_total` bigint(20) DEFAULT NULL,
  `download_total` bigint(20) DEFAULT NULL,
  `register_total` bigint(20) DEFAULT NULL,
  `pay_total` bigint(20) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `exchange_total` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for shield_user_model
-- ----------------------------
DROP TABLE IF EXISTS `shield_user_model`;
CREATE TABLE `shield_user_model` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `shield_id` bigint(20) DEFAULT NULL,
  `info` varchar(255) DEFAULT NULL,
  `status` int(4) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for shiyan
-- ----------------------------
DROP TABLE IF EXISTS `shiyan`;
CREATE TABLE `shiyan` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `cont` mediumtext,
  `time` smallint(50) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for short_video
-- ----------------------------
DROP TABLE IF EXISTS `short_video`;
CREATE TABLE `short_video` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL COMMENT '用户id',
  `forward` int(11) DEFAULT NULL,
  `thumb_down` bigint(20) DEFAULT NULL,
  `thumb_up` bigint(20) DEFAULT NULL,
  `comment_count` int(11) DEFAULT NULL,
  `douyin_user_id` varchar(255) DEFAULT NULL,
  `douyin_challenge` varchar(255) DEFAULT NULL,
  `download_url` varchar(255) DEFAULT NULL,
  `cover_url` varchar(255) DEFAULT NULL,
  `one_sentence` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `height` int(11) DEFAULT NULL,
  `width` int(11) DEFAULT NULL,
  `created` bigint(20) DEFAULT NULL,
  `size` bigint(20) DEFAULT NULL,
  `topic_tag` varchar(255) DEFAULT NULL,
  `origin_url` varchar(255) DEFAULT NULL,
  `approve_status` int(11) DEFAULT NULL,
  `approve_time` bigint(20) DEFAULT NULL,
  `top` tinyint(1) DEFAULT NULL,
  `top_time` bigint(20) DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT NULL,
  `source` varchar(255) DEFAULT NULL COMMENT 'xingxing, huoshan, xiaokaxiu, douyin',
  `o_thumb_up` bigint(20) DEFAULT '0' COMMENT '原始点赞数',
  `o_thumb_down` bigint(20) DEFAULT '0' COMMENT '原始点踩数',
  `o_comment_count` bigint(20) DEFAULT '0',
  `o_share_count` bigint(20) DEFAULT '0',
  `o_created` bigint(20) DEFAULT '0',
  `o_play_count` bigint(20) DEFAULT '0',
  `x_thumb_up` bigint(20) DEFAULT '0',
  `x_thumb_down` bigint(20) DEFAULT '0',
  `x_comment_count` bigint(20) DEFAULT '0',
  `x_share_count` bigint(20) DEFAULT '0',
  `x_play_count` bigint(20) DEFAULT '0',
  `op_top` bigint(20) DEFAULT '0' COMMENT '0, 1, 1表示置顶',
  `op_top_time` bigint(20) DEFAULT '0',
  `op_weight_level` bigint(20) DEFAULT '0',
  `o_user_id` varchar(255) DEFAULT NULL,
  `o_theme` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `o_url` varchar(255) DEFAULT NULL,
  `theme_tag` varchar(255) DEFAULT NULL COMMENT '主题标签',
  `deleted` tinyint(1) DEFAULT '0',
  `o_nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `o_user_gender` int(11) DEFAULT NULL,
  `o_user_age` int(11) DEFAULT NULL,
  `o_user_avatar` varchar(255) DEFAULT NULL,
  `op_weight` bigint(20) DEFAULT '0' COMMENT '运营添加权重热度',
  `op_weight_levle` bigint(20) DEFAULT '0',
  `op_weight_time` bigint(20) DEFAULT '0',
  `approve_reason` varchar(255) DEFAULT NULL,
  `updated` bigint(20) DEFAULT '0',
  `approve_decision` int(11) DEFAULT '0',
  `latitude` double DEFAULT NULL,
  `longitude` double DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='短视频';

-- ----------------------------
-- Table structure for short_video_activity
-- ----------------------------
DROP TABLE IF EXISTS `short_video_activity`;
CREATE TABLE `short_video_activity` (
  `id` bigint(20) NOT NULL,
  `name` varchar(255) NOT NULL,
  `cover` varchar(255) NOT NULL,
  `start` bigint(20) NOT NULL,
  `end` bigint(20) NOT NULL,
  `sort` int(11) DEFAULT NULL,
  `introduction` varchar(255) NOT NULL DEFAULT '',
  `view_count` bigint(20) NOT NULL DEFAULT '0',
  `created` bigint(20) NOT NULL,
  `position` int(11) DEFAULT NULL,
  `on_shelf` int(11) DEFAULT '0',
  `updated` bigint(20) DEFAULT '0',
  `recommend_count` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for short_video_activity_relation
-- ----------------------------
DROP TABLE IF EXISTS `short_video_activity_relation`;
CREATE TABLE `short_video_activity_relation` (
  `id` bigint(20) NOT NULL,
  `video_id` bigint(20) NOT NULL,
  `activity_id` bigint(20) NOT NULL,
  `recommended` tinyint(1) DEFAULT NULL,
  `created` bigint(20) DEFAULT NULL,
  `updated` bigint(20) DEFAULT '0',
  `recommended_time` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for short_video_comment
-- ----------------------------
DROP TABLE IF EXISTS `short_video_comment`;
CREATE TABLE `short_video_comment` (
  `id` bigint(20) NOT NULL,
  `video_id` bigint(20) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `x_thumb_up` bigint(20) DEFAULT '0',
  `x_thumb_down` bigint(20) DEFAULT '0',
  `o_thumb_up` bigint(20) DEFAULT '0',
  `o_thumb_down` bigint(20) DEFAULT NULL,
  `o_created` bigint(20) DEFAULT NULL,
  `created` bigint(20) DEFAULT NULL,
  `o_nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `o_user_avatar` varchar(255) DEFAULT NULL,
  `o_user_gender` bigint(20) DEFAULT NULL,
  `o_user_age` bigint(20) DEFAULT NULL,
  `source` varchar(255) DEFAULT NULL,
  `approve_status` int(11) DEFAULT NULL,
  `latitude` double NOT NULL DEFAULT '0',
  `longitude` double NOT NULL DEFAULT '0',
  `to_user_id` bigint(20) DEFAULT '0',
  `to_user_o_name` varchar(255) DEFAULT '',
  `to_comment_id` bigint(20) DEFAULT NULL,
  `type` int(11) DEFAULT '1',
  `to_user_name` varchar(255) DEFAULT NULL,
  `duration` int(11) DEFAULT NULL,
  `del` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for short_video_report
-- ----------------------------
DROP TABLE IF EXISTS `short_video_report`;
CREATE TABLE `short_video_report` (
  `id` bigint(20) NOT NULL,
  `video_id` bigint(20) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `decision` int(11) DEFAULT NULL,
  `decide_time` bigint(20) DEFAULT NULL,
  `created` bigint(20) NOT NULL,
  `updated` bigint(20) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for short_video_tag_relation
-- ----------------------------
DROP TABLE IF EXISTS `short_video_tag_relation`;
CREATE TABLE `short_video_tag_relation` (
  `id` bigint(20) NOT NULL,
  `video_id` bigint(20) DEFAULT NULL,
  `tag_id` bigint(20) DEFAULT NULL,
  `tag_level` int(11) DEFAULT NULL,
  `recommended` tinyint(1) DEFAULT NULL,
  `related_time` bigint(20) DEFAULT NULL,
  `del` int(11) NOT NULL DEFAULT '0',
  `created` bigint(20) DEFAULT NULL,
  `updated` bigint(20) DEFAULT '0',
  `recommended_time` bigint(20) DEFAULT NULL,
  `tag_name` varchar(255) DEFAULT NULL,
  `shelves` int(11) DEFAULT '1' COMMENT '1--上架 ， 2 -- 下架',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for short_video_theme_tag
-- ----------------------------
DROP TABLE IF EXISTS `short_video_theme_tag`;
CREATE TABLE `short_video_theme_tag` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for short_video_thumb
-- ----------------------------
DROP TABLE IF EXISTS `short_video_thumb`;
CREATE TABLE `short_video_thumb` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `target_type` int(11) DEFAULT NULL,
  `target_id` bigint(20) DEFAULT NULL,
  `up_or_down` bigint(20) DEFAULT NULL,
  `created` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for sms_checked_log
-- ----------------------------
DROP TABLE IF EXISTS `sms_checked_log`;
CREATE TABLE `sms_checked_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号码',
  `ip` varchar(255) NOT NULL DEFAULT '' COMMENT '发送者的IP',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='短信日志检查表';

-- ----------------------------
-- Table structure for sms_log
-- ----------------------------
DROP TABLE IF EXISTS `sms_log`;
CREATE TABLE `sms_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `phone` text NOT NULL COMMENT '发送短信手机号码,多个以,分割，最多200个',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '短信内容',
  `send_time` int(11) NOT NULL DEFAULT '0' COMMENT '发送时间',
  `send_log` text NOT NULL COMMENT '发送短信的日志，以序列化字符串形式存放',
  `send_id` int(11) NOT NULL DEFAULT '0' COMMENT '操作人ID',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='后台发送短信日志表';

-- ----------------------------
-- Table structure for star
-- ----------------------------
DROP TABLE IF EXISTS `star`;
CREATE TABLE `star` (
  `id` bigint(20) NOT NULL,
  `user_label_id` bigint(20) NOT NULL DEFAULT '0',
  `user_id` bigint(20) NOT NULL,
  `topic_id` bigint(20) NOT NULL DEFAULT '0',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user` (
  `id` bigint(20) NOT NULL,
  `phone` bigint(20) NOT NULL,
  `nick_name` varchar(45) DEFAULT NULL,
  `login_name` varchar(45) DEFAULT NULL,
  `login_password` varchar(45) DEFAULT NULL,
  `role_id` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `delete` smallint(1) DEFAULT '0' COMMENT '删除键',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for system_activity_message
-- ----------------------------
DROP TABLE IF EXISTS `system_activity_message`;
CREATE TABLE `system_activity_message` (
  `id` bigint(20) NOT NULL,
  `start` bigint(20) DEFAULT '0',
  `end` bigint(20) DEFAULT '0',
  `title` varchar(255) DEFAULT '',
  `imag` varchar(255) DEFAULT '',
  `created` bigint(20) DEFAULT '0',
  `status` int(11) DEFAULT '2',
  `jump_type` int(11) DEFAULT '0',
  `jump_dest` varchar(255) DEFAULT '',
  `jump_id` bigint(20) DEFAULT '0',
  `img` varchar(255) DEFAULT '',
  `share_subtitle` varchar(255) DEFAULT '',
  `share_icon` varchar(255) DEFAULT NULL,
  `jump_title` varchar(255) DEFAULT '',
  `push_status` tinyint(3) DEFAULT '0' COMMENT '????״̬??״̬ 0-δ???ͣ?1-???',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for system_voice_index
-- ----------------------------
DROP TABLE IF EXISTS `system_voice_index`;
CREATE TABLE `system_voice_index` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `gender` int(11) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `max` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for system_voice_sound
-- ----------------------------
DROP TABLE IF EXISTS `system_voice_sound`;
CREATE TABLE `system_voice_sound` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `gender` int(11) DEFAULT NULL,
  `template_link` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UQE_system_voice_sound_name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=132 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for tag_level_first
-- ----------------------------
DROP TABLE IF EXISTS `tag_level_first`;
CREATE TABLE `tag_level_first` (
  `id` bigint(20) NOT NULL,
  `name` varchar(256) DEFAULT NULL,
  `expire` int(11) DEFAULT '0',
  `del` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `is_special` tinyint(1) DEFAULT '0',
  `user_participate` tinyint(1) DEFAULT '1',
  `recommend_count` bigint(20) DEFAULT NULL,
  `start` bigint(20) DEFAULT NULL,
  `sort` int(11) DEFAULT NULL,
  `cover` varchar(255) DEFAULT NULL,
  `is_default_tag` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tag_level_second
-- ----------------------------
DROP TABLE IF EXISTS `tag_level_second`;
CREATE TABLE `tag_level_second` (
  `id` bigint(20) NOT NULL,
  `name` varchar(256) DEFAULT NULL,
  `expire` int(11) DEFAULT NULL,
  `del` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `level_first_id` bigint(20) DEFAULT NULL,
  `is_special` tinyint(1) DEFAULT '0',
  `user_participate` tinyint(1) DEFAULT '1',
  `recommend_count` bigint(20) DEFAULT NULL,
  `start` bigint(20) DEFAULT NULL,
  `sort` int(11) DEFAULT NULL,
  `cover` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tag_level_third
-- ----------------------------
DROP TABLE IF EXISTS `tag_level_third`;
CREATE TABLE `tag_level_third` (
  `id` bigint(20) NOT NULL,
  `name` varchar(256) DEFAULT NULL,
  `expire` int(11) DEFAULT NULL,
  `del` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `level_second_id` bigint(20) DEFAULT NULL,
  `is_special` tinyint(1) DEFAULT '0',
  `user_participate` tinyint(1) DEFAULT '1',
  `recommend_count` bigint(20) DEFAULT NULL,
  `start` bigint(20) DEFAULT NULL,
  `sort` int(11) DEFAULT NULL,
  `cover` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tag_module_connection
-- ----------------------------
DROP TABLE IF EXISTS `tag_module_connection`;
CREATE TABLE `tag_module_connection` (
  `id` bigint(20) NOT NULL,
  `tag_id` bigint(20) NOT NULL,
  `module_id` bigint(20) NOT NULL,
  `connected` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for task
-- ----------------------------
DROP TABLE IF EXISTS `task`;
CREATE TABLE `task` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `task_type` tinyint(1) DEFAULT NULL,
  `task_url` varchar(256) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `duration` int(11) unsigned DEFAULT '0',
  `cover_path` varchar(256) DEFAULT NULL,
  `del` tinyint(1) DEFAULT '0',
  `min_tip` int(11) unsigned DEFAULT '0',
  `call_time` varchar(256) DEFAULT NULL,
  `comment_total` int(11) DEFAULT '0',
  `check` tinyint(1) DEFAULT '1',
  `is_check` tinyint(1) DEFAULT '1',
  `user_label_id` bigint(20) NOT NULL,
  `user_label_name` varchar(256) DEFAULT NULL,
  `pricing` int(11) DEFAULT NULL,
  `content` varchar(256) DEFAULT NULL,
  `label_type` int(11) DEFAULT NULL,
  `gender` bigint(20) DEFAULT NULL,
  `cause` varchar(124) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for test
-- ----------------------------
DROP TABLE IF EXISTS `test`;
CREATE TABLE `test` (
  `id` int(11) NOT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `user_lable_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for threshold
-- ----------------------------
DROP TABLE IF EXISTS `threshold`;
CREATE TABLE `threshold` (
  `id` bigint(20) NOT NULL,
  `star` bigint(6) DEFAULT NULL,
  `comment` bigint(6) DEFAULT NULL,
  `watch` bigint(6) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for topic
-- ----------------------------
DROP TABLE IF EXISTS `topic`;
CREATE TABLE `topic` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `is_check` tinyint(1) DEFAULT '0',
  `comment_total` int(11) DEFAULT '0',
  `star_times` int(11) DEFAULT '0',
  `play_times` int(11) DEFAULT '0',
  `op_weight` bigint(20) DEFAULT '0',
  `label_id` bigint(20) DEFAULT NULL,
  `del` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `invitation_total` int(11) DEFAULT '0',
  `title` varchar(256) DEFAULT NULL,
  `collect_times` int(11) DEFAULT '0',
  `invitation_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for topic_banner
-- ----------------------------
DROP TABLE IF EXISTS `topic_banner`;
CREATE TABLE `topic_banner` (
  `id` bigint(20) NOT NULL,
  `cover` varchar(255) DEFAULT '',
  `jump_dest` varchar(255) DEFAULT '',
  `sort` int(11) DEFAULT NULL,
  `del` int(11) DEFAULT '0',
  `on_shelf` int(11) DEFAULT '1',
  `click` int(11) DEFAULT '0',
  `start` bigint(20) DEFAULT '0',
  `end` bigint(20) DEFAULT '0',
  `created` bigint(20) DEFAULT '0',
  `name` varchar(255) DEFAULT '',
  `updated` bigint(20) DEFAULT '0',
  `position` int(11) DEFAULT '0',
  `remark` varchar(255) DEFAULT NULL,
  `share_url` varchar(255) DEFAULT NULL COMMENT '分享活动链接地址',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for topic_banner_click
-- ----------------------------
DROP TABLE IF EXISTS `topic_banner_click`;
CREATE TABLE `topic_banner_click` (
  `id` bigint(20) NOT NULL,
  `uid` bigint(20) DEFAULT NULL,
  `banner_id` bigint(20) DEFAULT NULL,
  `platform` varchar(255) DEFAULT NULL,
  `model` varchar(255) DEFAULT NULL,
  `created` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for topic_label
-- ----------------------------
DROP TABLE IF EXISTS `topic_label`;
CREATE TABLE `topic_label` (
  `id` bigint(20) NOT NULL,
  `name` varchar(60) DEFAULT NULL,
  `info` varchar(200) DEFAULT NULL,
  `label_cover_path` varchar(256) DEFAULT NULL,
  `label_out_cover_path` varchar(256) DEFAULT NULL,
  `del` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `sort` bigint(20) DEFAULT NULL,
  `putaway_time` int(11) DEFAULT NULL,
  `on_lookers` bigint(20) NOT NULL,
  `height` int(11) DEFAULT NULL,
  `width` int(11) DEFAULT NULL,
  `shelves` int(11) DEFAULT '1' COMMENT '1--上架，2--下架',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for topic_label_copy
-- ----------------------------
DROP TABLE IF EXISTS `topic_label_copy`;
CREATE TABLE `topic_label_copy` (
  `id` bigint(20) NOT NULL,
  `name` varchar(60) DEFAULT NULL,
  `info` varchar(200) DEFAULT NULL,
  `label_cover_path` varchar(256) DEFAULT NULL,
  `label_out_cover_path` varchar(256) DEFAULT NULL,
  `del` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `sort` bigint(20) DEFAULT NULL,
  `putaway_time` int(11) DEFAULT NULL,
  `on_lookers` bigint(20) NOT NULL,
  `height` int(11) DEFAULT NULL,
  `width` int(11) DEFAULT NULL,
  `shelves` int(11) DEFAULT '1' COMMENT '1--上架，2--下架',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for topic_recommend
-- ----------------------------
DROP TABLE IF EXISTS `topic_recommend`;
CREATE TABLE `topic_recommend` (
  `id` bigint(20) NOT NULL,
  `expire` int(11) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for unlock_reason
-- ----------------------------
DROP TABLE IF EXISTS `unlock_reason`;
CREATE TABLE `unlock_reason` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `other_id` bigint(20) DEFAULT NULL COMMENT '其他id  type 1：房间id     type2：用户id',
  `type` tinyint(1) DEFAULT NULL COMMENT '1 直播间，2用户',
  `cause` varchar(255) DEFAULT NULL COMMENT '原因',
  `operator_id` int(11) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL,
  `wx_union_id` varchar(100) CHARACTER SET utf8 DEFAULT NULL,
  `nick_name` varchar(45) DEFAULT NULL,
  `phone` bigint(20) NOT NULL,
  `avator` varchar(200) CHARACTER SET utf8 DEFAULT NULL,
  `system` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `user_type` tinyint(1) DEFAULT '0',
  `gender` tinyint(1) DEFAULT NULL,
  `gold` int(11) unsigned DEFAULT '0',
  `income_gold` int(11) unsigned DEFAULT '0',
  `hot_degree` int(11) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  `wy_user` varchar(128) DEFAULT NULL,
  `wy_password` varchar(256) CHARACTER SET utf8 DEFAULT NULL,
  `role` tinyint(1) DEFAULT '1',
  `shield_user` text CHARACTER SET utf8,
  `call_duration` int(11) DEFAULT '0',
  `last_login_time` int(11) DEFAULT NULL,
  `last_login_history` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `signature` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `device_id` varchar(200) CHARACTER SET utf8 DEFAULT NULL,
  `device_token` varchar(200) CHARACTER SET utf8 DEFAULT NULL,
  `device_type` int(11) DEFAULT NULL,
  `fans_num` int(11) DEFAULT '0',
  `focus_num` int(11) NOT NULL DEFAULT '0',
  `used_duration` int(11) DEFAULT '0',
  `no_like_user` text CHARACTER SET utf8,
  `q_q_openid` varchar(256) CHARACTER SET utf8 DEFAULT NULL,
  `wx_g_z_open_id` varchar(60) CHARACTER SET utf8 DEFAULT NULL,
  `weibo_uid` bigint(20) DEFAULT NULL,
  `membership_date` bigint(20) DEFAULT '0',
  `city_num` int(11) DEFAULT NULL,
  `first_login_city` varchar(200) CHARACTER SET utf8 DEFAULT NULL,
  `match_gender` int(11) DEFAULT '0',
  `is_accept_member` int(11) DEFAULT '0',
  `facebook_id` bigint(20) DEFAULT NULL,
  `google_id` varchar(32) CHARACTER SET utf8 DEFAULT NULL,
  `play_gold` int(11) DEFAULT '0',
  `play_be_gold` int(11) DEFAULT '0',
  `c_code` varchar(256) CHARACTER SET utf8 DEFAULT NULL,
  `is_first_recharge` tinyint(1) DEFAULT '1',
  `check_time` int(11) DEFAULT '0',
  `recharge_total` int(11) DEFAULT '0',
  `password` varchar(60) CHARACTER SET utf8 DEFAULT NULL,
  `user_extra` text CHARACTER SET utf8,
  `v_c_code` varchar(256) CHARACTER SET utf8 DEFAULT NULL,
  `order` int(11) DEFAULT '0',
  `vc_code` varchar(256) CHARACTER SET utf8 DEFAULT NULL,
  `like_num` int(11) DEFAULT '0',
  `liked_num` int(11) DEFAULT '0',
  `province` varchar(60) CHARACTER SET utf8 DEFAULT NULL,
  `born` bigint(20) DEFAULT NULL,
  `age` bigint(20) DEFAULT NULL,
  `vc_income_gold` int(11) unsigned DEFAULT '0',
  `is_popular` int(11) DEFAULT '0',
  `last_active_time` int(11) DEFAULT NULL,
  `income_yuan` varchar(256) CHARACTER SET utf8 DEFAULT NULL,
  `constellation` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `city` varchar(64) CHARACTER SET utf8 DEFAULT NULL,
  `location_x` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `location_y` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `affective_state` int(11) DEFAULT '0',
  `profession` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `school` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `home_town` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `is_gay` int(11) DEFAULT '0',
  `ali_account` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `friends_num` int(11) DEFAULT NULL,
  `ali_pay_user_id` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `ali_pay_phone` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `ali_pay_real_name` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `xingxingid` int(8) DEFAULT NULL,
  `gay_time` int(11) DEFAULT '0',
  `is_pass` int(11) DEFAULT '0',
  `gay_overtime` int(11) DEFAULT NULL,
  `q_q_unionid` varchar(256) CHARACTER SET utf8 DEFAULT NULL,
  `user_info_update` int(11) DEFAULT NULL,
  `chan_type` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `job_position` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `has_col` tinyint(1) DEFAULT '0',
  `operator` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `notpass_cause` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `delete` smallint(1) unsigned DEFAULT '0' COMMENT '删除字段',
  `device_listm` text CHARACTER SET utf8,
  `version` varchar(256) CHARACTER SET utf8 DEFAULT NULL,
  `check_header_status` tinyint(3) DEFAULT '0' COMMENT '头像审核状态，状态 0-未审核，1-审核通过，2-审核不通过',
  `check_header_comment` varchar(256) CHARACTER SET utf8 DEFAULT '' COMMENT '头像审核备注，json字符串，{涉黄,涉暴}',
  `check_nickname_status` tinyint(3) DEFAULT '0' COMMENT '昵称审核状态，状态 0-未审核，1-审核通过，2-审核不通过',
  `check_nickname_comment` varchar(256) CHARACTER SET utf8 DEFAULT '' COMMENT '昵称审核备注，json字符串，{涉黄,涉暴}',
  `check_sign_status` tinyint(3) DEFAULT '0' COMMENT '签名审核状态，状态 0-未审核，1-审核通过，2-审核不通过',
  `check_sign_comment` varchar(256) CHARACTER SET utf8 DEFAULT '' COMMENT '签名审核备注，json字符串，{涉黄,涉暴}',
  `check_header_updated` int(11) DEFAULT '0' COMMENT 'ͷ??????ʱ?',
  `check_nickname_updated` int(11) DEFAULT '0' COMMENT '?ǳ?????ʱ?',
  `check_sign_updated` int(11) DEFAULT '0' COMMENT 'ǩ??????ʱ?',
  `score` int(11) NOT NULL DEFAULT '0' COMMENT '?û??ܻ??',
  `version_reg` varchar(256) DEFAULT NULL,
  `examine` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '审核(0：未操作，1：通过，2：不通过)',
  `tx_user` varchar(128) DEFAULT NULL COMMENT '腾讯用户\n',
  `tx_password` text CHARACTER SET utf8 COMMENT '腾讯用户登陆密码',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for user_antispam_log
-- ----------------------------
DROP TABLE IF EXISTS `user_antispam_log`;
CREATE TABLE `user_antispam_log` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `label` tinyint(3) NOT NULL COMMENT '违禁分类',
  `label_msg` text,
  `date` varchar(20) DEFAULT NULL COMMENT '日期',
  `score` int(11) NOT NULL DEFAULT '0' COMMENT '分数，得出风险值',
  `created` bigint(20) DEFAULT '0' COMMENT '创建时间',
  `is_reader` tinyint(3) DEFAULT '0' COMMENT '是否阅读，0-未读，1-已读',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='用户反垃圾信息表';

-- ----------------------------
-- Table structure for user_copy
-- ----------------------------
DROP TABLE IF EXISTS `user_copy`;
CREATE TABLE `user_copy` (
  `id` bigint(20) NOT NULL,
  `wx_union_id` varchar(100) DEFAULT NULL,
  `nick_name` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `phone` bigint(20) NOT NULL,
  `avator` varchar(200) DEFAULT NULL,
  `system` varchar(128) DEFAULT NULL,
  `user_type` tinyint(1) DEFAULT '0',
  `gender` tinyint(1) DEFAULT NULL,
  `gold` int(11) unsigned DEFAULT '0',
  `income_gold` int(11) unsigned DEFAULT '0',
  `hot_degree` int(11) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  `wy_user` varchar(128) DEFAULT NULL,
  `wy_password` varchar(256) DEFAULT NULL,
  `role` tinyint(1) DEFAULT '1',
  `shield_user` text,
  `call_duration` int(11) DEFAULT '0',
  `last_login_time` int(11) DEFAULT NULL,
  `last_login_history` varchar(128) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `signature` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `device_id` varchar(200) DEFAULT NULL,
  `device_token` varchar(200) DEFAULT NULL,
  `device_type` int(11) DEFAULT NULL,
  `fans_num` int(11) DEFAULT '0',
  `focus_num` int(11) NOT NULL DEFAULT '0',
  `used_duration` int(11) DEFAULT '0',
  `no_like_user` text,
  `q_q_openid` varchar(256) DEFAULT NULL,
  `wx_g_z_open_id` varchar(60) DEFAULT NULL,
  `weibo_uid` bigint(20) DEFAULT NULL,
  `membership_date` bigint(20) DEFAULT '0',
  `city_num` int(11) DEFAULT NULL,
  `first_login_city` varchar(200) DEFAULT NULL,
  `match_gender` int(11) DEFAULT '0',
  `is_accept_member` int(11) DEFAULT '0',
  `facebook_id` bigint(20) DEFAULT NULL,
  `google_id` varchar(32) DEFAULT NULL,
  `play_gold` int(11) DEFAULT '0',
  `play_be_gold` int(11) DEFAULT '0',
  `c_code` varchar(256) DEFAULT NULL,
  `is_first_recharge` tinyint(1) DEFAULT '1',
  `check_time` int(11) DEFAULT '0',
  `recharge_total` int(11) DEFAULT '0',
  `password` varchar(60) DEFAULT NULL,
  `user_extra` text,
  `v_c_code` varchar(256) DEFAULT NULL,
  `order` int(11) DEFAULT '0',
  `vc_code` varchar(256) DEFAULT NULL,
  `like_num` int(11) DEFAULT '0',
  `liked_num` int(11) DEFAULT '0',
  `province` varchar(60) DEFAULT NULL,
  `born` bigint(20) DEFAULT NULL,
  `age` bigint(20) DEFAULT NULL,
  `vc_income_gold` int(11) unsigned DEFAULT '0',
  `is_popular` int(11) DEFAULT '0',
  `last_active_time` int(11) DEFAULT NULL,
  `income_yuan` varchar(256) DEFAULT NULL,
  `constellation` varchar(128) DEFAULT NULL,
  `city` varchar(64) DEFAULT NULL,
  `location_x` varchar(128) DEFAULT NULL,
  `location_y` varchar(128) DEFAULT NULL,
  `affective_state` int(11) DEFAULT '0',
  `profession` varchar(128) DEFAULT NULL,
  `school` varchar(128) DEFAULT NULL,
  `home_town` varchar(128) DEFAULT NULL,
  `is_gay` int(11) DEFAULT '0',
  `ali_account` varchar(128) DEFAULT NULL,
  `friends_num` int(11) DEFAULT NULL,
  `ali_pay_user_id` varchar(128) DEFAULT NULL,
  `ali_pay_phone` varchar(128) DEFAULT NULL,
  `ali_pay_real_name` varchar(128) DEFAULT NULL,
  `xingxingid` int(8) DEFAULT NULL,
  `gay_time` int(11) DEFAULT '0',
  `is_pass` int(11) DEFAULT '0',
  `gay_overtime` int(11) DEFAULT NULL,
  `q_q_unionid` varchar(256) DEFAULT NULL,
  `user_info_update` int(11) DEFAULT NULL,
  `chan_type` varchar(255) DEFAULT NULL,
  `job_position` varchar(255) DEFAULT NULL,
  `has_col` tinyint(1) DEFAULT '0',
  `operator` varchar(50) DEFAULT NULL,
  `notpass_cause` varchar(255) DEFAULT NULL,
  `delete` smallint(1) unsigned DEFAULT '0' COMMENT '删除字段',
  `device_listm` text,
  `version` varchar(256) DEFAULT NULL,
  `check_header_status` tinyint(3) DEFAULT '0' COMMENT '头像审核状态，状态 0-未审核，1-审核通过，2-审核不通过',
  `check_header_comment` varchar(256) DEFAULT '' COMMENT '头像审核备注，json字符串，{涉黄,涉暴}',
  `check_nickname_status` tinyint(3) DEFAULT '0' COMMENT '昵称审核状态，状态 0-未审核，1-审核通过，2-审核不通过',
  `check_nickname_comment` varchar(256) DEFAULT '' COMMENT '昵称审核备注，json字符串，{涉黄,涉暴}',
  `check_sign_status` tinyint(3) DEFAULT '0' COMMENT '签名审核状态，状态 0-未审核，1-审核通过，2-审核不通过',
  `check_sign_comment` varchar(256) DEFAULT '' COMMENT '签名审核备注，json字符串，{涉黄,涉暴}',
  `check_header_updated` int(11) DEFAULT '0' COMMENT 'ͷ??????ʱ?',
  `check_nickname_updated` int(11) DEFAULT '0' COMMENT '?ǳ?????ʱ?',
  `check_sign_updated` int(11) DEFAULT '0' COMMENT 'ǩ??????ʱ?',
  `score` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `wx_union_id` (`wx_union_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_expansion
-- ----------------------------
DROP TABLE IF EXISTS `user_expansion`;
CREATE TABLE `user_expansion` (
  `user_id` bigint(20) NOT NULL,
  `pay_type` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '支付方式（1：微信，2：支付宝）',
  `id_number` varchar(20) DEFAULT NULL COMMENT '身份证号',
  `positive_path` varchar(255) DEFAULT NULL COMMENT '身份证正面照片',
  `back_path` varchar(255) DEFAULT NULL COMMENT '身份证反面照片',
  `address` varchar(255) DEFAULT NULL COMMENT '详细地址',
  `state` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '状态（0：关闭，1：开启）',
  `provinces` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '省id',
  `citys` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '市id',
  `charge_person` varchar(30) DEFAULT '' COMMENT '负责人',
  `full_name` varchar(30) DEFAULT '' COMMENT '姓名',
  `tel_val` varchar(25) DEFAULT '' COMMENT '联系电话',
  `examine` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '审核(0：未操作，1：通过，2：不通过)',
  `updated` int(30) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `exam_cont` mediumtext COMMENT '审核说明',
  `join_guild_time` int(30) unsigned NOT NULL DEFAULT '0' COMMENT '加入公会时间',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `user_id` (`user_id`) USING BTREE,
  KEY `provinces` (`provinces`) USING BTREE,
  KEY `city` (`citys`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='官方主播内容附表';

-- ----------------------------
-- Table structure for user_extra
-- ----------------------------
DROP TABLE IF EXISTS `user_extra`;
CREATE TABLE `user_extra` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) DEFAULT NULL,
  `url` varchar(255) DEFAULT NULL,
  `type` tinyint(1) DEFAULT '0',
  `created` int(11) DEFAULT NULL,
  `del` tinyint(1) DEFAULT '0',
  `flag` int(11) DEFAULT NULL,
  `cover_url` varchar(255) DEFAULT NULL,
  `cover_u_r_l` varchar(255) DEFAULT NULL,
  `pos` int(11) DEFAULT NULL,
  `checked_status` tinyint(3) DEFAULT '0' COMMENT '推送状态，状态 0-未审核，1-审核通过，2-审核不通过',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3067 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_label
-- ----------------------------
DROP TABLE IF EXISTS `user_label`;
CREATE TABLE `user_label` (
  `id` bigint(20) NOT NULL,
  `extend_type` tinyint(1) DEFAULT NULL,
  `extend_label_id` bigint(20) DEFAULT NULL,
  `label_name` varchar(256) DEFAULT NULL,
  `user_id` bigint(20) NOT NULL,
  `comment_total` int(11) DEFAULT '0',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `label_id` bigint(20) DEFAULT NULL,
  `cover_path` varchar(256) DEFAULT NULL,
  `pricing` int(11) DEFAULT NULL,
  `content` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `label_type` tinyint(2) DEFAULT '0',
  `is_check` tinyint(1) DEFAULT '1',
  `is_publish` tinyint(1) DEFAULT '1',
  `watch_num` int(11) DEFAULT '0',
  `star_times` int(11) DEFAULT '0',
  `play_times` int(11) DEFAULT '0',
  `op_weight` bigint(20) DEFAULT '0',
  `chat_times` bigint(20) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_mission
-- ----------------------------
DROP TABLE IF EXISTS `user_mission`;
CREATE TABLE `user_mission` (
  `id` bigint(20) NOT NULL COMMENT 'ID',
  `user_id` bigint(20) NOT NULL COMMENT 'user???û?ID',
  `code` int(6) DEFAULT '0' COMMENT '???????',
  `title` varchar(255) DEFAULT NULL,
  `date` varchar(50) DEFAULT '0000-00-00' COMMENT 'ÿ??????????: 2017-10-30',
  `type` int(6) DEFAULT '0' COMMENT '???????',
  `score` int(11) NOT NULL DEFAULT '0' COMMENT '????????',
  `times` int(11) NOT NULL DEFAULT '0' COMMENT '???ۼ????ɴ????? Ĭ?ϣ?1',
  `progress` int(11) DEFAULT NULL COMMENT '???????',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '????״̬ 1: ?????ɣ? 2: ?????ɿ???ȡ?????? 3: ????ȡ?????ر????',
  `created` int(11) DEFAULT '0' COMMENT '????ʱ?',
  `updated` int(11) DEFAULT '0' COMMENT '?޸?ʱ?',
  `property` int(6) DEFAULT '0' COMMENT '???????ʣ?ÿ?գ????'
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='?û???????¼?';

-- ----------------------------
-- Table structure for user_play_detail_data
-- ----------------------------
DROP TABLE IF EXISTS `user_play_detail_data`;
CREATE TABLE `user_play_detail_data` (
  `id` bigint(20) NOT NULL,
  `code` varchar(256) DEFAULT NULL,
  `share_record_id` bigint(20) DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `ip` varchar(256) DEFAULT NULL,
  `user_agent` varchar(256) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `reward_code` varchar(256) DEFAULT NULL,
  `active_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_voice_property
-- ----------------------------
DROP TABLE IF EXISTS `user_voice_property`;
CREATE TABLE `user_voice_property` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) DEFAULT NULL,
  `sound` text,
  `channel_code` varchar(255) DEFAULT NULL,
  `channel_name` varchar(255) DEFAULT NULL,
  `index` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_welfare
-- ----------------------------
DROP TABLE IF EXISTS `user_welfare`;
CREATE TABLE `user_welfare` (
  `id` bigint(20) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `uid` bigint(20) DEFAULT NULL,
  `start` bigint(20) DEFAULT NULL,
  `end` bigint(20) DEFAULT NULL,
  `created` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for user_welfare_talk
-- ----------------------------
DROP TABLE IF EXISTS `user_welfare_talk`;
CREATE TABLE `user_welfare_talk` (
  `id` bigint(20) NOT NULL,
  `uid` bigint(20) DEFAULT NULL,
  `created` bigint(20) DEFAULT NULL,
  `welfare_id` bigint(20) DEFAULT NULL,
  `talk_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for user_with_draw
-- ----------------------------
DROP TABLE IF EXISTS `user_with_draw`;
CREATE TABLE `user_with_draw` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `user_gold` bigint(20) DEFAULT NULL,
  `user_income_gold` bigint(20) DEFAULT NULL,
  `user_income_gold_yuan` bigint(20) DEFAULT NULL,
  `user_vc_income_gold` bigint(20) DEFAULT NULL,
  `user_vc_income_gold_yuan` bigint(20) DEFAULT NULL,
  `all_income_yuan` bigint(20) DEFAULT NULL,
  `scale` varchar(256) DEFAULT NULL,
  `banana_scale` varchar(256) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `all_income_real` int(11) DEFAULT NULL,
  `all_income_fee` int(11) DEFAULT NULL,
  `order_id` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for vc_channel
-- ----------------------------
DROP TABLE IF EXISTS `vc_channel`;
CREATE TABLE `vc_channel` (
  `id` bigint(20) NOT NULL,
  `phone` bigint(20) NOT NULL,
  `level` int(11) DEFAULT '1',
  `last_channel_code` varchar(256) DEFAULT NULL,
  `code` varchar(256) DEFAULT NULL,
  `name` varchar(256) DEFAULT NULL,
  `ali_pay_account` varchar(256) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `scale` varchar(256) DEFAULT NULL,
  `live_time` int(10) DEFAULT '20' COMMENT '明星时长',
  `user_id` bigint(20) DEFAULT NULL,
  `delete` smallint(1) DEFAULT '0' COMMENT '删除键',
  `logo_path` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '工会logo',
  `pay_type` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '支付方式',
  `charge_man` varchar(30) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '负责人',
  `id_corde` varchar(20) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '身份证号',
  `positive_path` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '身份证正面',
  `back_path` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '身份证反面',
  `license_numb` varchar(100) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '营业执照号',
  `license_path` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '营业执照照片',
  `license_level` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '评级',
  `province` double unsigned DEFAULT '1' COMMENT '省',
  `city` double unsigned DEFAULT '2' COMMENT '市',
  `addressval` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '详细地址',
  `pwd_val` varchar(50) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '渠道密码',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for vote_activity
-- ----------------------------
DROP TABLE IF EXISTS `vote_activity`;
CREATE TABLE `vote_activity` (
  `vote_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `vote_cont` varchar(255) NOT NULL DEFAULT '-1' COMMENT '活动内容',
  `vote_a` varchar(100) DEFAULT '-1' COMMENT 'a方说明',
  `vote_b` varchar(100) DEFAULT '-1' COMMENT 'b方说明',
  `vote_begin_time` int(50) unsigned NOT NULL DEFAULT '0' COMMENT '开始时间',
  `vote_end_time` int(50) unsigned NOT NULL DEFAULT '0' COMMENT '结束时间',
  `vote_onup` smallint(1) NOT NULL DEFAULT '1' COMMENT '是否上架',
  `vote_created` int(50) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `vote_update` int(50) NOT NULL DEFAULT '0' COMMENT '修改时间',
  `vote_click` smallint(5) NOT NULL DEFAULT '1' COMMENT '每天同一账户点击次数',
  `vote_del` smallint(1) NOT NULL DEFAULT '0' COMMENT '删除(1：删除，0：未删除）',
  `vote_sum` int(50) NOT NULL DEFAULT '0' COMMENT '分享出去打开次数',
  `vote_down` int(50) NOT NULL DEFAULT '0' COMMENT '分享下载次数',
  PRIMARY KEY (`vote_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COMMENT='活动投票表';

-- ----------------------------
-- Table structure for vote_activity_click
-- ----------------------------
DROP TABLE IF EXISTS `vote_activity_click`;
CREATE TABLE `vote_activity_click` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `vote_activity_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '活动id',
  `vote_activity_sum` smallint(11) NOT NULL DEFAULT '0' COMMENT '点击次数',
  `vote_activity_click_time` int(50) NOT NULL DEFAULT '0' COMMENT '点击时间',
  `vote_user_openid` varchar(50) DEFAULT '-1' COMMENT '微信开发识别符',
  `click_del` smallint(1) NOT NULL DEFAULT '0' COMMENT '删除（0：未删除，1：删除）',
  `vote_click_ip` varchar(20) DEFAULT '-1' COMMENT 'ip',
  `vote_activity_state` smallint(3) NOT NULL DEFAULT '0' COMMENT '支持状态',
  `vote_disting` varchar(100) DEFAULT '-1' COMMENT '识别服',
  `vote_disting_basctr` varchar(100) DEFAULT '-1' COMMENT '识别符',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8mb4 COMMENT='活动评论支持表';

-- ----------------------------
-- Table structure for vote_comment
-- ----------------------------
DROP TABLE IF EXISTS `vote_comment`;
CREATE TABLE `vote_comment` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `com_cont` varchar(255) DEFAULT '-1' COMMENT '评论内容',
  `vote_activity_id` int(11) NOT NULL DEFAULT '0' COMMENT '投票活动id',
  `com_time` int(50) NOT NULL DEFAULT '0' COMMENT '评论时间',
  `com_state` smallint(1) NOT NULL DEFAULT '0' COMMENT '状态（0：未审核,1：通过，2：审核不通过）',
  `com_del` smallint(1) NOT NULL DEFAULT '0' COMMENT '删除（0：未删除，1：已删除）',
  `com_uid` bigint(50) DEFAULT NULL,
  `com_type` smallint(1) NOT NULL DEFAULT '1' COMMENT '评论类型（1：文字评论，2：语音评论）',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=191 DEFAULT CHARSET=utf8mb4 COMMENT='活动投票评论表';

-- ----------------------------
-- Table structure for vote_down
-- ----------------------------
DROP TABLE IF EXISTS `vote_down`;
CREATE TABLE `vote_down` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for vote_fabulous
-- ----------------------------
DROP TABLE IF EXISTS `vote_fabulous`;
CREATE TABLE `vote_fabulous` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `com_id` int(11) NOT NULL DEFAULT '0' COMMENT '评论id',
  `user_openid` varchar(50) NOT NULL DEFAULT '-1' COMMENT '微信唯一识别符',
  `fab_time` int(50) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=718 DEFAULT CHARSET=utf8mb4 COMMENT='评论点赞表';

-- ----------------------------
-- Table structure for weibo_info
-- ----------------------------
DROP TABLE IF EXISTS `weibo_info`;
CREATE TABLE `weibo_info` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `access_token` varchar(100) DEFAULT NULL,
  `expire_in` bigint(20) DEFAULT NULL,
  `id_str` varchar(100) DEFAULT NULL,
  `screen_name` varchar(100) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `province` bigint(20) DEFAULT NULL,
  `city` bigint(20) DEFAULT NULL,
  `location` varchar(100) DEFAULT NULL,
  `description` varchar(1000) DEFAULT NULL,
  `url` varchar(100) DEFAULT NULL,
  `profile_image_url` varchar(100) DEFAULT NULL,
  `profile_url` varchar(100) DEFAULT NULL,
  `domain` varchar(100) DEFAULT NULL,
  `weihao` varchar(100) DEFAULT NULL,
  `gender` varchar(10) DEFAULT NULL,
  `followers_count` bigint(20) DEFAULT NULL,
  `friends_count` bigint(20) DEFAULT NULL,
  `statuses_count` bigint(20) DEFAULT NULL,
  `favourites_count` bigint(20) DEFAULT NULL,
  `created_at` varchar(100) DEFAULT NULL,
  `following` bigint(20) DEFAULT NULL,
  `allow_all_act_msg` bigint(20) DEFAULT NULL,
  `beo_enabled` bigint(20) DEFAULT NULL,
  `verified` bigint(20) DEFAULT NULL,
  `verified_type` bigint(20) DEFAULT NULL,
  `remark` varchar(100) DEFAULT NULL,
  `status` text,
  `allow_all_comment` bigint(20) DEFAULT NULL,
  `avatar_large` varchar(100) DEFAULT NULL,
  `avatar_hd` varchar(100) DEFAULT NULL,
  `verified_reason` varchar(100) DEFAULT NULL,
  `follow_me` bigint(20) DEFAULT NULL,
  `online_status` bigint(20) DEFAULT NULL,
  `bi_followers_count` bigint(20) DEFAULT NULL,
  `lang` varchar(100) DEFAULT NULL,
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  `geo_enabled` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for western_union
-- ----------------------------
DROP TABLE IF EXISTS `western_union`;
CREATE TABLE `western_union` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `remit_id` bigint(20) NOT NULL COMMENT '汇款ID',
  `user_id` bigint(20) NOT NULL,
  `user_with_draw_id` bigint(20) NOT NULL COMMENT '用户提现ID',
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `bank_addr` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '银行地址',
  `address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '地址',
  `money` float DEFAULT NULL,
  `country` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `service_charge` float DEFAULT NULL,
  `status` int(1) NOT NULL COMMENT '1未审核 2审核通过 3审核失败 4已汇款 5汇款失败',
  `currency` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '币种',
  `created` bigint(20) DEFAULT NULL,
  `updated` bigint(20) DEFAULT NULL,
  `mtcn` bigint(20) DEFAULT NULL COMMENT '汇款监控号',
  `remitter_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '汇款人姓名',
  `remitter_addr` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '汇款人地址',
  `remitter_country` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '汇款人姓名',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for xingxingid
-- ----------------------------
DROP TABLE IF EXISTS `xingxingid`;
CREATE TABLE `xingxingid` (
  `id` int(8) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10002555 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for zhima_certification
-- ----------------------------
DROP TABLE IF EXISTS `zhima_certification`;
CREATE TABLE `zhima_certification` (
  `id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL COMMENT '用户id',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1未认证 2已经认证',
  `id_card_code` varchar(128) DEFAULT NULL COMMENT '身份证号',
  `name` varchar(50) DEFAULT NULL COMMENT '姓名',
  `created` int(11) DEFAULT NULL,
  `updated` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
