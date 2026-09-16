import { Node } from '@tiptap/core'
import katex from 'katex'

const InlineMath = Node.create({
  name: 'inlineMath',
  group: 'inline',
  inline: true,
  atom: true,

  addAttributes() {
    return {
      formula: {
        default: ''
      }
    }
  },

  parseHTML() {
    return [{
      tag: 'span[data-inline-math]',
      getAttrs: (dom) => ({
        formula: dom.getAttribute('data-formula') || dom.textContent || ''
      })
    }]
  },

  renderHTML({ node }) {
    return ['span', {
      'data-inline-math': '',
      'data-formula': node.attrs.formula,
      class: 'katex-inline'
    }]
  },

  addNodeView() {
    return ({ node }) => {
      const dom = document.createElement('span')
      dom.className = 'katex-inline'
      dom.setAttribute('data-inline-math', '')
      dom.setAttribute('data-formula', node.attrs.formula)

      try {
        katex.render(node.attrs.formula, dom, {
          displayMode: false,
          throwOnError: false
        })
      } catch {
        dom.textContent = node.attrs.formula
      }

      return { dom }
    }
  },

  addCommands() {
    return {
      insertInlineMath: (formula) => ({ chain }) => {
        return chain()
          .insertContent({ type: this.name, attrs: { formula } })
          .run()
      }
    }
  }
})

export default InlineMath