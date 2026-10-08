# done-class-a-1 起手锚（只登记，不改任何东西）

- 腿名：`done-class-a-1`（只读盘点腿；工作语言中文）
- 任务：甲类 10 张（125/196/211/213/214/215/216/219/258/261）全勾未 `-done` 票，逐张核"今天能不能收口"
- 禁区：零 Go 命令、零翻框、不改票/台账/索引、不 push、不新建票；写面＝只新建本目录 `*.md`

## 现量（逐字登记）

`date`：

```
Thu Oct  8 19:15:27 CST 2026
```

`git log -1 --format='%H %ad %s'`：

```
45305086a2b94fa8c2c4cc8e15471e9aeb346fb7 Thu Oct 8 19:15:03 2026 +0800 A740 落账：收 done-key-sweep-1（50f220b6/700f25ca）全池收口键体检——计数 工单263枚/含README264个.md/-done 95；甲类10张（125/196/211/213/214/215/216/219/258/261）全勾未收口，我抽验 258=0未勾4已勾、261=0未勾5已勾，⇒ 不当场批量收口（规则4还要 Status 行＋索引更新，收口前逐张核凭据）排下一波；乙类8张（07 L51/92 L68L73L79/97 L55/104 L52L57/105 L56/110 L39L47/113 L54L56/115 L58L64）已 -done 却有未勾框，我抽验 07-done=1未勾、113-done=2未勾，★正是 README L189-191 记过的票62先例⇒⛔不批量撤名，下一波逐张核（真未做⇒撤名记一笔/漏勾⇒补凭据留名）。★README 口径纠正（纠我转述）：我不该说'完成度只认 -done'——原文 L15/L30-31 规则4＝全勾→Status done→rename→更新索引，L201 明写'别因为看到 -done 就以为它零残余'⇒定式 -done 是收口动作的产物不是零残余证明。在飞＝comment-fix-land-1/242-status-1；零翻框零 push
```

`git status --porcelain | head -20`（只登记，属共享工作树他腿状态）：

```
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 M .scratch/wisp/probes/242/r3/logs/probe-routed.txt
 M .scratch/wisp/probes/268/v1/evidence.md
 M cmd/wisp/config_readers_255.go
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 D design/doubao/README.md
 D design/doubao/demo/app.js
 D design/doubao/demo/index.html
```

## 十张票的盘上文件名（`ls .scratch/wisp/issues/` 现量，rc_final=0）

- `125-posix-c26-does-not-install-when-the-temp-dir-is-a-symlink-and-nothing-pins-it.md`
- `196-two-task-state-vocabularies-already-coexist-and-the-schema-has-no-check-so-the-d43-names-are-not-the-ones-in-use.md`
- `211-subagent-pool-cannot-exceed-the-d38d-tool-ceiling.md`
- `213-composer-plus-menu-slash-command-catalog.md`
- `214-composer-plus-menu-attachments-and-context-wired.md`
- `215-composer-plus-menu-skills-plugins-inventory.md`
- `216-menu-display-base-scrub-control-chars-and-two-sentences.md`
- `219-approval-card-three-reply-buttons-and-a-reason-box.md`
- `258-the-schema-calls-hotkey-hot-tier-but-the-resident-leg-builds-the-ball-from-default-hotkeys-and-the-only-reloader-caller-is-balldebug.md`
- `261-the-model-enabled-key-promises-removal-from-discovery-and-selection-but-no-production-code-reads-it-while-the-default-true-tag-is-inert-for-map-entries.md`
