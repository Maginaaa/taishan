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

 Date: 17/05/2025 17:36:23
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for parameter_files
-- ----------------------------
DROP TABLE IF EXISTS `parameter_files`;
CREATE TABLE `parameter_files` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `plan_id` int(11) NOT NULL DEFAULT '0' COMMENT '测试计划ID',
  `file_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件名',
  `size` int(11) NOT NULL COMMENT '文件大小',
  `rows` int(64) NOT NULL DEFAULT '0' COMMENT '文件行数',
  `column` json NOT NULL COMMENT '文件列信息',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用',
  `create_user_id` int(11) NOT NULL COMMENT '创建人id',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_user_id` int(11) DEFAULT NULL COMMENT '最后修改人id',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `is_delete` tinyint(1) unsigned zerofill NOT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=975 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
