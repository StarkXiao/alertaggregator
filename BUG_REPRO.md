# Bug Reproduction

基线分支：`base_bug_008`

问题类型：other

复现命令：

```bash
go test ./internal/api -count=1 -run '^TestBug008_EventIDIsAssignedByServer$'
```

预期结果：测试失败，客户端提交的事件 ID 被服务端原样保存，可能干扰去重和审计。

对比基线：`test_mode_bugfix_008` 保留 `base_bug_008` 的实现和测试，仅补充本复现说明。
