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

 Date: 17/05/2025 17:35:20
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for debug_record
-- ----------------------------
DROP TABLE IF EXISTS `debug_record`;
CREATE TABLE `debug_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '变量id',
  `plan_id` int(11) NOT NULL COMMENT '计划id',
  `status` tinyint(1) NOT NULL COMMENT '调试结果',
  `result_info` json NOT NULL COMMENT 'Case调试详情',
  `create_user_id` int(11) NOT NULL COMMENT '创建人id',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `is_delete` tinyint(1) NOT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2509 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
