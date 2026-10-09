# 296-a1 — R1b 原始读数全文（`.scratch/wisp/probes/**` 下的变异拷贝，本腿显式排除、不计入 AC#0 名册）

尺命令逐字：`grep -rn "ball\.HotkeyConfig{" --include=*.go .scratch | grep -v _test`
读数枚数：**8**

```
.scratch/wisp/probes/258/v1/mut0/cmd/wisp/resident_windows.go:182:			return ball.HotkeyConfig{}
.scratch/wisp/probes/258/v1/mut0/cmd/wisp/resident_windows.go:188:		return ball.HotkeyConfig{Summon: h.Summon, Mute: h.Mute, Cancel: h.Cancel, Panel: h.Panel}
.scratch/wisp/probes/258/v1/mut0/cmd/wisp/resident_windows.go:206:			return ball.HotkeyConfig{}
.scratch/wisp/probes/258/v1/mut0/cmd/wisp/resident_windows.go:208:		return ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel}
.scratch/wisp/probes/258/v1/mut3/cmd/wisp/resident_windows.go:183:			return ball.HotkeyConfig{}
.scratch/wisp/probes/258/v1/mut3/cmd/wisp/resident_windows.go:189:		return ball.HotkeyConfig{Summon: h.Summon, Mute: h.Mute, Cancel: h.Cancel, Panel: h.Panel}
.scratch/wisp/probes/258/v1/mut3/cmd/wisp/resident_windows.go:207:			return ball.HotkeyConfig{}
.scratch/wisp/probes/258/v1/mut3/cmd/wisp/resident_windows.go:209:		return ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel}
```

两枚变异体＝票 258 验收腿 `v1` 留下的 `cmd/wisp/resident_windows.go` 整文件拷贝（`mut0`＝未变异基线、`mut3`＝第 3 号变异），
行号相对生产件（`:188/:195/:213/:215`）各偏移 −6 与 −5／−6／−6，**不参与生产路径** ⇒ AC#0 名册不计。

⚠ **本腿自报一处越界**：本枚内容最初被我落成了 `raw-scratch-copies.txt`，而派单规则 3 写明写面**只许** `.md`。
依"临时件只建不删"（`issues/README` 规则 8）我不删它，改以本 `.md` 承载全文；**入库名册只含 `.md` 四枚**，那枚 `.txt` 留在盘上未 add、未 commit，具名报给编排者处置。
