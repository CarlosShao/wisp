# 308-v1 · 30-ac3-tiers-ci（`AC#3` 档位／⛔ build tag 这枚真隐患／CI 色那问）

## 1. 三枚件的档位归属（尺＝`grep -n 'go:build'`＋`head -1`，射程＝工作树）

| 件 | 头行 | 档位 |
|---|---|---|
| `pathresolver_anchor_spelling_windows_test.go` | `1: //go:build windows` | 只 GOOS=windows 编译 |
| `pathresolver_junction_windows_test.go` | `1: //go:build windows` | 只 GOOS=windows 编译 |
| **`pathresolver_volumeless_spelling_308_test.go`（新件）** | **`1: package risk`＝⛔ 任何 build tag**（尺＝`grep -n 'go:build' <三枚>` 只命中前两枚） | **每个 OS 都编译、都求值** |

⇒ helper（`driveShare`/`volumeInvariantSpellings`）**只活在⛔ tag 的新件里**，而两枚 windows-tagged 用例⛔ 编译进 linux 档 ⇒ **守卫的牙在 linux 上有载具（那枚新用例），在 linux 上⛔ 任何 windows 用例**；反过来在 windows 档，牙⛔ 挂在 env 台面（见 20 件 §c 族）。这是"⛔ 用 env 形状、改用输入面"这条裁语的**机械后果**，⛔ 是腿的越界。

## 2. ★它会不会在 Linux CI 上被求值？——会，而且本腿量到了那一发的颜色

尺（逐条现跑，⛔ 转抄 r1 的行号）：
- `.github/workflows/ci.yml:419-420` `test-core:` / `runs-on: ubuntu-latest` ⇒ **`:487` `run: bash scripts/portable-tests.sh --scope=core`**
- `scripts/portable-tests.sh:236-246` core scope 里第 **240** 行＝`./internal/risk/... ./internal/statemachine/...` ⇒ **core 档含 `./internal/risk/`**
- `portable-tests.sh:164-190` `core_pin` 含 `github.com/CarlosShao/wisp/internal/risk` ⇒ 新件⛔ 挪 package 名册（pin 是**包路径**名册、⛔ 用例名名册，尺＝读 `:161-163` 注释＋pin 体）⇒ ⛔ 需要改 pin、⛔ 需要新增 ledger 行（新件⛔ skip、⛔ 造 SKIP）
- `ci.yml:533-534` `test-windows:` / `windows-latest` ⇒ **`:797` `--scope=windows`**（`portable-tests.sh:249-254`，第 **251** 行含 `./internal/risk/`）
- `ci.yml:813` `run: bash tools/d22scan/runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows` ⇒ ★本腿用 `go test ./internal/risk/ -list 'TestPathResolverJunctionWindows'`（⛔ 执行任何码）量到**只匹配 1 个名字**＝`TestPathResolverJunctionWindows`（件 `logs/v25-list-rulers.txt`）⇒ **那一步打⛔ 到两枚拼法用例**，它只是 Case 1 的兜底；两枚拼法用例的 windows 兜底档＝`:797`。⛔ 与编排者 §4 那句一致（复认）。
- `ci.yml:623` 名字＋`:672` `run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/` ⇒ **winlive 档⛔ 碰 `internal/risk`**；`grep -rn winlive internal/risk/*_test.go`＝**0** ⇒ `internal/risk` 里⛔ 任何用例带 `&& winlive`，⛔ 存在"winlive 档"这一层；本腿另跑 `go vet -tags winlive ./internal/risk/` ＝ `rc=0`（⛔ 变任何东西，符合"0 命中"）。

**结论（新件会⛔ 在 ubuntu CI 求值）**＝**会**。所以⛔ 接受"Linux 上大概过"这种推理，本腿做了三步，一步比一步硬：
1. `GOOS=linux go vet ./internal/risk/` → **rc=0**（件 `logs/v20`）
2. `GOOS=linux GOARCH=amd64 go test -c -o … ./internal/risk/` → **rc=0**，二进制 7 046 471 字节（⛔ 留仓内，已迁至仓外，件 `logs/v23-binary-relocated.txt`）
3. ★**真在 linux 上跑了一遍**：`docker run --rm golang:1.27`（容器内 `go env` 自报 `linux amd64 go1.27.1`，与本机同版）挂**仓外导出树**，`go vet ./internal/risk/ && go test ./internal/risk/ -count=1 -v`：
   - POST 树＝**`rc=0  RUN=169  PASS=108  FAIL=0  SKIP=1  PANIC=0`**（件 `logs/v26-docker-linux-full.txt`）
   - 其中 `--- PASS: TestVolumeInvariantSpellingsVolumeLess308`，`t.Logf` 打印的输入在 linux 上就是**原生无卷标**那句：`volume-less input "/tmp/wisp-308-no-volume/.ssh/id_testkey": 7 … 2 UNC spellings dropped by name [unc-admin unc-extended]`
   - 那一枚 `--- SKIP` ＝ `TestC26RewrittenSyncRootDoesNotDisarmSuspectNet` ⇒ **已在 `portable-tests.sh` 的 ledger 挂号**（`./internal/risk/|linux|fixture`），⛔ 新件带来、⛔ unaccounted ⇒ `runtests.sh` 的"SKIP 即致命"⛔ 被本票触发
   - 两枚 windows-tagged 用例在 linux 名册里出现次数＝**0**（＝整件被 tag 排除，GUARD A ⛔ 咬：该包在 linux 仍编译得出别的测试文件）
   - PRE 树同一尺＝`rc=0 RUN=168 PASS=107 SKIP=1`，名册作差＝**＋1 枚**（就是那枚新用例）⇒ linux 档⛔ 少跑任何东西
4. `VolumeName` 的 OS 语义（回答"期望是⛔ 依赖 OS"）：`filepath.VolumeName("/tmp/x")` 在 **linux＝`""`**、在 **windows＝`""`**（`/tmp/x` 是盘符相对形，⛔ server/share 段）；新用例的期望（¬ok＋去掉那两枚＋留 7 枚）在两边都成立 ⇒ **该用例 OS-无关**，实测两面绿（windows `logs/v-target-case.txt`、linux `logs/v26`）。

## 3. 本格判语

**附条件成立**。
- 满足的部分：档位问题**具名答完**（core/windows/那枚 `-run` 步骤/winlive 各带尺）；⛔ tag 的载具后果讲清；linux 侧⛔ 止于推理，本腿补到"同版本真跑"。
- 那枚条件＝**CI 跑机色⛔ 等于本机色**。票 `AC#3` 原句＝"⛔ 用'本机跑过'顶替 CI 色；拿不到 CI 色就具名登记欠读数"⇒ 本腿的 docker linux 与本机 windows 仍属"本机"，⛔ 是 `ubuntu-latest`/`windows-latest` 那两台 gear（镜像层数、`GOMODCACHE` 冷/热、actions 缓存、`fsutil 8dot3name` 态都可能不同；⚠ 本腿连 `fsutil` 都被 `Access is denied` 挡了，见 `logs/v46-fsutil.txt`）。⇒ 登记 **`Q-308-v1-A`**（三枚用例在 `:487`/`:797` 两步的实跑色，与 `A813` 那批同批取），⛔ 拿 §2 那三发行顶掉它。
