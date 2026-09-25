Screens.register('config', {
  nav: { icon: 'settings', label: '设置' },
  html: `
    <div id="config-root">
      <div class="flex gap-4 items-center mb-5">
        <div class="relative flex-1 max-w-sm">
          <i data-lucide="search" class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground pointer-events-none"></i>
          <input id="cfg-search" class="input-field pl-9" placeholder="搜索配置项，如 api_base / 采样率 / 主题" />
        </div>
        <span class="text-xs text-muted-foreground ml-auto">config.toml 为唯一真相源 · GUI 仅作编辑器</span>
      </div>

      <div class="flex gap-6">
        <aside class="w-28 shrink-0">
          <div class="py-2 px-3 rounded-lg text-sm cursor-pointer transition-colors text-muted-foreground hover:bg-muted cfg-nav-item" data-target="general">通用</div>
          <div class="py-2 px-3 rounded-lg text-sm cursor-pointer transition-colors text-muted-foreground hover:bg-muted cfg-nav-item" data-target="voice">语音</div>
          <div class="py-2 px-3 rounded-lg text-sm cursor-pointer transition-colors text-muted-foreground hover:bg-muted cfg-nav-item" data-target="audio">音频</div>
          <div class="py-2 px-3 rounded-lg text-sm cursor-pointer transition-colors bg-accent text-accent-foreground font-medium cfg-nav-item" data-target="llm">大模型</div>
          <div class="py-2 px-3 rounded-lg text-sm cursor-pointer transition-colors text-muted-foreground hover:bg-muted cfg-nav-item" data-target="tools">工具</div>
          <div class="py-2 px-3 rounded-lg text-sm cursor-pointer transition-colors text-muted-foreground hover:bg-muted cfg-nav-item" data-target="appearance">外观</div>
          <div class="py-2 px-3 rounded-lg text-sm cursor-pointer transition-colors text-muted-foreground hover:bg-muted flex items-center gap-1.5 cfg-nav-item" data-target="privacy"><i data-lucide="lock" class="w-3 h-3 inline"></i>隐私</div>
        </aside>

        <section class="flex-1 min-w-0">
          <h1 class="page-title mb-1" id="cfg-title">大模型</h1>
          <p class="page-subtitle mb-4" id="cfg-sub">[llm] · 热加载即时生效</p>

          <!-- 加载骨架：节切换 / 搜索过渡 -->
          <div id="cfg-skeleton" style="display:none" aria-hidden="true">
            <div class="card p-4">
              <div class="space-y-3">
                <div class="flex items-center justify-between gap-4"><div class="shimmer-line w-40"></div><div class="shimmer-line w-24" style="height:32px"></div></div>
                <div class="flex items-center justify-between gap-4"><div class="shimmer-line w-32"></div><div class="shimmer-line w-36" style="height:32px"></div></div>
                <div class="flex items-center justify-between gap-4"><div class="shimmer-line w-48"></div><div class="shimmer-line w-28" style="height:32px"></div></div>
                <div class="flex items-center justify-between gap-4"><div class="shimmer-line w-36"></div><div class="shimmer-line w-24" style="height:32px"></div></div>
              </div>
            </div>
          </div>

          <div id="cfg-sections">
          <!-- ============ 通用 [app]/[session]/[hotkey] ============ -->
          <div class="cfg-section card" data-section="general" style="display:none; padding:4px 16px">
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="language 语言 zh-CN">
              <div class="min-w-0"><div class="font-mono text-sm">language</div><div class="text-xs text-muted-foreground mt-0.5">界面语言</div></div>
              <div class="flex items-center gap-2 shrink-0">
                <div class="bu-select" data-value="zh-CN">
                  <button type="button" class="bu-select-trigger" style="min-width:140px">
                    <span class="bu-select-value">zh-CN</span>
                    <i data-lucide="chevron-down" class="w-3.5 h-3.5 text-muted-foreground"></i>
                  </button>
                  <div class="bu-select-menu hidden">
                <button type="button" class="bu-select-option selected" data-value="zh-CN">zh-CN<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                <button type="button" class="bu-select-option" data-value="zh-TW">zh-TW<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                <button type="button" class="bu-select-option" data-value="en-US">en-US<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                  </div>
                </div>
                <span class="badge">重启</span>
              </div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="autostart 开机自启">
              <div class="min-w-0"><div class="font-mono text-sm">autostart</div><div class="text-xs text-muted-foreground mt-0.5">开机自动启动常驻进程</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch" data-static="1"></div><span class="badge">重启</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="single_instance 单实例">
              <div class="min-w-0"><div class="font-mono text-sm">single_instance</div><div class="text-xs text-muted-foreground mt-0.5">仅允许一个实例运行</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch on" data-static="1"></div><span class="badge">重启</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="warm_timeout_sec 预热 空闲">
              <div class="min-w-0"><div class="font-mono text-sm">warm_timeout_sec</div><div class="text-xs text-muted-foreground mt-0.5">语音常驻预热空闲超时（秒）</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-20 text-center" value="90" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="conversation_idle_sec 会话空闲">
              <div class="min-w-0"><div class="font-mono text-sm">conversation_idle_sec</div><div class="text-xs text-muted-foreground mt-0.5">对话空闲多久后进入收尾（秒）</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-20 text-center" value="30" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3" data-search="summon 呼出 快捷键 hotkey">
              <div class="min-w-0"><div class="font-mono text-sm">hotkey.summon</div><div class="text-xs text-muted-foreground mt-0.5">呼出悬浮球的全局快捷键</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-40 text-center" value="Alt + Space" /><span class="badge badge-warn">重载</span></div>
            </div>
          </div>

          <!-- ============ 语音 [voice] ============ -->
          <div class="cfg-section card" data-section="voice" style="display:none; padding:4px 16px">
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="enabled 语音 启用">
              <div class="min-w-0"><div class="font-mono text-sm">voice.enabled</div><div class="text-xs text-muted-foreground mt-0.5">启用语音链路</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch on"></div><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="asr provider 语音识别 sherpa">
              <div class="min-w-0"><div class="font-mono text-sm">voice.asr.provider</div><div class="text-xs text-muted-foreground mt-0.5">本地 ASR 引擎</div></div>
              <div class="flex items-center gap-2 shrink-0">
                <div class="bu-select" data-value="local-sherpa">
                  <button type="button" class="bu-select-trigger" style="min-width:150px">
                    <span class="bu-select-value">local-sherpa</span>
                    <i data-lucide="chevron-down" class="w-3.5 h-3.5 text-muted-foreground"></i>
                  </button>
                  <div class="bu-select-menu hidden">
                <button type="button" class="bu-select-option selected" data-value="local-sherpa">local-sherpa<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                <button type="button" class="bu-select-option" data-value="cloud-asr">cloud-asr<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                  </div>
                </div>
                <span class="badge badge-warn">重载</span>
              </div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="asr model 语音识别 模型">
              <div class="min-w-0"><div class="font-mono text-sm">voice.asr.model</div><div class="text-xs text-muted-foreground mt-0.5">ASR 模型 ID（更换需重载模型）</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-52 font-mono" value="sherpa-paraformer-zh-20240323" /><span class="badge badge-warn">重载</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="tts voice 语音合成 音色">
              <div class="min-w-0"><div class="font-mono text-sm">voice.tts.voice</div><div class="text-xs text-muted-foreground mt-0.5">TTS 音色</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-44" value="aina-low" /><span class="badge badge-warn">重载</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="tts speed 语速">
              <div class="min-w-0"><div class="font-mono text-sm">voice.tts.speed</div><div class="text-xs text-muted-foreground mt-0.5">TTS 语速倍率</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-20 text-center" value="1.0" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="punctuation 标点">
              <div class="min-w-0"><div class="font-mono text-sm">voice.punctuation</div><div class="text-xs text-muted-foreground mt-0.5">ASR 结果自动加标点</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch on"></div><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="conversation_mode 全双工 陪聊">
              <div class="min-w-0"><div class="font-mono text-sm">voice.conversation_mode</div><div class="text-xs text-muted-foreground mt-0.5">陪聊路径全双工（干活路径仍半双工）</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch"></div><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3" data-search="aec 回声消除">
              <div class="min-w-0"><div class="font-mono text-sm">voice.aec.enabled</div><div class="text-xs text-muted-foreground mt-0.5">回声消除（自渲染音频作回采参考）</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch on"></div><span class="badge badge-warn">重载</span></div>
            </div>
          </div>

          <!-- ============ 音频 [audio] ============ -->
          <div class="cfg-section card" data-section="audio" style="display:none; padding:4px 16px">
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="input_device 输入设备 麦克风 采样">
              <div class="min-w-0"><div class="font-mono text-sm">audio.input_device</div><div class="text-xs text-muted-foreground mt-0.5">录音输入设备</div></div>
              <div class="flex items-center gap-2 shrink-0">
                <div class="bu-select" data-value="default · 麦克风阵列">
                  <button type="button" class="bu-select-trigger" style="min-width:180px">
                    <span class="bu-select-value">default · 麦克风阵列</span>
                    <i data-lucide="chevron-down" class="w-3.5 h-3.5 text-muted-foreground"></i>
                  </button>
                  <div class="bu-select-menu hidden">
                <button type="button" class="bu-select-option selected" data-value="default · 麦克风阵列">default · 麦克风阵列<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                <button type="button" class="bu-select-option" data-value="麦克风 (Realtek Audio)">麦克风 (Realtek Audio)<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                  </div>
                </div>
                <span class="badge badge-warn">重载</span>
              </div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="sample_rate 采样率">
              <div class="min-w-0"><div class="font-mono text-sm">audio.sample_rate</div><div class="text-xs text-muted-foreground mt-0.5">录音采样率（Hz）</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-24 text-center" value="16000" /><span class="badge badge-warn">重载</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="vad_threshold 静音检测 阈值">
              <div class="min-w-0"><div class="font-mono text-sm">audio.vad_threshold</div><div class="text-xs text-muted-foreground mt-0.5">VAD 语音端点检测阈值（0–1）</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-20 text-center" value="0.5" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="half_duplex 半双工">
              <div class="min-w-0"><div class="font-mono text-sm">audio.half_duplex</div><div class="text-xs text-muted-foreground mt-0.5">干活路径半双工（硬编码只读）</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch on" data-static="1" style="opacity:0.6"></div><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3" data-search="mic_muted_default 默认静音">
              <div class="min-w-0"><div class="font-mono text-sm">audio.mic_muted_default</div><div class="text-xs text-muted-foreground mt-0.5">会话开始时麦克风默认静音</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch on"></div><span class="badge badge-success">即时</span></div>
            </div>
          </div>

          <!-- ============ 大模型 [llm] ============ -->
          <div class="cfg-section card" data-section="llm" style="padding:4px 16px">
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="provider 服务商 大模型">
              <div class="min-w-0"><div class="font-mono text-sm">provider</div><div class="text-xs text-muted-foreground mt-0.5">LLM 服务商预设</div></div>
              <div class="flex items-center gap-2 shrink-0">
                <div class="bu-select" data-value="deepseek">
                  <button type="button" class="bu-select-trigger" style="min-width:150px">
                    <span class="bu-select-value">deepseek</span>
                    <i data-lucide="chevron-down" class="w-3.5 h-3.5 text-muted-foreground"></i>
                  </button>
                  <div class="bu-select-menu hidden">
                <button type="button" class="bu-select-option selected" data-value="deepseek">deepseek<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                <button type="button" class="bu-select-option" data-value="openai">openai<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                <button type="button" class="bu-select-option" data-value="anthropic">anthropic<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                <button type="button" class="bu-select-option" data-value="qwen">qwen<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                <button type="button" class="bu-select-option" data-value="moonshot">moonshot<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                <button type="button" class="bu-select-option" data-value="ollama">ollama<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                  </div>
                </div>
                <span class="badge badge-success">即时</span>
              </div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="model 模型名 deepseek-chat">
              <div class="min-w-0"><div class="font-mono text-sm">model</div><div class="text-xs text-muted-foreground mt-0.5">对话模型 ID</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-44 font-mono" value="deepseek-chat" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="api_base 接口地址 base_url">
              <div class="min-w-0">
                <div class="font-mono text-sm">api_base</div>
                <div class="text-xs text-muted-foreground mt-0.5">OpenAI 兼容端点</div>
                <div id="cfg-api-base-err" class="text-xs text-destructive mt-1" style="display:none">api_base 不能为空</div>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <input id="cfg-api-base" class="input-field w-56 font-mono" placeholder="https://api.deepseek.com/v1" value="https://api.deepseek.com/v1" />
                <span class="badge badge-success">即时</span>
              </div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="api_key 密钥 密码">
              <div class="min-w-0"><div class="font-mono text-sm">api_key</div><div class="text-xs text-muted-foreground mt-0.5">写一次性输入，明文不落 TOML（仅存后 4 位）</div></div>
              <div class="flex items-center gap-2 shrink-0"><input type="password" class="input-field w-56 font-mono" placeholder="已保存 ····8f3a" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="temperature 温度 创造性 slider">
              <div class="min-w-0"><div class="font-mono text-sm">temperature</div><div class="text-xs text-muted-foreground mt-0.5">采样温度 0–2</div></div>
              <div class="flex items-center gap-3 shrink-0">
                <input type="range" min="0" max="2" step="0.1" value="0.7" id="cfg-temp" class="w-32">
                <span class="font-mono text-xs w-7 text-center" id="cfg-temp-val">0.7</span>
                <span class="badge badge-success">即时</span>
              </div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="max_tokens 最大输出 长度">
              <div class="min-w-0"><div class="font-mono text-sm">max_tokens</div><div class="text-xs text-muted-foreground mt-0.5">单次最大输出 token</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-24 text-center" value="4096" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="top_p 核采样">
              <div class="min-w-0"><div class="font-mono text-sm">top_p</div><div class="text-xs text-muted-foreground mt-0.5">核采样阈值</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-20 text-center" value="1.0" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="timeout_seconds 超时">
              <div class="min-w-0"><div class="font-mono text-sm">timeout_seconds</div><div class="text-xs text-muted-foreground mt-0.5">单次请求超时（秒）</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-20 text-center" value="30" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="stream 流式输出">
              <div class="min-w-0"><div class="font-mono text-sm">stream</div><div class="text-xs text-muted-foreground mt-0.5">流式输出增量 token</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch on"></div><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3" data-search="enable_cache 上下文缓存">
              <div class="min-w-0"><div class="font-mono text-sm">enable_cache</div><div class="text-xs text-muted-foreground mt-0.5">启用提示缓存（按 cached 价计费）</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch on"></div><span class="badge badge-success">即时</span></div>
            </div>
          </div>

          <!-- ============ 工具 [agent] + 工具默认级 ============ -->
          <div class="cfg-section card" data-section="tools" style="display:none; padding:4px 16px">
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="max_rounds 最大轮次 agent">
              <div class="min-w-0"><div class="font-mono text-sm">agent.max_rounds</div><div class="text-xs text-muted-foreground mt-0.5">单任务最大 Agent 轮次</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-20 text-center" value="50" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="token_budget 上下文预算">
              <div class="min-w-0"><div class="font-mono text-sm">agent.token_budget</div><div class="text-xs text-muted-foreground mt-0.5">上下文窗口预算</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-28 text-center" value="200000" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="per_tool_timeout_ms 工具超时">
              <div class="min-w-0"><div class="font-mono text-sm">agent.per_tool_timeout_ms</div><div class="text-xs text-muted-foreground mt-0.5">单个工具调用超时（毫秒）</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-28 text-center" value="30000" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3" data-search="steering 引导 纠偏">
              <div class="min-w-0"><div class="font-mono text-sm">agent.steering_enabled</div><div class="text-xs text-muted-foreground mt-0.5">允许用户中途纠偏（steering）</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch on"></div><span class="badge badge-success">即时</span></div>
            </div>

            <div class="section-title mt-4 mb-2">内置工具默认风险级（R1 下界，非结论）</div>
            <div class="card p-0 mb-3" style="background:transparent">
              <table class="data-table">
                <thead><tr><th>工具</th><th>默认级</th><th>说明</th></tr></thead>
                <tbody>
                  <tr><td class="font-mono">fs.read</td><td><span class="badge badge-success">L0</span></td><td>只读无副作用</td></tr>
                  <tr><td class="font-mono">fs.write</td><td><span class="badge badge-warn">L1</span></td><td>可逆写（temp + 原子 rename）</td></tr>
                  <tr><td class="font-mono">web.fetch</td><td><span class="badge badge-warn">L1</span></td><td>出网读取，受 R5/R4 约束</td></tr>
                  <tr><td class="font-mono">fs.delete</td><td><span class="badge badge-destructive">L2</span></td><td>不可逆，永不进会话授权</td></tr>
                  <tr><td class="font-mono">shell.exec</td><td><span class="badge badge-destructive">L2</span></td><td>强制 argv 向量，元字符升 L2</td></tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- ============ 外观 [ball]/[panel] ============ -->
          <div class="cfg-section card" data-section="appearance" style="display:none; padding:4px 16px">
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="theme 主题 浅色 深色">
              <div class="min-w-0"><div class="font-mono text-sm">app.theme</div><div class="text-xs text-muted-foreground mt-0.5">界面主题</div></div>
              <div class="flex items-center gap-2 shrink-0">
                <div class="bu-select" data-value="light">
                  <button type="button" class="bu-select-trigger" style="min-width:120px">
                    <span class="bu-select-value">light</span>
                    <i data-lucide="chevron-down" class="w-3.5 h-3.5 text-muted-foreground"></i>
                  </button>
                  <div class="bu-select-menu hidden">
                <button type="button" class="bu-select-option selected" data-value="light">light<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                <button type="button" class="bu-select-option" data-value="dark">dark<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                <button type="button" class="bu-select-option" data-value="auto">auto<i data-lucide="check" class="w-3.5 h-3.5"></i></button>
                  </div>
                </div>
                <span class="badge badge-success">即时</span>
              </div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="font_size 字体 字号">
              <div class="min-w-0"><div class="font-mono text-sm">panel.font_size</div><div class="text-xs text-muted-foreground mt-0.5">面板正文字号（px）</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-20 text-center" value="13" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="ball size 球大小 悬浮球">
              <div class="min-w-0"><div class="font-mono text-sm">ball.size</div><div class="text-xs text-muted-foreground mt-0.5">悬浮球直径（44–72）</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-20 text-center" value="56" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="opacity_idle 空闲透明度">
              <div class="min-w-0"><div class="font-mono text-sm">ball.opacity_idle</div><div class="text-xs text-muted-foreground mt-0.5">空闲时悬浮球不透明度</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-20 text-center" value="0.35" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="click_through 鼠标穿透">
              <div class="min-w-0"><div class="font-mono text-sm">ball.click_through</div><div class="text-xs text-muted-foreground mt-0.5">空闲时鼠标穿透悬浮球</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch on"></div><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3" data-search="panel width 面板宽度">
              <div class="min-w-0"><div class="font-mono text-sm">panel.width</div><div class="text-xs text-muted-foreground mt-0.5">面板宽度（px）</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-24 text-center" value="640" /><span class="badge badge-success">即时</span></div>
            </div>
          </div>

          <!-- ============ 隐私 [privacy]/[memory] ============ -->
          <div class="cfg-section card" data-section="privacy" style="display:none; padding:4px 16px">
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="redact_paths 脱敏 路径">
              <div class="min-w-0"><div class="font-mono text-sm">privacy.redact_paths</div><div class="text-xs text-muted-foreground mt-0.5">导出/诊断时脱敏本地路径</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch"></div><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="retention_days 留存 保留天数">
              <div class="min-w-0"><div class="font-mono text-sm">privacy.retention_days</div><div class="text-xs text-muted-foreground mt-0.5">任务日志/记忆保留天数</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-20 text-center" value="30" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="diagnostics_opt_in 诊断上报">
              <div class="min-w-0"><div class="font-mono text-sm">privacy.diagnostics_opt_in</div><div class="text-xs text-muted-foreground mt-0.5">同意上传匿名诊断数据</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch"></div><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="l1_max 记忆 条目">
              <div class="min-w-0"><div class="font-mono text-sm">memory.l1_max</div><div class="text-xs text-muted-foreground mt-0.5">L1 工作记忆最大条目数</div></div>
              <div class="flex items-center gap-2 shrink-0"><input class="input-field w-20 text-center" value="20" /><span class="badge badge-success">即时</span></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3 border-b border-border" data-search="keep_transcript 记录 转写">
              <div class="min-w-0"><div class="font-mono text-sm">privacy.keep_transcript</div><div class="text-xs text-muted-foreground mt-0.5">持久化对话转写（硬编码关闭，只读）</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch" data-static="1" style="opacity:0.4"></div><i data-lucide="lock" class="w-3.5 h-3.5 text-muted-foreground"></i></div>
            </div>
            <div class="cfg-row flex items-center justify-between gap-4 py-3" data-search="keep_audio 录音 音频">
              <div class="min-w-0"><div class="font-mono text-sm">privacy.keep_audio</div><div class="text-xs text-muted-foreground mt-0.5">持久化录音音频（硬编码关闭，只读）</div></div>
              <div class="flex items-center gap-2 shrink-0"><div class="switch" data-static="1" style="opacity:0.4"></div><i data-lucide="lock" class="w-3.5 h-3.5 text-muted-foreground"></i></div>
            </div>
          </div>
          </div><!-- /#cfg-sections -->

          <!-- config.toml 预览（Code Block，随节切换） -->
          <div class="mt-5">
            <button id="cfg-code-toggle" class="flex items-center gap-2 text-xs font-medium text-muted-foreground hover:text-foreground transition-colors">
              <i data-lucide="chevron-down" id="cfg-code-chevron" class="w-3.5 h-3.5 transition-transform"></i>
              <i data-lucide="file-code-2" class="w-3.5 h-3.5"></i>
              config.toml 预览 · 当前节
            </button>
            <div class="code-block mt-2" id="cfg-code-body"></div>
          </div>

          <!-- 未保存改动条 -->
          <div id="dirty-bar" style="display:none" class="mt-5 flex items-center gap-3 rounded-lg border border-warn/30 bg-warn/5 px-4 py-2.5 sticky bottom-2">
            <i data-lucide="circle-alert" class="w-4 h-4 text-warn shrink-0"></i>
            <span class="text-sm flex-1">有未保存改动</span>
            <button id="cfg-reset" class="btn btn-ghost btn-sm">恢复默认</button>
            <button id="cfg-save" class="btn btn-primary btn-sm">保存更改</button>
          </div>
        </section>
      </div>
    </div>
  `,
  onMount: function (app) {
    var root = document.getElementById('config-root');
    var titleMap = {
      general: ['通用', '[app] / [session] / [hotkey] · theme 外多为重启生效'],
      voice: ['语音', '[voice] · 换模型需重载子系统'],
      audio: ['音频', '[audio] · 设备与采样率需重开设备'],
      llm: ['大模型', '[llm] · 热加载即时生效'],
      tools: ['工具', '[agent] 预算与内置工具默认风险级'],
      appearance: ['外观', '[ball] / [panel] · 均热加载'],
      privacy: ['隐私', '[privacy] / [memory] · 硬编码项只读']
    };

    // 每节 toml 预览行（innerHTML，已带语法色 token）
    var CODE = {
      general: [
        '<span class="tok-com"># Wisp 配置 · 通用</span>',
        '<span class="tok-key">[app]</span>',
        '<span class="tok-key">language</span> = <span class="tok-str">"zh-CN"</span>  <span class="tok-com"># 界面语言</span>',
        '<span class="tok-key">autostart</span> = true',
        '<span class="tok-key">single_instance</span> = true',
        '',
        '<span class="tok-key">[session]</span>',
        '<span class="tok-key">warm_timeout_sec</span> = <span class="tok-num">90</span>',
        '<span class="tok-key">conversation_idle_sec</span> = <span class="tok-num">30</span>',
        '',
        '<span class="tok-key">[hotkey]</span>',
        '<span class="tok-key">summon</span> = <span class="tok-str">"Alt+Space"</span>'
      ],
      voice: [
        '<span class="tok-com"># Wisp 配置 · 语音链路</span>',
        '<span class="tok-key">[voice]</span>',
        '<span class="tok-key">enabled</span> = true',
        '<span class="tok-key">conversation_mode</span> = false  <span class="tok-com"># 陪聊全双工</span>',
        '',
        '<span class="tok-key">[voice.asr]</span>',
        '<span class="tok-key">provider</span> = <span class="tok-str">"local-sherpa"</span>',
        '<span class="tok-key">model</span> = <span class="tok-str">"sherpa-paraformer-zh-20240323"</span>',
        '<span class="tok-key">punctuation</span> = true',
        '',
        '<span class="tok-key">[voice.tts]</span>',
        '<span class="tok-key">voice</span> = <span class="tok-str">"aina-low"</span>',
        '<span class="tok-key">speed</span> = <span class="tok-num">1.0</span>',
        '',
        '<span class="tok-key">[voice.aec]</span>',
        '<span class="tok-key">enabled</span> = true'
      ],
      audio: [
        '<span class="tok-com"># Wisp 配置 · 音频设备</span>',
        '<span class="tok-key">[audio]</span>',
        '<span class="tok-key">input_device</span> = <span class="tok-str">"default"</span>',
        '<span class="tok-key">sample_rate</span> = <span class="tok-num">16000</span>',
        '<span class="tok-key">vad_threshold</span> = <span class="tok-num">0.5</span>',
        '<span class="tok-key">half_duplex</span> = true  <span class="tok-com"># 干活路径硬编码</span>',
        '<span class="tok-key">mic_muted_default</span> = true'
      ],
      llm: [
        '<span class="tok-com"># Wisp 配置 · 大模型（热加载）</span>',
        '<span class="tok-key">[llm]</span>',
        '<span class="tok-key">provider</span> = <span class="tok-str">"deepseek"</span>',
        '<span class="tok-key">model</span> = <span class="tok-str">"deepseek-chat"</span>',
        '<span class="tok-key">api_base</span> = <span class="tok-str">"https://api.deepseek.com/v1"</span>',
        '<span class="tok-com"># api_key 不落 TOML，仅存 DPAPI 凭据</span>',
        '<span class="tok-key">temperature</span> = <span class="tok-num">0.7</span>',
        '<span class="tok-key">max_tokens</span> = <span class="tok-num">4096</span>',
        '<span class="tok-key">top_p</span> = <span class="tok-num">1.0</span>',
        '<span class="tok-key">timeout_seconds</span> = <span class="tok-num">30</span>',
        '<span class="tok-key">stream</span> = true',
        '<span class="tok-key">enable_cache</span> = true'
      ],
      tools: [
        '<span class="tok-com"># Wisp 配置 · Agent 预算</span>',
        '<span class="tok-key">[agent]</span>',
        '<span class="tok-key">max_rounds</span> = <span class="tok-num">50</span>',
        '<span class="tok-key">token_budget</span> = <span class="tok-num">200000</span>',
        '<span class="tok-key">per_tool_timeout_ms</span> = <span class="tok-num">30000</span>',
        '<span class="tok-key">steering_enabled</span> = true',
        '<span class="tok-com"># 工具默认风险级见 SPEC-06 矩阵</span>'
      ],
      appearance: [
        '<span class="tok-com"># Wisp 配置 · 外观</span>',
        '<span class="tok-key">[app]</span>',
        '<span class="tok-key">theme</span> = <span class="tok-str">"auto"</span>',
        '',
        '<span class="tok-key">[ball]</span>',
        '<span class="tok-key">size</span> = <span class="tok-num">56</span>',
        '<span class="tok-key">opacity_idle</span> = <span class="tok-num">0.35</span>',
        '<span class="tok-key">click_through</span> = true',
        '',
        '<span class="tok-key">[panel]</span>',
        '<span class="tok-key">font_size</span> = <span class="tok-num">13</span>',
        '<span class="tok-key">width</span> = <span class="tok-num">640</span>'
      ],
      privacy: [
        '<span class="tok-com"># Wisp 配置 · 隐私与记忆</span>',
        '<span class="tok-key">[privacy]</span>',
        '<span class="tok-key">redact_paths</span> = false',
        '<span class="tok-key">retention_days</span> = <span class="tok-num">30</span>',
        '<span class="tok-key">diagnostics_opt_in</span> = false',
        '<span class="tok-com"># keep_transcript / keep_audio 硬编码关闭</span>',
        '',
        '<span class="tok-key">[memory]</span>',
        '<span class="tok-key">l1_max</span> = <span class="tok-num">20</span>'
      ]
    };

    var skeleton = root.querySelector('#cfg-skeleton');
    var sectionsWrap = root.querySelector('#cfg-sections');
    var codeBody = root.querySelector('#cfg-code-body');
    var search = root.querySelector('#cfg-search');
    var currentTarget = 'llm';
    var revealTimer = null;
    var SEARCH_MS = 120;
    var codeOpen = true;

    function paintNav(target) {
      root.querySelectorAll('.cfg-nav-item').forEach(function (n) {
        var on = n.dataset.target === target;
        n.classList.toggle('bg-accent', on);
        n.classList.toggle('text-accent-foreground', on);
        n.classList.toggle('font-medium', on);
        n.classList.toggle('text-muted-foreground', !on);
      });
    }

    function setTitle(target) {
      var t = titleMap[target];
      if (t) {
        document.getElementById('cfg-title').textContent = t[0];
        document.getElementById('cfg-sub').textContent = t[1];
      }
    }

    function renderCode(target) {
      var lines = CODE[target] || [];
      var h = '';
      lines.forEach(function (inner, i) {
        h += '<div class="code-line"><span class="line-num">' + (i + 1) + '</span><span>' + inner + '</span></div>';
      });
      codeBody.innerHTML = h;
    }

    function animateRows(rows) {
      rows.forEach(function (r, i) {
        r.classList.remove('animated-list-item');
        void r.offsetWidth;
        r.style.animationDelay = (i * 45) + 'ms';
        r.classList.add('animated-list-item');
      });
    }

    function reveal(target, q) {
      skeleton.style.display = 'none';
      sectionsWrap.style.display = '';
      var anim = [];
      root.querySelectorAll('.cfg-section').forEach(function (sec) {
        var visibleRows = 0;
        var secAnim = [];
        sec.querySelectorAll('.cfg-row').forEach(function (row) {
          var hit = !q || (row.dataset.search || '').toLowerCase().indexOf(q) !== -1;
          row.style.display = hit ? '' : 'none';
          if (hit) { visibleRows++; secAnim.push(row); }
        });
        var secVisible;
        if (q) {
          secVisible = !!visibleRows;
          sec.style.display = secVisible ? '' : 'none';
        } else {
          secVisible = (sec.dataset.section === target);
          sec.style.display = secVisible ? '' : 'none';
        }
        // 仅对当前可见节做入场动画，延迟按本节独立编号（隐藏节不参与）
        if (secVisible) {
          animateRows(secAnim);
          secAnim.forEach(function (r) { anim.push(r); });
        }
      });
      renderCode(target);
      app.refreshIcons();
    }

    function flashTo(target, q) {
      if (revealTimer) clearTimeout(revealTimer);
      sectionsWrap.style.display = 'none';
      skeleton.style.display = 'block';
      revealTimer = setTimeout(function () { reveal(target, q); }, SEARCH_MS);
    }

    function switchTo(target) {
      currentTarget = target;
      paintNav(target);
      setTitle(target);
      var q = search.value.trim().toLowerCase();
      reveal(target, q);
    }

    // 左侧节导航切换
    root.querySelectorAll('.cfg-nav-item').forEach(function (item) {
      item.addEventListener('click', function () { switchTo(item.dataset.target); });
    });

    // 配置搜索过滤（带骨架过渡）
    var searchDebounce = null;
    search.addEventListener('input', function () {
      if (searchDebounce) clearTimeout(searchDebounce);
      searchDebounce = setTimeout(function () {
        var q = search.value.trim().toLowerCase();
        flashTo(currentTarget, q);
      }, 200);
    });

    // toml 预览折叠
    var chevron = root.querySelector('#cfg-code-chevron');
    root.querySelector('#cfg-code-toggle').addEventListener('click', function () {
      codeOpen = !codeOpen;
      codeBody.style.display = codeOpen ? '' : 'none';
      chevron.style.transform = codeOpen ? '' : 'rotate(-90deg)';
    });

    // 未保存改动追踪
    var bar = root.querySelector('#dirty-bar');
    function markDirty() { bar.style.display = 'flex'; }
    root.querySelectorAll('.cfg-row input, .cfg-row select').forEach(function (el) {
      el.addEventListener('input', markDirty);
      el.addEventListener('change', markDirty);
    });
    root.querySelectorAll('.switch').forEach(function (sw) {
      if (sw.dataset.static === '1') return;
      sw.addEventListener('click', function () { sw.classList.toggle('on'); markDirty(); });
    });

    // temperature slider 读数联动
    var temp = root.querySelector('#cfg-temp');
    var tempVal = root.querySelector('#cfg-temp-val');
    temp.addEventListener('input', function () {
      tempVal.textContent = parseFloat(temp.value).toFixed(1);
      markDirty();
    });

    // api_base 校验错误态
    var apiBase = root.querySelector('#cfg-api-base');
    var err = root.querySelector('#cfg-api-base-err');
    function validate() {
      if (!apiBase.value.trim()) {
        apiBase.style.borderColor = 'hsl(var(--destructive))';
        err.style.display = 'block';
        return false;
      }
      apiBase.style.borderColor = '';
      err.style.display = 'none';
      return true;
    }
    apiBase.addEventListener('input', function () { validate(); markDirty(); });
    validate();

    // 保存 / 恢复默认
    root.querySelector('#cfg-save').addEventListener('click', function () {
      if (!validate()) { app.toast('api_base 不能为空，未保存'); return; }
      bar.style.display = 'none';
      app.toast('配置已写回 config.toml（演示）');
    });
    root.querySelector('#cfg-reset').addEventListener('click', function () {
      bar.style.display = 'none';
      app.toast('已从 Go 结构体重载默认值（演示）');
    });

    // 首屏：立即渲染 llm 节（无人为延迟）
    reveal('llm', '');
    BUSelect.init(root);
  }
});
