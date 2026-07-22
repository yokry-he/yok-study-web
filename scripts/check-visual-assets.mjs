import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const rootDir = path.resolve(process.argv[2] ?? process.cwd())
const docsDir = path.join(rootDir, 'docs')
const publicDir = path.join(docsDir, 'public')
const registerPath = path.join(docsDir, 'contribute', 'visual-asset-register.md')

const componentRe = /<DocFigure\s+([\s\S]*?)\/>/g
const propRe = /(?:^|\s)(src|alt|caption|source-url)="([^"]*)"/g
const markdownImageRe = /!\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)/g
const assetRe = /<!--\s*asset:\s*(\/images\/[^|\s]+)\s*\|([\s\S]*?)-->/g
const allowedExtensionRe = /\.(?:png|jpe?g|webp)$/i
const validFilenameRe = /^[a-z0-9]+(?:-[a-z0-9]+)*\.(?:png|jpe?g|webp)$/
const normalSizeLimit = 500 * 1024
const approvedLargeSizeLimit = 1536 * 1024

export function checkVisualAssets({ docsDir, publicDir, registerPath }) {
  const errors = []
  const referenced = new Map()
  const registered = readRegister(registerPath, errors)

  for (const markdownPath of walkFiles(docsDir, filePath => filePath.endsWith('.md'))) {
    const relativePath = path.relative(docsDir, markdownPath)
    const markdown = stripFencedCode(fs.readFileSync(markdownPath, 'utf8'))

    for (const match of markdown.matchAll(markdownImageRe)) {
      errors.push(`禁止使用裸 Markdown 图片：${relativePath} -> ${match[2]}`)
    }

    for (const match of markdown.matchAll(componentRe)) {
      const props = parseProps(match[1])
      for (const required of ['src', 'alt', 'caption']) {
        if (!props[required]?.trim()) {
          errors.push(`DocFigure 缺少非空 ${required}：${relativePath}`)
        }
      }

      const source = props.src?.trim()
      if (!source || isExternalSource(source)) continue
      if (!source.startsWith('/images/')) {
        errors.push(`本地图片路径必须以 /images/ 开头：${relativePath} -> ${source}`)
        continue
      }
      validateAssetPath(source, relativePath, errors)
      referenced.set(source, relativePath)

      const diskPath = publicPathFor(source, publicDir)
      if (!fs.existsSync(diskPath) || !fs.statSync(diskPath).isFile()) {
        errors.push(`图片文件不存在：${source}`)
      }
      if (!registered.has(source)) {
        errors.push(`图片未登记：${source}`)
      }
    }
  }

  for (const imagePath of walkFiles(path.join(publicDir, 'images'), filePath => allowedExtensionRe.test(filePath))) {
    const source = `/${path.relative(publicDir, imagePath).split(path.sep).join('/')}`
    validateAssetPath(source, path.relative(publicDir, imagePath), errors)
    if (!registered.has(source)) {
      errors.push(`图片未登记：${source}`)
      continue
    }
    validateFileSize(source, imagePath, registered.get(source), errors)
  }

  for (const [source, metadata] of registered) {
    const diskPath = publicPathFor(source, publicDir)
    if (!fs.existsSync(diskPath) || !fs.statSync(diskPath).isFile()) {
      errors.push(`登记图片文件不存在：${source}`)
      continue
    }
    validateAssetPath(source, registerPath, errors)
    validateFileSize(source, diskPath, metadata, errors)
  }

  return [...new Set(errors)].sort((left, right) => left.localeCompare(right, 'zh-CN'))
}

function walkFiles(root, include) {
  if (!fs.existsSync(root)) return []
  const files = []
  const entries = fs.readdirSync(root, { withFileTypes: true })
  for (const entry of entries) {
    const entryPath = path.join(root, entry.name)
    if (entry.isDirectory()) {
      files.push(...walkFiles(entryPath, include))
    } else if (entry.isFile() && include(entryPath)) {
      files.push(entryPath)
    }
  }
  return files.sort()
}

function stripFencedCode(markdown) {
  const lines = markdown.split('\n')
  let fence = ''
  return lines.map(line => {
    const opening = line.match(/^\s*(`{3,}|~{3,})/)
    if (!fence && opening) {
      fence = opening[1][0]
      return ''
    }
    if (fence && new RegExp(`^\\s*\\${fence}{3,}\\s*$`).test(line)) {
      fence = ''
      return ''
    }
    return fence ? '' : line
  }).join('\n')
}

function parseProps(raw) {
  return Object.fromEntries([...raw.matchAll(propRe)].map(match => [match[1], match[2]]))
}

function readRegister(registerPath, errors) {
  const registered = new Map()
  if (!fs.existsSync(registerPath)) return registered
  const markdown = stripFencedCode(fs.readFileSync(registerPath, 'utf8'))
  for (const match of markdown.matchAll(assetRe)) {
    const source = match[1].trim()
    if (registered.has(source)) {
      errors.push(`资产重复登记：${source}`)
      continue
    }
    registered.set(source, match[2].trim())
  }
  return registered
}

function isExternalSource(source) {
  return /^https?:\/\//i.test(source)
}

function publicPathFor(source, publicDir) {
  return path.join(publicDir, source.replace(/^\//, ''))
}

function validateAssetPath(source, location, errors) {
  if (!allowedExtensionRe.test(source)) {
    errors.push(`图片扩展名不受支持：${source}`)
  }
  const filename = path.posix.basename(source)
  if (!validFilenameRe.test(filename)) {
    errors.push(`图片文件名必须使用小写 kebab-case：${location} -> ${filename}`)
  }
}

function validateFileSize(source, diskPath, metadata, errors) {
  const size = fs.statSync(diskPath).size
  const largeApproved = metadata?.includes('large-approved') ?? false
  const limit = largeApproved ? approvedLargeSizeLimit : normalSizeLimit
  if (size > limit) {
    errors.push(`图片超过大小限制：${source} (${size} > ${limit} bytes)`)
  }
}

if (fileURLToPath(import.meta.url) === path.resolve(process.argv[1])) {
  const errors = checkVisualAssets({ docsDir, publicDir, registerPath })
  if (errors.length > 0) {
    console.error(errors.join('\n'))
    process.exit(1)
  }
  console.log('视觉资产检查通过')
}
