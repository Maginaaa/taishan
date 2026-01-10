# Requirements Document

## Introduction

本功能为压测平台添加测试报告对比功能。当压测结束（status=4）时，用户可以选择当前测试计划的历史压测报告，与当前报告进行性能指标对比分析，帮助用户快速识别性能变化趋势。

## Glossary

- **Report**: 压测报告，包含压测执行的完整结果数据
- **Plan**: 测试计划，包含多个场景和用例的压测配置
- **Report_Comparison_System**: 报告对比系统，负责获取、对比和展示报告数据
- **Comparison_Metrics**: 对比指标，包括RPS、RT、成功率等核心性能指标
- **Historical_Report**: 历史报告，同一测试计划下已完成的其他压测报告

## Requirements

### Requirement 1: 历史报告选择

**User Story:** As a 测试工程师, I want to 在压测结束后选择历史报告进行对比, so that 我可以分析不同压测之间的性能差异。

#### Acceptance Criteria

1. WHEN press_status equals 4 (压测结束), THE Report_Comparison_System SHALL display a comparison button in the report detail page
2. WHEN user clicks the comparison button, THE Report_Comparison_System SHALL open a dialog showing historical reports from the same plan
3. WHEN displaying historical reports, THE Report_Comparison_System SHALL show report_id, start_time, end_time, and create_user_name for each report
4. WHEN user selects a historical report, THE Report_Comparison_System SHALL enable the confirm button
5. THE Report_Comparison_System SHALL exclude the current report from the historical report list

### Requirement 2: 报告对比展示

**User Story:** As a 测试工程师, I want to 查看两份报告的对比结果, so that 我可以直观地了解性能变化。

#### Acceptance Criteria

1. WHEN user confirms the comparison selection, THE Report_Comparison_System SHALL display a comparison view with both reports' data
2. THE Report_Comparison_System SHALL display comparison metrics including: total_request_num, total_success_rate, total_rps, and total_avg_rt
3. WHEN displaying comparison metrics, THE Report_Comparison_System SHALL show the difference value and percentage change between two reports
4. WHEN a metric improves (higher RPS, lower RT, higher success rate), THE Report_Comparison_System SHALL highlight it in green
5. WHEN a metric degrades (lower RPS, higher RT, lower success rate), THE Report_Comparison_System SHALL highlight it in red

### Requirement 3: 场景级别对比

**User Story:** As a 测试工程师, I want to 对比各场景的性能指标, so that 我可以定位具体哪个场景的性能发生了变化。

#### Acceptance Criteria

1. THE Report_Comparison_System SHALL display scene-level comparison for matching scenes between two reports
2. WHEN scenes exist in both reports, THE Report_Comparison_System SHALL compare their RPS, RT, and success_rate
3. WHEN a scene exists only in one report, THE Report_Comparison_System SHALL mark it as "新增" or "已移除"
4. THE Report_Comparison_System SHALL display scene comparison in a table format with clear visual distinction

### Requirement 4: 接口级别对比

**User Story:** As a 测试工程师, I want to 对比各接口的性能指标, so that 我可以精确定位性能变化的接口。

#### Acceptance Criteria

1. THE Report_Comparison_System SHALL display case-level comparison within each scene
2. WHEN cases exist in both reports, THE Report_Comparison_System SHALL compare their total_rps, total_avg_rt, total_success_rate, and percentile values (90/95/99线)
3. WHEN a case exists only in one report, THE Report_Comparison_System SHALL mark it as "新增" or "已移除"
4. THE Report_Comparison_System SHALL allow expanding/collapsing scene sections to view case details

### Requirement 5: 对比结果关闭

**User Story:** As a 测试工程师, I want to 关闭对比视图返回原报告, so that 我可以继续查看当前报告详情。

#### Acceptance Criteria

1. WHEN comparison view is displayed, THE Report_Comparison_System SHALL provide a close button
2. WHEN user clicks the close button, THE Report_Comparison_System SHALL hide the comparison view and restore the original report view
3. THE Report_Comparison_System SHALL preserve the original report data when switching between views
