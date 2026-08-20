# Bug Reproduction

基线分支：`base_bug_002`

问题类型：nil

复现命令：

```bash
go test ./internal/api -count=1 -run '^TestBug002_EventWithoutLabelsIsAccepted$'
```

预期结果：测试失败，未携带 `labels` 的合法事件触发 nil map 写入 panic。

对比基线：`test_mode_bugfix_002` 保留 `base_bug_002` 的实现和测试，仅补充本复现说明。
