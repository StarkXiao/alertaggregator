# Bug Reproduction

基线分支：`base_bug_007`

问题类型：other

复现命令：

```bash
go test ./internal/api -count=1 -run '^TestBug007_AlertsLimitZeroReturnsNoItems$'
```

预期结果：测试失败，显式 `limit=0` 被错误当作缺省值并返回默认页大小的数据。

对比基线：`test_mode_bugfix_007` 保留 `base_bug_007` 的实现和测试，仅补充本复现说明。
