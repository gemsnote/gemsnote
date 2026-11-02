# Markdown Mermaid 预览

Markdown 笔记使用标记为 `mermaid` 的围栏代码块。预览显示图表，编辑、同步、导入和 Markdown 导出保留原文。富文本中的代码不会被当作 Mermaid 执行。

````markdown
```mermaid
flowchart LR
    A[记录想法] --> B[整理笔记]
    B --> C[积累知识]
```
````

图表适应正文宽度，可点击放大、调节缩放，并展开“图表源码”选择复制。语法错误或超出限制时显示提示与源码，不影响其余正文。内置中文和英文提示。

## 渲染与安全边界

- 依赖锁定 Mermaid 11.17.2，资源全部随构建发布，不访问 CDN，不向服务端上传渲染内容。
- `frontend/src/components/NotePreview.vue` 处理普通正文和 Mermaid 槽位；`MermaidDiagram.vue` 负责图表展示，`mermaidPreview.ts` 管理隔离 iframe 和请求生命周期。
- `frontend/mermaid/renderer.ts` 是离线渲染器源码，使用严格安全模式、纯 SVG 标签和 DOMPurify 清理。最终图表通过 SVG 图片显示，不把图表生成的 HTML 注入主页面。
- 渲染 iframe 只有 `allow-scripts`，不开放同源权限。资源文档的 CSP 拒绝网络连接、外部图片/字体、表单提交与外部脚本。渲染器不接收账号、token 或客户端 bridge。
- 禁用笔记内的 Mermaid 配置指令和配置 frontmatter，不开放 HTML 标签、外部图像或点击链接。
- 每个预览实例最多处理 32 个图块，每块最多 20,000 字符、500 行，图的连边上限为 200；超过限制保留源码并提示。加载超时 20 秒，单次绘图等待上限 30 秒。复杂图仍可能较慢，这些限制不能等同于严格的 CPU 隔离。
- 图表背景固定白色，主题固定为 default；三端使用同一渲染版本与配置，系统字体差异可能导致文字布局略有不同。

## 构建与维护

```bash
npm ci --prefix frontend
npm run build --prefix frontend
npm test --prefix frontend
node frontend/tests/mermaid.mjs
```

`prebuild` / `predev` 自动运行 `frontend/scripts/build-mermaid.mjs`，生成
`frontend/public/assets/mermaid/index.html`、`LICENSE` 和 `manifest.json`。
Vite 将它们复制到 `frontend/dist/assets/mermaid/`；现有 `/assets/` 路由直接提供这些文件。
生成文件不提交到服务端 Git。Docker / Release 构建使用同一流程，无需增加服务端进程、API2 或数据库迁移。

构建器把依赖和延迟加载模块打包成单文件，避免离线加载额外资源。`manifest.json` 记录 Mermaid 版本与 HTML 的 SHA-256，包内保留许可证注释。
升级时固定依赖版本并更新 lockfile，重新生成资源后运行浏览器回归，不能只依据 Node 构建成功判断浏览器兼容。

客户端的构建与接入分别见[桌面开发文档](https://github.com/gemsnote/gemsnote-desktop/blob/main/docs/DEVELOPMENT.md)及移动端仓库的 `docs/DEVELOPMENT.md`。客户端不需要新增服务端接口。
