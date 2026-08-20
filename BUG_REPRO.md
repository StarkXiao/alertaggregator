# Bug Reproduction

基线分支：`base_bug_001`

问题类型：concurrency

复现命令：

```bash
go test -race ./internal/store -count=1 -run '^TestBug001_ConcurrentEventsReadIsSynchronized$'
```

预期结果：测试失败并报告 `Store.Events` 与并发更新之间存在数据竞争。

对比基线：`test_mode_bugfix_001` 保留 `base_bug_001` 的实现和测试，仅补充本复现说明。
