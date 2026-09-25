Screens.register('tasks', {
  nav: { icon: 'list-checks', label: '任务' },
  html: `
    <h1 class="page-title">任务</h1>
    <p class="page-subtitle">并发执行、路径锁与定时调度 · 审批等待不阻塞其他任务的非冲突工作</p>

    <!-- Filter Table 状态 chip 行 -->
    <div class="flex items-center justify-between mb-5 mt-5 flex-wrap gap-3">
      <div class="filter-chips" id="task-filters">
        <button class="filter-chip active" data-filter="all">全部</button>
        <button class="filter-chip" data-filter="running">运行中</button>
        <button class="filter-chip" data-filter="approval">等待</button>
        <button class="filter-chip" data-filter="blocked">阻塞</button>
        <button class="filter-chip" data-filter="queued">排队</button>
        <button class="filter-chip" data-filter="done">完成</button>
        <button class="filter-chip" data-filter="failed">失败</button>
      </div>
      <button class="btn btn-primary btn-sm" id="open-task-modal">
        <i data-lucide="plus" class="w-3.5 h-3.5"></i>新建任务
      </button>
    </div>

    <!-- 任务行列表：初始为 shimmer 骨架，onMount 0.8s 后替换为真实行 -->
    <div class="card card-spotlight !p-2 mb-8" id="running-tasks">
      <div class="list-row" style="opacity:0.7">
        <span class="shimmer-block" style="width:24px;height:24px;border-radius:50%;flex-shrink:0"></span>
        <div class="flex-1" style="display:flex;flex-direction:column;gap:7px">
          <div class="shimmer-line" style="width:42%"></div>
          <div class="shimmer-line" style="width:68%;height:10px"></div>
        </div>
      </div>
      <div class="list-row" style="opacity:0.7">
        <span class="shimmer-block" style="width:24px;height:24px;border-radius:50%;flex-shrink:0"></span>
        <div class="flex-1" style="display:flex;flex-direction:column;gap:7px">
          <div class="shimmer-line" style="width:34%"></div>
          <div class="shimmer-line" style="width:60%;height:10px"></div>
        </div>
      </div>
      <div class="list-row" style="opacity:0.7">
        <span class="shimmer-block" style="width:24px;height:24px;border-radius:50%;flex-shrink:0"></span>
        <div class="flex-1" style="display:flex;flex-direction:column;gap:7px">
          <div class="shimmer-line" style="width:48%"></div>
          <div class="shimmer-line" style="width:55%;height:10px"></div>
        </div>
      </div>
      <div class="list-row" style="opacity:0.7">
        <span class="shimmer-block" style="width:24px;height:24px;border-radius:50%;flex-shrink:0"></span>
        <div class="flex-1" style="display:flex;flex-direction:column;gap:7px">
          <div class="shimmer-line" style="width:30%"></div>
          <div class="shimmer-line" style="width:64%;height:10px"></div>
        </div>
      </div>
    </div>

    <span class="section-title">定时任务（交给系统调度）</span>
    <div class="card !p-2 mb-8">
      <table class="records-table">
        <thead>
          <tr><th>任务名</th><th>触发</th><th>电源条件</th><th>工作目录</th><th>最近结果</th></tr>
        </thead>
        <tbody>
          <tr>
            <td class="font-medium">整理下载目录</td>
            <td>每日 09:00</td>
            <td>接通电源时</td>
            <td class="font-mono text-xs">C:\\Users\\swq\\Downloads</td>
            <td><span class="badge badge-success">成功</span></td>
          </tr>
          <tr>
            <td class="font-medium">生成周报草稿</td>
            <td>每周一 10:00</td>
            <td>解锁时</td>
            <td class="font-mono text-xs">C:\\Users\\swq\\Documents</td>
            <td><span class="badge badge-success">成功</span></td>
          </tr>
          <tr>
            <td class="font-medium">清理临时文件</td>
            <td>每月 1 日 03:00</td>
            <td>空闲时</td>
            <td class="font-mono text-xs">%TEMP%</td>
            <td><span class="badge badge-warn">失败</span></td>
          </tr>
          <tr>
            <td class="font-medium">归档项目快照</td>
            <td>每周五 18:00</td>
            <td>接通电源时</td>
            <td class="font-mono text-xs">D:\\work\\projects</td>
            <td><span class="badge">未运行</span></td>
          </tr>
        </tbody>
      </table>
      <div class="text-[11px] text-muted-foreground px-2 pt-2 pb-1">定时任务通过系统任务计划调用 CLI 入口执行，无人值守；电源条件与触发时间由系统调度器保证。</div>
    </div>

    <span class="section-title">定时提醒（进程内一次性）</span>
    <div class="card !p-2" id="reminder-list">
      <div class="list-row" style="border-bottom:1px solid hsl(var(--border)/0.5)">
        <i data-lucide="bell" class="w-4 h-4 text-muted-foreground"></i>
        <div class="flex-1">
          <div class="text-sm">今天 15:00 例会前准备好演示资料</div>
          <div class="text-[11px] text-muted-foreground">一次性 · 到点仅通知，不执行操作</div>
        </div>
        <div class="switch on" data-switch></div>
      </div>
      <div class="list-row" style="border-bottom:1px solid hsl(var(--border)/0.5)">
        <i data-lucide="bell" class="w-4 h-4 text-muted-foreground"></i>
        <div class="flex-1">
          <div class="text-sm">周五 18:00 前发出周报</div>
          <div class="text-[11px] text-muted-foreground">一次性 · 到点仅通知，不执行操作</div>
        </div>
        <div class="switch on" data-switch></div>
      </div>
      <div class="list-row">
        <i data-lucide="bell" class="w-4 h-4 text-muted-foreground"></i>
        <div class="flex-1">
          <div class="text-sm text-muted-foreground">起身活动一下</div>
          <div class="text-[11px] text-muted-foreground">已关闭</div>
        </div>
        <div class="switch" data-switch></div>
      </div>
      <div class="text-[11px] text-muted-foreground px-2 pt-2 pb-1">提醒为进程内一次性计时：进程未运行则不触发；重启后仅提示「错过 N 条提醒」，不自动补触发。</div>
    </div>

    <!-- 新建任务 Modal（本屏内建） -->
    <div id="task-modal" class="modal-overlay" style="display:none">
      <div class="firstrun-dialog" style="width:460px">
        <div class="flex items-center justify-between mb-4">
          <span class="text-base font-semibold">新建定时任务</span>
          <button class="title-bar-btn" id="close-task-modal"><i data-lucide="x" class="w-4 h-4"></i></button>
        </div>

        <label class="block text-xs text-muted-foreground mb-1.5">任务名</label>
        <input id="tf-name" class="input-field mb-3" placeholder="如：整理下载目录">

        <label class="block text-xs text-muted-foreground mb-1.5">类型</label>
        <div class="flex gap-2 mb-3" id="tf-type">
          <button class="btn btn-outline btn-sm flex-1" data-type="一次性">一次性</button>
          <button class="btn btn-primary btn-sm flex-1" data-type="周期">周期</button>
        </div>

        <label class="block text-xs text-muted-foreground mb-1.5">触发时间</label>
        <input id="tf-time" class="input-field mb-3" placeholder="如：每周一 10:00 / 每日 09:00">

        <label class="block text-xs text-muted-foreground mb-1.5">电源条件</label>
        <div class="bu-select mb-4" data-value="空闲时">
          <button type="button" class="bu-select-trigger" style="min-width:160px">
            <span class="bu-select-value">空闲时</span>
            <i data-lucide="chevron-down" class="w-3.5 h-3.5 text-muted-foreground"></i>
          </button>
          <div class="bu-select-menu hidden">
            <button type="button" class="bu-select-option selected" data-value="空闲时">空闲时<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
            <button type="button" class="bu-select-option" data-value="解锁时">解锁时<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
            <button type="button" class="bu-select-option" data-value="接通电源时">接通电源时<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
          </div>
        </div>

        <label class="block text-xs text-muted-foreground mb-1.5">工作目录</label>
        <input id="tf-cwd" class="input-field mb-5" placeholder="如：C:\\Users\\swq\\Downloads">

        <div class="flex gap-2 justify-end">
          <button class="btn btn-ghost" id="cancel-task-modal">取消</button>
          <button class="btn btn-primary" id="submit-task">创建任务</button>
        </div>
      </div>
    </div>
  `,
  onMount: function (app) {
    /* SpinnerRing —— 对照 task-rows.tsx：双 circle，活动弧 dasharray = c*0.28。
       全部内联样式，不依赖 Tailwind 运行时为延迟注入节点生成工具类。 */
    function spinnerRing(active, label) {
      var size = 24, stroke = 2;
      var r = (size - stroke) / 2;
      var c = 2 * Math.PI * r;
      var spin = active ? 'animation:buSpin 1.1s linear infinite;' : '';
      var arc = active
        ? '<circle cx="' + (size / 2) + '" cy="' + (size / 2) + '" r="' + r + '" fill="none" stroke="hsl(var(--primary))" stroke-width="' + stroke + '" stroke-linecap="round" stroke-dasharray="' + (c * 0.28) + ' ' + (c * 0.72) + '"/>'
        : '';
      return '<span style="position:relative;display:inline-flex;flex-shrink:0;align-items:center;justify-content:center;width:' + size + 'px;height:' + size + 'px">'
        + '<svg width="' + size + '" height="' + size + '" style="position:absolute;inset:0;' + spin + '">'
        + '<circle cx="' + (size / 2) + '" cy="' + (size / 2) + '" r="' + r + '" fill="none" stroke="hsl(var(--border))" stroke-width="' + stroke + '"/>'
        + arc
        + '</svg>'
        + '<span style="position:relative;font-size:10.5px;font-weight:600;color:hsl(var(--foreground));font-variant-numeric:tabular-nums">' + label + '</span>'
        + '</span>';
    }

    /* 左侧 40px：运行中 = spinner 环 + 序号；其余状态 = 状态图标 */
    function avatarHtml(t) {
      if (t.state === 'running') return spinnerRing(true, t.n);
      var icons = { approval: 'clock', blocked: 'lock', queued: 'hourglass', failed: 'x', done: 'check' };
      return '<i data-lucide="' + icons[t.state] + '" class="bu-av-ic bu-av-' + t.state + '"></i>';
    }

    /* 右侧状态 chip：运行中=绿 / 等待=青 / 阻塞=琥珀 / 排队=灰 / 失败=红 / 完成=灰勾 */
    function chipHtml(state) {
      switch (state) {
        case 'running': return '<span class="bu-chip bu-chip-running"><span class="bu-dot"></span>运行中</span>';
        case 'approval': return '<span class="bu-chip bu-chip-approval">等待审批</span>';
        case 'blocked':  return '<span class="bu-chip bu-chip-blocked">路径阻塞</span>';
        case 'queued':   return '<span class="bu-chip bu-chip-queued">排队</span>';
        case 'failed':   return '<span class="bu-chip bu-chip-failed">失败</span>';
        default:         return '<span class="bu-chip bu-chip-done"><i data-lucide="check" class="bu-chip-ic"></i>完成</span>';
      }
    }

    /* 任务数据：六态齐全，task-6 运行中展开子任务进度 */
    var tasks = [
      { id: 'task-6', n: '6', name: '整理桌面文件', state: 'running', lock: false, cancel: true, clock: '2/3',
        note: '正在扫描桌面目录 · 并发槽位 1/4',
        subs: [
          { label: '扫描文件', meta: '128 项', done: true },
          { label: '分类归档', meta: '图片 / 文档', done: true },
          { label: '移动到归档目录', meta: '进行中', done: false }
        ] },
      { id: 'task-3', n: '3', name: '下载语音识别模型', state: 'approval', lock: false, cancel: false, clock: 'L2',
        note: '等待 L2 确认：fs.write 写入模型权重（队列第 1 位）' },
      { id: 'task-9', n: '9', name: '归类本周文档', state: 'blocked', lock: true, cancel: false, clock: '—',
        note: '在等 task-6 释放 Desktop · 释放后自动恢复' },
      { id: 'task-2', n: '2', name: '备份周报附件', state: 'queued', lock: false, cancel: false, clock: '#2',
        note: '队位 #2 · 工具并发上限 4，等待空闲槽位' },
      { id: 'task-5', n: '5', name: '同步相册到网盘', state: 'failed', lock: false, cancel: false, clock: 'R2',
        note: 'R2 路径越界：目标在授权目录外 · 已终止，未产生副作用' },
      { id: 'task-1', n: '1', name: '清理临时目录', state: 'done', lock: false, cancel: false, clock: '1.2s',
        note: '已清理 142 个临时文件 · 耗时 1.2s' }
    ];

    /* 子任务块：缩进 24px 起、行高 36px，左侧 mini progress bar(120px) + 文字 */
    function subsBlock(subs) {
      var h = '<div class="bu-task-subs">';
      subs.forEach(function (s, j) {
        var icon = s.done
          ? '<i data-lucide="check-circle-2" class="bu-sub-ic bu-sub-done"></i>'
          : '<span class="bu-sub-spin"></span>';
        var pct = s.done ? 100 : 55;
        h += '<div class="bu-sub-row" style="animation-delay:' + (120 + j * 90) + 'ms">'
          + icon
          + '<span class="bu-sub-bar"><span class="bu-sub-bar-fill" style="width:' + pct + '%"></span></span>'
          + '<span class="bu-sub-label' + (s.done ? ' bu-sub-label-done' : '') + '">' + s.label + '</span>'
          + '<span class="bu-sub-meta">' + s.meta + '</span>'
          + '</div>';
      });
      h += '</div>';
      return h;
    }

    /* 单行结构：[avatar 40px] [mid flex-1 标题/副标题单行截断] [side chip+计时+取消] */
    function rowHtml(t, i) {
      var g = '<div class="bu-task-group" data-state="' + t.state + '">';
      g += '<div class="bu-task-row' + (t.state === 'blocked' ? ' bu-row-blocked' : '') + '" style="animation-delay:' + (i * 80) + 'ms">';
      g += '<span class="bu-task-avatar">' + avatarHtml(t) + '</span>';
      g += '<div class="bu-task-mid">';
      g += '<span class="bu-task-title">'
        + (t.lock ? '<i data-lucide="lock" class="bu-title-lock"></i>' : '')
        + '<span class="bu-task-title-text">' + t.name + '</span>'
        + '</span>';
      g += '<span class="bu-task-sub">' + t.note + '</span>';
      g += '</div>';
      g += '<div class="bu-task-side">';
      g += chipHtml(t.state);
      g += '<span class="bu-task-clock">' + t.clock + '</span>';
      if (t.cancel) g += '<button class="bu-cancel" data-cancel="' + t.id + '" title="取消任务"><i data-lucide="x" class="bu-cancel-ic"></i>取消</button>';
      g += '</div>';
      g += '</div>';
      if (t.subs) g += subsBlock(t.subs);
      g += '</div>';
      return g;
    }

    var runningEl = document.getElementById('running-tasks');

    function renderRunning() {
      var rows = '';
      tasks.forEach(function (t, i) { rows += rowHtml(t, i); });
      runningEl.innerHTML = rows;
      app.refreshIcons();

      runningEl.querySelectorAll('[data-cancel]').forEach(function (btn) {
        btn.onclick = function () {
          app.toast('已取消 ' + btn.dataset.cancel + '，已执行步骤将列入副作用报告');
        };
      });
    }

    /* onMount：骨架加载 0.8s 后替换为真实内容 */
    setTimeout(renderRunning, 800);
    BUSelect.init(document);

    /* Filter chips：点击实时重组下方任务行（按 .bu-task-group 的 data-state） */
    var chips = document.querySelectorAll('#task-filters .filter-chip');
    chips.forEach(function (chip) {
      chip.onclick = function () {
        chips.forEach(function (c) { c.classList.remove('active'); });
        chip.classList.add('active');
        var f = chip.dataset.filter;
        document.querySelectorAll('#running-tasks [data-state]').forEach(function (el) {
          el.style.display = (f === 'all' || el.dataset.state === f) ? '' : 'none';
        });
      };
    });

    // 提醒开关翻转
    document.querySelectorAll('[data-switch]').forEach(function (sw) {
      sw.onclick = function () {
        sw.classList.toggle('on');
        app.toast(sw.classList.contains('on') ? '提醒已开启' : '提醒已关闭');
      };
    });

    // Modal 显隐
    var modal = document.getElementById('task-modal');
    function openModal() { modal.style.display = 'flex'; }
    function closeModal() { modal.style.display = 'none'; }

    document.getElementById('open-task-modal').onclick = openModal;
    document.getElementById('close-task-modal').onclick = closeModal;
    document.getElementById('cancel-task-modal').onclick = closeModal;
    modal.addEventListener('click', function (e) { if (e.target === modal) closeModal(); });

    // 类型切换
    var typeBtns = document.querySelectorAll('#tf-type [data-type]');
    typeBtns.forEach(function (b) {
      b.onclick = function () {
        typeBtns.forEach(function (x) {
          x.classList.remove('btn-primary');
          x.classList.add('btn-outline');
        });
        b.classList.remove('btn-outline');
        b.classList.add('btn-primary');
      };
    });

    // 提交仅 toast
    document.getElementById('submit-task').onclick = function () {
      var name = document.getElementById('tf-name').value.trim();
      var type = document.querySelector('#tf-type .btn-primary').dataset.type;
      app.toast('已创建' + (type === '周期' ? '周期' : '一次性') + '任务「' + (name || '未命名') + '」（演示）');
      closeModal();
    };
  }
});
