export const extractPath = (url: string): string => {
  if (!url) return '/'
  const pathStart =
    url.indexOf('://') !== -1
      ? url.indexOf('/', url.indexOf('://') + 3)
      : url.startsWith('//')
        ? url.indexOf('/', 2)
        : url.startsWith('/')
          ? 0
          : -1

  if (pathStart === -1) return '/'

  let pathEnd = url.length
  const queryIndex = url.indexOf('?', pathStart)
  const hashIndex = url.indexOf('#', pathStart)

  if (queryIndex > 0) pathEnd = Math.min(pathEnd, queryIndex)
  if (hashIndex > 0) pathEnd = Math.min(pathEnd, hashIndex)

  const path = url.slice(pathStart, pathEnd)
  return path || '/'
}
