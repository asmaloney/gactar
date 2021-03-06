import { Issue, IssueList } from './api'

function commentString(language: string, text: string): string {
  let comment = ''

  switch (language) {
    case 'commonlisp':
      comment = ';;;'
      break

    case 'python':
      comment = '#'
      break

    default:
      return text
  }

  return `${comment} ${text}`
}

// issueToText takes an Issue and formats it for output.
function issueToText(issue: Issue): string {
  let text = `${issue.level}: ${issue.text}`

  if (issue.location) {
    text += `  (line ${issue.location.line}`
    if (issue.location.columnStart != 0 || issue.location.columnEnd != 0) {
      text += `, col ${issue.location.columnStart}`
      if (issue.location.columnEnd != issue.location.columnStart) {
        text += `-${issue.location.columnEnd}`
      }
    }
    text += ')'
  }

  return text
}

// issuesToArray takes an IssueList, formats the issues, and returns them as an array or strings.
function issuesToArray(list: IssueList): string[] {
  const issueTexts: string[] = []

  list.forEach((issue: Issue) => {
    const text = issueToText(issue)

    issueTexts.push(text)
  })

  return issueTexts
}

export { commentString, issuesToArray, issueToText }
