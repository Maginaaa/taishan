/*
 Navicat Premium Data Transfer

 Source Server         : taishan-qa-mysql
 Source Server Type    : MySQL
 Source Server Version : 50743
 Source Host           : 10.72.248.28:3306
 Source Schema         : press_dev

 Target Server Type    : MySQL
 Target Server Version : 50743
 File Encoding         : 65001

 Date: 17/05/2025 17:36:12
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for operation_log
-- ----------------------------
DROP TABLE IF EXISTS `operation_log`;
CREATE TABLE `operation_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `source_name` varchar(300) CHARACTER SET utf8 NOT NULL COMMENT '调用方系统或模块的名字',
  `source_id` int(16) NOT NULL COMMENT '操作的对象ID',
  `operation_type` int(8) NOT NULL COMMENT '操作别名',
  `operator_id` int(16) NOT NULL COMMENT '操作人',
  `value_before` json NOT NULL COMMENT '变更前',
  `value_after` json NOT NULL COMMENT '变更后',
  `value_diff` json NOT NULL COMMENT '差异变更',
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `record_source_id` (`source_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1089 DEFAULT CHARSET=utf8mb4 COMMENT='操作记录表';

SET FOREIGN_KEY_CHECKS = 1;
