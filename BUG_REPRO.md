# Bug Reproduction

基线分支：`base_bug_004`

问题类型：error

复现命令：

```bash
go test ./internal/validation -count=1 -run '^TestBug004_InvalidLevelPreservesSentinel$'
```

预期结果：测试失败，非法日志级别错误包含哨兵文本但无法通过 `errors.Is` 识别。

对比基线：`test_mode_bugfix_004` 保留 `base_bug_004` 的实现和测试，仅补充本复现说明。
