import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const p0Modules = new Set(['frontend', 'css', 'vue', 'browser'])
const p1Modules = new Set([
  'javascript', 'typescript', 'react', 'meta-frameworks', 'node', 'java',
  'go', 'database', 'engineering', 'devops', 'ai-engineering'
])

export function auditDocs(docsDir) {
  const rows = []
  for (const filePath of walkMarkdown(docsDir)) {
    const relativePath = path.relative(docsDir, filePath).split(path.sep).join('/')
    if (relativePath.startsWith('superpowers/') || relativePath === 'contribute/visual-audit.md') {
      continue
    }
    const markdown = fs.readFileSync(filePath, 'utf8')
    const moduleName = relativePath.includes('/') ? relativePath.split('/')[0] : 'root'
    rows.push({
      route: routeFor(relativePath),
      moduleName,
      mermaidCount: countMatches(markdown, /^```mermaid\s*$/gm),
      figureCount: countMatches(markdown, /<DocFigure\b/g),
      codeCount: countMatches(markdown, /^```[^`\n]*$/gm),
      priority: priorityFor(moduleName),
      status: 'review-required'
    })
  }
  return rows.sort((left, right) => left.route.localeCompare(right.route, 'en'))
}
export function renderAudit(rows) {
  const mermaidPages = rows.filter(row => row.mermaidCount > 0).length
  const figurePages = rows.filter(row => row.figureCount > 0).length
  const imageReferences = rows.reduce((total, row) => total + row.figureCount, 0)
  const lines = [
    '# 全站视觉讲解审计',
    '',
    '> 本文件由 `node scripts/audit-doc-visuals.mjs --write docs/contribute/visual-audit.md` 生成。正式内容发生增删后重新运行；人工结论应登记到资产表或对应页面，不直接修改自动表格。',
    '',
    '## 基线',
    '',
    `- 核对日期：2026-07-21。`,
    `- 正式内容文档：${rows.length} 篇；不含 \`docs/superpowers/\` 与本审计文件。`,
    `- 包含 Mermaid 的页面：${mermaidPages} 篇。`,
    `- 包含 DocFigure 的页面：${figurePages} 篇，共 ${imageReferences} 个引用。`,
    '- 初始状态统一为 `review-required`，优先级只表示人工审阅顺序，不代表每页都必须增加位图。',
    '',
    '## 媒介选择规则',
    '',
    '| 结论 | 适用场景 | 完成定义 |',
    '| --- | --- | --- |',
    '| `diagram-sufficient` | 状态、流程、依赖或时序可被 Mermaid 准确表达 | 图与正文一致，步骤解释、误区和自测完整 |',
    '| `needs-live-screenshot` | 必须观察浏览器、DevTools 或真实 UI 的最终结果 | 可复现场景、固定视口、脱敏截图、alt 与图注齐全 |',
    '| `needs-annotated-screenshot` | 截图中必须指出区域、变化或因果关系 | 标注不遮挡关键信息，正文逐一解释标注 |',
    '| `needs-generated-visual` | 抽象心智模型需要类比，且不承担精确规则 | prompt、工具、日期和人工核对结果已登记 |',
    '| `needs-official-source` | 工具界面或标准事实必须引用官方材料 | 来源 URL、许可、核对日期和版本明确 |',
    '| `needs-mermaid-refactor` | 现有图过密、重复或无法在移动端阅读 | 拆图后 SVG 非空、无横向溢出、正文仍可独立理解 |',
    '',
    '## 路由清单',
    '',
    '| Route | 模块 | Mermaid | DocFigure | 代码块 | 优先级 | 状态 |',
    '| --- | --- | ---: | ---: | ---: | --- | --- |'
  ]

  for (const row of rows) {
    lines.push(`| \`${row.route}\` | ${row.moduleName} | ${row.mermaidCount} | ${row.figureCount} | ${row.codeCount} | ${row.priority} | ${row.status} |`)
  }
  lines.push('')
  return lines.join('\n')
}

function walkMarkdown(root) {
  const files = []
  for (const entry of fs.readdirSync(root, { withFileTypes: true })) {
    const entryPath = path.join(root, entry.name)
    if (entry.isDirectory()) {
      files.push(...walkMarkdown(entryPath))
    } else if (entry.isFile() && entry.name.endsWith('.md')) {
      files.push(entryPath)
    }
  }
  return files
}

function routeFor(relativePath) {
  const withoutExtension = relativePath.replace(/\.md$/, '')
  if (withoutExtension === 'index') return '/'
  if (withoutExtension.endsWith('/index')) {
    return `/${withoutExtension.slice(0, -'/index'.length)}`
  }
  return `/${withoutExtension}`
}

function priorityFor(moduleName) {
  if (p0Modules.has(moduleName)) return 'P0'
  if (p1Modules.has(moduleName)) return 'P1'
  return 'P2'
}

function countMatches(markdown, pattern) {
  return [...markdown.matchAll(pattern)].length
}

if (fileURLToPath(import.meta.url) === path.resolve(process.argv[1])) {
  const rootDir = process.cwd()
  const docsDir = path.join(rootDir, 'docs')
  const output = renderAudit(auditDocs(docsDir))
  const writeIndex = process.argv.indexOf('--write')
  if (writeIndex >= 0) {
    const outputPath = process.argv[writeIndex + 1]
    if (!outputPath) {
      console.error('--write 必须提供输出路径')
      process.exit(1)
    }
    fs.writeFileSync(path.resolve(outputPath), output)
  } else {
    process.stdout.write(output)
  }
}
