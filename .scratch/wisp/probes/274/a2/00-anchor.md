# 274-a2 起手锚（只读普查腿）

> 本件是**起手锚**，先于本腿一切取数写出并单独 commit。本腿只读：⛔ 不改产码/测试/脚本/CI、
> ⛔ 不跑 `go build|vet|test`、⛔ 不跑 `npm`、⛔ 不跑 `build.ps1`（`33-v4` 在飞，构建读数会被染色）、
> ⛔ 不 push、⛔ 不进 `D:/wt/fe` 工作树、⛔ 不动 `build/wisp.exe`（只 stat 留证）。

## 1. 现量锚（全部为本腿命令当场输出，非引用派单）

命令与时刻逐字记录：

```
$ date "+%Y-%m-%d %H:%M %z"
2026-10-08 10:54 +0800

$ git log -1 --format=%H
7a367a08292ef94908d86656608974e81c42f0d4

$ git rev-parse --abbrev-ref HEAD
dev

$ git log -1 --format='%ad%n%s' --date=iso
2026-10-08 10:49:00 +0800
A700 落账：收只读前置普查腿 255-a2 ⇒ …（首行截断，完整 subject 见 git log）

$ git status --porcelain -- scripts .github docs internal frontend | wc -l
0

$ git status --porcelain | wc -l          # 全仓（共享工作树，别的腿在写）
755

$ git status --porcelain -- .scratch/wisp/probes/274 | wc -l
0
```

HEAD **＝ `7a367a08`**（现场量，非派单里假设的号）。派单没有写死 HEAD，故不存在顶回；
但派单第 10:49 那一发的 subject 表明 **`33-v4` 仍在飞**，本腿据此自查：不跑任何构建/仪器。

## 2. 本腿落点与既有件盘点（写锚时已存在的 a2 内容）

```
$ ls -d .scratch/wisp/probes/274/*/ | wc -l
8

$ find .scratch/wisp/probes/274 -maxdepth 2 -name '*.md' -printf '%p %s %TY-%Tm-%Td %TH:%TM\n' | sort
.scratch/wisp/probes/274/a1/census.md            29577 2026-10-07 10:39   （普查腿 a1）
.scratch/wisp/probes/274/a2/shipping-chain.md    30884 2026-10-07 11:12   （普查腿 a2 既有件）
.scratch/wisp/probes/274/r1/impl.md              33439 2026-10-07 12:03
.scratch/wisp/probes/274/r2/impl.md              19272 2026-10-07 15:42
.scratch/wisp/probes/274/v1/verdict.md           29240 2026-10-07 14:57
.scratch/wisp/probes/274/v2/verdict.md          27314 2026-10-07 17:07
.scratch/wisp/probes/274/v3/verdict.md           16280 2026-10-07 17:46

$ for d in a1 a2 r1 r2 r3 v1 v2 v3; do echo "$d: $(find .scratch/wisp/probes/274/$d -type f | wc -l) files"; done
a1: 1 files / a2: 19 files / r1: 20 files / r2: 22 files / r3: 10 files / v1: 18 files / v2: 29 files / v3: 16 files
```

**⚠ 顶回派单一处**：派单说"票 274 前面有 `274-a1` 等普查件"，把 `a2` 当成待填的空格；
盘上实测 **`.scratch/wisp/probes/274/a2/shipping-chain.md`（30884 B，Oct 7 11:12）连同 17 份 logs 已存在且已跟踪**。
⇒ 本腿把这份 a2 既有件**也当作"前面已答过的格"**处理（查重见 `costs.md` §1），⛔ 不重做、⛔ 不改写它，
新落点只用 `00-anchor.md`（本件）＋`costs.md`。另：`r3/` 目录下**没有 `.md` 交付件**（只有 logs 与桩件）。

## 3. 两条"只读留证"（派单要求跑 exe 前先记，本腿只 stat 不跑构建）

```
$ ls -la build/wisp.exe
-rwxr-xr-x 1 swq 197609 31074357 Oct  7 11:57 build/wisp.exe

$ stat -c '%n size=%s mtime=%y' build/wisp.exe
build/wisp.exe size=31074357 mtime=2026-10-07 11:57:46.007254300 +0800
```

来历**不在本腿证据面内**（本腿不重算 md5、不对拉源码；机主口令「还原 274 前 exe」意味着这份是特意留的），
`costs.md` §2 按"通道通 / 带的是这份源码"两格分述，⛔ 不拿它的输出当"exe 带了当前页面"的反证。

```
$ git ls-files frontend/dist
frontend/dist/.gitkeep

$ git check-ignore -v frontend/dist/index.html frontend/dist/assets
frontend/.gitignore:12:dist/*	frontend/dist/index.html
frontend/.gitignore:12:dist/*	frontend/dist/assets
```

⇒ 跟踪名册**只有 1 枚**（`.gitkeep`）；盘上 `frontend/dist/` 另有**被忽略的历史构建件**
（`index.html` 1044 B / `assets/`，dir mtime Oct 7 11:58）。派单"dev 分支里 frontend/dist 只有 .gitkeep"
**对跟踪名册成立**，对**盘上现状不成立**——那几枚是旧源码落的件，属"可能过期"一类，`costs.md` §2(a) 具名分开写。

## 4. 本腿件册与自报

- 本件 `00-anchor.md`：起手锚，单独 commit（显式 pathspec）。
- `costs.md`：六问正文（含每把尺的精确命令逐字）。
- 追加节：`.scratch/wisp/issues/274-*.md` 末追加一节，⛔ AC 框一字不勾不动（前后 `grep -cE '^[[:space:]]*- \[ \]'` 同数写进件里）。
- 只读 `go env` / `go list`：若用到，逐条自报于 `costs.md` §7。

rc=0
