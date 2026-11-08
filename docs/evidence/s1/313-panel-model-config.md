# 313 页面侧落地 — 模型供应商配置的界面半（AC#1 补上承载物 ＋ 全格活体证据）

**本腿**：票 313 的页面侧收尾（写面 `frontend/**`，机主 2026-11-07 授权前端归本编队）。
**日期**：2026-11-08。**性质**：写腿交付件；非实现者对抗验收另行。

## 0. 接手时的现量（起手三核，全部现读）

读面核过四处：票 313 全文（AC#1/#2/#3 原文＋编排者判语节）、`.scratch/wisp/probes/313/a1/` 四表、`docs/evidence/s1/313-panel-ask.md`（页面读键权威表）、`cmd/wisp/panel_config_store.go`（`ReadSettings`/`ApplySetting` 真键现读）。

**重要现量修正（对本腿派单的更正，如实记）**：派单写「页面侧半格未做」，但 git 现量显示页面侧**主体早已由前腿落地**——`a3fdb50dd`（65 外观腿顺带落设置门类型）、`33b264d70`（329 接收器）、`9b820f8e8`/`f82b909e9`/`7c80ac032`（编排者亲改的模型页 v2/v3：供应商行＋拉取清单＋计费＋密钥＋删除）、`e62877c60`（219）。本腿逐件现读核实：

- **AC#2（两个模型都能改）**＝已落：`model-settings.tsx` 默认聊天模型卡走 `role_chat_pair`；`voice-settings.tsx` 听写/朗读链编辑器走 `voice_cloud_asr_chain`/`voice_cloud_tts_chain`、实时语音走 `voice_realtime_provider`+`voice_realtime_model` 两枚；全部写后 `onReload()` 重读。零缺口。
- **AC#3（录 key 不回显）**＝已落：`type="password"` 无控值、`saveKey` 发送前 `input.value = ""`、值只走 `credentialValue` 键（`config-actions.ts`→`panel.ts` `requestCredentialWrite`），屏上只显示 `已保存（dpapi:<id>）` 引用名。零缺口。
- **AC#1（默认选中 stepfun）**＝**真缺口，本腿补**：`ConfigView.configured/configuredKnown` 两个键 Go 侧早已送（`panel_config_store.go:103-104`），但全前端树**零消费**——添加表单的供应商 picker 初始恒为空，聊天对恒从 `chatProvider`（空）起。「默认高亮 stepfun」没有任何承载物。

## 1. 本腿改动（三件产码＋一件仪器）

1. `frontend/src/components/model-settings.tsx`：
   - 新增导出常量 `OWNER_DEFAULT_PROVIDER = "stepfun"`（出处＝机主原话「可以把 stepfun 作为我们默认的模型来使用」＋`313-panel-ask.md` §2 A877 裁定③）。
   - `ProvidersCard` 计算 `defaultProvider`，**两道门**：`view.configuredKnown && !view.configured`（用户没配过任何东西）**且** `availablePresets.includes(OWNER_DEFAULT_PROVIDER)`（名字必须真的在 Go 侧送来的预设清单里）⇒ 不是盲硬编码：用户改过（`configured=true`）连预选都不给；Go 预设表变了名字就自然不出现。
   - 预选只作 `useState(defaultProvider)` 的**初始值**——不点「添加」不写任何东西，用户第一下点击起就完全归用户（无 effect 再强加）。
   - 导出 `AddProviderForm`/`ProviderRow` 两个纯件（仅仪器挂载用；静态 markup 点不了「添加供应商」按钮）。
2. `frontend/scripts/render-model-settings.tsx`（新仪器，`render:*` 形状）＋ `frontend/package.json` 挂 `render:model-settings`。四牙全绿：
   - 牙① AC#1：预选 markup 落在 `aria-label="供应商"` 的 picker 按钮文字上（`>stepfun<`），且全 markup 无 `value="stepfun"` 隐藏携带；源级钉两道门＋`useState(defaultProvider)` 恰用一次。
   - 牙② 用户为准：`defaultProvider=""` 时无预选（placeholder 出现）；聊天对 markup 显示用户已存对 `stepfun / 示例模型-1`；源级钉草稿起自 `useState(view.chatProvider)/useState(view.chatModel)`。
   - 牙③ 密钥不上 DOM：密码框 `type="password"`；markup 无非空 `value="…"` 属性；canary 不出现；源级钉 `input.value = ""` 清空先于发送、值只走 `performCredential`、两个设置页文件均无 `credentialValue`/`configValue` 字样。
   - 牙④ 两个模型可改写回：聊天对＋语音三卡（ASR/TTS/实时）picker 全挂载（≥3 供应商选择器）；源级钉五个白名单字段名字面（`role_chat_pair`、两链 `field="voice_cloud_*"` 挂点拼写、realtime 两枚）＋`setBusy(null); onReload();` 写回重读。
   - 另有全 markup 零 emoji 检查（ban #8 仪器集）。
3. 证据截图两枚（真 Chromium 活体，见 §3）：`313-panel-ac1-stepfun-default.png`、`313-panel-voice-chains.png`（同目录）。

## 2. 门禁读数（现跑）

| 门 | 读数 |
|---|---|
| `npm run typecheck` | rc=0 |
| `npm run render:model-settings`（新仪器） | rc=0（OK 句见仪器尾注） |
| `npm run build` | rc=0（dist 入口 `index-Dk6PIr5X.js`） |
| `npm run lint` | rc=0，0 errors；36 warnings 全为既有类（`model-settings.tsx:371` 那枚 `react(refs)` 是上轮 auto-fetch 既有代码，本腿 diff 未触） |
| `sh scripts/d22scan.sh` | clean（frontend/=121 files 过 ban #6/#8） |
| LF | 本腿四件（2 产码 1 仪器 1 package.json）`git ls-files --eol`＝`w/lf`；新件 file 无 CRLF |

## 3. 活体证据（真 Chromium，同一 bundle）

方法（照 219 AC#6 腿先例）：`npm run build` 的 dist 经 `vite preview`（127.0.0.1:4189）供出，Tabbit 真浏览器加载**同一嵌入 bundle**，页面挂载后注入一枚夹具宿主门 `window.wispDispatch`（照 Go 侧线形：散文＋`wisp-config-view:` 标记行，载荷键集照 `config_313_test.go`；夹具模型名沿用仓内占位「示例模型-N」，⛔ 编真名）。产品 bundle 本身零改动、未挂任何 fixture 路由（P9 红线未碰）——假门只在浏览器控制面，不进 exe。

读数（全部现跑现读）：

1. **AC#1 默认选中**（夹具＝未配置 fresh）：点开「添加供应商」表单，供应商 picker 按钮文字现读＝**`stepfun`**（`document.querySelector('button[aria-label="供应商"]').textContent`）；屏显「未配置」徽标。截图＝`313-panel-ac1-stepfun-default.png`。
2. **AC#2/AC#3 已配置态**（夹具切 configured）：行内现读「已保存（dpapi:<id>）」真（引用名上屏，值不上屏）；聊天卡「当前：stepfun / 示例模型-1」真。
   **密钥活体牙**：点「更换」→ 向密钥框填 canary `sk-canary-313-live-check` → 点「保存密钥」→ 现读：**整棵 DOM `innerHTML` grep canary＝false**、密钥输入框已消失（返回已保存态）、夹具门记录到的封套＝`configField:"provider_credential"` 且 **`credentialValue` 有值＝canary、`configValue` 键不存在**（值只走专属键，共享封套干净）。
3. **语音页**（夹具带两链＋实时对）：「听写（ASR）」「朗读（TTS）」「实时语音」三卡全挂载；「当前：stepfun/示例模型-1」与「当前：stepfun / 示例模型-2」真；供应商 picker×3、模型 picker×4。截图＝`313-panel-voice-chains.png`。

## 4. 端到端（bundle 嵌入重编）

`npm run build` 后 `go build -o wisp.exe ./cmd/wisp` rc=0（51,008,954 B）；**嵌入核验**＝二进制内现读 `index-Dk6PIr5X.js`（本腿 16:33 构建的 dist 入口 hash，`grep -a -o` 唯一命中）⇒ exe 携带的是改后页面。本机 wisp 的 WebView2 面板窗白窗为已知环境账（A1631/329-v1 记录在案，旧二进制同症），活体证据因此走真 Chromium 通路如上；真机双击端到端仍归机主签收。

## 5. 边界与具名没做

- 夹具门是浏览器侧的测试载具：**config.set 的真实落盘链（DPAPI blob、config.toml）本腿零触碰**——Go 半已由 313-r1/r2/r3/320 系判成立；本腿只证页面半。
- `voice.realtime.enabled` 无界面口＝`313-panel-ask.md` §5 缺口②（A877 名册无 enabled），页面照实不画开关。
- 清空云链无界面口＝§5 缺口①，页面照实（链编辑器至少留一段）。
- 仪器的「用户点击后」行为（打开 picker、输入、点保存）静态 markup 不可达，以源级钉＋§3 活体读数共同覆盖。
- 临时件：`%TEMP%\ms313-p{1,2,3}.js`（浏览器程序）、`/tmp/{ms,build,lint,d22}.log`、根目录 `wisp.exe`——只建不删。
- ⛔ 未 push；commit pathspec 逐件点名。
