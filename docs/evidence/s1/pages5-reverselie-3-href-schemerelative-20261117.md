# 票 36 消毒收紧 · scheme-relative href 负控牙（红/绿）

- **修复腿**：`agent-105`（写腿，frontend 写面；实现者，⛔ 非裁决者）
- **时间**：2026-11-17 +08
- **来路**：`docs/evidence/s1/pages5-adversarial-20261117-0953.md` §5（审计发现 `TIGHT_SCHEMA` 的协议清单收不进**无冒号 URL**）
- **修复提交**：见派单回话首行（本文件与代码同一枚 commit）

## 缺口的机理（根因在库代码，非本项目实现误用）

`hast-util-sanitize`（frontend 实装 v5.0.2）`lib/index.js` 的 `safeProtocol`（:684-697）：
`colon = url.indexOf(':')`，**`colon < 0` 直接 `return true`** ⇒ 协议清单（`protocols.href = ["http","https","mailto"]`）只审「声明了 scheme」的 URL，`//evil.example/x` 与 `\evil.example/x` 这类 scheme-relative 形状**逐字放行**，浏览器把它们按页面 origin 解析 ⇒ 开放重定向面。

## 修了哪里

- `frontend/src/components/result-history.tsx:113-126` — 新增 `tightUrlTransform`（只放行 `http/https/mailto` 绝对 URL 与 `#` 锚；空串/其余 scheme 一律返回 null＝剥 href、留文本）。
- `frontend/src/components/result-history.tsx:139` — `MarkdownText` 的 `ReactMarkdown` 挂 `urlTransform={tightUrlTransform}`（react-markdown v10 在 rehype 插件**之后**、JSX 之前对全树 url 属性套它 ⇒ 白名单之上的最后一道）。
- `frontend/scripts/render-history-36.tsx` — `XSS_CASES` 加两枚负控牙：
  - `protocol-relative href`：`[点这里](//evil.example/x)` ⇒ hostile 断言 `/evil\.example|href="\/\//i` 不出现在 markup，`kept: "点这里"`（文本保留）；
  - `backslash scheme-relative href`：`[点这里](\evil.example/x)` ⇒ 同上（浏览器把 `\` 归一成 `/`，同一个把戏）。
  - 仪器总数 4 → 6（`render-history36: OK, 6 XSS families neutralized`）。

## 红/绿读数（现跑）

**红**（修复前形状＝schema-only，无 urlTransform；探针为临时件 `frontend/node_modules/.tmp/hrefprobe/`，gitignore 内、不入库）：

```
RED  (schema only, pre-fix state): href SURVIVES -> <p><a href="//evil.example/x">点这里</a></p>
```

——审计报告的存活面逐字节复现。

**绿**（修复后，生产 `MarkdownText` 同路，仪器 `render:history36` rc=0，6 族全灭）：

```
GREEN (production gate mounted): href stripped -> <p><a>点这里</a></p>
text kept on green: yes
```

## 硬约束自证

- 未动 `docs/PLAN.md`／`docs/specs/**`／thresholds／golden／allowlist／D43／design 契约文字。
- **未放宽断言**：负控是**加严**（新增两条 hostile 断言），原有 4 族断言逐字未动。
- 未 push；未提交 `frontend/dist` 字节（构建后删回 `.gitkeep`）；UI 零 emoji。

## 串行读数（node = D:\work\server\node24，v24.21.0 / npm 11.19.0）

```
$ npm --prefix frontend run build          rc=0
$ npm --prefix frontend run typecheck      rc=0
$ npm --prefix frontend run lint           rc=0（38 warnings, 0 errors——存量告警，非本次引入）
$ npm --prefix frontend run render:history36  rc=0（6 XSS families neutralized，含新负控红/绿）
```
