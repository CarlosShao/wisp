Screens.register('cost', {
  nav: { icon: 'chart-column', label: '成本' },
  html: `
  <div class="px-8 py-8">
    <h1 class="page-title">用量与成本</h1>
    <p class="page-subtitle">C23 计量 · 整数微单位结算 · 缓存命中按缓存费率计价，未知价格仅记 token 不估费用</p>

    <div class="grid grid-cols-3 gap-3">
      <div class="insight-card animated-list-item" style="animation-delay:0ms">
        <div class="text-xs text-muted-foreground">今日</div>
        <div class="insight-num mt-1" id="cost-today">¥0.00</div>
        <div class="text-xs text-muted-foreground mt-1.5">12 次调用 · 6,760 in / 2,030 out</div>
      </div>
      <div class="insight-card animated-list-item" style="animation-delay:60ms;border-color:rgba(220,38,38,0.4)">
        <div class="flex items-center justify-between">
          <div class="text-xs text-muted-foreground">本月</div>
          <span class="badge badge-destructive">超预算</span>
        </div>
        <div class="insight-num mt-1" id="cost-month" style="color:hsl(var(--destructive))">¥0.00</div>
        <div class="text-xs text-muted-foreground mt-1.5">342 次调用 · 预算 ¥50.00 · 117%</div>
      </div>
      <div class="insight-card animated-list-item" style="animation-delay:120ms">
        <div class="text-xs text-muted-foreground">累计</div>
        <div class="insight-num mt-1" id="cost-total">¥0.00</div>
        <div class="text-xs text-muted-foreground mt-1.5">自 2026-08-11 起 · cost_daily 已回滚</div>
      </div>
    </div>

    <div class="card mt-4" style="border-color:rgba(220,38,38,0.4)">
      <div class="flex justify-between items-center mb-2">
        <span class="text-sm font-medium">本月预算</span>
        <span class="badge badge-destructive"><i data-lucide="alert-triangle"></i>&nbsp;已超支 · 新任务暂停</span>
      </div>
      <div class="relative pt-5">
        <div class="progress-bar" style="height:8px">
          <div class="progress-fill" style="width:100%;background-color:hsl(var(--destructive))"></div>
        </div>
        <div style="position:absolute;left:80%;top:-2px;bottom:-2px;width:1px;background-color:#b8860b"></div>
        <div style="position:absolute;left:80%;top:-18px;transform:translateX(-50%);font-size:10px;color:#b8860b;white-space:nowrap">80% 告警 ¥40</div>
      </div>
      <div class="text-xs mt-2" style="color:hsl(var(--destructive))">已用 ¥58.60 / ¥50.00 · 超支 ¥8.60。进行中的任务不受影响，次日零点自动恢复新任务派发。</div>
    </div>

    <div class="section-title mt-6">本周洞察</div>
    <div class="insight-card">
      <div class="flex items-start justify-between gap-4">
        <div class="flex-1 min-w-0" id="insight-body">
          <div class="text-xs text-muted-foreground mb-1" id="insight-label">峰值日</div>
          <div class="text-sm leading-relaxed" id="insight-text"></div>
        </div>
        <div id="insight-chart" class="flex-shrink-0"></div>
      </div>
      <div class="pagination-dots mt-4" id="insight-dots">
        <span class="active" data-i="0"></span>
        <span data-i="1"></span>
        <span data-i="2"></span>
      </div>
    </div>

    <div class="section-title mt-6">近 7 日费用</div>
    <div class="card p-4">
      <div class="flex items-end justify-between gap-2" style="height:150px">
        <div class="flex-1 flex flex-col items-center justify-end gap-1">
          <div class="text-[10px] font-mono text-muted-foreground">0.42</div>
          <div class="w-full rounded-t animated-list-item" style="height:34px;background-color:#d6d9dc;animation-delay:0ms"></div>
          <div class="text-[10px] text-muted-foreground">周一</div>
        </div>
        <div class="flex-1 flex flex-col items-center justify-end gap-1">
          <div class="text-[10px] font-mono text-muted-foreground">0.78</div>
          <div class="w-full rounded-t animated-list-item" style="height:64px;background-color:#d6d9dc;animation-delay:50ms"></div>
          <div class="text-[10px] text-muted-foreground">周二</div>
        </div>
        <div class="flex-1 flex flex-col items-center justify-end gap-1">
          <div class="text-[10px] font-mono" style="color:hsl(var(--primary))">1.35</div>
          <div class="w-full rounded-t animated-list-item" style="height:110px;background-color:#86C2B9;animation-delay:100ms"></div>
          <div class="text-[10px]" style="color:hsl(var(--primary))">周三</div>
        </div>
        <div class="flex-1 flex flex-col items-center justify-end gap-1">
          <div class="text-[10px] font-mono text-muted-foreground">0.87</div>
          <div class="w-full rounded-t animated-list-item" style="height:71px;background-color:#c2c9cc;animation-delay:150ms"></div>
          <div class="text-[10px] text-muted-foreground">今天</div>
        </div>
        <div class="flex-1 flex flex-col items-center justify-end gap-1">
          <div class="text-[10px] font-mono text-muted-foreground">0.95</div>
          <div class="w-full rounded-t animated-list-item" style="height:77px;background-color:#d6d9dc;animation-delay:200ms"></div>
          <div class="text-[10px] text-muted-foreground">周五</div>
        </div>
        <div class="flex-1 flex flex-col items-center justify-end gap-1">
          <div class="text-[10px] font-mono text-muted-foreground">0.31</div>
          <div class="w-full rounded-t animated-list-item" style="height:25px;background-color:#d6d9dc;animation-delay:250ms"></div>
          <div class="text-[10px] text-muted-foreground">周六</div>
        </div>
        <div class="flex-1 flex flex-col items-center justify-end gap-1">
          <div class="text-[10px] font-mono text-muted-foreground">0.22</div>
          <div class="w-full rounded-t animated-list-item" style="height:18px;background-color:#d6d9dc;animation-delay:300ms"></div>
          <div class="text-[10px] text-muted-foreground">周日</div>
        </div>
      </div>
    </div>

    <div class="tabs mt-6" id="cost-tabs">
      <div class="tab active" data-view="day">日视图</div>
      <div class="tab" data-view="month">月视图</div>
    </div>

    <div id="cost-content"></div>

    <p class="text-xs text-muted-foreground mt-4">价格表 v3（2026-09-01 生效）· 缓存输入按缓存费率折算 · 本地模型不计费 · 本页数字与 task_log / cost_daily 对账到微单位。</p>
  </div>
  `,
  onMount: function (app) {
    var state = { view: 'day', insight: 0 };
    var shimmerTimer = null;

    var insights = [
      { label: '峰值日', text: '本周三花费最高 ¥1.35，因批量归档任务触发 4 次 L2 审批',
        bars: [0.42, 0.78, 1.35, 0.87, 0.95, 0.31, 0.22], hi: 2 },
      { label: '缓存效率', text: '缓存命中占输入 token 的 68%，本周较缓存前节省约 ¥2.10',
        bars: [0.40, 0.55, 0.72, 0.60, 0.80, 0.50, 0.45], hi: 4 },
      { label: '本地模型', text: '本地 qwen3:7b 承接了 23% 的闲聊任务，本月该部分成本 ¥0.00',
        bars: [0.15, 0.20, 0.18, 0.30, 0.25, 0.40, 0.50], hi: 6 }
    ];

    function modelTable(rows) {
      var h = '<div class="card p-0 overflow-hidden"><table class="records-table"><thead><tr>' +
        '<th>模型</th><th style="text-align:right">IN tokens</th><th style="text-align:right">OUT tokens</th><th style="text-align:right">缓存命中</th><th style="text-align:right">费用</th>' +
        '</tr></thead><tbody>';
      rows.forEach(function (r) {
        h += '<tr class="animated-list-item">' +
          '<td><span class="tool-chip">' + r.m + '</span></td>' +
          '<td style="text-align:right;font-family:monospace;font-size:12px">' + r.in + '</td>' +
          '<td style="text-align:right;font-family:monospace;font-size:12px">' + r.out + '</td>' +
          '<td style="text-align:right;font-family:monospace;font-size:12px;color:hsl(var(--muted-foreground))">' + r.cached + '</td>' +
          '<td style="text-align:right;font-family:monospace;font-size:12px;font-weight:600">' + r.cost + '</td></tr>';
      });
      return h + '</tbody></table></div>';
    }

    function miniChart(ins) {
      var max = Math.max.apply(null, ins.bars);
      var h = '<div class="flex items-end gap-1" style="height:48px">';
      ins.bars.forEach(function (v, i) {
        var bh = Math.max(4, Math.round(v / max * 44));
        var col = (i === ins.hi) ? 'hsl(var(--primary))' : '#c9ced2';
        h += '<div style="width:10px;height:' + bh + 'px;background-color:' + col + ';border-radius:2px"></div>';
      });
      return h + '</div>';
    }

    function renderInsight() {
      var ins = insights[state.insight];
      var body = document.getElementById('insight-body');
      body.classList.remove('animated-list-item');
      void body.offsetHeight;
      body.classList.add('animated-list-item');
      document.getElementById('insight-label').textContent = ins.label;
      document.getElementById('insight-text').textContent = ins.text;
      document.getElementById('insight-chart').innerHTML = miniChart(ins);
      document.querySelectorAll('#insight-dots span').forEach(function (d) {
        d.classList.toggle('active', parseInt(d.dataset.i) === state.insight);
      });
    }

    function shimmerSkeleton() {
      return '<div class="card p-0">' +
        '<div class="list-row"><div class="shimmer-line" style="width:30%"></div></div>' +
        '<div class="list-row"><div class="shimmer-line" style="width:78%"></div></div>' +
        '<div class="list-row"><div class="shimmer-line" style="width:55%"></div></div>' +
        '<div class="list-row"><div class="shimmer-line" style="width:42%"></div></div>' +
        '</div>';
    }

    function render() {
      var c = document.getElementById('cost-content');
      var h = '';
      if (state.view === 'day') {
        h += '<div class="section-title mt-2">按模型分解（今日）</div>';
        h += modelTable([
          { m: 'deepseek-chat', in: '4,820', out: '1,260', cached: '12,400', cost: '¥0.62' },
          { m: 'qwen3:7b · 本地', in: '1,940', out: '680', cached: '0', cost: '¥0.00' },
          { m: 'realtime · 陪聊', in: '—', out: '—', cached: '—', cost: '¥0.25' }
        ]);
        h += '<div class="section-title mt-5">今日任务明细</div>';
        h += '<div class="card p-0 overflow-hidden"><table class="records-table"><thead><tr>' +
          '<th>任务</th><th style="text-align:right">IN</th><th style="text-align:right">OUT</th><th style="text-align:right">缓存</th><th style="text-align:right">费用</th><th style="text-align:right">耗时</th>' +
          '</tr></thead><tbody>';
        [
          ['整理本周截图归档', '3,240', '860', '9,800', '¥0.28', '12.4s'],
          ['读取日志并诊断报错', '2,480', '520', '6,100', '¥0.21', '8.8s'],
          ['起草周报邮件', '1,860', '740', '4,200', '¥0.19', '9.1s'],
          ['翻译英文文档段落', '960', '480', '0', '¥0.08', '3.2s'],
          ['查找并清理重复文件', '1,120', '310', '0', '¥0.06', '4.6s'],
          ['随口闲聊', '140', '90', '0', '¥0.05', '1.1s']
        ].forEach(function (r, i) {
          h += '<tr class="animated-list-item" style="animation-delay:' + (i * 45) + 'ms">' +
            '<td>' + r[0] + '</td>' +
            '<td style="text-align:right;font-family:monospace;font-size:12px">' + r[1] + '</td>' +
            '<td style="text-align:right;font-family:monospace;font-size:12px">' + r[2] + '</td>' +
            '<td style="text-align:right;font-family:monospace;font-size:12px;color:hsl(var(--muted-foreground))">' + r[3] + '</td>' +
            '<td style="text-align:right;font-family:monospace;font-size:12px;font-weight:600">' + r[4] + '</td>' +
            '<td style="text-align:right;font-family:monospace;font-size:12px;color:hsl(var(--muted-foreground))">' + r[5] + '</td></tr>';
        });
        h += '</tbody></table></div>';
      } else {
        h += '<div class="section-title mt-2">花费最高任务 Top 5（本月）</div>';
        h += '<div class="card p-0 overflow-hidden"><table class="records-table"><thead><tr>' +
          '<th style="width:44px">#</th><th>任务</th><th style="text-align:right">Tokens</th><th style="text-align:right">费用</th>' +
          '</tr></thead><tbody>';
        [
          ['重构导出模块并补测试', '86.0k', '¥3.12'],
          ['整理桌面与文件归档', '52.4k', '¥2.40'],
          ['周报与邮件批量起草', '41.2k', '¥1.96'],
          ['日志诊断与根因分析', '38.7k', '¥1.74'],
          ['翻译技术文档 3 篇', '27.1k', '¥1.31']
        ].forEach(function (r, i) {
          h += '<tr class="animated-list-item" style="animation-delay:' + (i * 50) + 'ms">' +
            '<td><span class="badge" style="justify-content:center">' + (i + 1) + '</span></td>' +
            '<td>' + r[0] + '</td>' +
            '<td style="text-align:right;font-family:monospace;font-size:12px;color:hsl(var(--muted-foreground))">' + r[1] + '</td>' +
            '<td style="text-align:right;font-family:monospace;font-size:12px;font-weight:600;color:hsl(var(--destructive))">' + r[2] + '</td></tr>';
        });
        h += '</tbody></table></div>';
        h += '<div class="section-title mt-5">按模型分解（本月）</div>';
        h += modelTable([
          { m: 'deepseek-chat', in: '96,400', out: '28,300', cached: '214,000', cost: '¥47.20' },
          { m: 'qwen3:7b · 本地', in: '31,200', out: '11,800', cached: '0', cost: '¥0.00' },
          { m: 'realtime · 陪聊', in: '—', out: '—', cached: '—', cost: '¥11.40' }
        ]);
      }
      c.innerHTML = h;
      app.refreshIcons();
    }

    function showShimmerThenRender() {
      var c = document.getElementById('cost-content');
      c.innerHTML = shimmerSkeleton();
      if (shimmerTimer) clearTimeout(shimmerTimer);
      shimmerTimer = setTimeout(function () { render(); }, 300);
    }

    function countUp(el, target, duration) {
      if (!el) return;
      var start = null;
      function step(ts) {
        if (!start) start = ts;
        var progress = Math.min((ts - start) / duration, 1);
        var eased = 1 - Math.pow(1 - progress, 3);
        el.textContent = '¥' + (target * eased).toFixed(2);
        if (progress < 1) requestAnimationFrame(step);
        else el.textContent = '¥' + target.toFixed(2);
      }
      requestAnimationFrame(step);
    }

    // 数字滚动动效（onMount）
    countUp(document.getElementById('cost-today'), 0.87, 600);
    countUp(document.getElementById('cost-month'), 58.60, 600);
    countUp(document.getElementById('cost-total'), 156.82, 600);

    // 洞察分页切换
    document.querySelectorAll('#insight-dots span').forEach(function (d) {
      d.onclick = function () {
        state.insight = parseInt(d.dataset.i);
        renderInsight();
      };
    });
    renderInsight();

    // 日 / 月视图切换
    document.querySelectorAll('#cost-tabs .tab').forEach(function (t) {
      t.onclick = function () {
        document.querySelectorAll('#cost-tabs .tab').forEach(function (x) { x.classList.remove('active'); });
        t.classList.add('active');
        state.view = t.dataset.view;
        showShimmerThenRender();
      };
    });

    render();
  }
});
