# Implementation Plan: Report Comparison Feature

## Overview

本实现计划将报告对比功能分解为可执行的编码任务。采用自底向上的方式，先实现工具函数和基础组件，再组装成完整功能。

## Tasks

- [x] 1. 创建数据类型定义和工具函数
  - [x] 1.1 创建报告对比相关的 TypeScript 类型定义
    - 在 `taishan-fe/src/views/Report/report-data.ts` 中添加对比相关类型
    - 包括 MetricComparison, SceneComparisonItem, CaseComparisonItem 等
    - _Requirements: 2.2, 2.3, 3.1, 4.1_

  - [x] 1.2 创建对比计算工具函数
    - 创建 `taishan-fe/src/views/Report/utils/comparison.ts`
    - 实现 calculateMetricComparison 函数
    - 实现 matchScenes 和 matchCases 函数
    - 实现 filterHistoricalReports 函数
    - _Requirements: 1.5, 2.3, 3.1, 3.3, 4.1, 4.3_

  - [ ]* 1.3 编写工具函数的属性测试
    - **Property 1: Historical Report Exclusion**
    - **Property 2: Metric Comparison Calculation**
    - **Property 3: Metric Improvement Detection**
    - **Property 4: Scene/Case Matching Completeness**
    - **Validates: Requirements 1.5, 2.3, 2.4, 2.5, 3.1, 3.3, 4.1, 4.3**

- [x] 2. 扩展 Report API 服务
  - [x] 2.1 添加获取历史报告列表的 API 方法
    - 在 `taishan-fe/src/api/scene/report/index.ts` 中添加 getHistoricalReports 方法
    - 参数: plan_id, exclude_report_id
    - _Requirements: 1.2, 1.3, 1.5_

- [x] 3. 实现基础 UI 组件
  - [x] 3.1 创建 ComparisonMetricCard 组件
    - 创建 `taishan-fe/src/views/Report/component/ComparisonMetricCard.vue`
    - 显示当前值、对比值、差异值和百分比
    - 根据 higherIsBetter 属性决定颜色（绿色=改进，红色=退化）
    - _Requirements: 2.2, 2.3, 2.4, 2.5_

  - [x] 3.2 创建 HistoricalReportDialog 组件
    - 创建 `taishan-fe/src/views/Report/component/HistoricalReportDialog.vue`
    - 显示历史报告列表（report_id, start_time, end_time, create_user_name）
    - 支持单选报告
    - 确认/取消按钮
    - _Requirements: 1.2, 1.3, 1.4_

- [x] 4. Checkpoint - 确保基础组件正常工作
  - 确保所有测试通过，如有问题请询问用户

- [x] 5. 实现对比展示组件
  - [x] 5.1 创建 SceneComparisonTable 组件
    - 创建 `taishan-fe/src/views/Report/component/SceneComparisonTable.vue`
    - 显示场景级别对比（RPS, RT, 成功率）
    - 标记新增/已移除的场景
    - 支持展开/折叠查看接口详情
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 4.1, 4.2, 4.3, 4.4_

  - [x] 5.2 创建 ReportComparisonView 组件
    - 创建 `taishan-fe/src/views/Report/component/ReportComparisonView.vue`
    - 顶部显示两份报告的基本信息
    - 使用 ComparisonMetricCard 展示总体指标对比
    - 使用 SceneComparisonTable 展示场景/接口对比
    - 提供关闭按钮
    - _Requirements: 2.1, 2.2, 5.1, 5.2_

- [x] 6. 集成到 ReportDetail 页面
  - [x] 6.1 在 ReportDetail.vue 中添加对比功能入口
    - 在 press_status === 4 的区域添加"报告对比"按钮
    - 引入 HistoricalReportDialog 组件
    - 引入 ReportComparisonView 组件
    - 管理对比状态（是否显示对比视图）
    - _Requirements: 1.1, 5.2, 5.3_

  - [x] 6.2 实现对比流程逻辑
    - 点击对比按钮 → 打开历史报告选择对话框
    - 选择报告并确认 → 获取目标报告数据
    - 显示对比视图
    - 关闭对比视图 → 恢复原报告视图
    - _Requirements: 1.1, 1.2, 1.4, 2.1, 5.1, 5.2, 5.3_

- [ ] 7. Final Checkpoint - 完整功能验证
  - 确保所有测试通过，如有问题请询问用户
  - 验证完整的对比流程

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties
- 使用 TypeScript 进行类型安全的开发
- 使用 fast-check 进行属性测试
