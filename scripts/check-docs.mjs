import fs from 'node:fs'
import path from 'node:path'

const rootDir = process.cwd()
const docsDir = path.join(rootDir, 'docs')
const configPath = path.join(docsDir, '.vitepress', 'config.ts')

const ignoredDirs = new Set(['.vitepress', 'public'])
const moduleDirsRequiringIntro = new Set([
  'ai-engineering',
  'browser',
  'css',
  'database',
  'devops',
  'engineering',
  'javascript',
  'meta-frameworks',
  'node',
  'react',
  'roadmap',
  'typescript',
  'vue'
])

const moduleDirsRequiringTroubleshooting = new Set([
  'ai-engineering',
  'browser',
  'css',
  'database',
  'devops',
  'engineering',
  'meta-frameworks',
  'node',
  'react',
  'typescript',
  'vue'
])

const requiredLearningSections = [
  ['## 适合谁看'],
  ['## 下一步学习', '## 下一步']
]

function walk(dir) {
  const entries = fs.readdirSync(dir, { withFileTypes: true })
  const files = []

  for (const entry of entries) {
    if (entry.isDirectory()) {
      if (ignoredDirs.has(entry.name)) continue
      files.push(...walk(path.join(dir, entry.name)))
      continue
    }

    if (entry.isFile() && entry.name.endsWith('.md')) {
      files.push(path.join(dir, entry.name))
    }
  }

  return files
}

function toRoute(filePath) {
  const relative = path.relative(docsDir, filePath).replaceAll(path.sep, '/')

  if (relative === 'index.md') return '/'
  if (relative.endsWith('/index.md')) return `/${relative.slice(0, -'index.md'.length)}`

  return `/${relative.slice(0, -'.md'.length)}`
}

function stripHashAndQuery(link) {
  return link.split('#')[0].split('?')[0]
}

function routeExists(route, routes) {
  if (route === '') return true
  if (route === '/') return routes.has('/')

  const normalized = route.endsWith('/') ? route : `${route}/`
  const withoutSlash = route.endsWith('/') ? route.slice(0, -1) : route

  return routes.has(route) || routes.has(normalized) || routes.has(withoutSlash)
}

function extractInternalLinks(content) {
  const links = []
  const markdownLinkRe = /\[[^\]]+\]\((\/[^)\s]+)\)/g
  const componentLinkRe = /link:\s*['"`](\/[^'"`]+)['"`]/g

  for (const match of content.matchAll(markdownLinkRe)) {
    links.push(match[1])
  }

  for (const match of content.matchAll(componentLinkRe)) {
    links.push(match[1])
  }

  return links
}

function reportError(errors, filePath, message) {
  const relative = path.relative(rootDir, filePath).replaceAll(path.sep, '/')
  errors.push(`${relative}: ${message}`)
}

const markdownFiles = walk(docsDir)
const routes = new Set(markdownFiles.map(toRoute))
const errors = []
const warnings = []

for (const filePath of markdownFiles) {
  const content = fs.readFileSync(filePath, 'utf8')
  const route = toRoute(filePath)
  const relative = path.relative(docsDir, filePath).replaceAll(path.sep, '/')
  const topLevel = relative.split('/')[0]

  for (const link of extractInternalLinks(content)) {
    const target = stripHashAndQuery(link)

    if (!routeExists(target, routes)) {
      reportError(errors, filePath, `内部链接不存在：${link}`)
    }
  }

  if (
    moduleDirsRequiringIntro.has(topLevel) &&
    !relative.endsWith('/introduction.md') &&
    !relative.endsWith('/troubleshooting.md')
  ) {
    for (const sectionOptions of requiredLearningSections) {
      if (!sectionOptions.some(section => content.includes(section))) {
        reportError(errors, filePath, `缺少必备章节：${sectionOptions.join(' 或 ')}`)
      }
    }
  }

  if (route !== '/' && !content.startsWith('# ') && !content.startsWith('---')) {
    reportError(errors, filePath, '文档应以一级标题或 frontmatter 开头')
  }
}

for (const moduleName of moduleDirsRequiringIntro) {
  const introRoute = `/${moduleName}/introduction`
  if (!routeExists(introRoute, routes)) {
    errors.push(`docs/${moduleName}: 缺少 introduction.md`)
  }
}

for (const moduleName of moduleDirsRequiringTroubleshooting) {
  const troubleshootingRoute = `/${moduleName}/troubleshooting`
  if (!routeExists(troubleshootingRoute, routes)) {
    errors.push(`docs/${moduleName}: 缺少 troubleshooting.md`)
  }
}

const configContent = fs.readFileSync(configPath, 'utf8')
const configuredRoutes = [...configContent.matchAll(/link:\s*['"`](\/[^'"`]+)['"`]/g)]
  .map(match => stripHashAndQuery(match[1]))
  .filter(link => !link.startsWith('http'))

for (const route of configuredRoutes) {
  if (!routeExists(route, routes)) {
    reportError(errors, configPath, `配置中的路由不存在：${route}`)
  }
}

const moduleCounts = new Map()
for (const filePath of markdownFiles) {
  const relative = path.relative(docsDir, filePath).replaceAll(path.sep, '/')
  const [topLevel] = relative.split('/')
  moduleCounts.set(topLevel, (moduleCounts.get(topLevel) ?? 0) + 1)
}

for (const [moduleName, count] of [...moduleCounts.entries()].sort()) {
  if (count === 1 && !moduleName.endsWith('.md')) {
    warnings.push(`${moduleName}: 当前只有 ${count} 篇文档，请确认是否需要导览和问题库`)
  }
}

if (errors.length > 0) {
  console.error('文档检查失败：')
  for (const error of errors) {
    console.error(`- ${error}`)
  }
  process.exit(1)
}

console.log('文档检查通过')
console.log(`- Markdown 文档：${markdownFiles.length} 篇`)
console.log(`- 内部路由：${routes.size} 个`)
console.log(`- 配置路由：${configuredRoutes.length} 个`)

if (warnings.length > 0) {
  console.log('提示：')
  for (const warning of warnings) {
    console.log(`- ${warning}`)
  }
}
