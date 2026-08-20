# Bug Reproduction

基线分支：`base_bug_003`

问题类型：slice

复现命令：

```bash
go test ./internal/store -count=1 -run '^TestBug003_AlertsReturnsIndependentSamples$'
```

预期结果：测试失败，修改查询返回的 `SampleEventIDs` 会污染 Store 内部告警数据。

对比基线：`test_mode_bugfix_003` 保留 `base_bug_003` 的实现和测试，仅补充本复现说明。
