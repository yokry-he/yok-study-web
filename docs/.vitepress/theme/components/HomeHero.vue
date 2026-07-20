<template>
  <main class="custom-home">
    <section class="custom-home__hero">
      <div>
        <span class="custom-home__eyebrow">Learning Path · Docs · Practice · Troubleshooting</span>
        <h1 class="custom-home__title">
          一座面向程序员的
          <span class="custom-home__title-accent">技术学习工作台</span>
        </h1>
        <p class="custom-home__lead">
          用图解、核心章节、从零到项目、专项练习和真实问题库打通学习闭环，覆盖前端、框架、后端、数据、工程化与 AI 工程。
        </p>
        <div class="custom-home__actions">
          <a class="custom-home__action custom-home__action--primary" href="/roadmap/introduction">开始学习路线</a>
          <a class="custom-home__action custom-home__action--secondary" href="/technologies/">浏览技术库</a>
        </div>
      </div>

      <aside class="custom-home__visual" aria-label="技术学习路线预览">
        <div class="custom-home__visual-header">
          <span class="custom-home__dot custom-home__dot--mint"></span>
          <span class="custom-home__dot custom-home__dot--lemon"></span>
          <span class="custom-home__dot custom-home__dot--rose"></span>
          <span class="custom-home__panel-title">learning-roadmap.md</span>
        </div>
        <div class="custom-home__route">
          <div class="custom-home__route-item">
            <span class="custom-home__route-index">01</span>
            <span>
              <strong class="custom-home__route-title">基础能力</strong>
              <span class="custom-home__route-copy">HTML、CSS、浏览器、网络</span>
            </span>
          </div>
          <div class="custom-home__route-item">
            <span class="custom-home__route-index">02</span>
            <span>
              <strong class="custom-home__route-title">语言与框架</strong>
              <span class="custom-home__route-copy">JavaScript、TypeScript、Vue</span>
            </span>
          </div>
          <div class="custom-home__route-item">
            <span class="custom-home__route-index">03</span>
            <span>
              <strong class="custom-home__route-title">工程实践</strong>
              <span class="custom-home__route-copy">Vite、规范、环境、构建部署</span>
            </span>
          </div>
          <div class="custom-home__route-item">
            <span class="custom-home__route-index">04</span>
            <span>
              <strong class="custom-home__route-title">问题沉淀</strong>
              <span class="custom-home__route-copy">实战项目、问题库、速查手册</span>
            </span>
          </div>
        </div>
      </aside>
    </section>

    <section class="custom-home__section">
      <h2 class="custom-home__section-title">当前核心模块</h2>
      <TechGrid :items="[
        { title: '阅读顺序与使用方法', description: '说明如何按路线、图解、项目、练习、问题库和速查手册组合使用本站。', link: '/roadmap/reading-guide', level: '导读' },
        { title: '学习路线总览', description: '按前端、后端、全栈、AI 工程等目标组织可执行学习路径。', link: '/roadmap/introduction', level: '路线' },
        { title: '前端基础模块', description: '图解 HTML、表单、图片、无障碍和浏览器页面链路，并用原生项目与故障注入完成验收。', link: '/frontend/introduction', level: '基础' },
        { title: 'JavaScript 模块', description: '覆盖类型、函数、数组对象、异步和模块化工程实践。', link: '/javascript/fundamentals', level: '语言' },
        { title: 'Vue 核心文档', description: '覆盖模板、响应式、组件、路由、Pinia、表单、性能和测试。', link: '/vue/introduction', level: '框架' },
        { title: '学习路径练习包', description: '用静态页面、JS 列表、TS 类型、Vue Admin、数据库和上线练习验证能力。', link: '/roadmap/practice-labs', level: '练习' },
        { title: '工程化文档', description: '整理 Vite、规范、环境配置、构建部署和上线排错。', link: '/engineering/vite', level: '工程' },
        { title: 'Vue Admin 实战', description: '串联登录、权限、菜单、表格、表单、接口封装和部署。', link: '/projects/vue-admin', level: '实战' },
        { title: '真实项目问题库', description: '按现象、根因、解决方案和预防方式沉淀项目问题。', link: '/projects/real-world-issues', level: '排错' }
      ]" />
    </section>
  </main>
</template>
