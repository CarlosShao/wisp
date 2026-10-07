# 35-r4 起手锚（票 35 `:75` 的 (a) 支：夹具两面盲区修成会咬人的）

现量时刻：2026-10-07 22:14:16+0800

```
r4 anchor 2026-10-07 22:14:00+0800
HEAD=6c7dd5ed348b8f070da1b344c54359b97509132b
--- forbidden zone diff 1d70fb2b..HEAD ---
rc=0
--- working tree vs HEAD on prod faces ---
rc=0
--- blob hashes ---
7185ab56
bebe8e70
--- CR counts (autocrlf phantom check) ---
cmd/wisp/panel_transport_35r2_test.go wt=0 head=0
cmd/wisp/panel_inbound_guards_35r3_test.go wt=0 head=0
--- r2 fixture blind-spot baseline ---
0
rc=1
```

读法：禁区三方差 1d70fb2b..HEAD 为空（rc=0）；工作树对 HEAD 在生产面/两枚禁区测试面上为空（rc=0）；
r2 夹具 blob 7185ab56（与 A686 记的同名）＝本程要改的对象；bridge.go blob bebe8e70 未变；
CR 幻影尺：两枚测试件工作树与 HEAD 同为 0，一致；
夹具 writable|configurable|defineProperty 命中 0 处（grep -c rc=1＝零命中，A684 那两枚恒真的根因）。
注意：起手时 HEAD 已从派单里的 af68866c 前进到 6c7dd5ed（共享工作树，别人已落笔），禁区差在新 HEAD 上仍为空。
rc=0
