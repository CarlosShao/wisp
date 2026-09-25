Screens.register('firstrun', {
  nav: { icon: 'life-buoy', label: '首次引导' },
  html: `
    <div style="width:400px;" class="flex flex-col">
      <div class="mb-5">
        <div class="page-title" style="font-size:18px;">欢迎使用一缕</div>
        <div class="page-subtitle" style="margin-bottom:0;">三步完成首次配置 · 每步均可跳过，稍后在设置里补</div>
      </div>

      <!-- 三步进度轨 -->
      <div class="flex items-center mb-6" id="rail">
        <div class="flex flex-col items-center">
          <div class="rail-dot w-6 h-6 rounded-full bg-primary text-white flex items-center justify-center text-xs font-medium" data-rail="0">
            <i data-lucide="check" class="w-3.5 h-3.5"></i>
          </div>
          <span class="rail-label text-[11px] font-medium mt-1.5 text-primary" data-rail-label="0">模型下载</span>
        </div>
        <div class="rail-line flex-1 h-0.5 bg-muted mx-2 mb-4" data-rail-line="0"></div>
        <div class="flex flex-col items-center">
          <div class="rail-dot w-6 h-6 rounded-full bg-muted text-muted-foreground flex items-center justify-center text-xs" data-rail="1">2</div>
          <span class="rail-label text-[11px] text-muted-foreground mt-1.5" data-rail-label="1">目录授权</span>
        </div>
        <div class="rail-line flex-1 h-0.5 bg-muted mx-2 mb-4" data-rail-line="1"></div>
        <div class="flex flex-col items-center">
          <div class="rail-dot w-6 h-6 rounded-full bg-muted text-muted-foreground flex items-center justify-center text-xs" data-rail="2">3</div>
          <span class="rail-label text-[11px] text-muted-foreground mt-1.5" data-rail-label="2">API Key</span>
        </div>
      </div>

      <!-- 第 1 步：模型下载 -->
      <div class="step-panel" data-step="0">
        <!-- Preloader 式状态：进度环 + 流光语义标签 + 呼吸点 -->
        <div class="card mb-3">
          <div class="flex items-center gap-3">
            <!-- 进度环（onMount 里 JS 模拟 0 -> 68%） -->
            <div class="relative w-14 h-14 flex-shrink-0" id="dl-ring" style="background: conic-gradient(hsl(var(--primary)) 0 0%, hsl(var(--muted)) 0% 100%); border-radius: 999px;">
              <div class="absolute inset-[5px] rounded-full bg-card flex items-center justify-center">
                <span class="text-[11px] font-semibold tabular-nums" id="dl-pct">0%</span>
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <span class="text-sm font-medium shimmer-text" id="dl-status">在装模型</span>
                <div class="churning-grid v-dots" id="dl-dots"><span></span><span></span><span></span><span></span><span></span><span></span><span></span><span></span><span></span></div>
              </div>
              <div class="text-xs text-muted-foreground mt-1">离线语音模型 · 已下载 <span id="dl-mb" class="tabular-nums">0</span> MB / 475 MB · minisign 清单已离线验签</div>
            </div>
          </div>
        </div>

        <div class="card mb-3 max-h-[220px] overflow-y-auto">
          <div class="flex flex-col">
            <!-- 当前下载中的主模型：spinner-ring 跟踪进度环 -->
            <div class="list-row flex items-center justify-between">
              <div class="flex items-center gap-2 min-w-0">
                <i data-lucide="audio-lines" class="w-4 h-4 text-muted-foreground flex-shrink-0"></i>
                <div class="min-w-0">
                  <div class="text-xs font-medium font-mono truncate">asr-streaming-paraformer-zh-en</div>
                  <div class="text-[10px] text-muted-foreground">流式识别 · int8 · 226 MB</div>
                </div>
              </div>
              <span class="flex items-center gap-2 flex-shrink-0" data-badge="paraformer">
                <span class="spinner-ring" style="width:14px;height:14px;border-width:1.5px;"></span>
                <span class="badge badge-primary">下载中 <span data-badge-pct>0</span>%</span>
              </span>
            </div>
            <!-- 其余模型：初始 shimmer-line 骨架，onMount 错峰翻成 check + 已校验 -->
            <div class="list-row flex items-center justify-between">
              <div class="flex items-center gap-2 min-w-0">
                <i data-lucide="scan-text" class="w-4 h-4 text-muted-foreground flex-shrink-0"></i>
                <div class="min-w-0">
                  <div class="text-xs font-medium font-mono truncate">asr-offline-sensevoice-zh</div>
                  <div class="text-[10px] text-muted-foreground">离线识别 · int8 · 155 MB</div>
                </div>
              </div>
              <span class="shimmer-line flex-shrink-0" data-badge="sensevoice" style="width:64px;height:18px;"></span>
            </div>
            <div class="list-row flex items-center justify-between">
              <div class="flex items-center gap-2 min-w-0">
                <i data-lucide="quote" class="w-4 h-4 text-muted-foreground flex-shrink-0"></i>
                <div class="min-w-0">
                  <div class="text-xs font-medium font-mono truncate">punc-ct-transformer-zh</div>
                  <div class="text-[10px] text-muted-foreground">标点恢复 · int8 · 62 MB</div>
                </div>
              </div>
              <span class="shimmer-line flex-shrink-0" data-badge="punc" style="width:64px;height:18px;"></span>
            </div>
            <div class="list-row flex items-center justify-between">
              <div class="flex items-center gap-2 min-w-0">
                <i data-lucide="mic-vocal" class="w-4 h-4 text-muted-foreground flex-shrink-0"></i>
                <div class="min-w-0">
                  <div class="text-xs font-medium font-mono truncate">kws-zipformer-wenet-3.3M</div>
                  <div class="text-[10px] text-muted-foreground">唤醒词 · int8 · 31 MB</div>
                </div>
              </div>
              <span class="shimmer-line flex-shrink-0" data-badge="kws" style="width:64px;height:18px;"></span>
            </div>
            <div class="list-row flex items-center justify-between">
              <div class="flex items-center gap-2 min-w-0">
                <i data-lucide="audio-lines" class="w-4 h-4 text-muted-foreground flex-shrink-0"></i>
                <div class="min-w-0">
                  <div class="text-xs font-medium font-mono truncate">vad-silero</div>
                  <div class="text-[10px] text-muted-foreground">语音活动检测 · fp32 · 0.6 MB</div>
                </div>
              </div>
              <span class="shimmer-line flex-shrink-0" data-badge="vad" style="width:64px;height:18px;"></span>
            </div>
            <div class="list-row flex items-center justify-between">
              <div class="flex items-center gap-2 min-w-0">
                <i data-lucide="speaker" class="w-4 h-4 text-muted-foreground flex-shrink-0"></i>
                <div class="min-w-0">
                  <div class="text-xs font-medium font-mono truncate">tts-matcha-zh-baker</div>
                  <div class="text-[10px] text-muted-foreground">语音合成 · fp32 · 123 MB</div>
                </div>
              </div>
              <span class="badge badge-warn flex-shrink-0"><i data-lucide="hourglass" class="w-3 h-3"></i>&nbsp;待替换</span>
            </div>
          </div>
        </div>
        <div class="text-[11px] text-muted-foreground leading-relaxed mb-1">
          哈希一律取自已验签的清单，不从镜像取；下载支持断点续传，可随时跳过，后台继续。
        </div>
      </div>

      <!-- 第 2 步：目录授权 -->
      <div class="step-panel hidden" data-step="1">
        <div class="card mb-3">
          <div class="text-sm font-medium mb-1">授权目录</div>
          <div class="text-xs text-muted-foreground mb-2">一缕只能读写你授权的目录；工作区必须落在其中（C26 拒绝符号链接外指）。</div>
          <div id="dir-list" class="flex flex-col">
            <div class="list-row flex items-center justify-between">
              <span class="font-mono text-xs truncate">C:\\Users\\swq\\Desktop</span>
              <button type="button" class="dir-remove text-muted-foreground hover:text-destructive flex-shrink-0"><i data-lucide="x" class="w-3.5 h-3.5"></i></button>
            </div>
            <div class="list-row flex items-center justify-between">
              <span class="font-mono text-xs truncate">D:\\work</span>
              <button type="button" class="dir-remove text-muted-foreground hover:text-destructive flex-shrink-0"><i data-lucide="x" class="w-3.5 h-3.5"></i></button>
            </div>
            <div class="list-row flex items-center justify-between">
              <span class="font-mono text-xs truncate">C:\\Users\\swq\\Documents</span>
              <button type="button" class="dir-remove text-muted-foreground hover:text-destructive flex-shrink-0"><i data-lucide="x" class="w-3.5 h-3.5"></i></button>
            </div>
          </div>
          <div class="flex gap-2 mt-2">
            <input id="dir-input" type="text" class="input-field flex-1" placeholder="粘贴或输入目录绝对路径…" />
            <button type="button" id="dir-add" class="btn btn-secondary btn-sm">添加</button>
          </div>
          <div class="flex flex-wrap gap-1.5 mt-2">
            <span class="text-[11px] text-muted-foreground self-center mr-1">推荐：</span>
            <button type="button" class="dir-suggest badge hover:bg-muted cursor-pointer" data-path="C:\\Users\\swq\\Downloads">下载</button>
            <button type="button" class="dir-suggest badge hover:bg-muted cursor-pointer" data-path="D:\\work\\workspace">工作区</button>
            <button type="button" class="dir-suggest badge hover:bg-muted cursor-pointer" data-path="C:\\Users\\swq\\Pictures">图片</button>
          </div>
        </div>
      </div>

      <!-- 第 3 步：API Key -->
      <div class="step-panel hidden" data-step="2">
        <div class="card mb-3">
          <div class="text-sm font-medium mb-1">LLM 凭据</div>
          <div class="text-xs text-muted-foreground mb-3">用于文本与工具调用；不填也能用离线语音，但对话能力不可用。</div>
          <div class="relative mb-2">
            <input id="key-input" type="password" class="input-field font-mono pr-9" placeholder="sk-..." value="sk-********************" />
            <button type="button" id="key-toggle" class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground">
              <i data-lucide="eye" class="w-4 h-4"></i>
            </button>
          </div>
          <div class="flex items-center gap-1.5 text-[11px] text-muted-foreground mb-1">
            <i data-lucide="lock" class="w-3.5 h-3.5"></i>
            <span>DPAPI 用户域加密存储，仅本机可读；不写入 config 明文，日志零回显（票63）</span>
          </div>
          <div class="flex items-center gap-1.5 text-[11px] text-muted-foreground">
            <i data-lucide="server" class="w-3.5 h-3.5"></i>
            <span>当前作用环境：prod（WISP_ENV 与 dev/test 隔离）</span>
          </div>
        </div>
      </div>

      <!-- 底部按钮 -->
      <div class="flex gap-2 justify-end mt-2">
        <button type="button" class="btn btn-ghost" id="skip-btn">跳过</button>
        <button type="button" class="btn btn-primary" id="next-btn">下一步</button>
      </div>
    </div>
  `,
  onMount: function (app) {
    var step = 0;
    var total = 3;
    var panels = document.querySelectorAll('.step-panel');
    var railDots = document.querySelectorAll('[data-rail]');
    var railLabels = document.querySelectorAll('[data-rail-label]');
    var railLines = document.querySelectorAll('[data-rail-line]');
    var nextBtn = document.getElementById('next-btn');
    var skipBtn = document.getElementById('skip-btn');

    function render() {
      panels.forEach(function (p) {
        p.classList.toggle('hidden', parseInt(p.dataset.step) !== step);
      });
      railDots.forEach(function (dot, i) {
        if (i < step) {
          dot.className = 'rail-dot w-6 h-6 rounded-full bg-primary text-white flex items-center justify-center text-xs font-medium flex-shrink-0';
          dot.innerHTML = '<i data-lucide="check" class="w-3.5 h-3.5"></i>';
        } else if (i === step) {
          dot.className = 'rail-dot w-6 h-6 rounded-full bg-primary text-white flex items-center justify-center text-xs font-medium flex-shrink-0';
          dot.textContent = String(i + 1);
        } else {
          dot.className = 'rail-dot w-6 h-6 rounded-full bg-muted text-muted-foreground flex items-center justify-center text-xs flex-shrink-0';
          dot.textContent = String(i + 1);
        }
      });
      railLabels.forEach(function (lb, i) {
        lb.className = 'rail-label text-[11px] mt-1.5 ' + (i <= step ? 'font-medium text-primary' : 'text-muted-foreground');
      });
      railLines.forEach(function (ln, i) {
        ln.className = 'rail-line flex-1 h-0.5 mx-2 mb-4 ' + (i < step ? 'bg-primary' : 'bg-muted');
      });
      nextBtn.textContent = step === total - 1 ? '完成' : '下一步';
      app.refreshIcons();
    }

    nextBtn.addEventListener('click', function () {
      if (step < total - 1) {
        step++;
        render();
      } else {
        App.closeFirstRun();
        app.toast('首次配置完成，进入待命态');
      }
    });
    skipBtn.addEventListener('click', function () {
      App.closeFirstRun();
      app.toast('已跳过引导，可在设置中补全');
    });

    // 目录授权：删除（事件委托，覆盖新增行）
    var dirList = document.getElementById('dir-list');
    dirList.addEventListener('click', function (e) {
      var rm = e.target.closest('.dir-remove');
      if (rm) {
        var row = rm.closest('.list-row');
        if (row) row.remove();
      }
    });
    function addDir(path) {
      if (!path) return;
      var row = document.createElement('div');
      row.className = 'list-row flex items-center justify-between';
      row.innerHTML = '<span class="font-mono text-xs truncate"></span>' +
        '<button type="button" class="dir-remove text-muted-foreground hover:text-destructive flex-shrink-0"><i data-lucide="x" class="w-3.5 h-3.5"></i></button>';
      row.querySelector('span').textContent = path;
      dirList.appendChild(row);
      app.refreshIcons();
    }
    document.getElementById('dir-add').addEventListener('click', function () {
      var inp = document.getElementById('dir-input');
      addDir(inp.value.trim());
      inp.value = '';
    });
    document.querySelectorAll('.dir-suggest').forEach(function (s) {
      s.addEventListener('click', function () { addDir(s.dataset.path); });
    });

    // API Key 掩码切换
    var keyInput = document.getElementById('key-input');
    var keyToggle = document.getElementById('key-toggle');
    keyToggle.addEventListener('click', function () {
      var show = keyInput.type === 'password';
      keyInput.type = show ? 'text' : 'password';
      keyToggle.innerHTML = '<i data-lucide="' + (show ? 'eye-off' : 'eye') + '" class="w-4 h-4"></i>';
      app.refreshIcons();
    });

    // —— Preloader 式模型下载动效：进度环 JS 模拟 0 -> 68% ——
    var dlRing = document.getElementById('dl-ring');
    var dlPct = document.getElementById('dl-pct');
    var dlMb = document.getElementById('dl-mb');
    var dlStatus = document.getElementById('dl-status');
    var paraPct = document.querySelector('[data-badge-pct]');
    var pct = 0;
    function paintDl() {
      var p = Math.round(pct);
      dlRing.style.background = 'conic-gradient(hsl(var(--primary)) 0 ' + p + '%, hsl(var(--muted)) ' + p + '% 100%)';
      dlPct.textContent = p + '%';
      dlMb.textContent = Math.round(475 * pct / 100);
      if (paraPct) paraPct.textContent = p;
      dlStatus.textContent = pct >= 45 ? '在起引擎' : '在装模型';
    }
    var dlTimer = setInterval(function () {
      pct += 0.5 + Math.random() * 1.3;
      if (pct >= 68) { pct = 68; clearInterval(dlTimer); }
      paintDl();
    }, 120);
    paintDl();

    // —— 已校验行错峰落定：shimmer-line 骨架 -> check 徽章 ——
    var DONE_BADGE = '<span class="badge badge-success flex-shrink-0"><i data-lucide="check" class="w-3 h-3"></i>&nbsp;已校验</span>';
    var settleOrder = ['sensevoice', 'punc', 'kws', 'vad'];
    settleOrder.forEach(function (key, i) {
      setTimeout(function () {
        var slot = document.querySelector('[data-badge="' + key + '"]');
        if (slot) {
          slot.outerHTML = DONE_BADGE;
          app.refreshIcons();
        }
      }, 600 + i * 700);
    });

    render();
  }
});
