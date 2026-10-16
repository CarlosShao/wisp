收 166-r2 补齐三处写面：证据件（94 行）＋读数件 23 枚＋票尾一条 bullet ⇒ ★更正 c2b1ade5 那句"盘上唯一一枚中间态"不实

代提清单（三处写面逐枚打勾，规约＝A999 那条）：
1) 产码＝internal/session/manage_test.go（14 加／2 删）── 已在 c2b1ade5，本笔不动。
2) 证据件＝docs/evidence/s1/ticket166-ac6-166r2.md（现量 94 行／9363 字节，占位符 0）── 盘上有、我当时漏提交。
3) 读数目录＝.scratch/wisp/probes/166/r2/（盘上 23 枚＝暂存 23 枚，逐枚对数）＋本票末尾追加一条。

★对我自己上一笔的更正：c2b1ade5 正文写"盘上唯一一枚中间态"，实测不实── 该腿的证据件与 23 枚读数件都在盘上，
我当时只 git add 了产码一枚就交那笔。按⛔ amend 规矩，原笔一字不改，就地追加这一笔更正与补齐。
同一枚坑今天第二次抓到（A999 那次＝165-r2 的 165-wiring.md），定式已写但执行仍漏 ⇒ 新增尺：
代提之后必跑一把 "git status --porcelain 对自家票号目录＋docs/evidence/s1 的残留计数"，非 0 就当场补。

读数摘要（腿自陈，⛔ 我裁决；AC 框六枚全未碰）：
AC#6 ⓐ 五发齐（改前 FAIL_TOP=7／第二发 12 证非恒定／改后 20/20×2＋Search 同形＋TestSession 族 100/100＋整包 720/0
／摘步进的突变态 FAIL_TOP=22＋还原绿且 cmp 0 差异）。
AC#6 ⓑ 答＝看不见（git grep session.NewManager rc=1 零命中；Manager.List/Search 9 枚调用点全在测试件）。
两处与派单前提不符由腿具名上报：:84 的 Search 实测不是硬币（改平仍 20/20 绿，其断言只比命中枚数）⇒ :84 算防未来；
三形 tie-break 料上交，其中 AppendSessionTurn 不刷 session.updated_at 属行为缺陷，归我裁、⛔ 本腿动。

next＝派 166-v1（非实现者一次裁 AC#1–AC#6）。本笔零产码改动、⛔ push。
