# Bug Reproduction

基线分支：`base_bug_009`

问题类型：error

复现命令：

```bash
go test ./internal/api -count=1 -run '^TestBug009_MissingAlertAckReturnsNotFound$'
```

预期结果：测试失败，不存在告警的确认请求因错误链断裂返回 500，而不是 404。

对比基线：`test_mode_bugfix_009` 保留 `base_bug_009` 的实现和测试，仅补充本复现说明。
