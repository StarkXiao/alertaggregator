# Bug Reproduction

基线分支：`base_bug_006`

问题类型：slice

复现命令：

```bash
go test ./internal/store -count=1 -run '^TestBug006_EventsReturnsIndependentLabels$'
```

预期结果：测试失败，修改查询返回的事件标签 map 会污染 Store 内部事件。

对比基线：`test_mode_bugfix_006` 保留 `base_bug_006` 的实现和测试，仅补充本复现说明。
