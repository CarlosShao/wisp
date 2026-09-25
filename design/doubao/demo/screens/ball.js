Screens.register('ball', {
  nav: { icon: 'orbit', label: '球状态' },
  html: `
<style>
  @keyframes wbPulse { 0%,100%{transform:scale(1)} 50%{transform:scale(1.1)} }
  @keyframes wbBreathe { 0%,100%{opacity:.55} 50%{opacity:.82} }
  @keyframes wbSpin { to { transform:rotate(360deg) } }
  @keyframes wbFlow { 0%{background-position:0% 50%} 100%{background-position:200% 50%} }
  @keyframes wbRing { 0%{box-shadow:0 0 0 0 rgba(220,80,80,.45)} 100%{box-shadow:0 0 0 14px rgba(220,80,80,0)} }
  #demo-ball {
    width: 72px; height: 72px; border-radius: 50%; position: relative;
    background: radial-gradient(circle at 32% 28%, rgba(255,255,255,.75), rgba(134,194,185,.9) 46%, rgba(108,168,160,.95));
    box-shadow: 0 0 26px rgba(134,194,185,.5), inset 0 2px 6px rgba(255,255,255,.5), inset 0 -4px 10px rgba(0,0,0,.08);
    transition: width 260ms cubic-bezier(.32,.72,0,1), height 260ms cubic-bezier(.32,.72,0,1), opacity 260ms ease, box-shadow 260ms ease, background 260ms ease;
  }
  .state-row.selected { background: hsl(var(--accent)); }
  .state-row.selected .state-en { color: hsl(var(--accent-foreground)); }
</style>

<div class="p-6 max-w-4xl mx-auto">
  <div class="page-title">悬浮球状态机</div>
  <div class="page-subtitle mb-6">20 态 · 42 条权威转移 · 液态玻璃材质（点击下方任意一行，上方预览球随之变化）</div>

  <!-- 选中态预览球 -->
  <div class="card mb-6 flex items-center gap-6">
    <div class="flex items-center justify-center" style="width:140px;height:140px;flex-shrink:0;position:relative;">
      <div id="demo-ball"></div>
      <!-- Thinking 态：球下呼吸三点（选中 Thinking 时才显示，不动球本身变色逻辑） -->
      <div class="churning-grid v-dots" id="ball-thinking" style="position:absolute;bottom:14px;display:none;"><span></span><span></span><span></span><span></span><span></span><span></span><span></span><span></span><span></span></div>
    </div>
    <div class="flex-1 min-w-0">
      <div class="flex items-center gap-2 flex-wrap">
        <span id="demo-en" class="text-lg font-semibold font-mono"></span>
        <span id="demo-cn" class="badge badge-primary"></span>
        <span id="demo-timer" class="badge"></span>
      </div>
      <div class="text-sm mt-2"><span class="text-muted-foreground">触发：</span><span id="demo-trigger"></span></div>
      <div class="text-xs text-muted-foreground mt-1.5 leading-relaxed"><span class="font-medium">视觉/动效：</span><span id="demo-visual"></span></div>
      <div class="text-xs text-muted-foreground mt-1"><span class="font-medium">点球行为：</span><span id="demo-click"></span></div>
    </div>
  </div>

  <!-- 参考图 -->
  <div class="card mb-6 p-0 overflow-hidden">
    <img src="../01-ball-states.jpg" alt="球状态定稿参考图" class="w-full h-auto block" />
    <div class="px-4 py-2 text-[11px] text-muted-foreground border-t border-border">参考图仅作质感参照；本屏主体为 20 态交互列表</div>
  </div>

  <div class="section-title">20 态清单（点击行联动上方预览球）</div>
  <div id="state-list" class="card p-0 overflow-hidden"></div>

  <!-- 关键路径转移图 -->
  <div class="card mt-6">
    <div class="section-title">关键路径转移（D43 权威转移表子集）</div>
    <div class="flex flex-wrap items-center gap-y-2 text-[13px]">
      <span class="badge badge-primary">FirstRun</span><i data-lucide="arrow-right" class="w-3.5 h-3.5 text-muted-foreground"></i>
      <span class="badge">Sleeping</span><i data-lucide="arrow-right" class="w-3.5 h-3.5 text-muted-foreground"></i>
      <span class="badge">Armed</span><i data-lucide="arrow-right" class="w-3.5 h-3.5 text-muted-foreground"></i>
      <span class="badge badge-primary">Listening</span><i data-lucide="arrow-right" class="w-3.5 h-3.5 text-muted-foreground"></i>
      <span class="badge">Thinking</span><i data-lucide="arrow-right" class="w-3.5 h-3.5 text-muted-foreground"></i>
      <span class="badge">Acting</span><i data-lucide="arrow-right" class="w-3.5 h-3.5 text-muted-foreground"></i>
      <span class="badge badge-destructive">Confirming</span><i data-lucide="arrow-right" class="w-3.5 h-3.5 text-muted-foreground"></i>
      <span class="badge badge-success">Speaking</span><i data-lucide="arrow-right" class="w-3.5 h-3.5 text-muted-foreground"></i>
      <span class="badge badge-warn">Warm</span><i data-lucide="arrow-right" class="w-3.5 h-3.5 text-muted-foreground"></i>
      <span class="badge">Settling</span><i data-lucide="arrow-right" class="w-3.5 h-3.5 text-muted-foreground"></i>
      <span class="badge">Sleeping</span>
    </div>
    <div class="mt-3 flex items-center gap-2 text-xs">
      <i data-lucide="zap" class="w-3.5 h-3.5 text-muted-foreground"></i>
      <span class="text-muted-foreground">热路径二次唤起：</span>
      <span class="badge badge-warn">Warm</span>
      <i data-lucide="arrow-right" class="w-3.5 h-3.5 text-muted-foreground"></i>
      <span class="badge badge-primary">Listening</span>
      <span class="text-muted-foreground">——零模型加载（B4：保留 SessionScope，免冷启动 1–3s）</span>
    </div>
    <div class="mt-2 text-[11px] text-muted-foreground leading-relaxed">
      其它分支：Acting 高风险工具入 AwaitingApproval 队列；路径锁冲突转入 Queued；重复调用达阈值 8 转入 Stuck；
      网络探测失败转入 NoNetwork（保留任务 ctx）；90s 无交互进入 Settling 回落；拨报中检出人声（AEC）切回 Listening。
    </div>
  </div>

  <!-- 热键与交互 -->
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-6">
    <div class="card">
      <div class="section-title">全局热键</div>
      <div class="space-y-2 text-[13px]">
        <div class="flex items-center justify-between"><span>唤起语音聆听</span><span class="kbd">Ctrl+Alt+Q</span></div>
        <div class="flex items-center justify-between"><span>静音切换</span><span class="kbd">Ctrl+Alt+M</span></div>
        <div class="flex items-center justify-between"><span>命令面板</span><span class="kbd">Ctrl+K</span></div>
        <div class="flex items-center justify-between"><span>打开设置</span><span class="kbd">Ctrl+,</span></div>
        <div class="flex items-center justify-between"><span>Confirming 期间取消</span><span class="kbd">Esc</span></div>
        <div class="mt-2 text-[11px] text-muted-foreground leading-relaxed">
          Ctrl+Alt+Space 可配置但不作默认（中文输入法会抢占）；原 Ctrl+Alt+W 在本机被第三方程序占用，故改 Q。
        </div>
      </div>
    </div>
    <div class="card">
      <div class="section-title">球体交互</div>
      <div class="space-y-2 text-[13px]">
        <div class="flex gap-2"><span class="kbd-xs">单击</span><span>Sleeping/Armed/Warm 唤起 Listening；Confirming 期间单击 = 否决该调用</span></div>
        <div class="flex gap-2"><span class="kbd-xs">拖拽</span><span>移动位置；拖到左/右/上边缘收缩吸附半隐，悬停或单击弹回</span></div>
        <div class="flex gap-2"><span class="kbd-xs">双击</span><span class="text-muted-foreground">当前版本未绑定（预留）</span></div>
        <div class="flex gap-2"><span class="kbd-xs">右键</span><span>托盘右键菜单：打开面板 / 静音 / 暂停唤醒 / 退出</span></div>
        <div class="flex gap-2"><span class="kbd-xs">穿透</span><span>透明区点击穿透到桌面；球永不抢焦点</span></div>
      </div>
    </div>
  </div>
</div>
`,
  onMount: function (app) {
    // 20 态权威数据：视觉列对齐 SPEC-08 §2.1，触发/转移对齐 D43
    var STATES = [
      { key: 'FirstRun', cn: '首次引导', trigger: '首次启动，未写 config.toml', visual: 'accent 底 + 引导脉冲，引导模型/目录授权/填 Key', click: '进入引导流程，完成后写库落配置', glow: '134,194,185', size: 72, anim: 'pulse', timer: '动画' },
      { key: 'Sleeping', cn: '休眠', trigger: '空闲回落 / 进程常驻但无会话', visual: '静态玻璃体 34.72px（56×0.62），零动画零定时器，CPU约0', click: '唤起 Listening（冷启动加载 VAD+ASR 1–3s）', glow: '150,160,170', size: 44, anim: 'none', timer: '零定时器' },
      { key: 'Armed', cn: '待命(唤醒词)', trigger: 'KWS 开启，已加载唤醒词模型', visual: 'opacity 0.6 静态，等待唤醒词命中', click: '单击球唤起 Listening；唤醒词亦可', glow: '134,194,185', size: 60, anim: 'none', timer: '静态' },
      { key: 'Muted', cn: '静音', trigger: '按静音键，停 KWS 推理', visual: 'opacity 0.4 + 1.5px 斜杠', click: '再按静音键回 Armed（或 Sleeping）', glow: '150,160,170', size: 56, anim: 'none', timer: '静态' },
      { key: 'Listening', cn: '聆听', trigger: '单击球 / 唤起键 / 唤醒词命中', visual: '外环随音量波动 + accent 描边，液体随音频电平旋转起伏', click: '单击球 = 否决，丢弃缓冲回 Warm（不落盘）', glow: '134,194,185', size: 72, anim: 'pulse', timer: '动画' },
      { key: 'Thinking', cn: '思考', trigger: 'VAD 判停且语音 >=300ms', visual: '球体底部 2px accent 光带流动', click: '无（流式进行中）', glow: '134,194,185', size: 72, anim: 'flow', timer: '动画' },
      { key: 'Acting', cn: '执行', trigger: '收到首 token 且有 tool call', visual: 'accent 常亮实心，opacity 1', click: '无（工具执行中）', glow: '134,194,185', size: 72, anim: 'none', timer: '常亮' },
      { key: 'Confirming', cn: '待确认(L1)', trigger: '工具需 2–3s 倒计时确认', visual: 'danger 脉冲环（仅一次）+ 环下 micro 倒计时', click: '单击球 / Esc / 否决词 = 否决该调用回 LLM', glow: '220,90,90', size: 72, anim: 'ring', timer: '脉冲一次' },
      { key: 'AwaitingApproval', cn: '待审批(L2)', trigger: '高风险工具入 C18 审批队列', visual: 'danger 环 + 右上角标（12px 圆，白字=队列深度）', click: '取消该审批；面板侧只能拒绝，允许权在原生侧', glow: '220,90,90', size: 72, anim: 'none', timer: '静态' },
      { key: 'Speaking', cn: '播报', trigger: '工具完成 / 纯文本首 token', visual: 'success 呼吸（1.6s）', click: '单击 / 唤起键打断，停 TTS 回 Listening', glow: '90,190,130', size: 72, anim: 'breathe', timer: '呼吸' },
      { key: 'Settling', cn: '收束', trigger: '静默完成 / 90s 无交互', visual: 'opacity 1 渐隐至 0.35，260ms，回收 SessionScope', click: '新唤起则取消回落，直接 Listening', glow: '150,160,170', size: 60, anim: 'fade', timer: '渐隐' },
      { key: 'Warm', cn: '温热会话', trigger: '播报完 / 会话内热态', visual: '暖色微呼吸 2.4s（opacity 在 0.55 至 0.7 间起伏），日常最常见态', click: '单击球唤起 Listening，零模型加载', glow: '220,170,110', size: 72, anim: 'breathe', timer: '呼吸' },
      { key: 'Conversation', cn: '陪聊全双工', trigger: '开启陪聊模式（首次需 L2 隐私确认）', visual: 'danger 常亮环，不得渐隐、不得弱于 Confirming', click: '结束陪聊回文本循环', glow: '220,90,90', size: 72, anim: 'none', timer: '常亮' },
      { key: 'Downloading', cn: '下载模型', trigger: '缺模型文件', visual: '环形进度 + 百分比（sha256 + manifest 签名校验）', click: '无', glow: '134,194,185', size: 72, anim: 'spin', timer: '环形' },
      { key: 'Error', cn: '错误', trigger: '网络失败 / panic recover / 设备占用', visual: 'danger 底 + x 图标，明示错误类与人话', click: '确认或 10s 后回错误前会话态 / Warm', glow: '220,90,90', size: 72, anim: 'none', timer: '静态' },
      { key: 'NoNetwork', cn: '无网络', trigger: '任意态网络探测失败', visual: 'fg-tertiary 底 + wifi-off', click: '重试，保留当前任务 ctx 不中断', glow: '150,160,170', size: 64, anim: 'none', timer: '静态' },
      { key: 'Unconfigured', cn: '未配置', trigger: '缺 API Key / 配置非法', visual: 'warn 底 + key-round', click: '打开配置填入 Key', glow: '217,178,106', size: 64, anim: 'none', timer: '静态' },
      { key: 'WatchdogAlert', cn: '看门狗告警', trigger: '连续 3 次回落失败', visual: 'warn 底 + alert-triangle', click: '一键重启', glow: '217,178,106', size: 64, anim: 'pulse', timer: '告警' },
      { key: 'Queued', cn: '排队', trigger: '路径锁冲突（C20），等前序任务', visual: '主态 + 左下 6px info 小点叠加', click: '面板可见「在等谁」', glow: '134,194,185', size: 68, anim: 'none', timer: '静态' },
      { key: 'Stuck', cn: '卡住', trigger: '重复调用达阈值 [3,5,8] 中 8', visual: 'warn 底 + refresh-cw，明示重复了什么', click: '继续 / 换个方式 / 放弃', glow: '217,178,106', size: 64, anim: 'spin', timer: '告警' }
    ];

    var ball = document.getElementById('demo-ball');
    var list = document.getElementById('state-list');

    function applyState(s, rowEl) {
      // 球体变色/变尺寸/变动画
      ball.style.width = s.size + 'px';
      ball.style.height = s.size + 'px';
      ball.style.opacity = '';
      ball.style.animation = '';
      ball.style.background = 'radial-gradient(circle at 32% 28%, rgba(255,255,255,.75), rgba(' + s.glow + ',.9) 46%, rgba(' + s.glow + ',.95))';
      ball.style.boxShadow = '0 0 26px rgba(' + s.glow + ',.5), inset 0 2px 6px rgba(255,255,255,.5), inset 0 -4px 10px rgba(0,0,0,.08)';
      if (s.anim === 'pulse') ball.style.animation = 'wbPulse 1.6s ease-in-out infinite';
      else if (s.anim === 'breathe') ball.style.animation = 'wbBreathe 2.4s ease-in-out infinite';
      else if (s.anim === 'spin') ball.style.animation = 'wbSpin 2.4s linear infinite';
      else if (s.anim === 'flow') ball.style.animation = 'wbFlow 1.2s linear infinite';
      else if (s.anim === 'ring') ball.style.animation = 'wbRing 2s ease-out infinite';
      else if (s.anim === 'fade') { ball.style.animation = 'wbBreathe 260ms ease-in-out forwards'; ball.style.opacity = '0.4'; }

      document.getElementById('demo-en').textContent = s.key;
      document.getElementById('demo-cn').textContent = s.cn;
      document.getElementById('demo-timer').textContent = s.timer;
      document.getElementById('demo-trigger').textContent = s.trigger;
      document.getElementById('demo-visual').textContent = s.visual;
      document.getElementById('demo-click').textContent = s.click;

      // Thinking 态：球下呼吸三点（仅选中 Thinking 时显示，不改球变色/尺寸逻辑）
      var thinkDots = document.getElementById('ball-thinking');
      if (thinkDots) thinkDots.style.display = (s.key === 'Thinking') ? 'inline-flex' : 'none';

      list.querySelectorAll('.state-row').forEach(function (r) { r.classList.remove('selected'); });
      if (rowEl) rowEl.classList.add('selected');
      app.refreshIcons();
    }

    // 初始骨架：5-6 行 shimmer-line 占位，模拟列表加载
    function renderSkeleton() {
      list.innerHTML = '';
      for (var i = 0; i < 6; i++) {
        var sk = document.createElement('div');
        sk.className = 'flex items-center gap-3 px-4 py-3 border-b border-border/50';
        sk.innerHTML =
          '<span class="shimmer-block flex-shrink-0" style="width:10px;height:10px;border-radius:50%;"></span>' +
          '<span class="shimmer-line flex-shrink-0" style="width:120px;"></span>' +
          '<span class="shimmer-line flex-shrink-0" style="width:70px;"></span>' +
          '<span class="shimmer-line flex-1" style="height:12px;"></span>';
        list.appendChild(sk);
      }
    }

    // 真实 20 态行（错峰淡入上滑）
    function renderStates() {
      list.innerHTML = '';
      STATES.forEach(function (s, idx) {
        var row = document.createElement('div');
        row.className = 'state-row animated-list-item flex items-center gap-3 px-4 py-2.5 cursor-pointer border-b border-border/50';
        row.style.animationDelay = (idx * 40) + 'ms';
        row.innerHTML =
          '<span style="width:10px;height:10px;border-radius:50%;background:rgb(' + s.glow + ');flex-shrink:0;box-shadow:0 0 8px rgba(' + s.glow + ',.6)"></span>' +
          '<span class="state-en font-mono text-[13px] font-semibold w-[168px] flex-shrink-0">' + s.key + '</span>' +
          '<span class="text-[13px] w-[110px] flex-shrink-0 text-muted-foreground">' + s.cn + '</span>' +
          '<span class="text-xs text-muted-foreground flex-1 min-w-0 truncate">' + s.trigger + '</span>' +
          '<span class="kbd-xs flex-shrink-0">' + s.timer + '</span>';
        row.addEventListener('click', function () { applyState(s, row); });
        list.appendChild(row);
      });
      // 真实列表就位后，选中第二行 Sleeping
      applyState(STATES[1], list.children[1]);
    }

    // 预览区即时就绪（不带行选中），列表 0.6s 后由骨架替换为真实 20 态
    renderSkeleton();
    applyState(STATES[1], null);
    setTimeout(renderStates, 600);
  }
});
