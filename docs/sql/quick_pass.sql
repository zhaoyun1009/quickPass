/*
 Navicat Premium Data Transfer

 Source Server         : localhost_3306
 Source Server Type    : MySQL
 Source Server Version : 50724
 Source Host           : localhost:3306
 Source Schema         : quick_pass

 Target Server Type    : MySQL
 Target Server Version : 50724
 File Encoding         : 65001

 Date: 19/06/2020 09:17:29
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE `quick_pass`;
USE `quick_pass`;

-- ----------------------------
-- Table structure for abnormal_order
-- ----------------------------
DROP TABLE IF EXISTS `abnormal_order`;
CREATE TABLE `abnormal_order` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `agency` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '代理',
  `username` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '用户名',
  `nickname` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '昵称',
  `channel` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '收款通道({“BANK_CARD":"银行卡","ALIPAY":"支付宝","WECHAT":"微信"})',
  `original_order_no` varchar(23) COLLATE utf8mb4_bin NOT NULL COMMENT '原始订单号',
  `abnormal_order_no` varchar(23) COLLATE utf8mb4_bin NOT NULL DEFAULT '0' COMMENT '异常单号',
  `abnormal_order_type` int(4) NOT NULL DEFAULT '0' COMMENT '异常类型（1: 未知, 2：超时取消，3：未生成订单， 4：订单金额不符）',
  `abnormal_order_status` int(8) NOT NULL COMMENT '订单状态（1：未处理  2：已处理  3：处理中）',
  `accept_card_account` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '收款账户',
  `accept_card_no` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '收款卡号',
  `accept_card_bank` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '银行名称',
  `amount` bigint(20) NOT NULL COMMENT '金额',
  `append_info` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '摘要',
  `payment_date` datetime DEFAULT NULL COMMENT '付款日期',
  `finish_time` datetime DEFAULT NULL COMMENT '完成时间',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_abnormal_order_no` (`abnormal_order_no`) USING BTREE,
  KEY `idx_order_no` (`original_order_no`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for acceptor
-- ----------------------------
DROP TABLE IF EXISTS `acceptor`;
CREATE TABLE `acceptor` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `acceptor` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '承兑人账号',
  `agency` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '所属代理',
  `if_auto_accept` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否自动承兑（1：不自动，2：自动）',
  `accept_switch` tinyint(1) NOT NULL DEFAULT '1' COMMENT '承兑开关(1：关闭，2：开启)',
  `accept_status` tinyint(1) NOT NULL DEFAULT '2' COMMENT '承兑账户状态(1：关闭，2：开启)',
  `deposit` int(64) NOT NULL DEFAULT '0' COMMENT '保证金',
  `max_accepted_amount` bigint(20) NOT NULL DEFAULT '500000000' COMMENT '最大承兑金额',
  `min_accepted_amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '最小承兑金额',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_agency_acceptor` (`agency`,`acceptor`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for acceptor_card
-- ----------------------------
DROP TABLE IF EXISTS `acceptor_card`;
CREATE TABLE `acceptor_card` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `agency` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '代理',
  `if_open` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否开启={1:关闭，2：开启}',
  `acceptor` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '承兑人账号',
  `card_type` varchar(32) COLLATE utf8mb4_bin NOT NULL COMMENT '卡类型={“BANK_CARD":"银行卡","ALIPAY":"支付宝","WECHAT":"微信"}',
  `card_no` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '卡号',
  `card_account` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '卡账户名',
  `card_bank` varchar(64) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '银行名称',
  `card_sub_bank` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '银行支行',
  `card_img` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '图片地址',
  `day_available_amt` bigint(20) NOT NULL COMMENT '当日剩余可用流水',
  `day_frozen_amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '单日冻结流水',
  `day_max_amt` bigint(20) NOT NULL DEFAULT '500000000' COMMENT '当日最大流水',
  `delete_flag` tinyint(1) NOT NULL DEFAULT '1' COMMENT '删除标记（1：未删除，2：已删除）',
  `last_match_time` datetime DEFAULT NULL COMMENT '上一次的匹配时间',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_agency_acceptor_card` (`agency`,`acceptor`,`card_type`,`card_no`,`delete_flag`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for bill
-- ----------------------------
DROP TABLE IF EXISTS `bill`;
CREATE TABLE `bill` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `bill_no` varchar(23) COLLATE utf8mb4_bin NOT NULL COMMENT '账单号',
  `agency` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '代理',
  `own_user_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '已方交易账号',
  `own_role` smallint(8) NOT NULL COMMENT '角色（1:代理，2：承兑人、3：商家）',
  `opposite_user_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '对方交易账号',
  `opposite_role` smallint(8) NOT NULL COMMENT '对方角色（1:代理，2：承兑人、3：商家）',
  `amount` bigint(20) NOT NULL COMMENT '涉及金额',
  `usable_amount` bigint(20) NOT NULL COMMENT '当前可用金额',
  `frozen_amount` bigint(20) NOT NULL COMMENT '当前冻结金额',
  `income_expenses_type` tinyint(4) NOT NULL COMMENT '收支类型（1：支出，2：收入）',
  `business_type` int(8) NOT NULL COMMENT '会计科目（1：转账，2：买入手续费，3：买入，4：卖出）',
  `append_info` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '摘要信息',
  `current_rate` varchar(64) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '当前汇率',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_bill_no` (`bill_no`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for bill_statistics
-- ----------------------------
DROP TABLE IF EXISTS `bill_statistics`;
CREATE TABLE `bill_statistics` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `agency` varchar(255) NOT NULL COMMENT '代理',
  `username` varchar(255) NOT NULL COMMENT '用户名',
  `role` tinyint(1) NOT NULL COMMENT '角色',
  `statistics_date` varchar(255) NOT NULL DEFAULT '' COMMENT '统计日期',
  `left_amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '剩余金额',
  `accept_amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '承兑金额',
  `accept_count` int(11) NOT NULL DEFAULT '0' COMMENT '承兑次数',
  `withdrawal_amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '提现金额',
  `withdrawal_count` int(11) NOT NULL DEFAULT '0' COMMENT '提现次数',
  `recharge_amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '充值金额',
  `recharge_count` int(11) NOT NULL DEFAULT '0' COMMENT '充值次数',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_date_idx` (`agency`,`username`,`statistics_date`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for channel
-- ----------------------------
DROP TABLE IF EXISTS `channel`;
CREATE TABLE `channel` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `agency` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '代理',
  `channel` varchar(32) COLLATE utf8mb4_bin NOT NULL COMMENT '通道名称',
  `if_open` tinyint(4) NOT NULL COMMENT '通道开关(1: 关闭  2:开启)',
  `rate` int(64) NOT NULL COMMENT '通道汇率',
  `limit_amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '单账户每日承兑上限',
  `append_info` varchar(32) COLLATE utf8mb4_bin NOT NULL COMMENT '通道补充信息',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_channel` (`channel`,`agency`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for fund
-- ----------------------------
DROP TABLE IF EXISTS `fund`;
CREATE TABLE `fund` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `agency` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '代理',
  `user_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '用户名',
  `type` tinyint(1) NOT NULL DEFAULT '2' COMMENT '账户类型(1：系统资金账户  2：普通资金账户)',
  `available_amount` bigint(20) NOT NULL COMMENT '可用资金',
  `frozen_amount` bigint(20) NOT NULL COMMENT '冻结资金',
  `version` bigint(20) NOT NULL DEFAULT '0' COMMENT '版本号',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_agency_user` (`agency`,`user_name`) USING BTREE,
  KEY `agencyUserRefer` (`agency`,`user_name`,`type`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of fund
-- ----------------------------
BEGIN;
INSERT INTO `fund` VALUES (1, '', 'admin', 1, 100000000000000, 0, 0, '2020-06-19 09:16:31', '2020-06-19 09:16:38');
COMMIT;

-- ----------------------------
-- Table structure for management
-- ----------------------------
DROP TABLE IF EXISTS `management`;
CREATE TABLE `management` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `agency` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '代理名称',
  `user_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '后台账号',
  `password` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '密码',
  `full_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '姓名',
  `role` tinyint(8) NOT NULL COMMENT '1：管理员，2：客服，3：财务',
  `rules` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '权限列表(如：order,user,merchant,acceptor)',
  `phone_number` varchar(16) COLLATE utf8mb4_bin NOT NULL COMMENT '电话号码',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_agency_user` (`agency`,`user_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for match_cache
-- ----------------------------
DROP TABLE IF EXISTS `match_cache`;
CREATE TABLE `match_cache` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `agency` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '代理',
  `acceptor` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '承兑人账号',
  `card_id` bigint(20) NOT NULL COMMENT '卡id',
  `card_type` varchar(32) COLLATE utf8mb4_bin NOT NULL COMMENT '卡类型',
  `max_matched_amount` bigint(20) NOT NULL COMMENT '最大承兑金额',
  `min_matched_amount` bigint(20) NOT NULL COMMENT '最小承兑金额',
  `last_match_time` timestamp(6) NULL DEFAULT NULL COMMENT '上一次的匹配时间',
  `fund_version` bigint(20) NOT NULL COMMENT '资金版本号',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_card` (`card_id`) USING BTREE,
  KEY `idx_agency_match` (`agency`,`card_type`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for merchant
-- ----------------------------
DROP TABLE IF EXISTS `merchant`;
CREATE TABLE `merchant` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `merchant_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '商家账号',
  `agency` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '所属代理',
  `merchant_private_key` varchar(2048) COLLATE utf8mb4_bin NOT NULL COMMENT '商家私钥',
  `merchant_public_key` varchar(2048) COLLATE utf8mb4_bin NOT NULL COMMENT '商家公钥',
  `system_private_key` varchar(2048) COLLATE utf8mb4_bin NOT NULL COMMENT '系统私钥',
  `system_public_key` varchar(2048) COLLATE utf8mb4_bin NOT NULL COMMENT '系统公钥',
  `return_url` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '商家接口返回url',
  `notify_url` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '商家接口回调url',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_agency_merchant` (`agency`,`merchant_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for merchant_card
-- ----------------------------
DROP TABLE IF EXISTS `merchant_card`;
CREATE TABLE `merchant_card` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `agency` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '代理',
  `merchant` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '商家账号',
  `card_type` varchar(32) COLLATE utf8mb4_bin NOT NULL COMMENT '卡类型={“BANK_CARD":"银行卡","ALIPAY":"支付宝","WECHAT":"微信"}',
  `card_no` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '卡号',
  `card_account` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '卡账户名',
  `card_bank` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '银行名称',
  `card_sub_bank` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '银行支行',
  `card_img` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '图片地址',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_agency_acceptor_card` (`agency`,`merchant`,`card_type`,`card_no`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for merchant_order
-- ----------------------------
CREATE TABLE `merchant_order` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `agency` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '代理',
  `username` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '商家账户名',
  `order_no` varchar(23) COLLATE utf8mb4_bin NOT NULL COMMENT '系统订单号',
  `merchant_order_no` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '商家平台传过来的订单号',
  `submit_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '下单方式（1、接口，2、非接口）',
  `return_url` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '支付成功的返回URL',
  `callback_times` int(11) NOT NULL DEFAULT '0',
  `callback_url` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '回调接口地址',
  `callback_status` tinyint(1) NOT NULL COMMENT '回调结果（1、未处理，2、回调成功，3、回调失败）',
  `callback_info` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '回调失败的信息',
  `callback_lock` tinyint(1) NOT NULL DEFAULT '1' COMMENT '回调结果（1、未请求，2、请求中）',
  `append_info` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '备注信息',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_merchant_orderId` (`agency`,`username`,`order_no`) USING BTREE,
  UNIQUE KEY `idx_merchant_merchantOrderId` (`agency`,`username`,`merchant_order_no`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for order
-- ----------------------------
DROP TABLE IF EXISTS `order`;
CREATE TABLE `order` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `agency` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '代理',
  `order_no` varchar(23) COLLATE utf8mb4_bin NOT NULL COMMENT '订单号',
  `order_type` int(4) NOT NULL COMMENT '订单类型（1：转账，2：买入，3：卖出）',
  `order_status` int(8) NOT NULL COMMENT '订单状态（1：创建，2：待支付，3：待放行，4：已取消  5：已失败  6：已完成）',
  `abnormal_order_no` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '异常单号',
  `from_user_name` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '发起者账号',
  `to_user_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '接受者账号',
  `merchant_user_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '商家账号',
  `acceptor_nickname` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '订单所匹配的承兑人昵称',
  `card_no` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '卡号',
  `card_account` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '卡账户名',
  `card_bank` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '银行名称',
  `card_sub_bank` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '银行支行',
  `card_img` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '图片地址',
  `amount` bigint(20) NOT NULL COMMENT '金额',
  `submit_type` tinyint(4) DEFAULT NULL COMMENT '提交类型=（1：直接提交，2：接口提交）',
  `channel_type` varchar(32) COLLATE utf8mb4_bin NOT NULL COMMENT '通道类型= (BANK_CARD,ALIPAY,WECHAT)',
  `append_info` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '摘要',
  `client_ip` varchar(32) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'ip地址',
  `finish_time` datetime DEFAULT NULL COMMENT '订单完成时间',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_order_no` (`order_no`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_name` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '用户名',
  `password` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '登录密码',
  `trade_key` varchar(64) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '交易密码',
  `agency` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '所属代理',
  `full_name` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '姓名',
  `phone_number` varchar(16) COLLATE utf8mb4_bin NOT NULL COMMENT '电话号码',
  `address` varchar(1024) COLLATE utf8mb4_bin NOT NULL COMMENT '地址',
  `type` tinyint(1) NOT NULL DEFAULT '2' COMMENT '账户类型(1:系统账户 2:普通账户)',
  `role` smallint(8) NOT NULL COMMENT '角色（1:代理，2：承兑人、3：商家）',
  `secret_key` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'token盐值',
  `status` tinyint(4) NOT NULL DEFAULT '2' COMMENT '账户状态（1：暂停，2：启用）',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_agency_user` (`agency`,`user_name`) USING BTREE,
  UNIQUE KEY `idx_agency_full_name` (`agency`,`full_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of user
-- ----------------------------
-- password:13686229817   trade_key:232511
BEGIN;
INSERT INTO `user` VALUES (1, 'admin', 'bb264f4fde074348c6a9f93fd3fd9f0a', 'e53073d0f7f36959ad303d9b9045e881', '', '系统账户', '18817936112', '阿拉伯联合酋长国', 1, 4, '', 2, '2020-06-19 09:15:37', '2020-06-19 09:15:42');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
