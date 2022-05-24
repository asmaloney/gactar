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

export { commentString }
