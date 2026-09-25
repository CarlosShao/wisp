Screens.register('approval', {
  nav: { icon: 'shield-check', label: '审批' },
  html: `
    <h1 class="page-title">审批中心</h1>
    <p class="page-subtitle">不可逆与高危操作在原生侧批准，面板仅可拒绝与查看完整参数</p>

    <div class="flex items-center justify-between mb-3 mt-5">
      <div class="flex items-center gap-2 flex-wrap">
        <span class="section-title mb-0">待审批队列</span>
        <span class="badge badge-primary" id="queue-depth">9 / 10</span>
        <span class="badge badge-warn">接近容量上限</span>
      </div>
      <button class="btn btn-outline btn-sm" id="batch-reject" style="display:none">
        <i data-lucide="ban" class="w-3.5 h-3.5"></i>
        <span>批量拒绝 <span id="batch-n">0</span> 条</span>
      </button>
    </div>

    <div class="flex items-start gap-2 rounded-lg border border-amber-200 bg-amber-50 dark:bg-amber-950/20 dark:border-amber-800/50 px-3 py-2 mb-4 text-xs text-amber-800 dark:text-amber-300">
      <i data-lucide="alert-triangle" class="w-4 h-4 flex-shrink-0 mt-px"></i>
      <span>队列容量上限 10 条：满额后新的待确认请求按 fail-closed 自动拒绝，绝不默认放行；宿主不可达或队列出错同样判拒绝。</span>
    </div>

    <div class="grid gap-4" style="grid-template-columns: 248px 1fr; align-items:start">
      <div class="card !p-1.5" id="queue-list"></div>
      <div id="detail-slot"></div>
    </div>
  `,
  onMount: function (app) {
    var QUEUE = [
      {
        id: 'a1', level: 'L2', risk: 'Deny', icon: 'trash-2', tool: 'fs.delete', title: '删除桌面 3 张旧截图',
        corr: 'corr-7f3a91', task: 'task-6', age: '刚刚', countdown: '剩余 268s 后自动拒绝（已进入 270s 警告态）',
        rules: [
          { code: 'R3', text: '批量删除 3 个及以上文件' },
          { code: 'R8', text: '不可逆：永久删除不进回收站' }
        ],
        reasonKnown: true,
        reason: '该操作将永久删除 3 个文件，不进回收站，无法撤销。',
        params: {
          paths: [
            'C:\\Users\\swq\\Desktop\\shot-01.png',
            'C:\\Users\\swq\\Desktop\\shot-02.png',
            'C:\\Users\\swq\\Desktop\\shot-03.png'
          ],
          cwd: 'C:\\Users\\swq\\Desktop'
        },
        rec: { conf: 78, text: '检测到批量删除操作，建议先移动到待确认文件夹而非直接删除', alt1: '移动到待确认', alt2: '直接删除' }
      },
      {
        id: 'a2', level: 'L2', risk: 'Deny', icon: 'terminal', tool: 'shell.exec', title: '执行磁盘清理脚本 clean.ps1',
        corr: 'corr-2c8e47', task: 'task-3', age: '1 分钟前', countdown: '剩余 291s 后自动拒绝',
        rules: [
          { code: 'R6', text: 'shell argv 含管道与重定向' },
          { code: 'R1', text: '工具声明风险仅为下界' }
        ],
        reasonKnown: true,
        reason: '脚本包含管道、重定向与元字符，并写入系统临时目录，需逐条确认，不可聚合。',
        params: {
          cmd: 'powershell -NoProfile -File clean.ps1 -Target C:\\Windows\\Temp | Out-File clean.log',
          cwd: 'C:\\Users\\swq\\Scripts'
        },
        rec: { conf: 64, text: '脚本写入系统临时目录，建议先在沙箱目录试运行一次再全量执行。', alt1: '先沙箱试运行', alt2: '直接全量执行' }
      },
      {
        id: 'a3', level: 'L1', risk: 'Caution', icon: 'folder-tree', tool: 'fs.write', title: '批量归类 12 个下载文件',
        corr: 'corr-9d16b0', task: 'task-2', age: '刚刚',
        rules: [
          { code: 'D45-1', text: '单次调用 12 个同质 L1 操作已聚合' }
        ],
        reasonKnown: true,
        reason: '可逆写操作：移动 / 归类下载目录中的 12 个文件，一次确认覆盖整批；否决则整批取消。',
        params: {
          paths: [
            'C:\\Users\\swq\\Downloads\\report-final.pdf',
            'C:\\Users\\swq\\Downloads\\photo-1204.jpg'
          ],
          more: '… 另 10 个文件归入 图片 / 文档 / 安装包 三个子目录',
          cwd: 'C:\\Users\\swq\\Downloads'
        },
        rec: { conf: 82, text: '12 个文件已按类型聚合，建议确认归类规则后再一次性移动。', alt1: '按建议归类', alt2: '逐个人工移动' }
      },
      {
        id: 'a4', level: 'L2', risk: 'Deny', icon: 'file-key', tool: 'fs.write', title: '写入项目 .env 凭据文件',
        corr: 'corr-5a7f2c', task: 'task-6', age: '3 分钟前', countdown: '剩余 245s 后自动拒绝',
        rules: [
          { code: 'R3', text: '命中 B 档敏感路径 .env*（单文件豁免）' },
          { code: 'R8', text: '覆盖已存在文件' }
        ],
        reasonKnown: true,
        reason: '目标命中 B 档敏感路径 .env*，默认拒绝，需单文件豁免级确认；且将覆盖已存在文件。',
        params: {
          paths: [
            'C:\\Users\\swq\\dev\\wisp\\.env'
          ],
          cwd: 'C:\\Users\\swq\\dev\\wisp'
        },
        rec: { conf: 55, text: '覆盖 .env 会清掉现有凭据，建议先备份再写入。', alt1: '先备份 .env', alt2: '直接覆盖写入' }
      },
      {
        id: 'a6', level: 'L1', risk: 'Caution', icon: 'pencil-line', tool: 'fs.rename', title: '重命名本周录屏文件',
        corr: 'corr-6e3a19', task: 'task-1', age: '5 分钟前',
        rules: [
          { code: 'L1', text: '可逆写：重命名可手动还原' }
        ],
        reasonKnown: true,
        reason: '可逆操作：重命名单个文件，不删除或改动内容。',
        params: {
          paths: [
            '录制 2026-09-24.mp4  ->  周会-演示录屏.mp4'
          ],
          cwd: 'C:\\Users\\swq\\Videos'
        },
        rec: { conf: 90, text: '重命名可逆且无副作用，建议直接执行，无需额外处置。', alt1: '直接重命名', alt2: '取消操作' }
      }
    ];

    var Q_PREVIEWS = [
      '第 1/3 项 · 操作参数：核对工具名、路径 / 命令与工作目录',
      '第 2/3 项 · 风险原因：原生风险评估为何要求人工确认',
      '第 3/3 项 · 调用链：本调用来自哪个任务、命中了哪条规则'
    ];

    var selectedId = QUEUE[0].id;
    var checked = {};
    var expanded = {};
    var qIndex = 0;

    var listEl = document.getElementById('queue-list');
    var detailEl = document.getElementById('detail-slot');

    /* LevelBadge —— 对照 l2-approval-card：L2 红 / L1 青（C21 token） */
    function headerBadge(level) {
      if (level === 'L1') {
        return '<span class="badge badge-primary" style="font-size:10.5px">'
          + '<i data-lucide="shield-check" class="w-3 h-3"></i>'
          + '<span>L1 轻确认</span>'
          + '</span>';
      }
      return '<span class="badge badge-destructive" style="font-size:10.5px">'
        + '<i data-lucide="shield-alert" class="w-3 h-3"></i>'
        + '<span class="uppercase tracking-wide">L2 强确认</span>'
        + '</span>';
    }

    /* 队列行内紧凑风险 badge */
    function rowBadge(it) {
      if (it.level === 'L1') {
        return '<span class="badge badge-primary" style="font-size:10px; height:18px">L1</span>';
      }
      var cls = it.risk === 'Deny' ? 'badge-destructive' : 'badge-warn';
      return '<span class="badge ' + cls + '" style="font-size:10px; height:18px">L2 · ' + it.risk + '</span>';
    }

    /* RuleList outline badge —— 只渲染规则 id，hover 看原文 */
    function ruleBadges(rules) {
      if (!rules || rules.length === 0) {
        return '<span class="text-[11.5px] text-muted-foreground">no rule id recorded</span>';
      }
      return '<div class="flex flex-wrap gap-1" aria-label="matched risk rules">'
        + rules.map(function (r) {
            return '<span class="badge" title="' + r.text + '" style="font-family:monospace; font-size:10.5px; height:20px;">'
              + r.code + '</span>';
          }).join('')
        + '</div>';
    }

    function renderList() {
      var html = '';
      QUEUE.forEach(function (it, i) {
        var sel = it.id === selectedId;
        html += ''
          + '<div class="list-row animated-list-item" data-id="' + it.id + '" style="cursor:pointer;'
          +   'animation-delay:' + (i * 60) + 'ms;'
          +   'border:1px solid ' + (sel ? 'hsl(var(--primary)/0.45)' : 'transparent') + ';'
          +   'background:' + (sel ? 'hsl(var(--accent)/0.5)' : 'transparent') + '">'
          +   '<label class="flex items-center" onclick="event.stopPropagation()">'
          +     '<input type="checkbox" data-check="' + it.id + '"' + (checked[it.id] ? ' checked' : '')
          +       ' class="w-3.5 h-3.5" style="accent-color:hsl(var(--primary))">'
          +   '</label>'
          +   '<div class="flex-1 min-w-0">'
          +     '<div class="flex items-center gap-1.5">' + rowBadge(it) + '</div>'
          +     '<div class="text-sm font-medium truncate mt-1.5">' + it.title + '</div>'
          +     '<div class="text-[11px] text-muted-foreground mt-0.5 font-mono">' + it.corr + ' · ' + it.age + '</div>'
          +   '</div>'
          + '</div>';
      });
      html += '<div class="text-[11px] text-muted-foreground px-2 pt-2 pb-1 leading-relaxed">队列共 9 条待确认 · 悬浮球角标实时显示队头与深度</div>';
      listEl.innerHTML = html;

      listEl.querySelectorAll('.list-row').forEach(function (row) {
        row.onclick = function () {
          selectedId = row.dataset.id;
          qIndex = 0;
          renderList();
          renderDetail();
        };
      });
      listEl.querySelectorAll('input[data-check]').forEach(function (cb) {
        cb.onchange = function () {
          checked[cb.dataset.check] = cb.checked;
          updateBatch();
        };
      });
    }

    /* 完整参数区：mono 逐行显示 路径 / 命令 / 工作目录 */
    function renderParams(it, isExp) {
      if (it.params.cmd) {
        return '<div class="rounded-lg border border-border bg-muted/40 px-3 py-2.5 overflow-x-auto">'
          + '<pre class="font-mono text-[11.5px] leading-relaxed whitespace-pre text-foreground m-0">'
          + 'shell.exec\n' + it.params.cmd + '</pre>'
          + '<div class="font-mono text-[11px] text-muted-foreground mt-2">工作目录  ' + it.params.cwd + '</div>'
          + '</div>';
      }
      var lines = it.params.paths.join('\n');
      if (!isExp && it.params.paths.length > 1) {
        lines = it.params.paths.slice(0, 1).join('\n');
      }
      var tail = '';
      if (!isExp && it.params.more) tail = '\n' + it.params.more;
      else if (!isExp && it.params.paths.length > 1) tail = '\n… 另 ' + (it.params.paths.length - 1) + ' 项，展开查看完整参数';
      return '<div class="rounded-lg border border-border bg-muted/40 px-3 py-2.5 overflow-x-auto">'
        + '<pre class="font-mono text-[11.5px] leading-relaxed whitespace-pre text-foreground m-0">'
        + it.tool + '\n' + lines + tail + '</pre>'
        + '<div class="font-mono text-[11px] text-muted-foreground mt-2">工作目录  ' + it.params.cwd + '</div>'
        + '</div>';
    }

    function renderDetail() {
      var it = QUEUE.filter(function (q) { return q.id === selectedId; })[0];
      var isL2 = it.level === 'L2';
      var isExp = !!expanded[it.id];

      var topBar = isL2
        ? '<div class="absolute top-0 left-0 right-0 h-1" style="background:hsl(var(--destructive))"></div>'
        : '<div class="absolute top-0 left-0 right-0 h-1" style="background:hsl(var(--primary))"></div>';

      /* 原因分区：无原因文本时 fail-closed 文案，对照 ReasonLine */
      var reasonBlock;
      if (!it.reasonKnown || !it.reason || it.reason.trim() === '') {
        reasonBlock = '<p class="text-[12.5px] text-muted-foreground leading-relaxed mb-2">'
          + '<span class="font-medium" style="color:#b8860b">信息不足：</span>'
          + '风险评估未给出原因文本，本卡不能推断为"无风险"。</p>';
      } else {
        reasonBlock = '<p class="text-[12.5px] text-muted-foreground leading-relaxed mb-2">' + it.reason + '</p>';
      }

      /* 倒计时警告条：L2 300s/270s 警告 · L1 2-3s 进度条 */
      var midBlock;
      if (isL2) {
        midBlock = ''
          + '<div class="rounded-lg border border-amber-200 bg-amber-50 dark:bg-amber-950/20 dark:border-amber-800/50 px-3 py-2 flex items-center gap-2 text-xs text-amber-800 dark:text-amber-300 mt-3 mb-3">'
          +   '<i data-lucide="timer" class="w-4 h-4 flex-shrink-0"></i>'
          +   '<span>' + it.countdown + ' · 超时一律判拒绝；拒绝后任务继续运行，可一键重放该调用。</span>'
          + '</div>';
      } else {
        midBlock = ''
          + '<div class="rounded-lg border px-3 py-2.5 mt-3 mb-3" style="border-color:hsl(var(--primary)/0.35); background:hsl(var(--accent)/0.35)">'
          +   '<div class="flex items-center gap-2 text-sm font-medium mb-2">'
          +     '<i data-lucide="hourglass" class="w-4 h-4"></i>'
          +     '<span>执行前阻止窗口 2.3s · 窗口结束自动执行</span>'
          +   '</div>'
          +   '<div class="progress-bar"><div class="progress-fill" style="width:64%"></div></div>'
          +   '<div class="text-[11px] text-muted-foreground mt-2">否决通道：单击悬浮球 · <kbd class="kbd">Esc</kbd> · 语音否决词 · 面板取消（语音引擎未加载时显示「语音取消不可用」，不假装可用）。</div>'
          + '</div>';
      }

      var rejectLabel = isL2 ? '拒绝' : '取消执行';

      /* Approval Card 本体 —— 头部：工具图标 + mono 工具名 + LevelBadge + 分页点 */
      var approvalCard = ''
        + '<div class="card card-spotlight relative overflow-hidden animated-list-item" data-role="approval-card" style="animation-delay:0ms">'
        +   topBar
        +   '<div class="flex items-start justify-between gap-3">'
        +     '<div class="flex items-center gap-2.5 min-w-0">'
        +       '<div class="flex items-center justify-center w-8 h-8 rounded-lg flex-shrink-0" style="background:hsl(var(--muted)/0.6)">'
        +         '<i data-lucide="' + it.icon + '" class="w-4 h-4 text-muted-foreground"></i>'
        +       '</div>'
        +       '<div class="min-w-0">'
        +         '<div class="text-[13px] font-medium text-foreground">需要你的确认</div>'
        +         '<div class="text-[11.5px] text-muted-foreground font-mono mt-0.5">' + it.tool + '</div>'
        +       '</div>'
        +     '</div>'
        +     '<div class="flex items-center gap-2 flex-shrink-0">'
        +       headerBadge(it.level)
        +       '<div class="pagination-dots" data-role="pager" aria-label="问题分页">'
        +         '<span data-q="0" class="' + (qIndex === 0 ? 'active' : '') + '"></span>'
        +         '<span data-q="1" class="' + (qIndex === 1 ? 'active' : '') + '"></span>'
        +         '<span data-q="2" class="' + (qIndex === 2 ? 'active' : '') + '"></span>'
        +       '</div>'
        +     '</div>'
        +   '</div>'
        +   '<div class="text-[11.5px] text-muted-foreground mt-2.5 mb-3 flex items-center gap-1.5">'
        +     '<i data-lucide="file-search" class="w-3.5 h-3.5 flex-shrink-0"></i>'
        +     '<span>' + Q_PREVIEWS[qIndex] + '</span>'
        +   '</div>'
        +   '<div class="section-title mt-1">完整参数</div>'
        +   renderParams(it, isExp)
        +   '<div class="section-title mt-4">为什么需要确认</div>'
        +   reasonBlock
        +   ruleBadges(it.rules)
        +   '<div class="text-[11.5px] text-muted-foreground mt-3 font-mono">调用链：' + it.corr + ' &gt; ' + it.task + '</div>'
        +   midBlock
        +   '<div class="flex items-start gap-1.5 mt-4 text-[11px] text-muted-foreground leading-relaxed">'
        +     '<i data-lucide="info" class="w-3.5 h-3.5 flex-shrink-0 mt-px"></i>'
        +     '<span>批准权在原生侧：点击悬浮球或按 <kbd class="kbd">F2</kbd> 以批准 · 面板不提供允许按钮（面板来源的允许在服务端被结构拒绝）。</span>'
        +   '</div>'
        +   '<div class="flex items-center gap-2 mt-3 pt-3" style="border-top:1px solid hsl(var(--border))">'
        +     '<span class="mr-auto text-[11px] text-muted-foreground">判定来自原生风险评估</span>'
        +     '<button class="btn btn-destructive btn-sm" data-act="reject">' + rejectLabel + '</button>'
        +     '<button class="btn btn-outline btn-sm" data-act="toggle">' + (isExp ? '收起参数' : '查看完整参数') + '</button>'
        +   '</div>'
        + '</div>';

      /* Recommendation Card —— 队列详情末尾 */
      var rec = it.rec;
      var recCard = ''
        + '<div class="card card-spotlight animated-list-item" data-role="rec-card" style="animation-delay:60ms">'
        +   '<div class="flex items-center justify-between mb-2.5">'
        +     '<div class="flex items-center gap-1.5">'
        +       '<i data-lucide="lightbulb" class="w-4 h-4" style="color:hsl(var(--primary))"></i>'
        +       '<span class="section-title mb-0">智能建议</span>'
        +     '</div>'
        +     '<span class="text-[11px] font-mono text-muted-foreground">置信度 ' + rec.conf + '%</span>'
        +   '</div>'
        +   '<div class="confidence-bar mb-2.5"><div class="conf-fill" style="width:' + rec.conf + '%"></div></div>'
        +   '<p class="text-[12.5px] text-foreground leading-relaxed mb-3">' + rec.text + '</p>'
        +   '<div class="flex items-center gap-2 flex-wrap">'
        +     '<button class="btn btn-outline btn-sm" data-rec="alt1">' + rec.alt1 + '</button>'
        +     '<button class="btn btn-outline btn-sm" data-rec="alt2">' + rec.alt2 + '</button>'
        +     '<button class="btn btn-primary btn-sm ml-auto" data-rec="accept">'
        +       '<i data-lucide="check" class="w-3.5 h-3.5"></i>'
        +       '<span>接受建议</span>'
        +     '</button>'
        +   '</div>'
        + '</div>';

      detailEl.innerHTML = '<div class="flex flex-col gap-3">' + approvalCard + recCard + '</div>';

      app.refreshIcons();

      detailEl.querySelector('[data-act="reject"]').onclick = function () {
        app.toast(isL2 ? '已拒绝该操作，任务继续运行，可一键重放' : '已在阻止窗口内取消，操作未执行');
      };
      detailEl.querySelector('[data-act="toggle"]').onclick = function () {
        expanded[it.id] = !expanded[it.id];
        renderDetail();
      };

      detailEl.querySelectorAll('[data-role="pager"] span').forEach(function (s) {
        s.style.cursor = 'pointer';
        s.onclick = function () {
          qIndex = parseInt(s.dataset.q, 10);
          renderDetail();
        };
      });

      detailEl.querySelectorAll('[data-rec]').forEach(function (b) {
        b.onclick = function () {
          var k = b.dataset.rec;
          if (k === 'accept') app.toast('已采纳智能建议，将按建议方案处理');
          else if (k === 'alt1') app.toast('已选择备选：' + rec.alt1);
          else app.toast('已选择备选：' + rec.alt2);
        };
      });

      detailEl.querySelectorAll('.card-spotlight').forEach(function (card) {
        card.onmousemove = function (e) {
          var r = card.getBoundingClientRect();
          card.style.setProperty('--mx', (e.clientX - r.left) + 'px');
          card.style.setProperty('--my', (e.clientY - r.top) + 'px');
        };
      });
    }

    function updateBatch() {
      var n = Object.keys(checked).filter(function (k) { return checked[k]; }).length;
      document.getElementById('batch-n').textContent = n;
      document.getElementById('batch-reject').style.display = n > 0 ? '' : 'none';
    }

    document.getElementById('batch-reject').onclick = function () {
      var n = Object.keys(checked).filter(function (k) { return checked[k]; }).length;
      app.toast('已批量拒绝 ' + n + ' 条待确认操作');
      checked = {};
      renderList();
      updateBatch();
    };

    renderList();
    renderDetail();
  }
});
