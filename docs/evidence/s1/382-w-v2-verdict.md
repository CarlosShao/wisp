# 票 R382-1 接线终裁腿 `382-w-v2` —— 非实现者裁决表（五格逐条判）

日期 2026-11-01 10:02-10:5x +08 · 腿 `382-w-v2`（终裁、非实现者；重派——382-w-v1 过夜静默死亡判死）·
被裁件 **`8189f4a9b`**（4 枚 +340/-5：`cmd/wisp/resident_voice_windows.go`、`resident_voice_engines_windows.go`、
`resident_voice_ttsinherit_382wr2_windows_test.go`（新）、件 `382wr2-l2-wiring.md`）·
派单存档 `.scratch/wisp/probes/orch/20261101/txt/382-w-v1-dispatch.txt`。
**方法**：`git archive 8189f4a9b` 导出仓外树 `D:\tmp\wisp\20261101\382wv2\tree`（工作树零碰；在飞
`16-v2`/`26p-v2` 只读同行未触未归因）；二进制三枚全自建（no-tag／mut1／mut2＋tagged `wisp26p`，rc 各 0）；
staged 全链（A1326 三 DLL 同目录＋`WISP_*` env 指真仓 third_party 只读＋cwd=导出树包目录＋`unset GOTMPDIR`
＋rc 重定向现取）；音频全程错峰串行。**全部读数本腿现跑**，实现件日志只作形状对照。

## 0. 一句话

**五格全过，R382-1 机制面可记已验**：382-v1 修法四项逐字落地；三读数本腿复跑绿；M1/M2 两发反向红形
如述＋还原核 3/3 逐位；零回归 34/1/0 逐位复现 392 基线＋`262/640/902` 逐位同＋26p (b) 判据位逐字＋(a)
solo 全不超 400ms。**本腿主新发现＝票 394 两形不穷尽**：生产带内容 mark 共 4 处跨 3 包，
`tools.Options.OnMark` 只能盖 1 处——第三形（risk 层 mark 观察缝）覆盖全部且契约面更小，应并入 394
选型单（契约批前）。

## 1. 五格逐条判

| # | 派单要求 | 判 | 本腿独立读数 |
|---|---|---|---|
| 1 | 修法逐字核（对照 382-v1 四项） | **可翻** | 见 §2 |
| 2 | 三态读数复跑 | **可翻** | 见 §3 |
| 3 | mutation 两发复刻＋还原核 | **可翻** | 见 §4 |
| 4 | 零回归（34/1/0＋26p＋262/640/902） | **可翻** | 见 §5 |
| 5 | 停手②证据链复核＋394 归口 | **可翻（394 两形不穷尽＝新发现）** | 见 §6/§7 |

## 2. 修法逐字核（382-v1 §8 四项对照，全部码面现读）

1. **字段**：`residentVoice` 加 `ttsGate *risk.TTSGate`（`resident_voice_windows.go:361`，注释钉死唯一写者
   ＝attachProvenance、与 prov 同生同死）。全包 grep：`rv.prov =` 恰 1 处（:609）、`rv.ttsGate =` 恰 1 处
   （:615）——nil 等价成立：早退肢（:605-608 作用域打不开）prov 与 gate 一齐留 nil，speak() 跳检＝旧
   prov==nil 肢逐字节同义。
2. **初始化**：`attachProvenance` 内 `rv.ttsGate = risk.NewTTSGate(p, risk.DefaultTTSGateConfig())`
   （:615）。借**同一引擎**核过：生产挂点 `resident_windows.go:577-582` 交的是 `src.run.c25`（bridge 往里
   mark 的那枚），注释逐字拒绝"私有引擎＝check against marks production never made"。
   `DefaultTTSGateConfig`＝形 B 开＋L2 开＋深度 16（`internal/risk/tts_gate.go:47-53`）＝票姿势，零调参。
3. **换调用**：speak() HOP 1 `_, refused, source := rv.ttsGate.CheckTextWithInheritance(req.scope,
   risk.ChTTS, text)`（`resident_voice_engines_windows.go:904`）；refused ⇒ `ttsGated.Add(1)`＋noteProblem
   ＋**return，绝不合成**（:905-908；闸在合成前＝测试 `slot_loads=0` 位钉死）。382-v1 写 `hit, refused,
   source :=`、实现弃 hit 取 `_`——refused 是唯一门变量、source 进问题行，形状等价。
4. **source 进 log**：noteProblem 带 `source`（:907），读数实测「来源 search.content (from closed task
   scope task-382wr2-closed)」与「来源 fs.read D:\work\secret.txt」两形。

## 3. 三态读数复跑（本腿二进制，log 01／05 两遍）

| 读数 | 判据 | 实测 |
|---|---|---|
| INHERIT（形 B 污句拒播） | refused=1 spoken=0 | `INHERIT refused=1 passed=0 spoken=0 samples=0 slot_loads=0`＋问题行带 "from closed task scope task-382wr2-closed"＋`risk/C25: recorded taint inheritance link`（闸在合成前：slot_loads=0） |
| CLEAN（干净句照播） | 照播 | `CLEAN passed=1 refused=0 spoken=1 samples=75565`（真 matcha 合成 6.85s） |
| DIRECT（活作用域仍拒） | 仍拒 | `DIRECT refused=1 passed=0 spoken=0 samples=0`（直查半边未因换接口丢失） |

实现件读数（75431／75363）与本腿（75565／75436）：**计数位逐位同**；samples 是合成非确定量，带内浮动。

## 4. mutation 两发复刻＋还原核（变异全仓外 `-overlay`，锚点各恰 1 处；log 02-05）

| 突变 | 摘什么 | 本腿红相（rc／读数） |
|---|---|---|
| **M1** | speak() 闸调用（`if rv.ttsGate != nil` ⇒ `if false && ...`，恒 no-refusal） | rc=1：INHERIT `timed out waiting for the gate to refuse the closed-task sentence`＋**污句照播 samples=56651**；DIRECT 同红**照播 42620**；CLEAN PASS——闸两半（继承扫描＋直查）全死（实现件 56482 同形） |
| **M2** | attachProvenance 的 `NewTTSGate` 初始化行 | 全套 rc=2：**nil pointer panic**（INHERIT 0.12s 死于测试自身 `rv.ttsGate.RecordTaintInheritance`＝闸对象整个不存在的形状）；direct-only rc=1：**污句照播 42482**＋等拒超时（实现件 panic 形＋42635 同形） |
| **还原核** | pristine 二进制（变异前构建、全程未动）复跑 | rc=0 **3/3 PASS，计数位与首发逐位同**（INHERIT refused=1／CLEAN spoken=1／DIRECT refused=1） |

## 5. 零回归复核（本腿 8189f4a9b 纯导出树；log 09-12）

- **六族 staged 合跑**（`Test43|Test15e|Test41a2r1|Test381r2|Test390r1|Test26R3|Test384`）：**34 PASS /
  1 FAIL / 0 SKIP，逐位复现 392 基线**；唯一红 `Test384ResidentWatchdogEndToEnd`（3.32s，合跑时窗），
  **本腿单跑 PASS（1.22s）**＝A1329 既有归因自证。
- **`262/640/902` 逐位同**：`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms`（Test15e 内）。
- **26p (b)** PASS（19.31s）：`ECHO-ZERO episodes=1 all_below_cut=true triggers=0 finals_delta=0
  handed_delta=0 windows_untrusted=6`＝392 基线判据位逐字。
- **26p (a) solo** PASS（25.77s）：五轮 `36/89/85/82/96ms` 全不超 400ms（实现件 solo 26/83/50/50/34；两腿
  皆绿，392-v2 已钉 (a)＝摆动面、solo 绿是时段读数非稳定性质——照账引用）。
- **读数树位差（如实记账）**：本腿读数取于 `8189f4a9b` 纯导出树，实现件取于共享树合并态（393-r2 在
  飞）——两树同得 34/1/0，基线归因更强。
- **staging 两课（非码面，392-v2 同族）**：导出树缺未追踪 `build/models/manifest.json.minisig` ⇒ 15e
  e2e 红（本腿首跑 33/2/0 弃，补 sig 后逐位复现）；`third_party/model-fetch` 以 junction 指真仓只读后
  标点级正常挂载。后续导出树复跑者先补这两件。
- 工件：26p 测试 hardcode out 落 `D:\tmp\wisp\20261031\26pr1\out\26p-b-echo-zero.json`＝本腿复跑覆盖写
  （先前份在 `392v2\out-backup`，392-v2 §8 先例同形）。

## 6. 停手②证据链复核（三条逐条码面核实）

1. **mark 点普查**（生产非测试）：`internal/tools/bridge.go:687` `MarkWithHostPath(dec.TaskID, …, res.Text,
   hostPath)`——毒内容（res.Text）只在桥内跨界成立。**但"生产唯一 mark 点"表述不周延**：
   `internal/tools/subagent_197.go:512` `Mark(parentID, SrcTaskOutput, …, text)` 同是**带内容、进任务作用
   域**的生产 mark（子代理结论进父任务）；另有 `cmd/wisp/panel_assets.go:245`（argv 内容、私有引擎、非
   任务作用域）与 `internal/projctx/projctx.go:349`（文件内容进 loader 作用域）。
2. **Options 无观察钩**：`internal/tools/bridge.go:32-109` 全部字段现读——OnUpdate＝流式进度增量、
   OnDecision＝执行前裁决，无一携带 mark 内容。判"无 mark 观察钩"成立。
3. **ScopeTaints 只回元数据**：`internal/risk/provenance.go:654-669` 只回
   `TaintInfo{Tool,Origin,RuneLen,MarkedAt}`，注释逐字 "never the content itself"；`Scope.Close`
   （provenance.go:447）纯 `delete(p.scopes, s.id)`，无钩、无内容出口。判成立。
⇒ 判：admitTask 收尾闭包（`cmd/wisp/run.go:1185-1190`）只能发空内容（假接线）或拿 Origin 冒充内容
（改闸语义），停手②**成立**，实现腿不造假挂点正确。

## 7. 主新发现：票 394 两形不穷尽——risk 层 mark 观察缝为第三形

394 机理节两形：①`tools.Options` 加 OnMark 观察钩（扩 tools 契约）②闸／调用方下沉 internal/tools。
本腿判**不穷尽**，第三形＝**`risk.Provenance` 加 mark 观察缝（ProvOptions 钩或 TTSGate 借引擎时自注
册），闸在引擎侧自动累积**：

- **覆盖面**：§6.1 实测带内容 mark 共 4 处跨 3 包——Options.OnMark 只盖 bridge 1 处，**盖不住
  subagent_197.go:512（同是任务作用域毒）**；risk 层缝天然盖全部。
- **契约面更小**：只动 risk，tools 契约零扩。且形 B 契约"调用方在 mark 时累积**归一化**内容"对调用方
  根本无法履行——`normalizeTaint` 是 risk 未导出函数（`internal/risk/taintmatch.go:49`），382wr2 测试被
  迫手选已归一化常量即此痛；第三形把归一化义务收回唯一有权履行的层。
- **不属被否的形 A**：不动 `Scope.Close` 删除语义（provenance.go:447 一字不动）、不延长 scope 寿命，
  tri-v1 否决理由不适用。
- **代价（如实）**：链接留存策略从"调用方显式 opt-in"变"引擎侧姿势"（哪些 mark 进链需姿势界定——
  至少任务作用域 mark 必收、announce/session 常驻作用域建议不收）；`checkAncestorScopes` 第二循环全
  链扫描（tts_gate.go:181-189）使任务边界对匹配非实质，姿势面须写清。
- **处置建议**：394 契约批前由编排者把第三形并入选型单（批仍归机主）；三形并陈不改变停手②判语与本
  票机制面结论。

## 8. 五格判＋销账判语

五格全过。**R382-1 机制面可记已验**：R382-1 名目＝修法四项（字段／初始化／换调用／source 进问题行），
全部落地且经非实现者四重复核（逐字／三态／反向双发＋还原／零回归）。销账路径照派单：票 394 契约批
后由执行腿收生产积累钩；**L2 弹卡 UI 面仍归票 26 集成残余**（本票只接机制调用＋三态读数，382-v1 判
例同形）。残余三件：①394 生产积累钩（契约待批，第三形并入选型）②26p (a) 摆动面（A1329／392-v2 在
账）③L2 卡（26 面）。

## 9. 纪律与写面

仓内仅本件 `docs/evidence/s1/382-w-v2-verdict.md`；仓外 `D:\tmp\wisp\20261101\382wv2\**`（tree＝8189f4a9b
导出、bin-notag／bin-tag 自建二进制、mut 变异件＋overlay、logs 01-12、cases、run382.sh——只建不删）。
未落任何 `.go` 进仓（变异全 overlay 仓外）；票面／台账／thresholds／golden 零触碰；在飞 `16-v2`／
`26p-v2` 只读同行零触碰零归因。d22scan 本腿现跑＝**clean rc=0**（分母 bans#1-5 internal/=287、cmd/=52，
ban#8 internal/=730、cmd/=183——cmd 较实现腿 +1＝16-r1 已落库新测试件，非在飞脏；docs/ 不在 ban#8 扫
描分母，本件另以禁字符面自核 0 命中）。
