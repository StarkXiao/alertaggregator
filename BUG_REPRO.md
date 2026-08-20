# Bug Reproduction

基线分支：`base_bug_010`

问题类型：other

复现命令：

```bash
go test ./internal/notify -count=1 -run '^TestBug010_FirstNotificationUsesInitialBackoff$'
```

预期结果：测试失败，首次通知计算出的下一次提醒时间跳过了初始退避级别。

对比基线：`test_mode_bugfix_010` 保留 `base_bug_010` 的实现和测试，仅补充本复现说明。
