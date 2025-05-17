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

 Date: 17/05/2025 17:36:39
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for report
-- ----------------------------
DROP TABLE IF EXISTS `report`;
CREATE TABLE `report` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '测试报告id',
  `report_name` varchar(255) NOT NULL COMMENT '测试报告名',
  `plan_id` int(11) NOT NULL COMMENT '测试计划id',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '测试报告状态',
  `duration` int(32) NOT NULL COMMENT '预计持续时间',
  `press_type` int(1) NOT NULL COMMENT '压测模式',
  `engine_list` json NOT NULL COMMENT '施压机ip列表',
  `concurrency` int(32) NOT NULL COMMENT '起始并发数',
  `vum` int(64) DEFAULT NULL COMMENT 'Vum',
  `create_user_id` int(11) NOT NULL COMMENT '创建人id',
  `update_user_id` int(11) NOT NULL COMMENT '最后修改人id',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `is_delete` tinyint(1) NOT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6066 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
