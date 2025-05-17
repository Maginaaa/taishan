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

 Date: 17/05/2025 17:37:15
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for scene_case
-- ----------------------------
DROP TABLE IF EXISTS `scene_case`;
CREATE TABLE `scene_case` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Case ID',
  `parent_id` int(11) NOT NULL COMMENT '父节点id',
  `type` int(11) NOT NULL COMMENT 'case类型',
  `scene_id` int(11) NOT NULL COMMENT '所属场景id',
  `sort` int(11) NOT NULL COMMENT '排序',
  `disabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否禁用',
  `extend` json NOT NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3904 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
