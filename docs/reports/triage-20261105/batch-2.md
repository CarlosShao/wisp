# 批 2 只读分档（56 57 58 60 62 64 65 68 70 71 72 77 85 98 100 103 107 108 112 122）

> 只读分档，2026-11-05。票面与代码一字未改，未跑 git 写操作。判据取自票面 Status／Blocked by／AC 框／Progress log 现读，外加盘上现量（`-done` 存在性、`ci.yml` 的两处版本字面量、`scripts/wisp-cli-tests.sh` 的 PATH 行）。
> 分档尺：T1＝此刻有非 frontend／非真机桌面／非机主拍板的可动 Go/scripts 判据（残余格子照批 1 惯例写在依据里）；T2＝语义前置是另一张未结票；T3＝余下每一格都要 frontend owner／真机桌面 GUI／CI-runner colour／钱／契约文字／纯机主裁决／RESERVED 触发门。

| 票 | 档 | 依据(一行) | 阻断者/首动作 |
|---|---|---|---|
| 56 | T3 | DEFERRED 三件外部事：SignPath 资格申请＋CA 签名过 SmartScreen、clean VM 上 winget/Scoop e2e、P10 命名/商标/域名裁决（票面明写冲突只由机主定，不许自创名）；`Blocked by` 另含 54 | 阻断：签名资质＋机主命名拍板（54 未 done）；无可派首动作 |
| 57 | T2 | 第三方 SDK／契约测试套件／registry 的范例与 e2e 要先有票 53 的命令型插件范例，S8 又须 54；票面 `Blocked by: 54 + 53 + user go-ahead` 逐字 | 阻断：53（次 54） |
| 58 | T3 | 首格即"UI 全走 locale 文件、frontend/原生渲染零硬编码中文"，其余要英文 ASR/TTS/KWS 选型（P-new）与英文明测基线；`Blocked by` 另含机主放行 | 阻断：frontend owner＋英文语音模型选型（54 未 done） |
| 60 | T3 | 四个 `Blocked by`（59/16/28/44）均已 `-done`；真正卡的是 Gate condition #2＝机主给 realtime provider Key，缺则票面规定 AUTO-DEFER；首格还要真 S2S 多轮与 barge-in | 阻断：机主 Provider Key（账号/钱）；按票面自动 defer，无可派首动作 |
| 62 | T3 | review 且票面自述"从未被验收"，8 个未勾框全是真桌面差分截屏（浅/深桌布可见性、逐态动效、靠边吸附＋DPI/副屏）＋回填 SPEC-08 要 owner 签字；AC#5 的 C8 wav 注入半格虽是 Go，判读仍靠差分截屏 | 阻断：真机桌面视觉＋owner 签收；可切出 sidecar＝AC#5 的 wav→液面单调性数值半（无桌面的那半） |
| 64 | T3 | review；8 格里只剩"多显示器：真拖到第二屏验证并留证"，票面自定"本机无第二屏则保持未勾" | 阻断：第二屏真机硬件 |
| 65 | T3 | blocked-on-owner；第一步就是机主把参考图入库 `design/refs/ball-glass-*.png`，完成判据是机主第二次看 `-tour` 说"这次对了"——是人的判决不是测试绿 | 阻断：机主参考图＋人工判决 |
| 68 | T3 | AC#1（尺寸真相）与 AC#4（无桌面半格）已勾；AC#2 默认值翻转与 AC#3 签收面复测票面明写 blocked-on-owner，要 owner 答 R15#3/#4/#5 且需桌面 | 阻断：机主答 R15#3/#4/#5＋真桌面 |
| 70 | T3 | blocked-on-owner：剩下四件的判据都不在代理手里（`internal/risk` 的 ubuntu 归口、staticcheck pin-vs-35 顺序、放行 `internal/config` 那行注释、AC#6 要 push 后五 job 全 pass） | 阻断：机主/编排者归口裁决＋CI-runner colour（push 后才有 AC#6） |
| 71 | T1 | 门禁自身 Go 活：AC#1/AC#2 的 `cmd/balldebug -diff` 阈值门（px_delta 低于阈值/缺行/timers 非 no 即非零退出；阳性对照可用"把阈值调到不可达"造红，不必真桌面）；AC#4 的 d22scan 台账已由 `b551fef` 落完，AC#3 已勾 | 首动作＝AC#1：`cmd/balldebug/main.go` 加 -diff 阈值门＋AC#2 阳性对照测（约 1 文件＋1 测文件）；残余 AC#4 的 ban#6 处置待人裁、AC#5 的"真 CI 红"归编排者 push |
| 72 | T3 | 修复已在树（`f1033e1`），AC#1/#2/#3/#5 已勾；AC#4 的 runner 半边票面自写"runner 侧必须来自一次真 run，不是本地等价环境"，AC#6 的步骤改名在 `ci.yml` | 阻断：CI-runner colour（push）＋ci.yml 归口（71 在途） |
| 77 | T3 | 09-23 owner 把前端收走（A102、README 规则 7），AC#1/#3/#4/#6 的正解全在 `frontend/**`；唯一非前端格是 AC#4 里 `tools/d22scan` 的 ban#8 `frontend/` 作用域台账那一行 | 阻断：frontend owner；可切出 sidecar＝d22scan 那一行作用域自报（归编排者，纯扫描器配置） |
| 85 | T1 | 票面"排队等写码槽位，怕撞 81/83/77"的理由已消失（81/83 done、77 冻结）；AC#1 只需钉 `ci.yml` 的 `go install mvdan.cc/gofumpt@latest`（现读仍在），AC#2 用已钉的 `staticcheck@2026.2.2` 在本机出第一份判据清单 | 首动作＝AC#1 钉子半＋AC#2：ci.yml 一行版本字面量＋跑 staticcheck 出"包×检查码×file:line"清单入 log（不动任何断言/allowlist）；AC#1 尾巴"一次真 run 的日志行"归编排者 push |
| 98 | T1 | 09-28 更正后只剩两格：①那道"永远跑不了"已被 `scripts/wisp-cli-tests.sh:20` 的 PATH 前缀与 `ci.yml:475` 推翻；②ubuntu 腿 19 枚红至今〔未定性〕才是真账——归因是只读跑测，不改 cmd/wisp 包 | 首动作＝②：纯净快照/容器里跑 wisp-cli-tests.sh 的 ubuntu 腿，把 19 枚顶格红逐条归因（首批只归因不修）；①"取数口令"是否入 AGENTS 归编排者裁 |
| 100 | T3 | RESERVED（不是 open 也不是 blocked-on-owner）：AC#0 触发门＝生产路径真会起子进程，票面自报今天"门未开"，且明令门未开不许开工、不许用主进程 ACL 折中 | 阻断：AC#0 触发门（现状未开，保持 RESERVED）；无可派首动作 |
| 103 | T2 | `rejected-needs-fix`，退回单＝票 108（AC#1 的守卫被验收三枚探针绕过，重做归 108）；AC#2 的 `RemoveUnlinked` 拒 junction 与 108 不同面 | 阻断：108（AC#2 那一半可独立先落） |
| 107 | T1 | 107b 的放行侧修复已在 HEAD（`8e10095`：root 腿 `treeResolvedAsNamed`＋target 腿"解析后仍在同一 root"），AC#1/#2 有 POSIX 真跑证据；欠的是以当前 HEAD 重走双侧与出新裁决表，纯 `internal/tools` Go，Docker/WSL 可跑 | 首动作＝AC#3：POSIX 容器＋Windows 双侧 `-count=2 -v` 重跑三枚探针（A/B/C）并出 1:1 裁决表（放行侧不许比拒绝侧宽）；AC#1 的 CI run id 归编排者补 |
| 108 | T1 | AC#1/#2/#4/#5 验收已过（Windows 56/56、变异四向、48 线程混打 0 race），唯一红格 AC#3＝POSIX `platformVerifyPlacement` 缺链接腿，已由票 113（done）补；欠当前 HEAD 复跑＋1:1 表 | 首动作＝AC#3：容器里重跑 POSIX link 腿并复判整表（约 1 包＋测）；与 112 报的装机守卫 8.3 短名问题同包，需同波协同 |
| 112 | T3 | AC#1 三条已分开定因，②③在 CI 已 PASS；剩下的 ① 与四条新红同根＝通知 `Path` 用解析后长名而用例拿调用方拼写去比，正解要在"通知带调用方拼写"与"按树比"间选，后者正面撞 D22 ban#2；8.3 短名形状只有 runner 能造，AC#5 结案要 CI success | 阻断：D22 契约/设计裁决（待裁）＋CI-runner colour（push）；本地造不出 runner 那个形状 |
| 122 | T2 | 票面 `Blocked by: 票 85` 逐字（同改 `ci.yml`＋"那把 staticcheck 得先真能跑"）；AC#7 的版本钉也写明与主账同一次 ci.yml 编辑，别单开一票 | 阻断：85 |

T1-count=5
T1: 71, 85, 98, 107, 108
