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

 Date: 17/05/2025 17:36:32
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for plan
-- ----------------------------
DROP TABLE IF EXISTS `plan`;
CREATE TABLE `plan` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '测试计划id',
  `plan_name` varchar(255) CHARACTER SET utf8 NOT NULL COMMENT '测试计划名称',
  `engine_count` int(4) NOT NULL DEFAULT '1' COMMENT '压测机数量',
  `press_info` json NOT NULL COMMENT '施压策略',
  `global_variable` json NOT NULL COMMENT '全局变量',
  `default_header` json NOT NULL COMMENT '默认请求头',
  `sampling_info` json NOT NULL COMMENT '采样策略',
  `server_info` json NOT NULL COMMENT '管理服务',
  `break_type` int(11) NOT NULL COMMENT '熔断机制类型',
  `break_value` float(8,2) NOT NULL COMMENT '熔断阈值',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '定时任务id',
  `tag` json DEFAULT NULL,
  `remark` text CHARACTER SET utf8 COMMENT '备注',
  `debug_status` tinyint(1) NOT NULL COMMENT '是否调试通过',
  `create_user_id` int(11) NOT NULL COMMENT '创建人id',
  `update_user_id` int(11) NOT NULL COMMENT '最后修改人id',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `is_delete` tinyint(1) NOT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=526 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
