Screens.register('palette', {
  nav: { icon: 'command', label: '命令' },
  html: `
<div class="p-6 max-w-2xl mx-auto">
  <div class="page-title">命令面板</div>
  <div class="page-subtitle mb-6">命令面板静态全貌 · 键盘优先 · 模糊匹配命令、工具与页面 · 回车即把自由文本交给 Wisp</div>

  <!-- 搜索框：大字号，对齐 BeautifulUI Search -->
  <div class="card card-spotlight p-0 overflow-hidden">
    <div class="flex items-center gap-3 px-4 h-14 border-b border-border">
      <i data-lucide="search" class="w-5 h-5 text-muted-foreground"></i>
      <input id="pal-search" class="flex-1 bg-transparent text-[15px] outline-none placeholder:text-muted-foreground" placeholder="输入命令、工具名，或直接输入任务…" />
      <span class="kbd">Ctrl K</span>
    </div>

    <div id="pal-list" class="py-2 px-2 max-h-[420px] overflow-y-auto">

      <!-- 快捷指令 -->
      <div class="pal-group">
        <div class="flex items-center gap-1.5 px-2 pt-3 pb-1 text-[11px] font-semibold uppercase tracking-wide text-muted-foreground">
          <i data-lucide="zap" class="w-3.5 h-3.5"></i><span>快捷指令</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="唤起语音 聆听 summon listening ctrl alt q" data-toast="已唤起语音聆听（冷启动加载 VAD+ASR 约 1–3s）">
          <i data-lucide="audio-lines" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">唤起语音聆听</span>
          <span class="kbd">Ctrl+Alt+Q</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="切换静音 静音 mute ctrl alt m" data-toast="已切换静音状态（Muted / Armed）">
          <i data-lucide="mic-off" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">切换静音</span>
          <span class="kbd">Ctrl+Alt+M</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="整理桌面 归档 截图 organize" data-toast="已把今日桌面截图归档到 图片/Wisp/今日/">
          <i data-lucide="layout-grid" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">整理今日桌面截图</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="切换陪聊模式 全双工 conversation" data-toast="陪聊模式需 L2 级隐私确认后开启">
          <i data-lucide="message-circle" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">切换陪聊模式</span>
        </div>
      </div>

      <!-- 工具 -->
      <div class="pal-group">
        <div class="flex items-center gap-1.5 px-2 pt-3 pb-1 text-[11px] font-semibold uppercase tracking-wide text-muted-foreground">
          <i data-lucide="wrench" class="w-3.5 h-3.5"></i><span>工具</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="列出目录 文件夹 ls folder" data-toast="工具 directory.list 走同一宿主桥与风险门控">
          <i data-lucide="folder" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">列出目录</span>
          <span class="tool-chip">L0</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="读取文件 读文件 read file-text" data-toast="工具 file.read 已经 C26 路径解析收口">
          <i data-lucide="file-text" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">读取文件</span>
          <span class="tool-chip">L0</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="运行命令 终端 terminal shell exec" data-toast="工具 shell.exec 需按风险等级走 L1/L2 确认">
          <i data-lucide="terminal" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">运行命令</span>
          <span class="tool-chip">L2</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="网页搜索 搜索 web search 上网" data-toast="web.search 实现路径待定（抓结果页 vs API），S3 前不接线">
          <i data-lucide="globe" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">网页搜索</span>
          <span class="tool-chip">L1</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="新建提醒 定时 任务 闹钟 reminder bell" data-toast="已登记一条提醒任务（任务调度走 PathLock 路径锁）">
          <i data-lucide="bell" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">新建提醒</span>
        </div>
      </div>

      <!-- 设置 -->
      <div class="pal-group">
        <div class="flex items-center gap-1.5 px-2 pt-3 pb-1 text-[11px] font-semibold uppercase tracking-wide text-muted-foreground">
          <i data-lucide="sliders-horizontal" class="w-3.5 h-3.5"></i><span>设置</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="打开偏好设置 设置 config ctrl 逗号" data-jump="config">
          <i data-lucide="settings" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">打开偏好设置</span>
          <span class="kbd">Ctrl+,</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="管理 api 密钥 凭据 key secret" data-jump="config">
          <i data-lucide="key-round" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">管理 API 密钥</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="安全 隐私 权限 门控 security" data-jump="security">
          <i data-lucide="shield-check" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">安全门控与权限</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="成本 用量 花费 cost token" data-jump="cost">
          <i data-lucide="chart-column" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">查看今日成本</span>
        </div>
      </div>

      <!-- 跳转 -->
      <div class="pal-group">
        <div class="flex items-center gap-1.5 px-2 pt-3 pb-1 text-[11px] font-semibold uppercase tracking-wide text-muted-foreground">
          <i data-lucide="corner-up-right" class="w-3.5 h-3.5"></i><span>跳转</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="对话 聊天 chat 面板" data-jump="chat">
          <i data-lucide="message-square-text" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">对话面板</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="审批 队列 确认 approval" data-jump="approval">
          <i data-lucide="shield-check" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">审批队列</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="任务 列表 进行中 tasks" data-jump="tasks">
          <i data-lucide="list-checks" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">任务列表</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="球 状态 悬浮球 ball 状态机" data-jump="ball">
          <i data-lucide="orbit" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">悬浮球状态机</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="隐私 导出 privacy" data-jump="privacy">
          <i data-lucide="lock" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">隐私与数据</span>
        </div>
      </div>

      <!-- 最近任务（按时间倒序） -->
      <div class="pal-group">
        <div class="flex items-center gap-1.5 px-2 pt-3 pb-1 text-[11px] font-semibold uppercase tracking-wide text-muted-foreground">
          <i data-lucide="history" class="w-3.5 h-3.5"></i><span>最近任务</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="整理今日桌面截图 归档 图片" data-toast="任务已重放入队：整理今日桌面截图">
          <i data-lucide="layout-grid" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">整理今日桌面截图</span>
          <span class="text-[11px] text-muted-foreground">2 分钟前</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="生成本周工作周报 周报" data-toast="任务已重放入队：生成本周工作周报">
          <i data-lucide="file-text" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">生成本周工作周报</span>
          <span class="text-[11px] text-muted-foreground">18 分钟前</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="汇总会议纪要 纪要" data-toast="任务已重放入队：汇总 9/22 会议纪要">
          <i data-lucide="notebook-pen" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">汇总 9/22 会议纪要</span>
          <span class="text-[11px] text-muted-foreground">1 小时前</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="备份配置 onedrive 备份" data-toast="任务已重放入队：备份配置到 OneDrive">
          <i data-lucide="cloud-upload" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">备份配置到 OneDrive</span>
          <span class="text-[11px] text-muted-foreground">昨天 22:14</span>
        </div>
        <div class="pal-row flex items-center gap-2.5 px-2 py-2 rounded-md cursor-pointer hover:bg-[var(--accent)]" data-search="查询上周 api 成本 成本 用量" data-toast="任务已重放入队：查询上周 API 成本">
          <i data-lucide="chart-column" class="w-4 h-4 text-muted-foreground"></i>
          <span class="text-sm flex-1">查询上周 API 成本</span>
          <span class="text-[11px] text-muted-foreground">昨天 09:40</span>
        </div>
      </div>

      <!-- 过滤空态 -->
      <div id="pal-empty" class="hidden flex-col items-center py-10 text-center">
        <i data-lucide="search-x" class="w-8 h-8 text-muted-foreground mb-3"></i>
        <div class="text-sm font-medium">无匹配结果</div>
        <div class="text-xs text-muted-foreground mt-1">按 Enter 把这句话作为自由文本任务交给 Wisp（噪声环境下的文本兜底通道）</div>
      </div>
    </div>

    <div class="flex items-center gap-4 px-4 h-10 border-t border-border text-[11px] text-muted-foreground">
      <i data-lucide="arrow-up" class="w-3 h-3"></i><i data-lucide="arrow-down" class="w-3 h-3"></i><span>导航</span>
      <span class="kbd-xs">Enter</span><span>执行 / 提交任务</span>
      <span class="kbd-xs">Esc</span><span>关闭</span>
      <span class="flex-1"></span>
      <span class="kbd-xs">Ctrl+Alt+Q</span><span>唤起语音</span>
    </div>
  </div>

  <div class="mt-6 text-center text-xs text-muted-foreground">面板为无状态视图：每次打开由宿主侧 resync 全量推送，不缓存跨会话状态</div>
</div>
`,
  onMount: function (app) {
    var input = document.getElementById('pal-search');
    var empty = document.getElementById('pal-empty');
    var groups = document.querySelectorAll('#pal-list .pal-group');

    function applyFilter() {
      var q = input.value.trim().toLowerCase();
      var visibleTotal = 0;
      groups.forEach(function (g) {
        var vis = 0;
        g.querySelectorAll('.pal-row').forEach(function (row) {
          var hit = !q || row.dataset.search.indexOf(q) !== -1;
          row.style.display = hit ? '' : 'none';
          if (hit) vis++;
        });
        g.style.display = vis ? '' : 'none';
        visibleTotal += vis;
      });
      empty.style.display = visibleTotal ? 'none' : 'flex';
    }

    input.addEventListener('input', applyFilter);
    input.addEventListener('keydown', function (e) {
      if (e.key === 'Enter' && input.value.trim()) {
        app.toast('已提交为自由文本任务：' + input.value.trim());
      }
    });
    input.focus();

    document.querySelectorAll('#pal-list .pal-row').forEach(function (row) {
      row.addEventListener('click', function () {
        var jump = row.dataset.jump;
        var msg = row.dataset.toast;
        if (jump) {
          app.showScreen(jump);
        } else if (msg) {
          app.toast(msg);
        }
      });
    });
  }
});
