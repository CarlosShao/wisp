# 票 124 落地腿 124-r1 交件 part3（sec 7..10：残余红归因、门禁、交回项）

> 前接 `docs/evidence/s1/124-r1-impl.md`（sec 0..5）与 `124-r1-impl-part2.md`（sec 6）。

## 7. 残余 2 枚红：改前就在、两形都在、与形状无关（既有 flake，非本腿引入）

`RESTORED-LINK` 余 2 名：`TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent`、
`TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed`。
两枚在**改前普通形** `PRE-PLAIN` 也已红（`PRE-PLAIN FAIL=2` 正是这两名）=>  不属 only-in-link 那一档。

第二枚我单发 `-count=5` 在**纯净 HEAD 快照**上量到它是**摇摆**的（2 PASS / 3 FAIL），
与我的改动无关，归 **既有 flake**、**不在本票清零目标内**，本腿未碰其断言一个字。

## 8. 门禁（我自己现跑，逐条贴末行）

```
$ GOFLAGS= go build ./...            -> rc=0
$ go vet ./internal/config/          -> rc=0
$ go test -count=1 ./internal/config/-> rc=0   ok github.com/CarlosShao/wisp/internal/config 1.952s
      RUN=195  PASS=195  FAIL=0  SKIP=0      （四数分开量，未相减）
$ sh scripts/d22scan.sh              -> rc=0
      d22scan: clean - no D22 ban violations; examined 286 production Go files; ban #8 internal/=608 cmd/=143
      runtests.sh: OK - top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77
$ sh scripts/check-path-length-budget.sh -> rc=0
      VERDICT GREEN - every over-budget tracked path is rostered by name with a reason
      denominator: tracked paths=17488  over-budget=60  covered by roster=60  not in roster=0
$ gofmt -l <4 touched files>         -> 0 行 (rc=0)
$ gofmt -l internal/config           -> 0 行 (rc=0)
$ gofumpt v0.12.0 (go1.27.1) -l <4 touched files> -> 0 行 (rc=0)
```

`gofumpt` 版本按本机实测报 **v0.12.0**（`ci.yml` 是 `@latest` 未钉版本，前几批各报自己手上那枚）。

## 9. 红名册（only-in-link 清零，逐名对账）

```
内部对照（都在软链形，同一包）：
  only-in-link 改前 = 33 名
  only-in-link 改后 =  0 名        <- 差集为空
  改后残余 2 名 = 改前普通形也红的既有 flake（sec 7 已逐名归因）
新增红 = 0
```

## 10. 未验证 / 交回编排者

1. **本腿只做 `internal/config` 一包**。`internal/tools`（38 枚 / 16 文件）与 `cmd/wisp`（37 枚 / 13 文件）
   的逐名清单已量出但**未改**，请排后续批；`cmd/wisp` 是共享编译面，按 AGENTS.md 每波只许一枚写腿。
2. 那 34 枚**间接**红（经 `secret.NewStore` / `memory.Open` / config manager 递根）未逐枚探针复算其「另一种拒」，
   本腿只按红名差集判归零——若裁决方要求逐枚探针，请再放一格。
3. 票 123 那族三枚 300 s 腿（`TestL1Write...`、`TestLateVeto...`、`TestFSReadOnlyNeverOpensACard`）
   按账**豁免**，本腿对它们零需要、零贡献，也未顺手修。
4. 票面 AC 一格未翻（翻框归非实现者）。
