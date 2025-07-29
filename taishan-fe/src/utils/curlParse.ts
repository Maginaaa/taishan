import { HttpCaseExtend } from '@/views/Plan/scene/case-data'

export function curlParse(s: string): HttpCaseExtend | undefined {
  if (0 != s.indexOf('curl ')) return
  const args = rewrite(shellSplit(s))
  const out: HttpCaseExtend | any = {
    url: '',
    method_type: 'GET',
    headers_form: [],
    body: {
      body_type: 'json',
      body_value: ''
    }
  }
  let state = ''
  args.forEach((arg) => {
    switch (true) {
      case isURL(arg):
        out.url = arg
        break

      case arg == '-A' || arg == '--user-agent':
        state = 'user-agent'
        break

      case arg == '-H' || arg == '--header':
        state = 'header'
        break

      case arg == '--data-urlencode':
        state = 'data-urlencode'
        break

      case arg == '--form':
        state = 'form'
        break

      case arg == '-d' ||
        arg == '--data' ||
        arg == '--data-ascii' ||
        arg == '--data-raw' ||
        arg == '--data-binary':
        state = 'data'
        break

      case arg == '-u' || arg == '--user':
        state = 'user'
        break

      case arg == '-I' || arg == '--head':
        out.method_type = 'HEAD'
        break

      case arg == '-X' || arg == '--request':
        state = 'method'
        break

      case arg == '-b' || arg == '--cookie':
        state = 'cookie'
        break

      case arg == '--compressed':
        out.headers_form.push({
          enable: true,
          key: 'Accept-Encoding',
          value: 'deflate, gzip'
        })
        break

      case !!arg:
        switch (state) {
          case 'header':
            const field = parseField(arg)
            out.headers_form.push({
              enable: true,
              key: field[0],
              value: field[1]
            })
            state = ''
            break
          case 'data-urlencode':
            const index = arg.indexOf('=')
            if (index === -1) {
              state = ''
              break
            } else {
              // const left = arg.substring(0, index)
              // const right = arg.substring(index + 1)
              // if (!out.dataUrlencode) out.dataUrlencode = {}
              // out.dataUrlencode[left] = right
              // state = ''
              break
            }
          case 'form':
            const index1 = arg.indexOf('=')
            if (index1 === -1) {
              state = ''
              break
            } else {
              // const left = arg.substring(0, index1)
              // const right = arg.substring(index1 + 1)
              // console.log(left, right)
              // if (!out.form) out.form = {}
              // out.form[left] = right
              // state = ''
              break
            }

          case 'user-agent':
            out.headers_form.push({
              enable: true,
              key: 'User-Agent',
              value: arg
            })
            state = ''
            break
          case 'data':
            if (out.method_type === 'GET' || out.method_type === 'HEAD') out.method_type = 'POST'
            out.headers_form.push({
              enable: true,
              key: 'Content-Type',
              value: 'application/json'
            })
            out.body = {
              body_type: 'json',
              body_value: unescapeJsonString(arg)
            }
            state = ''
            break
          case 'user':
            out.headers_form.push({
              enable: true,
              key: 'Authorization',
              value: 'Basic ' + btoa(arg)
            })
            state = ''
            break
          case 'method':
            out.method_type = arg
            state = ''
            break
          case 'cookie':
            out.headers_form.push({
              enable: true,
              key: 'Cookie',
              value: arg
            })
            state = ''
            break
        }
        break
    }
  })
  return out
}

function rewrite(args) {
  return args.reduce((args, a) => {
    if (0 == a.indexOf('-X')) {
      args.push('-X')
      args.push(a.slice(2))
    } else {
      args.push(a)
    }

    return args
  }, [])
}

/**
 * Parse header field.
 */

function parseField(s) {
  return s.split(/: (.+)/)
}

/**
 * Check if `s` looks like a url.
 */

function isURL(s) {
  return /^https?:\/\//.test(s)
}

function shellSplit(line) {
  let field = ''
  if (line === null) {
    line = ''
  }
  const words: any[] = []
  field = ''
  scan(
    line,
    /\s*(?:([^\s\\\'\"]+)|'((?:[^\'\\]|\\.)*)'|"((?:[^\"\\]|\\.)*)"|(\\.?)|(\S))(\s|$)?/,
    function (match) {
      const word = match[1]
      const sq = match[2]
      const dq = match[3]
      const escape = match[4]
      const garbage = match[5]
      const seperator = match[6]
      if (garbage !== undefined) {
        throw new Error('Unmatched quote')
      }
      field += word || (sq || dq || escape).replace(/\\(?=.)/, '')
      if (seperator !== undefined) {
        words.push(field)
        field = ''
      }
    }
  )
  if (field) {
    words.push(field)
  }
  return words
}

function scan(string, pattern, callback) {
  let match, result
  result = ''
  while (string.length > 0) {
    match = string.match(pattern)
    if (match) {
      result += string.slice(0, match.index)
      result += callback(match)
      string = string.slice(match.index + match[0].length)
    } else {
      result += string
      string = ''
    }
  }
  return result
}

function unescapeJsonString(str) {
  if (!str) return str

  // 处理JSON字符串中的常见转义字符
  return str
    .replace(/\\"/g, '"') // 反转义双引号
    .replace(/\\'/g, "'") // 反转义单引号
    .replace(/\\\\/g, '\\') // 反转义反斜杠
    .replace(/\\n/g, '\n') // 反转义换行符
    .replace(/\\r/g, '\r') // 反转义回车符
    .replace(/\\t/g, '\t') // 反转义制表符
}
