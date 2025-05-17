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

 Date: 17/05/2025 17:37:23
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for scene_variable
-- ----------------------------
DROP TABLE IF EXISTS `scene_variable`;
CREATE TABLE `scene_variable` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '变量id',
  `scene_id` int(11) NOT NULL COMMENT '所属场景id',
  `variable_name` varchar(64) NOT NULL COMMENT '变量名',
  `variable_val` varchar(2048) NOT NULL COMMENT '变量值',
  `remark` varchar(255) DEFAULT NULL COMMENT '描述',
  `create_user_id` int(11) NOT NULL COMMENT '创建人id',
  `update_user_id` int(11) NOT NULL COMMENT '修改人id',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `is_delete` tinyint(1) NOT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=152 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
