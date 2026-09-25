Screens.register('privacy', {
  nav: { icon: 'lock', label: '隐私' },
  html: `
    <div class="px-8 py-8">
      <h1 class="page-title">隐私与数据</h1>
      <p class="page-subtitle">全部数据仅存于本地 SQLite，不上传任何服务器；转写中的敏感信息按尽力而为掩码，不保证完全脱敏</p>

      <div class="tabs" id="privacy-tabs">
        <div class="tab active" data-tab="l1">L1 画像</div>
        <div class="tab" data-tab="l2">L2 记忆</div>
        <div class="tab" data-tab="l3">L3 任务日志</div>
        <div class="tab" data-tab="tool">tool_call 取证</div>
        <div class="tab" data-tab="art">artifacts</div>
      </div>

      <div id="privacy-content"></div>

      <div class="grid grid-cols-2 gap-4 mt-6">
        <div class="card">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <i data-lucide="mic" class="w-4 h-4 text-muted-foreground"></i>
              <span class="text-sm font-medium">麦克风与录音</span>
            </div>
            <div class="switch" id="mic-switch"></div>
          </div>
          <p class="text-xs text-muted-foreground mt-2 mb-0" id="mic-status">当前未激活。默认不开启麦克风，唤醒词为 opt-in；音频缓冲永不落盘、不写日志。</p>
        </div>
        <div class="card">
          <div class="flex items-center gap-2 mb-2">
            <i data-lucide="clock" class="w-4 h-4 text-muted-foreground"></i>
            <span class="text-sm font-medium">记录保留期</span>
          </div>
          <div class="flex gap-2" id="retention-group">
            <button class="btn btn-secondary btn-sm retention-btn" data-days="30">30 天</button>
            <button class="btn btn-secondary btn-sm retention-btn" data-days="90">90 天</button>
            <button class="btn btn-secondary btn-sm retention-btn" data-days="0">永久</button>
          </div>
          <p class="text-xs text-muted-foreground mt-2 mb-0">作用于任务日志、tool_call 取证与转写记录；L1/L2 画像与记忆由你逐条管理。</p>
        </div>
      </div>

      <div class="card mt-4">
        <div class="flex items-center gap-2 mb-1">
          <i data-lucide="database" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm font-medium">存储位置</span>
        </div>
        <div class="list-row flex items-center gap-3">
          <i data-lucide="file-text" class="w-4 h-4 text-muted-foreground"></i>
          <span class="font-mono text-xs flex-1">%APPDATA%\wisp\wisp.db</span>
          <span class="badge">12.4 MB</span>
          <button class="btn btn-ghost btn-sm" data-open="db">打开目录</button>
        </div>
        <div class="list-row flex items-center gap-3">
          <i data-lucide="folder" class="w-4 h-4 text-muted-foreground"></i>
          <span class="font-mono text-xs flex-1">%APPDATA%\wisp\artifacts\</span>
          <span class="badge">86 MB / 500 MB</span>
          <button class="btn btn-ghost btn-sm" data-open="art">打开目录</button>
        </div>
      </div>

      <div class="flex justify-end mt-4">
        <button class="btn btn-outline" id="export-all"><i data-lucide="download"></i>导出全部数据 (JSON)</button>
      </div>
    </div>
  `,
  onMount: function (app) {
    var state = { tab: 'l1' };
    var shimmerTimer = null;

    var data = {
      l1: [
        { slot: 'pref.language', text: '偏好中文回复', src: 'extracted', at: '09-22 21:14', via: '会话 #132' },
        { slot: 'pref.tone', text: '语气简洁直接，少寒暄', src: 'extracted', at: '09-20 09:02', via: '会话 #129' },
        { slot: 'habit.project_dir', text: '常用工作目录 D:\\work\\workspace', src: 'extracted', at: '09-18 16:40', via: '会话 #127' },
        { slot: 'fact.pet', text: '家里养一只橘猫，叫年糕', src: 'manual', at: '09-15 22:31', via: '手动录入' },
        { slot: 'habit.work_hours', text: '主要在 9:00–18:00 使用', src: 'extracted', at: '09-12 10:05', via: '会话 #121' },
        { slot: 'dislikes', text: '不喜欢被反复催促确认', src: 'manual', at: '09-08 19:48', via: '手动录入' }
      ],
      l2: [
        { text: '周报模板放在 D:\\work\\templates\\weekly.docx', via: '会话 #48', hits: 3, at: '2026-09-18 22:10' },
        { text: '部署脚本要用 --no-cache，上次缓存导致旧镜像上线', via: '会话 #45', hits: 2, at: '2026-09-16 11:42' },
        { text: '年糕傍晚要喂化毛膏', via: '会话 #39', hits: 5, at: '2026-09-11 19:30' },
        { text: '周报每周五 18:00 前发给组长', via: '会话 #33', hits: 4, at: '2026-09-05 10:11' }
      ],
      l3: [
        { name: '整理本周截图归档', state: '成功', cls: 'badge-success', at: '09-24 14:28', dur: '12.4s', cost: '¥0.28' },
        { name: '读取日志并诊断报错', state: '成功', cls: 'badge-success', at: '09-24 11:03', dur: '8.8s', cost: '¥0.21' },
        { name: '随口闲聊', state: '成功', cls: 'badge-success', at: '09-24 09:47', dur: '1.1s', cost: '¥0.05' },
        { name: '查找并清理重复文件', state: '已中断', cls: 'badge-warn', at: '09-23 16:55', dur: '4.6s', cost: '¥0.06' },
        { name: '翻译英文文档段落', state: '成功', cls: 'badge-success', at: '09-23 15:12', dur: '3.2s', cost: '¥0.08' },
        { name: '重构导出模块', state: '失败', cls: 'badge-destructive', at: '09-22 20:31', dur: '26.0s', cost: '¥0.34' }
      ],
      tool: [
        { tool: 'read_file', risk: 'L0', args: '"D:\\work\\README.md"', dec: '允许', dcls: 'badge-success', out: '成功', ocls: 'badge-success', at: '14:28:03' },
        { tool: 'list_dir', risk: 'L0', args: '"D:\\work\\workspace"', dec: '允许', dcls: 'badge-success', out: '成功', ocls: 'badge-success', at: '14:28:04' },
        { tool: 'write_file', risk: 'L1', args: '"D:\\work\\out\\report.md" (1.2 KB)', dec: '会话授权', dcls: 'badge-primary', out: '成功', ocls: 'badge-success', at: '14:28:11' },
        { tool: 'bash', risk: 'L2', args: '"rm -rf node_modules"', dec: '已拒绝', dcls: 'badge-destructive', out: '已拒绝', ocls: 'badge-destructive', at: '14:29:40' },
        { tool: 'memory.save', risk: 'L1', args: '"周报模板路径…"', dec: '批量授权', dcls: 'badge-primary', out: '成功', ocls: 'badge-success', at: '22:10:05' },
        { tool: 'run_command', risk: 'L2', args: '"npm install"', dec: '允许', dcls: 'badge-success', out: '失败', ocls: 'badge-destructive', at: '11:03:52' }
      ],
      art: [
        { name: 'tool-output-8f3a.txt', size: '12 KB', at: '2026-09-24 14:28' },
        { name: 'screenshot-archive-2026w38.zip', size: '4.2 MB', at: '2026-09-24 14:31' },
        { name: 'report-draft-22c1.md', size: '8 KB', at: '2026-09-23 16:50' },
        { name: 'log-diagnosis-b7.json', size: '36 KB', at: '2026-09-23 11:04' },
        { name: 'translated-para-41en.md', size: '4 KB', at: '2026-09-23 15:13' }
      ]
    };

    function hdr(left, rightHtml) {
      return '<div class="flex justify-between items-center mb-3">' +
        '<span class="text-sm text-muted-foreground">' + left + '</span>' +
        '<div class="flex gap-2 items-center">' + rightHtml + '</div></div>';
    }
    function expBtn(label) {
      return '<button class="btn btn-outline btn-sm tab-export"><i data-lucide="download"></i>' + label + '</button>';
    }

    function shimmerSkeleton() {
      return '<div class="card p-0">' +
        '<div class="list-row"><div class="shimmer-line" style="width:26%"></div></div>' +
        '<div class="list-row"><div class="shimmer-line" style="width:72%"></div></div>' +
        '<div class="list-row"><div class="shimmer-line" style="width:52%"></div></div>' +
        '<div class="list-row"><div class="shimmer-line" style="width:38%"></div></div>' +
        '</div>';
    }

    function emptyState(icon, text) {
      return '<div class="card p-0" style="padding:44px 16px; text-align:center;">' +
        '<div style="display:flex;flex-direction:column;align-items:center;gap:10px;color:hsl(var(--muted-foreground));">' +
        '<i data-lucide="' + icon + '" style="width:32px;height:32px;stroke-width:1.5;"></i>' +
        '<div class="text-sm">' + text + '</div>' +
        '</div></div>';
    }

    function rowDelay(i) {
      return ' style="animation-delay:' + (i * 45) + 'ms"';
    }
    function monoNum(txt) {
      return '<td style="text-align:right;font-family:monospace;font-size:12px;color:hsl(var(--muted-foreground))">' + txt + '</td>';
    }

    function renderL1() {
      var h = hdr('共 ' + data.l1.length + ' 条画像 · 上限 20 条，超出按更新时间 LRU 淘汰', expBtn('导出画像'));
      h += '<div class="grid grid-cols-2 gap-3">';
      data.l1.forEach(function (r, i) {
        h += '<div class="context-card card-spotlight animated-list-item"' + rowDelay(i) + '>' +
          '<div class="flex items-center justify-between mb-1.5">' +
          '<span class="badge badge-primary font-mono" style="font-size:10px;font-weight:500">' + r.slot + '</span>' +
          '<span class="ctx-source">' + (r.src === 'manual' ? '手动' : '提取') + ' · ' + r.at + '</span>' +
          '</div>' +
          '<div class="text-sm leading-snug">' + r.text + '</div>' +
          '<div class="flex items-center justify-between mt-2.5 pt-2" style="border-top:1px solid hsl(var(--border)/0.6)">' +
          '<span class="text-[11px] text-muted-foreground">来源 ' + r.via + '</span>' +
          '<button class="btn btn-ghost btn-sm text-destructive ctx-del" title="删除"><i data-lucide="trash-2" class="w-3.5 h-3.5"></i></button>' +
          '</div>' +
          '</div>';
      });
      h += '</div>';
      return h;
    }

    function renderL2() {
      var h = hdr('共 ' + data.l2.length + ' 条显式记忆 · 不自动注入 prompt，按关键词+最近命中召回',
        '<label class="flex items-center gap-2 text-xs text-muted-foreground"><input type="checkbox" id="sel-all" style="accent-color:hsl(var(--primary))">全选</label>' +
        '<button class="btn btn-destructive btn-sm" id="batch-del"><i data-lucide="trash-2"></i>批量删除</button>' +
        expBtn('导出记忆'));
      h += '<div class="grid grid-cols-2 gap-3">';
      data.l2.forEach(function (r, i) {
        h += '<div class="context-card card-spotlight animated-list-item"' + rowDelay(i) + '>' +
          '<div class="flex items-start gap-2.5">' +
          '<input type="checkbox" class="row-check" style="accent-color:hsl(var(--primary));margin-top:3px">' +
          '<div class="flex-1 min-w-0">' +
          '<div class="flex items-start justify-between gap-2">' +
          '<div class="text-sm leading-snug">' + r.text + '</div>' +
          '<button class="btn btn-ghost btn-sm text-destructive ctx-del" title="删除"><i data-lucide="trash-2" class="w-3.5 h-3.5"></i></button>' +
          '</div>' +
          '<div class="flex items-center justify-between mt-2 pt-2" style="border-top:1px solid hsl(var(--border)/0.6)">' +
          '<span class="ctx-source">' + r.via + ' · 命中 ' + r.hits + ' 次</span>' +
          '<span class="text-[11px] text-muted-foreground">' + r.at + '</span>' +
          '</div>' +
          '</div>' +
          '</div>' +
          '</div>';
      });
      h += '</div>';
      return h;
    }

    function renderL3() {
      if (data.l3.length === 0) {
        return hdr('共 0 条任务 · 按保留期自动清理，query 已脱敏', expBtn('导出日志')) + emptyState('inbox', '暂无任务日志');
      }
      var h = hdr('共 ' + data.l3.length + ' 条任务 · 按保留期自动清理，query 已脱敏', expBtn('导出日志'));
      h += '<div class="card p-0 overflow-hidden">' +
        '<table class="records-table"><thead><tr>' +
        '<th>任务</th><th>时间</th><th>状态</th><th style="text-align:right">耗时</th><th style="text-align:right">费用</th><th style="width:44px"></th>' +
        '</tr></thead><tbody>';
      data.l3.forEach(function (r, i) {
        h += '<tr class="animated-list-item" data-idx="' + i + '"' + rowDelay(i) + '>' +
          '<td><div class="flex items-center gap-2"><i data-lucide="file-text" class="w-3.5 h-3.5 text-muted-foreground"></i>' + r.name + '</div></td>' +
          '<td style="font-size:12px;color:hsl(var(--muted-foreground))">' + r.at + '</td>' +
          '<td><span class="badge ' + r.cls + '">' + r.state + '</span></td>' +
          monoNum(r.dur) + monoNum(r.cost) +
          '<td style="text-align:right"><button class="btn btn-ghost btn-sm text-destructive l3-del" title="删除"><i data-lucide="trash-2" class="w-3.5 h-3.5"></i></button></td>' +
          '</tr>';
      });
      h += '</tbody></table></div>';
      return h;
    }

    function renderTool() {
      var h = hdr('取证记录 · 每条工具调用的参数摘要与审批决定（长参数已截断）', expBtn('导出取证'));
      h += '<div class="card p-0 overflow-hidden">' +
        '<table class="records-table"><thead><tr>' +
        '<th>工具</th><th>风险</th><th>参数</th><th>审批决定</th><th>结果</th><th>时间</th>' +
        '</tr></thead><tbody>';
      data.tool.forEach(function (r, i) {
        h += '<tr class="animated-list-item"' + rowDelay(i) + '>' +
          '<td><span class="tool-chip">' + r.tool + '</span></td>' +
          '<td><span class="badge">' + r.risk + '</span></td>' +
          '<td style="font-family:monospace;font-size:12px;color:hsl(var(--muted-foreground))">' + r.args + '</td>' +
          '<td><span class="badge ' + r.dcls + '">' + r.dec + '</span></td>' +
          '<td><span class="badge ' + r.ocls + '">' + r.out + '</span></td>' +
          '<td style="font-family:monospace;font-size:12px;color:hsl(var(--muted-foreground))">' + r.at + '</td>' +
          '</tr>';
      });
      h += '</tbody></table></div>';
      return h;
    }

    function renderArt() {
      var h = hdr('生成文件 · 目录配额 500 MB，超出按 LRU 自动清理', expBtn('导出清单'));
      h += '<div class="card p-0 overflow-hidden">' +
        '<table class="records-table"><thead><tr>' +
        '<th>文件</th><th>大小</th><th>时间</th><th>路径</th><th style="width:80px"></th>' +
        '</tr></thead><tbody>';
      data.art.forEach(function (r, i) {
        h += '<tr class="animated-list-item"' + rowDelay(i) + '>' +
          '<td><div class="flex items-center gap-2"><i data-lucide="file" class="w-3.5 h-3.5 text-muted-foreground"></i><span style="font-family:monospace;font-size:12px">' + r.name + '</span></div></td>' +
          '<td><span class="badge">' + r.size + '</span></td>' +
          '<td style="font-size:12px;color:hsl(var(--muted-foreground))">' + r.at + '</td>' +
          '<td style="font-family:monospace;font-size:11px;color:hsl(var(--muted-foreground))">%APPDATA%\\wisp\\artifacts\\</td>' +
          '<td style="text-align:right"><button class="btn btn-outline btn-sm row-open">打开</button></td>' +
          '</tr>';
      });
      h += '</tbody></table></div>';
      return h;
    }

    function render() {
      var c = document.getElementById('privacy-content');
      var h = '';
      if (state.tab === 'l1') h = renderL1();
      else if (state.tab === 'l2') h = renderL2();
      else if (state.tab === 'l3') h = renderL3();
      else if (state.tab === 'tool') h = renderTool();
      else if (state.tab === 'art') h = renderArt();
      c.innerHTML = h;
      bindRows();
      app.refreshIcons();
    }

    function showShimmerThenRender() {
      var c = document.getElementById('privacy-content');
      c.innerHTML = shimmerSkeleton();
      app.refreshIcons();
      if (shimmerTimer) clearTimeout(shimmerTimer);
      shimmerTimer = setTimeout(function () { render(); }, 300);
    }

    function bindRows() {
      document.querySelectorAll('.ctx-del').forEach(function (b) {
        b.onclick = function () {
          b.closest('.context-card').remove();
          app.toast('已删除 1 条记录');
        };
      });
      document.querySelectorAll('.l3-del').forEach(function (b) {
        b.onclick = function () {
          var row = b.closest('tr');
          var idx = parseInt(row.dataset.idx);
          if (!isNaN(idx) && data.l3[idx]) data.l3.splice(idx, 1);
          app.toast('已删除 1 条任务日志');
          render();
        };
      });
      document.querySelectorAll('.row-open').forEach(function (b) {
        b.onclick = function () { app.toast('已打开 %APPDATA%\\wisp\\artifacts\\ 目录'); };
      });
      document.querySelectorAll('.tab-export').forEach(function (b) {
        b.onclick = function () { app.toast('导出完成：该视图 JSON（不含密钥与转写原文）'); };
      });
      var selAll = document.getElementById('sel-all');
      if (selAll) selAll.onchange = function () {
        document.querySelectorAll('.row-check').forEach(function (cb) { cb.checked = selAll.checked; });
      };
      var batch = document.getElementById('batch-del');
      if (batch) batch.onclick = function () {
        var n = 0;
        document.querySelectorAll('.row-check:checked').forEach(function (cb) {
          cb.closest('.context-card').remove(); n++;
        });
        app.toast(n ? '已删除 ' + n + ' 条记忆' : '未勾选任何条目');
      };
    }

    // tab 切换
    document.querySelectorAll('#privacy-tabs .tab').forEach(function (t) {
      t.onclick = function () {
        document.querySelectorAll('#privacy-tabs .tab').forEach(function (x) { x.classList.remove('active'); });
        t.classList.add('active');
        state.tab = t.dataset.tab;
        showShimmerThenRender();
      };
    });

    // 麦克风开关
    var mic = document.getElementById('mic-switch');
    mic.onclick = function () {
      mic.classList.toggle('on');
      var on = mic.classList.contains('on');
      document.getElementById('mic-status').textContent = on
        ? '正在监听唤醒词，悬浮球常亮指示中。半双工：TTS 播报期间麦克风自动关闭；可一键静音。'
        : '当前未激活。默认不开启麦克风，唤醒词为 opt-in；音频缓冲永不落盘、不写日志。';
    };

    // 保留期
    document.querySelectorAll('.retention-btn').forEach(function (b) {
      b.onclick = function () {
        document.querySelectorAll('.retention-btn').forEach(function (x) {
          x.classList.remove('btn-primary'); x.classList.add('btn-secondary');
        });
        b.classList.remove('btn-secondary'); b.classList.add('btn-primary');
        app.toast('记录保留期已设为 ' + (b.dataset.days === '0' ? '永久' : b.dataset.days + ' 天'));
      };
    });

    // 打开目录
    document.querySelectorAll('[data-open]').forEach(function (b) {
      b.onclick = function () {
        app.toast(b.dataset.open === 'db'
          ? '已打开 %APPDATA%\\wisp\\ 目录（wisp.db 在此）'
          : '已打开 %APPDATA%\\wisp\\artifacts\\ 目录');
      };
    });

    // 全局导出
    document.getElementById('export-all').onclick = function () {
      app.toast('导出完成：wisp-export-2026-09-24.json（画像/记忆/日志/取证，不含密钥与转写原文）');
    };

    render();
  }
});
