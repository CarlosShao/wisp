# 334-r1 · 零字节件的读数名分（这九枚是缺席证明，⛔ 占位件）

量法＝wc -c ＋ 各自那条命令的 rc；生成时刻＝本腿跑那发的时候，文件此后未再被写过。

| 件 | 它证明什么 | 尺 |
|---|---|---|
| `ac3-after-full.err.log` | 那一发的 **stderr 为空**＝没有 build failed／no test files／setup failed／0xc0000135（这几形一旦响了就会写进这里，见派单凡红句含 walk／build failed…）| `wc -c`＝0，同发主件里 rc 见 `grep '^rc=' ac3-after-full.log` |
| `ac3-after-roster.err.log` | 那一发的 **stderr 为空**＝没有 build failed／no test files／setup failed／0xc0000135（这几形一旦响了就会写进这里，见派单凡红句含 walk／build failed…）| `wc -c`＝0，同发主件里 rc 见 `grep '^rc=' ac3-after-roster.log` |
| `ac3-before-full.err.log` | 那一发的 **stderr 为空**＝没有 build failed／no test files／setup failed／0xc0000135（这几形一旦响了就会写进这里，见派单凡红句含 walk／build failed…）| `wc -c`＝0，同发主件里 rc 见 `grep '^rc=' ac3-before-full.log` |
| `ac3-before-roster.err.log` | 那一发的 **stderr 为空**＝没有 build failed／no test files／setup failed／0xc0000135（这几形一旦响了就会写进这里，见派单凡红句含 walk／build failed…）| `wc -c`＝0，同发主件里 rc 见 `grep '^rc=' ac3-before-roster.log` |
| `mut-baseline.err.log` | 那一发的 **stderr 为空**＝没有 build failed／no test files／setup failed／0xc0000135（这几形一旦响了就会写进这里，见派单凡红句含 walk／build failed…）| `wc -c`＝0，同发 rc 见 mut-baseline.summary.log 的 MUT_TAG 行 |
| `mut-M-A.err.log` | 那一发的 **stderr 为空**＝没有 build failed／no test files／setup failed／0xc0000135（这几形一旦响了就会写进这里，见派单凡红句含 walk／build failed…）| `wc -c`＝0，同发 rc 见 mut-M-A.summary.log 的 MUT_TAG 行 |
| `mut-M-B.err.log` | 那一发的 **stderr 为空**＝没有 build failed／no test files／setup failed／0xc0000135（这几形一旦响了就会写进这里，见派单凡红句含 walk／build failed…）| `wc -c`＝0，同发 rc 见 mut-M-B.summary.log 的 MUT_TAG 行 |
| `mut-M-C.err.log` | 那一发的 **stderr 为空**＝没有 build failed／no test files／setup failed／0xc0000135（这几形一旦响了就会写进这里，见派单凡红句含 walk／build failed…）| `wc -c`＝0，同发 rc 见 mut-M-C.summary.log 的 MUT_TAG 行 |
| `ac3-after-roster.log.redlist` | 改后那枚名册 **红名册为空**（顶层＋子用例逐名 0 行），对照改前那枚 `.redlist`＝4 行 | `wc -l`＝0，主件末行 FOURCOUNT RUN=24 PASS_TOP=6 FAIL_TOP=0 SKIP=0 rc=0 |

非零同批对照（证明这把尺分得开空与非空）：
302 ac3-before-roster.log.redlist
  0 ac3-before-roster.err.log
302 total
mut-M-A.err.log 等八枚空件的兄弟主件字节数：
 5004 ac3-after-roster.log
40520 ac3-after-full.log
 4918 mut-M-A.log
 3099 mut-M-C.log
53541 total
