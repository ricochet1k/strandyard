// Frontmatter parsing utilities

export interface RoleMeta {
  description?: string
}

export interface TemplateMeta {
  role?: string
  priority?: string
  description?: string
  id_prefix?: string
}

export interface ParsedFile<T> {
  frontmatter: T
  body: string
}

export function parseFrontmatter<T>(content: string): ParsedFile<T> {
  const lines = content.split('\n')
  
  // Check if file starts with frontmatter delimiter
  if (lines.length === 0 || lines[0].trim() !== '---') {
    return { frontmatter: {} as T, body: content }
  }
  
  // Find end of frontmatter
  let endIndex = -1
  for (let i = 1; i < lines.length; i++) {
    if (lines[i].trim() === '---') {
      endIndex = i
      break
    }
  }
  
  if (endIndex === -1) {
    return { frontmatter: {} as T, body: content }
  }
  
  // Extract frontmatter and body
  const frontmatterLines = lines.slice(1, endIndex)
  const bodyLines = lines.slice(endIndex + 1)
  
  // Parse YAML frontmatter (simple key-value parsing)
  const frontmatter: any = {}
  for (const line of frontmatterLines) {
    const match = line.match(/^([a-z_]+):\s*(.*)$/i)
    if (match) {
      const [, key, value] = match
      // Remove quotes if present
      const cleanValue = value.replace(/^["']|["']$/g, '').trim()
      frontmatter[key] = cleanValue
    }
  }
  
  return {
    frontmatter: frontmatter as T,
    body: bodyLines.join('\n').trim()
  }
}

export function serializeFrontmatter<T extends Record<string, any>>(
  frontmatter: T,
  body: string
): string {
  const lines = ['---']
  
  for (const [key, value] of Object.entries(frontmatter)) {
    if (value !== undefined && value !== '') {
      // Add quotes if value contains special characters
      const needsQuotes = typeof value === 'string' && /[:#\n]/.test(value)
      const serialized = needsQuotes ? `"${value}"` : value
      lines.push(`${key}: ${serialized}`)
    }
  }
  
  lines.push('---')
  lines.push('')
  lines.push(body)
  
  return lines.join('\n')
}
