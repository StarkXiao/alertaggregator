# Bug Reproduction

基线分支：`base_bug_005`

问题类型：context

复现命令：

```bash
go test ./internal/worker -count=1 -run '^TestBug005_CancelledWorkerDoesNotProcessEvents$'
```

预期结果：测试失败，已取消的 Worker 仍执行首次 tick 并创建告警。

对比基线：`test_mode_bugfix_005` 保留 `base_bug_005` 的实现和测试，仅补充本复现说明。
