import {type SearchResult} from "minisearch";

const termRegExp = (term) => new RegExp(`(?<=^|[^\\p{L}])${term}(?=[^\\p{L}]|$)`, 'gui')

export function highlightMatches(result: SearchResult, classList = 'bg-primary-200 text-gray-900 dark:text-gray-900', surroundingLength = 50) {
  let fieldMatches = {}
  for (const term in result.match) {
    result.match[term].map(field =>
      (fieldMatches[field] ?? (fieldMatches[field] = [])).push(term)
    )
  }
  for (const field in fieldMatches) {
    let bestFragment = ''
    let bestFragmentTerms = []
    let bestFragmentMark = []
    const content = result[field]
    for (const term of fieldMatches[field]) {
      for (const position of [...content.matchAll(termRegExp(term))].map(match => match.index)) {
        let fragmentStart = Math.max(0, position - surroundingLength)
        if (fragmentStart > 0) {
          fragmentStart = content.lastIndexOf(' ', fragmentStart) + 1 || 0
        }
        let fragmentEnd = Math.min(content.length, position + term.length + surroundingLength)
        if (fragmentEnd < content.length) {
          const nextSpace = content.indexOf(' ', fragmentEnd);
          fragmentEnd = nextSpace === -1 ? content.length : nextSpace;
        }
        const fragment = content.slice(fragmentStart, fragmentEnd)
        let fragmentTerms = [term]
        let fragmentMark = [position - fragmentStart, position - fragmentStart + term.length]
        if (bestFragment === '') {
          bestFragment = fragment
          bestFragmentTerms = fragmentTerms
          bestFragmentMark = fragmentMark
        }
        for (const t of fieldMatches[field].filter(t => t !== term)) {
          for (const position of [...fragment.matchAll(termRegExp(t))].map(match => match.index)) {
            fragmentTerms.push(t)
            if (position < fragmentMark[0]) fragmentMark[0] = position
            if (position + t.length > fragmentMark[1]) fragmentMark[1] = position + t.length
          }
          if ((fragmentTerms.length > bestFragmentTerms.length) ||
            (fragmentTerms.length === bestFragmentTerms.length && fragmentTerms.length > bestFragmentTerms.length)) {
            bestFragment = fragment
            bestFragmentTerms = fragmentTerms
            bestFragmentMark = fragmentMark
          }
        }
      }
    }
    if (bestFragmentMark.length === 2) {
      const head = bestFragment.slice(0, bestFragmentMark[0]);
      const mark = bestFragment.slice(bestFragmentMark[0], bestFragmentMark[1]);
      const tail = bestFragment.slice(bestFragmentMark[1]);
      result[`${field}_match`] = `${head}<searchmark class="${classList}">${mark}</searchmark>${tail}`
    } else {
      if (field === 'content') {
        result[`${content}_match`] = result.content?.slice(0, 2 * surroundingLength)
      }
    }
  }
  return result
}