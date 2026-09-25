/* ============================================================
   Wisp Minimal Demo · app.js
   路由 / 导航 / Toast / 命令面板 / Modal / 主题切换
   ============================================================ */

// ---------- 屏幕注册表 ----------
window.Screens = {
  _registry: {},
  register(id, config) {
    this._registry[id] = config;
  },
  get(id) { return this._registry[id]; },
  all() { return this._registry; }
};

// ---------- 导航配置 ----------
const NAV_ITEMS = [
  { id: 'chat',     icon: 'message-square-text', label: '对话' },
  { id: 'approval', icon: 'shield-check',        label: '审批' },
  { id: 'palette',  icon: 'command',             label: '命令' },
  { id: 'tasks',    icon: 'list-checks',         label: '任务' },
  { id: 'ball',     icon: 'orbit',               label: '球状态' },
  { id: 'config',   icon: 'settings',            label: '设置' },
  { id: 'security', icon: 'shield-alert',        label: '安全' },
  { id: 'privacy',  icon: 'lock',                label: '隐私' },
  { id: 'cost',     icon: 'chart-column',        label: '成本' }
];

// ---------- App 主体 ----------
const App = {
  currentScreen: 'chat',
  theme: 'light',

  init() {
    this.loadTheme();
    this.buildNav();
    this.bindEvents();
    this.showScreen('chat');
    // 首次加载弹 firstrun
    if (!localStorage.getItem('wisp-firstrun-done')) {
      setTimeout(() => this.showFirstRun(), 400);
    }
    // 首载 Preloader：1.1s 后淡出（仅首次，切屏不闪）
    setTimeout(() => {
      const pl = document.getElementById('preloader');
      if (pl) pl.classList.add('hidden-loader');
    }, 1100);
  },

  // ---------- 主题 ----------
  loadTheme() {
    const saved = localStorage.getItem('wisp-theme');
    if (saved) this.theme = saved;
    this.applyTheme();
  },
  applyTheme() {
    document.documentElement.classList.toggle('dark', this.theme === 'dark');
  },
  toggleTheme() {
    this.theme = this.theme === 'dark' ? 'light' : 'dark';
    localStorage.setItem('wisp-theme', this.theme);
    this.applyTheme();
    this.refreshIcons();
  },

  // ---------- 导航 ----------
  buildNav() {
    const nav = document.getElementById('side-nav');
    nav.innerHTML = '';
    // 滑翔指示条
    const glide = document.createElement('div');
    glide.className = 'nav-glide';
    glide.id = 'nav-glide';
    nav.appendChild(glide);
    NAV_ITEMS.forEach(item => {
      const btn = document.createElement('button');
      btn.className = 'nav-item' + (item.id === this.currentScreen ? ' active' : '');
      btn.dataset.screen = item.id;
      btn.title = item.label;
      btn.innerHTML = `<i data-lucide="${item.icon}"></i>`;
      btn.onclick = () => this.showScreen(item.id);
      nav.appendChild(btn);
    });
    // 底部 firstrun 入口
    const spacer = document.createElement('div');
    spacer.className = 'nav-spacer';
    nav.appendChild(spacer);
    const helpBtn = document.createElement('button');
    helpBtn.className = 'nav-item';
    helpBtn.title = '首次引导';
    helpBtn.innerHTML = `<i data-lucide="life-buoy"></i>`;
    helpBtn.onclick = () => this.showFirstRun();
    nav.appendChild(helpBtn);
    this.updateNavGlide();
  },

  updateNavGlide() {
    const glide = document.getElementById('nav-glide');
    const active = document.querySelector('.nav-item.active');
    if (!glide || !active) return;
    const navRect = active.parentElement.getBoundingClientRect();
    const activeRect = active.getBoundingClientRect();
    // glide 锚定 top:0、left:10px(=(60-40)/2)，在 60px 导航列内水平居中；
    // 垂直平移使 40x40 glide 顶部对齐 active 40px 按钮顶部（同高 => 完美居中，不塌）。
    const y = activeRect.top - navRect.top;
    glide.style.transform = `translateY(${y}px)`;
  },

  updateNavActive(id) {
    document.querySelectorAll('.nav-item').forEach(el => {
      el.classList.toggle('active', el.dataset.screen === id);
    });
    this.updateNavGlide();
  },

  // ---------- 屏幕渲染 ----------
  showScreen(id) {
    const screen = Screens.get(id);
    if (!screen) {
      document.getElementById('app').innerHTML =
        `<div class="text-muted-foreground text-sm">屏幕 "${id}" 尚未实现</div>`;
      return;
    }
    this.currentScreen = id;
    this.updateNavActive(id);
    const app = document.getElementById('app');
    app.style.animation = 'none';
    app.offsetHeight; // reflow
    app.style.animation = '';
    app.innerHTML = screen.html;
    app.classList.remove('screen-enter');
    app.offsetHeight; // reflow
    app.classList.add('screen-enter');
    this.refreshIcons();
    // 屏幕级初始化钩子
    if (typeof screen.onMount === 'function') {
      try { screen.onMount(this); } catch(e) { console.warn('screen onMount error:', e); }
    }
  },

  refreshIcons() {
    if (window.lucide && lucide.createIcons) {
      lucide.createIcons();
    }
  },

  // ---------- Toast ----------
  toast(msg) {
    const container = document.getElementById('toast-container');
    const el = document.createElement('div');
    el.className = 'toast';
    el.textContent = msg;
    container.appendChild(el);
    setTimeout(() => el.remove(), 2600);
  },

  // ---------- 命令面板 ----------
  paletteItems: [],
  paletteIndex: 0,

  openPalette() {
    const modal = document.getElementById('command-palette');
    modal.classList.remove('hidden');
    const input = document.getElementById('palette-input');
    input.value = '';
    this.buildPaletteList('');
    setTimeout(() => input.focus(), 50);
  },
  closePalette() {
    document.getElementById('command-palette').classList.add('hidden');
  },
  buildPaletteList(query) {
    const list = document.getElementById('palette-list');
    const q = query.toLowerCase().trim();
    let items = [];

    // 快捷指令
    const quick = [
      { label: '整理桌面', icon: 'sparkles', action: () => { this.closePalette(); this.toast('演示原型，仅作展示'); } },
      { label: '切换陪聊模式', icon: 'message-circle', action: () => { this.closePalette(); this.toast('演示原型，仅作展示'); } },
      { label: '查看今日成本', icon: 'chart-column', action: () => { this.closePalette(); this.showScreen('cost'); } }
    ];
    // 工具
    const tools = [
      { label: '列出目录', icon: 'folder', action: () => { this.closePalette(); this.toast('演示原型，仅作展示'); } },
      { label: '读取文件', icon: 'file-text', action: () => { this.closePalette(); this.toast('演示原型，仅作展示'); } },
      { label: '运行命令', icon: 'terminal', action: () => { this.closePalette(); this.toast('演示原型，仅作展示'); } }
    ];
    // 屏幕跳转
    const navItems = NAV_ITEMS.map(n => ({
      label: n.label, icon: n.icon, action: () => { this.closePalette(); this.showScreen(n.id); }
    }));

    const groups = [
      { name: '快捷指令', items: quick },
      { name: '工具', items: tools },
      { name: '导航', items: navItems }
    ];

    let html = '';
    let flatIndex = 0;
    this.paletteItems = [];
    groups.forEach(g => {
      const filtered = g.items.filter(i => !q || i.label.toLowerCase().includes(q));
      if (filtered.length === 0) return;
      html += `<div class="palette-group-label">${g.name}</div>`;
      filtered.forEach(item => {
        const idx = this.paletteItems.length;
        this.paletteItems.push(item);
        html += `<div class="palette-item${idx === this.paletteIndex ? ' active' : ''}" data-idx="${idx}">
          <i data-lucide="${item.icon}"></i>${item.label}</div>`;
      });
    });
    if (!html) html = '<div class="palette-group-label">无匹配结果</div>';
    list.innerHTML = html;
    this.refreshIcons();

    // 绑定点击
    list.querySelectorAll('.palette-item').forEach(el => {
      el.onclick = () => {
        const idx = parseInt(el.dataset.idx);
        this.paletteItems[idx]?.action();
      };
    });
  },

  // ---------- 首次引导 ----------
  showFirstRun() {
    const screen = Screens.get('firstrun');
    if (!screen) return;
    const modal = document.getElementById('firstrun-modal');
    const content = document.getElementById('firstrun-content');
    content.innerHTML = screen.html;
    modal.classList.remove('hidden');
    this.refreshIcons();
    if (typeof screen.onMount === 'function') {
      try { screen.onMount(this); } catch(e) {}
    }
  },
  closeFirstRun() {
    document.getElementById('firstrun-modal').classList.add('hidden');
    localStorage.setItem('wisp-firstrun-done', '1');
  },

  // ---------- 窗口显隐 ----------
  toggleWindow() {
    const win = document.getElementById('main-window');
    win.classList.toggle('hidden-window');
  },

  // ---------- 事件绑定 ----------
  bindEvents() {
    // 主题切换
    document.getElementById('theme-toggle').onclick = () => this.toggleTheme();

    // 悬浮球
    document.getElementById('floating-ball').onclick = () => this.toggleWindow();

    // 命令面板 Ctrl+K
    document.addEventListener('keydown', (e) => {
      if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        this.openPalette();
      }
      if (e.key === 'Escape') {
        this.closePalette();
        this.closeFirstRun();
      }
      // 命令面板内导航
      const palette = document.getElementById('command-palette');
      if (!palette.classList.contains('hidden')) {
        if (e.key === 'ArrowDown') {
          e.preventDefault();
          this.paletteIndex = Math.min(this.paletteIndex + 1, this.paletteItems.length - 1);
          this.buildPaletteList(document.getElementById('palette-input').value);
        }
        if (e.key === 'ArrowUp') {
          e.preventDefault();
          this.paletteIndex = Math.max(this.paletteIndex - 1, 0);
          this.buildPaletteList(document.getElementById('palette-input').value);
        }
        if (e.key === 'Enter' && this.paletteItems[this.paletteIndex]) {
          e.preventDefault();
          this.paletteItems[this.paletteIndex].action();
        }
      }
    });

    // 命令面板搜索
    document.getElementById('palette-input').addEventListener('input', (e) => {
      this.paletteIndex = 0;
      this.buildPaletteList(e.target.value);
    });

    // 点击遮罩关闭 modal
    document.getElementById('command-palette').addEventListener('click', (e) => {
      if (e.target.id === 'command-palette') this.closePalette();
    });
    document.getElementById('firstrun-modal').addEventListener('click', (e) => {
      if (e.target.id === 'firstrun-modal') this.closeFirstRun();
    });
  }
};

// 全局暴露
window.App = App;
window.toast = (msg) => App.toast(msg);
window.toggleWindow = () => App.toggleWindow();

/* ===== bu-select: global custom dropdown component ===== */
window.BUSelect = (function () {
  var openMenu = null;

  function closeAll() {
    if (openMenu) {
      openMenu.classList.remove('open');
      var m = openMenu.querySelector('.bu-select-menu');
      if (m) m.classList.add('hidden');
      openMenu = null;
    }
  }

  function init(root) {
    if (!root) root = document;
    root.querySelectorAll('.bu-select').forEach(function (wrap) {
      if (wrap.dataset.buInit) return;
      wrap.dataset.buInit = '1';
      var trigger = wrap.querySelector('.bu-select-trigger');
      var menu = wrap.querySelector('.bu-select-menu');
      if (!trigger || !menu) return;

      trigger.addEventListener('click', function (e) {
        e.stopPropagation();
        var isOpen = wrap.classList.contains('open');
        closeAll();
        if (!isOpen) {
          wrap.classList.add('open');
          menu.classList.remove('hidden');
          openMenu = wrap;
        }
      });

      menu.querySelectorAll('.bu-select-option').forEach(function (opt) {
        opt.addEventListener('click', function (e) {
          e.stopPropagation();
          var val = opt.dataset.value || opt.textContent.trim();
          var valEl = wrap.querySelector('.bu-select-value');
          if (valEl) valEl.textContent = opt.textContent.trim();
          wrap.dataset.value = val;
          menu.querySelectorAll('.bu-select-option').forEach(function (o) { o.classList.remove('selected'); });
          opt.classList.add('selected');
          closeAll();
          if (window.app && window.app.refreshIcons) window.app.refreshIcons();
        });
      });
    });
  }

  document.addEventListener('click', function () { closeAll(); });
  document.addEventListener('keydown', function (e) { if (e.key === 'Escape') closeAll(); });

  return { init: init, closeAll: closeAll };
})();
