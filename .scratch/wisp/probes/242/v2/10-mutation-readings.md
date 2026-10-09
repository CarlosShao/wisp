# 242-v2 读数原件（命令逐字＋计数逐字，供复核）

CWD＝仓库根 `D:\work\workspace\projects plans\Wisp`。裁决与判语在
`docs/evidence/s1/242-grant-binding-v2.md`；这一枚只放命令与原始计数。

## 起手与基线

```
date                                          -> 2026-10-09 10:08 +0800
git rev-parse HEAD                            -> 2bdf385d3037b7063d797bd07428dc1d37cb2561
git log --oneline -3                          -> 2bdf385d / 6029b743 / 8ccc9568
git status --porcelain -- internal cmd        -> 空
ls -la .scratch/wisp/probes/242/v2/           -> anchor.md 2277 / logs/ / verdict.md 19358（10-08 前腿，未覆盖）
ls -la docs/evidence/s1/242-*.md              -> 只有 242-grant-binding-v1.md 29758
go test ./internal/agent/approval/ -count=1 -v -> rc=0 / RUN 91 / PASS 90 / FAIL 0 / ok 0.372s
```

## 载具（票 242 AC#1 那句"两枚 corr 不同"的盘上前置）

```
grep -n "CorrelationID: corrID|TaskID: taskID|CorrelationID: taskID" cmd/wisp/subagent_selfapproval_197_test.go
  -> 116:			TaskID: taskID, CorrelationID: corrID,     （唯一命中；票面 :14/:60 的 :109 同值句已过期）
grep -n "startChildWrite197(|-corr-1|-corr-2|foreignErr" 同一件
  -> 101 签名 / 396 corr-1 / 463 corr-2 / 473 rec.foreignErr = rt.gate.Native().Allow(…, card2.CorrelationID, card1.Grant)
  -> 252 {"别张卡的令牌（跨卡借证）", rec.foreignErr} / 262 一带 errors.Is(…, ErrBadGrant|ErrPanelAllow)
git log -1 --format 同一件 -> 77150dcc 10-08 19:01（票 259 AC#4）
PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./cmd/wisp/ -count=1 -run '…197…' -v
  -> rc=0 / PASS 2 / FAIL 0 / SKIP 0
grep -c "FORGED-OR-STALE|GRANT-DENY" <那次输出>  -> 0  （⇒ 路由级 denial token 不在 197 的断言面/打印面）
```

## 产码形状（本腿复量行号，⛔ 不抄票 259 §10 的 10-05 快照）

```
approval.go:558  func (s *grantStore) spend(nonce, bind string) grantDenial {
approval.go:569  	    if !equalSecret(stored, bind) {
approval.go:570  		    return denialMisbound
approval.go:574  	return denialSpentNonce
approval.go:477  //     property the ErrBadGrant comment states, and ticket 259 does not move it.   （仍成立）
ui.go:50-57      type PanelItem struct{ CorrelationID/Tool/Level/Reason/Paths/Depth }   （零 grant 字段）
ui.go:130        ErrBadGrant = errors.New("approval: 原生令牌无效（缺失/已用/与本次请求不绑定）")
ui.go:167-171    type PanelAPI interface { Reject / Head / View }
gate.go:626      func (g *Gate) Panel() PanelAPI { return panelAPI{g.q} }
gate.go:700-706  type panelAPI struct{ q *Queue } ＋ Reject/Head/View 三枚方法
queue.go:414     denial := it.grants.spend(nonce, it.bind)
queue.go:423/424 GRANT-DENY … denial=%s ／ FORGED-OR-STALE …(native grant missing/spent/misbound)
queue.go:571-577 func (q *Queue) viewLocked(it *qitem) PanelItem（投影）
cmd/wisp/approval_reply.go:229 / :274  合并文案第二、三份副本（票 259 §9 写的 :227/:272 已漂）
cmd/wisp/approval_reply_201_test.go:387  {"yes " + corr, "原生令牌无效"}   （子串钉）
```

## 七发突变（四件套逐发）

种前 hash：approval.go `67fb1468…b4fcf`／ui.go `55bcd1da…52b13`／gate.go `78e2d28a…4d3ce`／queue.go `66fec7ae…0cf41f`
（全值见裁决件表格）。⛔ 全部 `-overlay` 零使用；全部种在盘上。

```
MU-A   edit approval.go:569 -> `if false && !equalSecret(stored, bind) {`
       sed -n '569p' 复量 = 该行逐字
       go test ./internal/agent/approval/ -count=1 -v -> rc=1 / PASS 86 / FAIL 4
         ticket242_binding_test.go:40 / :57 / :171 ；ticket259_denial_rulers_test.go:172
       还原 -> hash=67fb146898dc7b235623bb8ad2fac0e0b2465fe1（等值）；porcelain 空
MU-M   edit approval.go:574 -> `return denialNone // MU-M planted by leg 242-v2: membership scan gutted`
       sed -n '574p' 复量 = 该行逐字
       go test … -v -> rc=1 / PASS 81 / FAIL 9
         名单：TestTicket242SpendRejectsForgedBindingAndConsumesTheNonce /
               TestTicket242ForgedBindingCannotSpendAnotherItemsGrant /
               TestTicket259R1DenialNamesSpentOrNeverLiveNonce / …NamesMisbound /
               TestTicket259R1OutwardAnswerStaysMergedAcrossDenials /
               TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer /
               TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis /
               TestGrantIsSingleUseAndBoundToItsItem / TestTicket224ForgedSessionAnswerRecordsNothing
         红句：ticket242_binding_test.go:46/:193；ticket259_denial_rulers_test.go:139/:180/:259；
               ticket259_panel_capability_rulers_test.go:301/:308；queue_test.go:167
       还原 -> hash 等值；go test ./internal/agent/approval/ -count=1 -> ok 0.415s
MU-D2a edit ui.go PanelItem 加 `Permitted     bool`
       sed -n '50,58p' 复量 = struct 含 Permitted bool
       go test … -v -> rc=1 / PASS 88 / FAIL 2
         ticket242_panelface_test.go:39 ；ticket259_panel_capability_rulers_test.go:174
MU-D2b 同字段改 `Permitted     int`
       sed -n '50,58p' 复量 = struct 含 Permitted int
       go test ./internal/agent/approval/ -count=1 -> rc=1 / FAIL 1（无 -v，PASS 枚数未枚举）
         唯一 FAIL = TestTicket242PanelItemReadFaceIsFullyDeclared ⇒ :174/:199/:229 三枚能力尺绿
MU-D2c MU-D2b ＋ queue.go viewLocked 加 `Permitted:     it.grants.live(),`
       sed -n '571,578p' 复量 = 该赋值在投影里
       go test … -v -> rc=1 / PASS 88 / FAIL 2
         ticket242_panelface_test.go:39 ；ticket259_panel_capability_rulers_test.go:229
MU-E   edit ui.go:171 加 `Allow(correlationID, grant string) error` ＋ gate.go:703
       加 `func (p panelAPI) Allow(corr, grant string) error   { return p.q.allow(corr, grant) }`
       sed -n '167,172p' / '702,703p' 复量 = 两行逐字
       go test … -v -> rc=1 / PASS 88 / FAIL 2（先打印 88 枚 PASS ⇒ 非编译失败冒充红）
         ticket259_panel_capability_rulers_test.go:73 / :82 / :101 / :111
MU-OUT edit ui.go:130 -> `ErrBadGrant = errors.New("approval: 原生令牌无效（与本次请求不绑定）")`
       sed -n '130p' 复量 = 该行逐字
       go test … -count=1 -v -run 'Ticket259R1Outward|Ticket259R1Denial' -> rc=1 / PASS 4 / FAIL 1
         ticket259_denial_rulers_test.go:252（同发四枚 Ticket259R1DenialNames* 全绿）
```

## 收工四件套（全发还原后）

```
git hash-object 四枚 = 种前四枚逐字等值
git status --porcelain -- internal cmd -> 空
go test ./internal/agent/approval/ -count=1 -> ok 0.380s
go test ./internal/tools/ -count=1 -v -run 'TestTicket283|Corr' -> 2 PASS / ok 0.063s（票 283 基线未被本腿洗）
git status --porcelain -- . 未取全仓（在飞改动不属本腿，见锚件）
```
