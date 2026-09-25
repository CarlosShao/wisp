Screens.register('chat', {
  nav: { icon: 'message-square-text', label: '对话' },
  html: `
    <div class="flex flex-col h-full">

      <!-- 顶部会话页签：会话名 + 下拉切换历史会话 + 右侧模型 chip -->
      <div class="flex items-center justify-between mb-3">
        <div class="relative">
          <button type="button" id="session-btn" class="flex items-center gap-1.5 rounded-md px-1.5 py-1 -mx-1.5 hover:bg-muted/60 transition-colors">
            <i data-lucide="messages-square" class="w-4 h-4 text-muted-foreground"></i>
            <span class="text-sm font-semibold text-foreground">归档桌面截图</span>
            <i data-lucide="chevron-down" class="w-3.5 h-3.5 text-muted-foreground"></i>
          </button>
          <div id="session-menu" class="hidden absolute left-0 top-full mt-1 w-52 rounded-lg border border-border bg-popover shadow-md p-1 z-30">
            <div class="px-2 py-1 text-[10px] font-medium uppercase tracking-wide text-muted-foreground">历史会话</div>
            <button type="button" class="session-opt w-full flex items-center gap-2 px-2 py-1.5 rounded text-xs hover:bg-muted text-foreground" data-s="归档桌面截图">
              <i data-lucide="folder-archive" class="w-3.5 h-3.5 text-muted-foreground"></i>归档桌面截图
            </button>
            <button type="button" class="session-opt w-full flex items-center gap-2 px-2 py-1.5 rounded text-xs hover:bg-muted text-foreground" data-s="周报生成">
              <i data-lucide="file-text" class="w-3.5 h-3.5 text-muted-foreground"></i>周报生成
            </button>
          </div>
        </div>
        <button type="button" id="model-chip-top" class="badge badge-primary">
          <i data-lucide="cpu" class="w-3 h-3"></i>&nbsp;deepseek-chat
        </button>
      </div>

      <!-- 消息流 -->
      <div class="flex-1 overflow-y-auto flex flex-col gap-5 pr-1">

        <!-- 日期分隔线 -->
        <div class="flex items-center gap-3 mt-1 animated-list-item" style="animation-delay:0ms">
          <div class="flex-1 h-px bg-border"></div>
          <span class="text-[11px] text-muted-foreground font-medium">今天 · 9月24日 周四</span>
          <div class="flex-1 h-px bg-border"></div>
        </div>

        <!-- 用户消息 1 -->
        <div class="flex justify-end animated-list-item" style="animation-delay:70ms">
          <div class="bg-secondary rounded-lg px-3.5 py-2.5 text-sm max-w-[75%] leading-relaxed">
            帮我把桌面上的截图按月份归档，重复的删掉
          </div>
        </div>

        <!-- Agent 完成回复 1 -->
        <div class="flex flex-col gap-2 animated-list-item" style="animation-delay:140ms">
          <!-- 推理过程（默认折叠） -->
          <div class="reasoning-wrap">
            <button type="button" class="reasoning-toggle flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors">
              <i data-lucide="chevron-right" class="w-3.5 h-3.5 reasoning-chevron transition-transform"></i>
              <span>推理过程 · 1.2s</span>
            </button>
            <div class="reasoning-body hidden mt-2 ml-4 pl-3 border-l-2 border-border text-xs text-muted-foreground leading-relaxed flex flex-col gap-1">
              <span>1. 识别意图：文件整理任务，列举 -> 比对 -> 移动。</span>
              <span>2. 桌面截图按修改时间归月；重复按 sha256 去重。</span>
              <span>3. 删除属 L2，本票未授权批量删除 => 改移到「待确认」而非直接删。</span>
            </div>
          </div>

          <!-- 工具调用 chips（四态 + Blur Highlight 焦点） -->
          <div class="blur-highlight flex flex-wrap items-center gap-1.5" id="tool-row-done">
            <div class="bh-item tool-item">
              <button type="button" class="tool-chip tool-toggle">
                <i data-lucide="folder-search" class="text-muted-foreground"></i>
                <span class="font-mono">fs.listdir</span>
                <span class="opacity-60 max-w-[110px] truncate">C:\\Users\\swq\\Desktop</span>
                <i data-lucide="check" class="text-[hsl(142_71%_38%)]"></i>
                <span class="badge badge-success">L0</span>
              </button>
              <div class="tool-detail hidden">
                <div class="tool-detail-line mono">{"path": "C:\\Users\\swq\\Desktop", "count": 12}</div>
                <div class="tool-detail-line">耗时 84ms · corrId a3f9c2</div>
              </div>
            </div>

            <div class="bh-item tool-item">
              <button type="button" class="tool-chip tool-toggle">
                <i data-lucide="folder-input" class="text-muted-foreground"></i>
                <span class="font-mono">fs.mv</span>
                <span class="opacity-60 max-w-[110px] truncate">截图 -> 2026-09</span>
                <i data-lucide="check" class="text-[hsl(142_71%_38%)]"></i>
                <span class="badge badge-warn">L1</span>
              </button>
              <div class="tool-detail hidden">
                <div class="tool-detail-line mono">{"moved": 9, "to": "Desktop\\截图\\2026-09"}</div>
                <div class="tool-detail-line">耗时 210ms · corrId b71e04 · L1 2.6s 内可撤销</div>
              </div>
            </div>

            <div class="bh-item tool-item">
              <button type="button" class="tool-chip tool-toggle" style="border-color: hsl(var(--destructive) / 0.4); color: hsl(var(--destructive));">
                <i data-lucide="ban"></i>
                <span class="font-mono">fs.delete</span>
                <span class="opacity-70">被拒：未授权批量删除</span>
                <i data-lucide="x"></i>
                <span class="badge badge-destructive">L2</span>
              </button>
              <div class="tool-detail hidden">
                <div class="tool-detail-line mono">{"refused": true, "reason": "批量删除 3 个及以上文件命中 R3，需 L2 原生批准"}</div>
                <div class="tool-detail-line">已改为移入「待确认」目录 · corrId c0d7ab</div>
              </div>
            </div>
          </div>

          <!-- 结果正文 -->
          <div class="text-sm text-foreground leading-relaxed">
            已归档 9 张截图到 <span class="font-mono text-xs">Desktop\\截图\\2026-09</span>；
            3 张疑似重复未直接删除，已移到 <span class="font-mono text-xs">Desktop\\截图\\待确认</span>，需要你在面板里逐条批准。
          </div>

          <!-- 结果分流标识（D10 tier-3：通知 + 面板） -->
          <div class="flex flex-wrap items-center gap-1.5 mt-0.5">
            <span class="badge badge-primary"><i data-lucide="bell" class="w-3 h-3"></i>&nbsp;系统通知</span>
            <span class="badge badge-primary"><i data-lucide="layout-dashboard" class="w-3 h-3"></i>&nbsp;结果面板</span>
            <span class="text-[11px] text-muted-foreground">播报首句不超过 60 字</span>
          </div>

          <!-- 成本行 -->
          <div class="text-[11px] text-muted-foreground font-mono tabular-nums mt-0.5">
            12,480 tok · ¥0.31 · IN 8,200 / OUT 4,280 / CACHED 3,100
          </div>
        </div>

        <!-- 用户消息 2（短结果，tier-4） -->
        <div class="flex justify-end animated-list-item" style="animation-delay:210ms">
          <div class="bg-secondary rounded-lg px-3.5 py-2.5 text-sm max-w-[75%]">现在几点了</div>
        </div>

        <div class="flex flex-col gap-1.5 animated-list-item" style="animation-delay:280ms">
          <div class="text-sm text-foreground">现在是下午 3 点 42 分。</div>
          <div class="flex flex-wrap items-center gap-1.5">
            <span class="badge badge-primary"><i data-lucide="volume-2" class="w-3 h-3"></i>&nbsp;TTS 播报</span>
            <span class="badge badge-primary"><i data-lucide="clipboard" class="w-3 h-3"></i>&nbsp;已复制剪贴板</span>
            <span class="text-[11px] text-muted-foreground">60 字以内短结果 · 不拉面板</span>
          </div>
        </div>

        <!-- 用户消息 3（触发流式 hero 回复） -->
        <div class="flex justify-end animated-list-item" style="animation-delay:350ms">
          <div class="bg-secondary rounded-lg px-3.5 py-2.5 text-sm max-w-[75%] leading-relaxed">
            会议纪要里有哪些待办？
          </div>
        </div>

        <!-- 流式进行中回复（hero：Thinking 子页签 + Streaming Text + 内联来源 + Follow-ups） -->
        <div class="flex flex-col gap-2 animated-list-item" style="animation-delay:420ms">

          <!-- Thinking 块：可折叠，头部 dots + shimmer + 耗时；展开为四子页签 -->
          <div class="thinking-wrap">
            <button type="button" class="thinking-toggle flex items-center gap-1.5 w-fit rounded-md px-1.5 py-1 -mx-1.5 hover:bg-muted/60 transition-colors">
              <i data-lucide="sparkles" class="w-3 h-3 text-muted-foreground"></i>
              <span class="text-[13px] font-medium text-foreground whitespace-nowrap">思考了 <span class="bu-think-time font-mono tabular-nums text-muted-foreground">0.0</span> 秒</span>
              <i data-lucide="chevron-down" class="w-3 h-3 text-muted-foreground thinking-chevron transition-transform duration-300"></i>
            </button>
            <div class="thinking-body hidden mt-1 ml-2 pl-3 border-l-2 border-border flex flex-col gap-2">
              <!-- Steps：3 步，前 2 完成 check，第 3 执行中 spinner -->
              <div class="trace-panel flex flex-col gap-1.5" data-panel="steps">
                <div class="flex items-center gap-2 text-xs">
                  <i data-lucide="check" class="w-3.5 h-3.5 text-[hsl(142_71%_38%)]"></i>
                  <span class="text-foreground">识别意图：文档整理 + 待办抽取（D11）</span>
                </div>
                <div class="flex items-center gap-2 text-xs">
                  <i data-lucide="check" class="w-3.5 h-3.5 text-[hsl(142_71%_38%)]"></i>
                  <span class="text-foreground">读取 artifacts 目录，定位会议纪要.md</span>
                </div>
                <div class="flex items-center gap-2 text-xs">
                  <span class="spinner-ring" style="width:12px;height:12px;border-width:1.5px"></span>
                  <span class="text-foreground">按议题分段，抽取待办事项</span>
                </div>
              </div>

              <!-- Reasoning：一段推理文本 -->
              <div class="trace-panel hidden flex-col gap-1.5" data-panel="reasoning">
                <span class="text-xs text-muted-foreground leading-relaxed max-w-md">
                  会议纪要按议题分为三段：版本规划、缺陷复盘、下周分工。待办通常落在「下周分工」段，需要抽出责任人与截止时间，再与 artifacts 目录里的任务台账去重，避免把已完成项重复登记。
                </span>
              </div>

              <!-- Search：2 条搜索来源 context-card -->
              <div class="trace-panel hidden flex-col gap-2 max-w-md" data-panel="search">
                <div class="context-card">
                  <div class="ctx-source">artifacts · 会议纪要.md</div>
                  <div class="text-xs text-foreground mt-1">9月24日站会 · 3 个议题，正文中标注待办 3 项</div>
                </div>
                <div class="context-card">
                  <div class="ctx-source">tasks · 待办台账.json</div>
                  <div class="text-xs text-foreground mt-1">已有未完成待办 5 条，抽取时需做名称去重</div>
                </div>
              </div>

              <!-- Coding：一小段 code-block -->
              <div class="trace-panel hidden flex-col gap-1.5" data-panel="coding">
                <div class="code-block max-w-md">
                  <div class="code-line"><span class="line-num">1</span><span><span class="tok-com"># 从纪要抽取待办</span></span></div>
                  <div class="code-line"><span class="line-num">2</span><span><span class="tok-key">todos</span> = extract(<span class="tok-str">"会议纪要.md"</span>)</span></div>
                  <div class="code-line"><span class="line-num">3</span><span>merge(todos, <span class="tok-str">"待办台账.json"</span>)</span></div>
                  <div class="code-line"><span class="line-num">4</span><span><span class="tok-key">return</span> dedup(todos)</span></div>
                </div>
              </div>

              <!-- 底部 pill 页签（Steps / Reasoning / Search / Coding） -->
              <div class="bu-pill-row" id="trace-tabs">
                <button type="button" class="bu-pill active" data-tab="steps">Steps</button>
                <button type="button" class="bu-pill" data-tab="reasoning">Reasoning</button>
                <button type="button" class="bu-pill" data-tab="search">Search</button>
                <button type="button" class="bu-pill" data-tab="coding">Coding</button>
              </div>
            </div>
          </div>

          <!-- 工具 chip：执行中（shimmer） + 等待（spinner） -->
          <div class="flex flex-wrap items-center gap-1.5">
            <div class="tool-item animated-list-item" style="animation-delay:60ms">
              <button type="button" class="tool-chip tool-toggle" style="border-color: hsl(var(--primary) / 0.5); color: hsl(var(--accent-foreground));">
                <i data-lucide="file-text" class="text-muted-foreground"></i>
                <span class="font-mono">doc.read</span>
                <span class="opacity-60 max-w-[120px] truncate">artifacts\\会议纪要.md</span>
                <span class="shimmer-text text-[10px]">执行中</span>
              </button>
              <div class="tool-detail hidden">
                <div class="tool-detail-line mono">{"file": "artifacts\\2026-09-24\\会议纪要.md", "bytes": 18420}</div>
                <div class="tool-detail-line">读取中 …</div>
              </div>
            </div>

            <div class="tool-item animated-list-item" style="animation-delay:180ms">
              <button type="button" class="tool-chip tool-toggle opacity-80">
                <span class="churning-grid v-drive" data-churning aria-hidden><i></i><i></i><i></i><i></i><i></i><i></i><i></i><i></i><i></i></span>
                <span class="font-mono">task.extract</span>
                <span class="bu-churn-label">运转中 <span class="bu-churn-time font-mono tabular-nums">0.0</span>s</span>
              </button>
              <div class="tool-detail hidden">
                <div class="tool-detail-line mono">排队等待 doc.read 完成后执行 · 3×3 churning 点阵</div>
                <div class="bu-variant-row">
                  <span class="text-[10px] text-muted-foreground">点阵样式</span>
                  <button type="button" class="bu-pill bu-variant active" data-v="drive">Drive</button>
                  <button type="button" class="bu-pill bu-variant" data-v="dots">Dots</button>
                  <button type="button" class="bu-pill bu-variant" data-v="orbit">Orbit</button>
                  <button type="button" class="bu-pill bu-variant" data-v="surfer">Surfer</button>
                </div>
              </div>
            </div>
          </div>

          <!-- 流式正文（逐词 Staggered Text） + 光标 + 重播 -->
          <div class="text-sm text-foreground leading-relaxed">
            <span id="stream-body"><span class="sw" style="opacity:0">会议纪要</span><span class="sw" style="opacity:0">已经</span><span class="sw" style="opacity:0">按议题</span><span class="sw" style="opacity:0">拆成</span><span class="sw" style="opacity:0">三段，</span><span class="sw" style="opacity:0">待办</span><span class="sw" style="opacity:0">抽出来</span><span class="sw" style="opacity:0">了，</span><span class="sw" style="opacity:0">一共</span><span class="sw" style="opacity:0">3 项，</span><span class="sw" style="opacity:0">下面</span><span class="sw" style="opacity:0">逐条</span><span class="sw" style="opacity:0">列给你。</span></span><span class="stream-cursor" id="stream-cursor"></span>
            <button type="button" id="stream-replay" class="hidden ml-1 inline-flex items-center gap-1 text-[11px] text-muted-foreground hover:text-foreground transition-colors align-middle">
              <i data-lucide="rotate-ccw" class="w-3 h-3"></i>重播
            </button>
          </div>

          <!-- 内联来源标记（2 个小 chip） -->
          <div class="flex flex-wrap items-center gap-1.5 mt-0.5">
            <span class="inline-flex items-center gap-1 rounded-md bg-muted px-1.5 py-0.5 font-mono text-[10.5px] text-muted-foreground">
              <i data-lucide="file-text" class="w-3 h-3"></i>会议纪要.md
            </span>
            <span class="inline-flex items-center gap-1 rounded-md bg-muted px-1.5 py-0.5 font-mono text-[10.5px] text-muted-foreground">
              <i data-lucide="file-text" class="w-3 h-3"></i>待办台账.json
            </span>
          </div>

          <!-- Follow-ups：3 个建议问题 -->
          <div class="mt-1">
            <div class="text-[12px] font-medium text-muted-foreground mb-1">Follow-ups</div>
            <div class="flex flex-col items-start">
              <button type="button" class="follow-up -mx-1.5 flex items-center gap-2 rounded-md px-1.5 py-1.5 text-left text-[12.5px] text-foreground hover:bg-muted transition-colors">
                <i data-lucide="corner-down-left" class="w-3 h-3.5 text-muted-foreground"></i>把这 3 项待办写进 tasks 屏
              </button>
              <button type="button" class="follow-up -mx-1.5 flex items-center gap-2 rounded-md px-1.5 py-1.5 text-left text-[12.5px] text-foreground hover:bg-muted transition-colors">
                <i data-lucide="corner-down-left" class="w-3 h-3.5 text-muted-foreground"></i>按人汇总下周分工
              </button>
              <button type="button" class="follow-up -mx-1.5 flex items-center gap-2 rounded-md px-1.5 py-1.5 text-left text-[12.5px] text-foreground hover:bg-muted transition-colors">
                <i data-lucide="corner-down-left" class="w-3 h-3.5 text-muted-foreground"></i>对比上月待办完成率
              </button>
            </div>
          </div>

          <!-- 预计分流（长/结构化 -> tier-1） -->
          <div class="flex flex-wrap items-center gap-1.5">
            <span class="badge"><i data-lucide="file-down" class="w-3 h-3"></i>&nbsp;预计落文件 + 面板</span>
            <span class="text-[11px] text-muted-foreground">含待办清单 3 项及以上 · 结构化优先于长度</span>
          </div>
        </div>

      </div>

      <!-- Conversation 红点标识（陪聊模式才显示，布局保留） -->
      <div id="conv-badge" class="hidden flex justify-end items-center gap-1.5 px-1 pb-1.5 text-xs mt-2">
        <span class="w-2 h-2 rounded-full bg-destructive"></span>
        <span class="text-destructive font-medium">Conversation</span>
        <span class="text-muted-foreground">· 全双工，麦克风常开，说话即插话（D47）</span>
      </div>

      <!-- Composer -->
      <div class="border-t border-border pt-3 px-1 pb-1 mt-2">
        <!-- 附件缩略图条 -->
        <div id="attachment-strip" class="flex flex-wrap items-center gap-2 mb-2">
          <div class="flex items-center gap-2 rounded-md border border-border bg-muted/40 px-2 py-1">
            <div class="w-8 h-8 rounded bg-primary/10 flex items-center justify-center">
              <i data-lucide="image" class="w-4 h-4 text-primary"></i>
            </div>
            <div class="flex flex-col leading-tight">
              <span class="text-xs font-medium">粘贴的截图.png</span>
              <span class="text-[10px] text-muted-foreground font-mono">1240×860 · 384 KB</span>
            </div>
            <button type="button" class="attach-remove text-muted-foreground hover:text-destructive" title="移除附件">
              <i data-lucide="x" class="w-3.5 h-3.5"></i>
            </button>
          </div>
        </div>

        <!-- 工作区路径 -->
        <div class="flex items-center gap-1.5 mb-2 text-xs text-muted-foreground">
          <i data-lucide="folder-cog" class="w-3.5 h-3.5"></i>
          <span class="font-mono">工作区：D:\\work\\workspace\\projects plans\\Wisp</span>
          <button type="button" id="ws-change" class="hover:text-foreground transition-colors flex items-center gap-0.5">
            <i data-lucide="chevron-right" class="w-3 h-3"></i>切换
          </button>
        </div>

        <!-- Prompt Bar（Rounded 变体） -->
        <div class="prompt-bar" id="prompt-bar">
          <button type="button" id="src-btn" class="text-muted-foreground hover:text-foreground p-1.5 rounded-md hover:bg-muted transition-colors" title="选择来源">
            <i data-lucide="at-sign" class="w-4 h-4"></i>
          </button>
          <button type="button" id="slash-btn" class="text-muted-foreground hover:text-foreground p-1.5 rounded-md hover:bg-muted transition-colors" title="命令">
            <i data-lucide="slash" class="w-4 h-4"></i>
          </button>
          <textarea rows="1" placeholder="输入消息或粘贴图片…"></textarea>
          <div class="prompt-actions">
            <div class="relative">
              <button type="button" id="model-btn" class="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors px-1.5 py-1 rounded-md hover:bg-muted">
                <i data-lucide="cpu" class="w-3.5 h-3.5"></i>
                <span>deepseek-chat</span>
                <i data-lucide="chevron-down" class="w-3 h-3"></i>
              </button>
              <div id="model-menu" class="hidden absolute bottom-full right-0 mb-1 w-44 rounded-lg border border-border bg-popover shadow-md p-1 z-30">
                <button type="button" class="model-opt w-full flex items-center gap-2 px-2 py-1.5 rounded text-xs hover:bg-muted text-foreground" data-m="deepseek-chat">deepseek-chat</button>
                <button type="button" class="model-opt w-full flex items-center gap-2 px-2 py-1.5 rounded text-xs hover:bg-muted text-foreground" data-m="deepseek-reasoner">deepseek-reasoner</button>
              </div>
            </div>
            <button type="button" id="mic-btn" class="text-muted-foreground hover:text-foreground p-1.5 rounded-md hover:bg-muted transition-colors" title="听写">
              <i data-lucide="mic" class="w-4 h-4"></i>
            </button>
            <button type="button" class="btn-destructive btn-sm" id="stop-btn" title="停止">
              <i data-lucide="square" class="w-3.5 h-3.5"></i><span>停止</span>
            </button>
            <button type="button" class="btn-primary btn-sm hidden" id="send-btn" title="发送">
              <i data-lucide="send" class="w-4 h-4"></i>
            </button>
          </div>
        </div>

        <!-- 底部行：模式下拉 + 提示 -->
        <div class="flex items-center justify-between mt-2">
          <div class="relative">
            <button type="button" class="btn-outline btn-sm" id="mode-btn">
              <i data-lucide="zap" class="w-3.5 h-3.5"></i>
              <span id="mode-label">自动</span>
              <i data-lucide="chevron-down" class="w-3.5 h-3.5"></i>
            </button>
            <div id="mode-menu" class="hidden absolute bottom-full left-0 mb-1 w-36 rounded-lg border border-border bg-popover shadow-md p-1 z-30">
              <button type="button" class="mode-opt w-full flex items-center gap-2 px-2 py-1.5 rounded text-xs hover:bg-muted text-foreground" data-mode="干活">
                <i data-lucide="hammer" class="w-3.5 h-3.5 text-muted-foreground"></i>干活
              </button>
              <button type="button" class="mode-opt w-full flex items-center gap-2 px-2 py-1.5 rounded text-xs hover:bg-muted text-foreground" data-mode="陪聊">
                <i data-lucide="message-circle" class="w-3.5 h-3.5 text-muted-foreground"></i>陪聊
              </button>
              <button type="button" class="mode-opt w-full flex items-center gap-2 px-2 py-1.5 rounded text-xs hover:bg-muted text-foreground" data-mode="自动">
                <i data-lucide="zap" class="w-3.5 h-3.5 text-muted-foreground"></i>自动
              </button>
            </div>
          </div>
          <span class="text-[11px] text-muted-foreground">Enter 发送 · Shift+Enter 换行</span>
        </div>
      </div>

    </div>
  `,
  onMount: function (app) {
    // ---- 顶部会话下拉 ----
    var sessionBtn = document.getElementById('session-btn');
    var sessionMenu = document.getElementById('session-menu');
    sessionBtn.addEventListener('click', function (e) {
      e.stopPropagation();
      sessionMenu.classList.toggle('hidden');
    });
    document.querySelectorAll('.session-opt').forEach(function (opt) {
      opt.addEventListener('click', function () {
        var name = opt.dataset.s;
        sessionMenu.classList.add('hidden');
        app.toast('切换到历史会话「' + name + '」（原型仅提示）');
        app.refreshIcons();
      });
    });

    // 顶部模型 chip
    document.getElementById('model-chip-top').addEventListener('click', function () {
      app.toast('当前模型 deepseek-chat（C8 抽象，原型仅提示）');
    });

    // ---- 完成态推理过程折叠 ----
    document.querySelectorAll('.reasoning-toggle').forEach(function (btn) {
      btn.addEventListener('click', function () {
        var wrap = btn.parentElement;
        var body = wrap.querySelector('.reasoning-body');
        var chev = btn.querySelector('.reasoning-chevron');
        body.classList.toggle('hidden');
        chev.style.transform = body.classList.contains('hidden') ? '' : 'rotate(90deg)';
      });
    });

    // ---- Thinking 块折叠展开 ----
    document.querySelectorAll('.thinking-toggle').forEach(function (btn) {
      btn.addEventListener('click', function () {
        var body = btn.parentElement.querySelector('.thinking-body');
        var chev = btn.querySelector('.thinking-chevron');
        body.classList.toggle('hidden');
        chev.style.transform = body.classList.contains('hidden') ? '' : 'rotate(180deg)';
      });
    });

    // ---- Thinking 子页签切换（Steps / Reasoning / Search / Coding） ----
    document.querySelectorAll('#trace-tabs .bu-pill').forEach(function (tab) {
      tab.addEventListener('click', function () {
        var key = tab.dataset.tab;
        document.querySelectorAll('#trace-tabs .bu-pill').forEach(function (t) {
          t.classList.toggle('active', t === tab);
        });
        document.querySelectorAll('.trace-panel').forEach(function (p) {
          var show = p.dataset.panel === key;
          p.classList.toggle('hidden', !show);
          p.classList.toggle('flex', show);
        });
        app.refreshIcons();
      });
    });

    // ---- Thinking 计时 + Churning 计时（onMount 起每秒 100ms 步进，保留 1 位小数） ----
    var thinkTime = document.querySelector('.bu-think-time');
    var churnTime = document.querySelector('.bu-churn-time');
    var t0 = Date.now();
    setInterval(function () {
      var s = (Date.now() - t0) / 1000;
      if (thinkTime) thinkTime.textContent = s.toFixed(1);
      if (churnTime) churnTime.textContent = s.toFixed(1);
    }, 100);

    // ---- Churning 点阵变体切换（Drive / Dots / Orbit / Surfer，纯视觉） ----
    var churnGrid = document.querySelector('[data-churning]');
    document.querySelectorAll('.bu-variant').forEach(function (v) {
      v.addEventListener('click', function (e) {
        e.stopPropagation();
        document.querySelectorAll('.bu-variant').forEach(function (x) {
          x.classList.toggle('active', x === v);
        });
        if (churnGrid) {
          ['v-drive', 'v-dots', 'v-orbit', 'v-surfer'].forEach(function (c) {
            churnGrid.classList.remove(c);
          });
          churnGrid.classList.add('v-' + v.dataset.v);
        }
      });
    });

    // ---- 工具 chip 展开/折叠；落在 blur-highlight 行里的 chip 切换焦点 ----
    document.querySelectorAll('.tool-toggle').forEach(function (chip) {
      chip.addEventListener('click', function () {
        var item = chip.closest('.tool-item');
        var detail = item.querySelector('.tool-detail');
        var row = chip.closest('.blur-highlight');
        if (!row) {
          if (detail) detail.classList.toggle('hidden');
          return;
        }
        var willOpen = detail.classList.contains('hidden');
        row.querySelectorAll('.bh-item').forEach(function (it) {
          it.classList.remove('bh-focused');
          var d = it.querySelector('.tool-detail');
          if (d) d.classList.add('hidden');
        });
        if (willOpen) {
          detail.classList.remove('hidden');
          item.classList.add('bh-focused');
          row.classList.add('has-focus');
        } else {
          row.classList.remove('has-focus');
        }
        app.refreshIcons();
      });
    });

    // ---- Follow-ups 点击 ----
    document.querySelectorAll('.follow-up').forEach(function (btn) {
      btn.addEventListener('click', function () {
        app.toast('已采纳建议：' + btn.textContent.trim());
      });
    });

    // ---- 附件删除 ----
    document.querySelectorAll('.attach-remove').forEach(function (btn) {
      btn.addEventListener('click', function () {
        var strip = document.getElementById('attachment-strip');
        var node = btn.closest('.flex.items-center.gap-2.rounded-md');
        if (node) node.remove();
        app.toast('已移除附件');
        if (strip && !strip.querySelector('.rounded-md')) strip.classList.add('hidden');
      });
    });

    // 工作区切换
    document.getElementById('ws-change').addEventListener('click', function () {
      app.toast('切换工作区需原生侧校验路径（C26），原型仅提示');
    });

    // ---- Prompt Bar 左右按钮 ----
    document.getElementById('src-btn').addEventListener('click', function () {
      app.toast('选择来源');
    });
    document.getElementById('slash-btn').addEventListener('click', function () {
      app.toast('命令面板（Ctrl+K）');
    });
    document.getElementById('mic-btn').addEventListener('click', function () {
      app.toast('语音输入（半双工干活路径 D16）');
    });

    // ---- Prompt Bar 模型下拉 ----
    var modelBtn = document.getElementById('model-btn');
    var modelMenu = document.getElementById('model-menu');
    modelBtn.addEventListener('click', function (e) {
      e.stopPropagation();
      modelMenu.classList.toggle('hidden');
    });
    document.querySelectorAll('.model-opt').forEach(function (opt) {
      opt.addEventListener('click', function () {
        var m = opt.dataset.m;
        modelMenu.classList.add('hidden');
        modelBtn.querySelector('span').textContent = m;
        app.toast('切换模型为 ' + m + '（C8 抽象，原型仅提示）');
        app.refreshIcons();
      });
    });

    // ---- 模式下拉（陪聊出红环） ----
    var modeBtn = document.getElementById('mode-btn');
    var modeMenu = document.getElementById('mode-menu');
    var modeLabel = document.getElementById('mode-label');
    var convBadge = document.getElementById('conv-badge');
    var promptBar = document.getElementById('prompt-bar');
    function resetRing() {
      promptBar.style.borderColor = '';
      promptBar.style.boxShadow = '';
    }
    modeBtn.addEventListener('click', function (e) {
      e.stopPropagation();
      modeMenu.classList.toggle('hidden');
    });
    document.querySelectorAll('.mode-opt').forEach(function (opt) {
      opt.addEventListener('click', function () {
        var mode = opt.dataset.mode;
        modeLabel.textContent = mode;
        modeMenu.classList.add('hidden');
        if (mode === '陪聊') {
          convBadge.classList.remove('hidden');
          convBadge.classList.add('flex');
          promptBar.style.borderColor = 'hsl(var(--destructive) / 0.5)';
          promptBar.style.boxShadow = '0 0 0 3px hsl(var(--destructive) / 0.12)';
          app.toast('陪聊模式：麦克风常开 + AEC，首次开启需 L2 隐私确认（D47）');
        } else {
          convBadge.classList.add('hidden');
          convBadge.classList.remove('flex');
          resetRing();
          if (mode === '自动') app.toast('切到全自动需 L2 强确认（M4），由原生侧批准');
        }
        app.refreshIcons();
      });
    });

    // 点击外部关闭下拉
    document.addEventListener('click', function closeMenus(e) {
      if (sessionMenu && !sessionMenu.classList.contains('hidden') && !sessionMenu.contains(e.target) && e.target !== sessionBtn) sessionMenu.classList.add('hidden');
      if (modelMenu && !modelMenu.classList.contains('hidden') && !modelMenu.contains(e.target) && e.target !== modelBtn) modelMenu.classList.add('hidden');
      if (modeMenu && !modeMenu.classList.contains('hidden') && !modeMenu.contains(e.target) && e.target !== modeBtn) modeMenu.classList.add('hidden');
    });

    // ---- 流式正文：逐词 Staggered Text（45ms/词 + 光标 + 重播） ----
    var streamBody = document.getElementById('stream-body');
    var streamCursor = document.getElementById('stream-cursor');
    var replayBtn = document.getElementById('stream-replay');
    var words = streamBody.querySelectorAll('.sw');
    var sendBtn = document.getElementById('send-btn');
    var stopBtn = document.getElementById('stop-btn');
    var WORD_MS = 45;
    var streamTimer = null;

    function clearStream() {
      if (streamTimer) { clearTimeout(streamTimer); streamTimer = null; }
    }
    function playWord(i) {
      if (i >= words.length) { finishStream(); return; }
      words[i].classList.add('stream-word');
      streamTimer = setTimeout(function () { playWord(i + 1); }, WORD_MS);
    }
    function playStream() {
      clearStream();
      words.forEach(function (w) { w.classList.remove('stream-word'); });
      void streamBody.offsetWidth; // 重排以重放 CSS 动画
      streamCursor.style.display = '';
      replayBtn.classList.add('hidden');
      sendBtn.classList.add('hidden');
      stopBtn.classList.remove('hidden');
      playWord(0);
    }
    function finishStream() {
      clearStream();
      streamCursor.style.display = 'none';
      replayBtn.classList.remove('hidden');
      sendBtn.classList.remove('hidden');
      stopBtn.classList.add('hidden');
      app.refreshIcons();
    }
    function stopStreamNow() {
      clearStream();
      words.forEach(function (w) { w.classList.add('stream-word'); });
      finishStream();
      app.toast('已停止生成（SSE 中断，取消不是错误）');
    }

    replayBtn.addEventListener('click', playStream);
    sendBtn.addEventListener('click', function () {
      app.toast('已发送，正在流式回复…');
      playStream();
    });
    stopBtn.addEventListener('click', stopStreamNow);

    // 进场：先停在思考态，约 1.1s 后开始逐词流式
    setTimeout(playStream, 1100);

    app.refreshIcons();
  }
});
