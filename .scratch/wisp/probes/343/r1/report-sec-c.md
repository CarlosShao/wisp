# 票 343 交件报告 C 半（交件自证三把尺＋票框终量尺）343-r1

A 半＝`report.md`（靶子／读准的两把尺／交件句子／AC#0 同族结论／AC#3 结论／未定案上报）。
B 半＝`report-sec-b.md`（起手三读数／成对四数／8 枚牙／门禁四数）。拆件口径＝公共块那句
"每件 ≤160 行"只管我撰写的正文件、写不下⇒拆具名分件（`322-r2` 先例）、⛔ 压正文。

## 九 交件自证（三把尺）
```
尺① 名册尺＝单次 `git log --all --format=%H%x09%s` 导出＋内存内比前缀（⛔ --grep、
  ⛔ 假设其后任何分隔符＝编排者 27cefc8f 那条改正）；逐枚 `git show --name-only` 取
  并集 ⇒ ⊆ 三处写面、并集 24 枚、★.go 命中恰 1（正控＝同尺打 339-r1 读出 10 笔）：
  31d4a476  internal/risk/blacklist.go                                    （.go＝1）
  2b6f8a9a  .scratch/wisp/probes/343/r1/sh/linux-343.sh                   （.go＝0）
  360a6af6  .scratch/wisp/probes/343/r1/logs/ 18 枚（4 正文件＋14 原始 stdout）（.go＝0）
  bf5977f4  .scratch/wisp/issues/343-…-in-the-package.md（票尾一条）      （.go＝0）
  ＋report.md／report-sec-b.md／report-sec-c.md（3 枚）                       （.go＝0）
  ★编排者自己那把名册尺已独立读到过我（27cefc8f 里"343-r1 的 2 笔为真、现验都在 HEAD
  祖先链"）＝正控在盘上、⛔ 我自己盖章。
尺② 证据件绝对路径＋字节＋行数＋占位符尺
  根＝`D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\343\r1\`；逐枚
  `wc -c` 现跑（下表字节＝一次现量快照；★行数⛔ 抄死＝抄死随每次补笔过期、尺现跑）：
    `report.md` ＋ `report-sec-b.md` ＋ 本件（三枚＝报告笔）｜`sh\linux-343.sh` 1506 B
    `logs\01` 7407 B｜`logs\02` 2960 B｜`logs\03` 10764 B｜`logs\04` 9674 B
    `logs\linux-{base,del}-343.log` 32297／32306 B、各 428 行
    `logs\teeth-sweep\ac2-sweep-runner.log` 15827 B／218｜`pair-{base,del}-risk.log`
      42912／42895 B、各 543｜`ac2-MU-1{a..h}-teeth.log` 8 枚 4484..5801 B、各 51
      ｜`ac2-00-pristine-teeth.log` 3386 B／45
  ★⛔ 0 字节件（逐枚 `wc -c` 见上、最小 1506 B）；★⛔ `.out` 扩展名（只用 .md/.txt/
  .log/.sh＝`.gitignore` 那条全仓 `*.out` 会静默跳过＝票 275 那枚腿自抓）
  ★占位符尺＝逐枚扫四枚 .txt＋本三笔 report＋台件＝命中 **0**；★同尺正控＝把那几个
  串种进一枚一次性字符串、尺读出 4 枚命中 ⇒ 那枚 0 是尺在工作。★另＝`logs/01`–`04`
  那四枚（ASCII 编码落的读数件）非 ASCII 残留逐枚 0＝⛔ 乱码名／乱码正文件；本三枚
  report 是中文正文＝真 UTF-8 无 BOM（那把"非 ASCII＝0"的尺⛔ 适用它们、⛔ 混用）；
  我用 ASCII 写中文那次确实造出过问号乱码、已重落＝另记的自我缺陷。★分节符那枚码位
  U+00A7 在我全部写面件内 **0** 枚（派单⛔ 那枚字符本身 ⇒ 这里只写码位名、⛔ 把它打出来；
  我原先在 A 半确实打出过 4 次、被我自己那把尺抓到后清零＝同记一笔自我缺陷）。
尺③ 决定性那把尺的逐字输出＝`logs/02`（起手侧⇄交件侧、带 file:line）＋`logs/01`
  （两侧四数⇄红名册⇄8 枚牙逐名⇄11 次还原哈希）＋`logs/teeth-sweep/` 原始 stdout。
```

## 十 票框终量尺现跑（六格一枚没翻）
```
git grep -c '^- \[ \]' HEAD -- <票343路径>  ->  6
'^- [x]' 命中                                ->  0
票名／Status:／文件名字面                     ->  未动（标题行现读＝"# 343 — …"）
★⛔ 翻框、⛔ 改票名／Status:、⛔ 加 -done、⛔ push（我只 commit、推送归编排者）
```

