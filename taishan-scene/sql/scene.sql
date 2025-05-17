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

 Date: 17/05/2025 17:37:05
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for scene
-- ----------------------------
DROP TABLE IF EXISTS `scene`;
CREATE TABLE `scene` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '场景id',
  `plan_id` int(11) NOT NULL COMMENT '场景id',
  `sort` int(11) NOT NULL DEFAULT '1' COMMENT '排序',
  `scene_type` int(8) NOT NULL COMMENT '场景类型',
  `scene_name` varchar(255) CHARACTER SET utf8 NOT NULL COMMENT '场景名',
  `export_info` json NOT NULL COMMENT '导出数据信息',
  `disabled` tinyint(1) NOT NULL COMMENT '是否禁用',
  `create_user_id` int(11) NOT NULL COMMENT '创建人id',
  `update_user_id` int(11) NOT NULL COMMENT '最后修改人id',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2352 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
