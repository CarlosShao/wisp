# 299-a1 · 10-runtime-asset-roster — 入口 HTML 运行时按 URL 取的资源名册

HEAD 锚＝`6ef14788`；钟＝`2026-10-10 09:43:39 +08`（起手）。

## 0. 先报一枚尺本身的不成立（具名，⛔ 不是绕过去）

派单尺①写的是 `git show HEAD:frontend/dist/index.html`。**这发跑不通**，现量：

```
$ git show HEAD:frontend/dist/index.html
fatal: path 'frontend/dist/index.html' exists on disk, but not in 'HEAD'
```

`frontend/dist/**` 在 HEAD 上跟踪的枚数（尺＝`git ls-tree -r --name-only HEAD -- frontend/dist`）＝**1 枚**，名册全文：

```
frontend/dist/.gitkeep
```

被忽略的证据（尺＝`git check-ignore -v frontend/dist/index.html`）＝`frontend/.gitignore:12:dist/*`。
`.gitignore` 里那段注释逐字写着为什么留那一枚锚：

> `# dist/ is a build product, but go:embed needs the directory to exist in a clean`
> `# checkout (ticket 77 AC#1), so one anchor file stays tracked; panel.Assets.`
> `# Built() reports false while the bundle is only that anchor, and the host shows`

### 这个差别对判据意味着什么（具名回答）

1. **"运行时名册"不是一份 HEAD 静态事实，而是"某一次构建的产物"的事实。** 同一枚 commit 上换一次 `npm run build`，hash 名就换（票面现量第 2 条用的 `index-B8yINMF1.js` 就是 08:51 那一发的名字）。⇒ 名册的可复现锚只能是**产物本身**（工作树 `frontend/dist/**`，或 `build/wisp.exe` 里的字节），⛔ 不可能写成 `HEAD:<path>`。
2. **exe 里那棵树是构建期才有的**（尺＝`git show HEAD:frontend/embed.go` 逐字 `//go:embed all:dist` ＋ `var distFS embed.FS` ＋ `func Dist() embed.FS`）。在 HEAD 的对象层里那棵树只有 `.gitkeep`，所以**从 HEAD 读 `Assets` 永远读到 `Built()==false`**（`internal/panel/assets.go:54-60` `newAssets` 逐字：只有 `fs.Stat(tree, EntryFile)` 成功才置 `a.built = true`）。⇒ 任何"Resolve 取不取得到"的判语都必须声明"以一枚构建过的树为前题"，否则一律是 `errNotBuilt`（`assets.go:35`）。
3. **本腿因此用两把尺并列**：名册内容取自工作树产物（08:51 那一发，mtime 实测＝`Oct 10 08:51`）；"取不取得到"取自 `Resolve`/`contentTypeOf` 的 HEAD 源码语义。工作树那一半⛔ 不能当契约锚，只能当"今天这台机器上这一发的读数"。

## 1. 第一层：入口 HTML 的引用（`entryRefRe` 同形正则）

正则本体逐字（`git show HEAD:internal/panel/assets.go` 第 126 行）：

```go
var entryRefRe = regexp.MustCompile(`(?i)\b(?:src|href)\s*=\s*["']([^"']+)["']`)
```

本腿那把同形尺（GNU grep 的 ERE 等价体，⛔ 不是同一个引擎，差别见 §1.3）：

```
grep -oiE '(src|href)[[:space:]]*=[[:space:]]*["'"'"'][^"'"'"']+["'"'"']' frontend/dist/index.html
```

**射程声明（派单要求逐枚写清"哪把尺＋射程文件＋含不含注释行"）**

- 尺＝上面那行 `grep -oiE`，`-i` 对应 `(?i)`，`-o` 取每个非重叠命中；
- 射程文件＝**只 `frontend/dist/index.html` 一枚**（工作树，1044 字节）；
- **含注释行**：`entryRefRe` 是对整段 entry 字节跑的、代码里**没有任何剥 HTML 注释的动作**，本尺同样不剥 ⇒ 那 7 行文件头注释在射程内。本发之所以没多计，是注释里没出现 `src=`／`href=`，⛔ 不是尺把它们排除了。
- 命中枚数＝**2**。

| # | 逐字命中 | `entryRefRe` 的捕获组① | `Assets.Resolve` 取不取得到 | Content-Type（`contentTypeOf` 逐字返回值） |
|---|---|---|---|---|
| 1 | `src="./assets/index-B8yINMF1.js"` | `./assets/index-B8yINMF1.js` | **取得到**（前题＝树已构建） | `text/javascript; charset=utf-8` |
| 2 | `href="./assets/index-yy8KMgdf.css"` | `./assets/index-yy8KMgdf.css` | **取得到**（前题＝树已构建） | `text/css; charset=utf-8` |

"取得到"的链条（HEAD 源码逐字锚）：`assets.go:79` `name := strings.TrimPrefix(strings.TrimPrefix(requestPath, "/"), "./")` ⇒ `./assets/…` 被削成 `assets/…`；`:83` 的三道拒绝（`..`／`\`／`:`）都不命中；`:86` `fs.ReadFile(a.tree, name)` 在 `//go:embed all:dist` 那棵树上读 `dist/assets/…`；`:93` `return data, contentTypeOf(name), nil`。工作树侧对得上：`ls frontend/dist/assets/` ⇒ `index-B8yINMF1.js`、`index-yy8KMgdf.css` 两枚，`wc -c` ⇒ **551989 / 49540**（与票面现量第 2 条逐字相同，本腿复尺复核）。

### 1.3 `Check()` 用同一枚正则做的不是同一件事（避免把 §1 当它的替身）

`assets.go:137-160` 的 `Check()` 除了跑 `entryRefRe`，还多了四道**跳过**（逐字 `:149-151`：`#`／`//`／`data:`／含 `://`／已见过），并先削 `?#` 尾巴（`:146-148`）。本尺**没有**这四道跳过 ⇒ 两把尺的枚数今天恰好都是 2，但**语义不同**：`Check()` 回答"引用在不在同一棵树里"，本尺回答"入口页面上写了几个 src/href"。票面现量第 4 条那句"`Check()` 只证字节在同一棵树里、证不了运行时取得到"与本节读数一致，本腿复核为**成立**。

## 2. 第二层：js 产物再扫一层动态形状

尺＝对 `frontend/dist/assets/index-B8yINMF1.js`（551989 字节，单行 minified）逐 pattern `grep -oa -- "$p" | grep -c .`。⚠ **`-a` 是必需的**（该文件被 grep 判为 binary，不加 `-a` 会只打印 `Binary file matches` 而给出 0 枚假读数）。

| pattern | 命中枚数 | 判语 |
|---|---|---|
| `import(` | **0** | 无动态 import ⇒ 运行时不会新增 js 分片 URL |
| `new URL(` | **0** | 无 `new URL(...)` 资源引用形状 |
| `import.meta.url` | **0** | 无 base-相对解析点 ⇒ Vite 没把产物拆成 `new URL(import.meta.url)` 那种静态资源入口 |
| `webpackAsync` | **0** | 非 webpack 分片器 |
| `\.wasm` | **0** | 无 wasm 取用（与"零 `'wasm-eval'`"的 CSP 注释对得上） |
| `Worker(` | **0** | 无 worker 脚本 |
| `XMLHttpRequest` | **0** | — |
| `fetch(` | **1** | **不是资源名册的一员**，见 §2.1 |
| `url(` | **1** | **假阳性**，见 §2.2 |

### 2.1 那枚 `fetch(`＝Vite 的 modulepreload 垫片

逐字片段（本腿 `grep -oaE '.{80}fetch\(.{60}'` 取的窗，⛔ 不是手写）：

```
onymous`?`omit`:`same-origin`,t}function n(e){if(e.ep)return;e.ep=!0;let n=t(e);fetch(e.href,n)}})();
```

它是"对 `<link rel="modulepreload">` 的 `href` 再发一发 fetch"的垫片。入口文档里**一枚 `rel="modulepreload"` 都没有**（§1 全文只有那 2 枚 src/href）⇒ 这一发 `fetch` 今天**不产生命中资源名册的请求**。⚠ 但如果将来入口页多出 modulepreload 链接（Vite 拆包后很常见），**这一层会自己长大**，名册的枚数判据就得把 `fetch(e.href,…)` 这一支算进去。本腿把它记成"今天 0 员、有明确的长大条件"。

### 2.2 那枚 `url(`＝CSS 单位判断，不是资源

逐字片段：`&&(hl.test(e)||e===`0`)&&!e.startsWith(`url(`))` —— 是样式值处理里对字符串 `url(` 的判断，⛔ 不取资源。

### 2.3 css 产物

尺＝`grep -oaE 'url\([^)]{0,60}' frontend/dist/assets/index-yy8KMgdf.css | sort -u` ⇒ **0 命中** ⇒ 今天 CSS 不引字体／图片 ⇒ 名册不会因为 `font-src`／`img-src` 再长。

## 3. 名册终态（本节终态）

**入口页运行时按 URL 去取的资源＝2 枚**：`./assets/index-B8yINMF1.js`、`./assets/index-yy8KMgdf.css`；两枚都在同一棵 `//go:embed all:dist` 树里、`Assets.Resolve` 在**构建过的前提下**都取得到、`contentTypeOf` 两者都给了正确的 `text/*`；第二层扫（`import(`/`new URL(`/`import.meta.url`/`.wasm`/`Worker(`）＝**0 枚新增**。
⚠ 三条限定，一枚都不能省：① 名册取自**工作树产物**，HEAD 对象层里 `frontend/dist/**` 只有 `.gitkeep` 一枚；② 枚数是**这一发构建**的读数，换一次 build 就换 hash；③ "取得到"说的是 `Resolve` 那棵 FS 树，⛔ 不是"宿主拿得到"——后者是件 `20`/`30` 的问题，本件不判。
