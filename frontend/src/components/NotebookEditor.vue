<template>
  <div class="editor-container">
    <div v-if="editor" class="editor-toolbar">
      <button @click="editor.chain().focus().toggleBold().run()" :class="{ active: editor.isActive('bold') }">
        <i class="fas fa-bold"></i>
      </button>
      <button @click="editor.chain().focus().toggleItalic().run()" :class="{ active: editor.isActive('italic') }">
        <i class="fas fa-italic"></i>
      </button>
      <button @click="editor.chain().focus().toggleHeading({ level: 2 }).run()" :class="{ active: editor.isActive('heading', { level: 2 }) }">
        <i class="fas fa-heading"></i>
      </button>
      <button @click="editor.chain().focus().toggleBulletList().run()" :class="{ active: editor.isActive('bulletList') }">
        <i class="fas fa-list-ul"></i>
      </button>
      <button @click="editor.chain().focus().toggleOrderedList().run()" :class="{ active: editor.isActive('orderedList') }">
        <i class="fas fa-list-ol"></i>
      </button>
      <button @click="editor.chain().focus().setTextAlign('left').run()" :class="{ active: editor.isActive({ textAlign: 'left' }) }">
        <i class="fas fa-align-left"></i>
      </button>
      <button @click="editor.chain().focus().setTextAlign('center').run()" :class="{ active: editor.isActive({ textAlign: 'center' }) }">
        <i class="fas fa-align-center"></i>
      </button>
      <button @click="editor.chain().focus().setTextAlign('right').run()" :class="{ active: editor.isActive({ textAlign: 'right' }) }">
        <i class="fas fa-align-right"></i>
      </button>
      <button @click="showTableDialog = true" title="Таблица">
        <i class="fas fa-table"></i>
      </button>
      
      <template v-if="isInTable">
        <button @click="editor.chain().focus().addRowAfter().run()" title="Строку ниже">
          <i class="fas fa-arrow-down"></i>
        </button>
        <button @click="editor.chain().focus().addColumnAfter().run()" title="Колонку справа">
          <i class="fas fa-arrow-right"></i>
        </button>
        <button @click="editor.chain().focus().deleteRow().run()" title="Удалить строку">
          <i class="fas fa-minus"></i> row
        </button>
        <button @click="editor.chain().focus().deleteColumn().run()" title="Удалить колонку">
          <i class="fas fa-minus"></i> col
        </button>
        <button @click="editor.chain().focus().deleteTable().run()" title="Удалить таблицу">
          <i class="fas fa-trash"></i>
        </button>
      </template>
      
      <button @click="showFormulaInput = !showFormulaInput" title="Формула">
        <i class="fas fa-square-root-variable"></i>
      </button>
      <button @click="showLinkInput = !showLinkInput" title="Ссылка">
        <i class="fas fa-link"></i>
      </button>
      <button @click="editor.chain().focus().undo().run()">
        <i class="fas fa-undo"></i>
      </button>
      <button @click="editor.chain().focus().redo().run()">
        <i class="fas fa-redo"></i>
      </button>
    </div>
    
    <div v-if="showTableDialog" class="table-dialog">
      <div class="dialog-header">
        <span>Вставка таблицы</span>
        <button @click="showTableDialog = false"><i class="fas fa-times"></i></button>
      </div>
      <div class="dialog-body">
        <label>Строки:</label>
        <input v-model="tableRows" type="number" min="1" max="20" />
        <label>Колонки:</label>
        <input v-model="tableCols" type="number" min="1" max="20" />
      </div>
      <button class="btn-insert" @click="insertTable">Вставить</button>
    </div>
    
    <div v-if="showFormulaInput" class="formula-input">
      <input v-model="formula" placeholder="Например: \frac{a}{b}" @keydown.enter="insertFormula" />
      <button class="btn-primary" @click="insertFormula">Вставить</button>
      <button class="btn-icon" @click="showFormulaInput = false"><i class="fas fa-times"></i></button>
    </div>

    <div v-if="showLinkInput" class="formula-input">
      <input v-model="linkUrl" placeholder="https://..." @keydown.enter="insertLink" />
      <button class="btn-primary" @click="insertLink">Вставить</button>
      <button class="btn-icon" @click="showLinkInput = false"><i class="fas fa-times"></i></button>
    </div>
    
    <editor-content :editor="editor" class="editor-content" />
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount, watch } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import Typography from '@tiptap/extension-typography'
import Table from '@tiptap/extension-table'
import TableRow from '@tiptap/extension-table-row'
import TableCell from '@tiptap/extension-table-cell'
import TableHeader from '@tiptap/extension-table-header'
import TipTapImage from '@tiptap/extension-image'
import TextAlign from '@tiptap/extension-text-align'
import InlineMath from '../extensions/InlineMath'
import Link from '@tiptap/extension-link'
import { supabase } from '../api/supabase'
import { useAuthStore } from '../stores/auth'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const auth = useAuthStore()

const showTableDialog = ref(false)
const tableRows = ref(3)
const tableCols = ref(3)
const showFormulaInput = ref(false)
const formula = ref('')
const showLinkInput = ref(false)
const linkUrl = ref('')
const isInTable = ref(false)

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit,
    InlineMath,
    Placeholder.configure({ placeholder: 'Пиши конспект...' }),
    Typography,
    Table.configure({ resizable: true }),
    TableRow,
    TableHeader,
    TableCell,
    TipTapImage,
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    Link.configure({ openOnClick: true })
  ],
  editorProps: {
    handlePaste: (view, event) => {
      const items = event.clipboardData?.items
      if (!items) return false
      
      for (const item of items) {
        if (item.type.startsWith('image/')) {
          event.preventDefault()
          uploadImage(item.getAsFile())
          return true
        }
      }
      return false
    }
  },
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
    setTimeout(addCopyButtons, 0)
  },
  onSelectionUpdate: ({ editor }) => {
    isInTable.value = editor.isActive('table')
  }
})

async function compressImage(file, maxWidth = 1200, quality = 0.8) {
  return new Promise((resolve) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const img = new Image()
      img.onload = () => {
        const canvas = document.createElement('canvas')
        let width = img.width
        let height = img.height
        
        if (width > maxWidth) {
          height = (height * maxWidth) / width
          width = maxWidth
        }
        
        canvas.width = width
        canvas.height = height
        const ctx = canvas.getContext('2d')
        ctx.drawImage(img, 0, 0, width, height)
        
        canvas.toBlob((blob) => {
          resolve(new File([blob], file.name, { type: 'image/jpeg' }))
        }, 'image/jpeg', quality)
      }
      img.src = e.target.result
    }
    reader.readAsDataURL(file)
  })
}

async function uploadImage(file) {
  if (!file) return
  
  try {
    const compressedFile = await compressImage(file)
    const userId = auth.user?.id || auth.profile?.id
    const filePath = `${userId}/${Date.now()}.jpg`
    
    const { error: uploadError } = await supabase.storage
      .from('notebook_images')
      .upload(filePath, compressedFile)
    
    if (uploadError) throw uploadError
    
    const { data: { publicUrl } } = supabase.storage
      .from('notebook_images')
      .getPublicUrl(filePath)
    
    editor.value.chain().focus().setImage({ src: publicUrl }).run()
  } catch (e) {
    console.error(e)
  }
}

function addCopyButtons() {
  document.querySelectorAll('.editor-content pre').forEach(pre => {
    if (pre.querySelector('.copy-code-btn')) return
    
    const btn = document.createElement('button')
    btn.className = 'copy-code-btn'
    btn.innerHTML = '<i class="fas fa-copy"></i>'
    btn.onclick = async () => {
      const code = pre.querySelector('code')?.innerText || ''
      await navigator.clipboard.writeText(code)
      btn.innerHTML = '<i class="fas fa-check"></i>'
      setTimeout(() => { btn.innerHTML = '<i class="fas fa-copy"></i>' }, 2000)
    }
    pre.style.position = 'relative'
    pre.appendChild(btn)
  })
}

function insertTable() {
  editor.value.chain().focus().insertTable({ rows: tableRows.value, cols: tableCols.value, withHeaderRow: true }).run()
  showTableDialog.value = false
}

function insertFormula() {
  if (formula.value) {
    const cleanFormula = formula.value.replace(/\\\\/g, '\\')
    editor.value.chain().focus().insertInlineMath(cleanFormula).run()
    formula.value = ''
    showFormulaInput.value = false
  }
}

function insertLink() {
  if (linkUrl.value) {
    const { state } = editor.value
    const { from, to } = state.selection
    
    if (from === to) {
      editor.value.chain().focus().insertContent(`<a href="${linkUrl.value}">${linkUrl.value}</a>`).run()
    } else {
      editor.value.chain().focus().setLink({ href: linkUrl.value }).run()
    }
    
    linkUrl.value = ''
    showLinkInput.value = false
  }
}

watch(() => props.modelValue, (newVal) => {
  if (editor.value && newVal !== editor.value.getHTML()) {
    editor.value.commands.setContent(newVal, false)
  }
})

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<style scoped>
.editor-container {
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 14px;
  overflow: hidden;
  position: relative;
}

.editor-toolbar {
  display: flex;
  gap: 2px;
  padding: 8px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  flex-wrap: wrap;
  background: #0a0a0a;
}

.editor-toolbar button {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: #737373;
  cursor: pointer;
  transition: all 0.15s ease;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.editor-toolbar button:hover {
  background: rgba(255,255,255,0.06);
  color: #e5e5e5;
}

.editor-toolbar button.active {
  background: rgba(255,255,255,0.08);
  color: #ffffff;
}

.table-dialog {
  position: absolute;
  top: 48px;
  left: 8px;
  background: #111111;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 14px;
  padding: 20px;
  z-index: 20;
  min-width: 220px;
  box-shadow: 0 16px 32px rgba(0,0,0,0.4);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  font-weight: 600;
  font-size: 0.9rem;
  color: #e5e5e5;
}

.dialog-header button {
  background: none;
  border: none;
  color: #525252;
  cursor: pointer;
  font-size: 0.85rem;
  padding: 4px;
  transition: color 0.15s ease;
}

.dialog-header button:hover {
  color: #e5e5e5;
}

.dialog-body {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 10px 14px;
  margin-bottom: 16px;
}

.dialog-body label {
  font-size: 0.8rem;
  color: #a3a3a3;
  align-self: center;
  font-weight: 500;
}

.dialog-body input {
  padding: 8px 12px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 8px;
  color: #e5e5e5;
  font-family: inherit;
  outline: none;
  width: 60px;
  font-size: 0.9rem;
  transition: all 0.2s ease;
}

.dialog-body input:focus {
  border-color: rgba(255,255,255,0.15);
  background: rgba(255,255,255,0.04);
}

.btn-insert {
  width: 100%;
  padding: 10px;
  background: #ffffff;
  color: #0a0a0a;
  border: 1px solid #ffffff;
  border-radius: 10px;
  font-family: inherit;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.85rem;
}

.btn-insert:hover {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.formula-input {
  padding: 12px 16px;
  background: #0a0a0a;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  display: flex;
  gap: 8px;
  align-items: center;
}

.formula-input input {
  flex: 1;
  padding: 10px 14px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  color: #e5e5e5;
  font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace;
  outline: none;
  font-size: 0.9rem;
  transition: all 0.2s ease;
}

.formula-input input:focus {
  border-color: rgba(255,255,255,0.15);
  background: rgba(255,255,255,0.04);
}

.formula-input input::placeholder {
  color: #525252;
}

.formula-input .btn-primary {
  padding: 10px 16px;
  background: #ffffff;
  color: #0a0a0a;
  border: 1px solid #ffffff;
  border-radius: 10px;
  cursor: pointer;
  font-family: inherit;
  font-weight: 500;
  font-size: 0.85rem;
  transition: all 0.2s ease;
}

.formula-input .btn-primary:hover {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.formula-input .btn-icon {
  padding: 10px;
  background: none;
  border: none;
  color: #525252;
  cursor: pointer;
  font-size: 0.85rem;
  transition: color 0.15s ease;
}

.formula-input .btn-icon:hover {
  color: #e5e5e5;
}

.editor-content {
  padding: 20px;
  min-height: 400px;
  color: #e5e5e5;
  line-height: 1.7;
  outline: none;
  background: #0a0a0a;
}

.editor-content :deep(.ProseMirror) {
  outline: none;
  min-height: 400px;
}

.editor-content :deep(.ProseMirror p) {
  margin-bottom: 12px;
}

.editor-content :deep(.ProseMirror p.is-editor-empty:first-child::before) {
  color: #525252;
  content: attr(data-placeholder);
  float: left;
  height: 0;
  pointer-events: none;
}

.editor-content :deep(.ProseMirror h2) {
  margin: 24px 0 12px;
  font-size: 1.3rem;
  font-weight: 600;
  color: #ffffff;
  letter-spacing: -0.01em;
}

.editor-content :deep(.ProseMirror h3) {
  margin: 20px 0 10px;
  font-size: 1.1rem;
  font-weight: 600;
  color: #e5e5e5;
  letter-spacing: -0.01em;
}

.editor-content :deep(.ProseMirror ul),
.editor-content :deep(.ProseMirror ol) {
  margin-left: 20px;
  margin-bottom: 12px;
}

.editor-content :deep(.ProseMirror li) {
  margin-bottom: 4px;
}

.editor-content :deep(.ProseMirror table) {
  border-collapse: collapse;
  width: 100%;
  margin: 20px 0;
  font-size: 0.9rem;
}

.editor-content :deep(.ProseMirror th),
.editor-content :deep(.ProseMirror td) {
  border: 1px solid rgba(255,255,255,0.08);
  padding: 8px 12px;
  min-width: 50px;
}

.editor-content :deep(.ProseMirror th) {
  background: rgba(255,255,255,0.04);
  font-weight: 600;
  color: #e5e5e5;
}

.editor-content :deep(.ProseMirror img) {
  max-width: 100%;
  border-radius: 10px;
  margin: 12px 0;
}

.editor-content :deep(.ProseMirror .katex) {
  font-size: 1.1em;
}

.editor-content :deep(.katex-mathml) {
  display: none;
}

.editor-content :deep(.ProseMirror blockquote) {
  border-left: 2px solid rgba(255,255,255,0.15);
  padding-left: 16px;
  margin: 16px 0;
  color: #737373;
  font-style: italic;
}

.editor-content :deep(.ProseMirror code) {
  background: rgba(255,255,255,0.06);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace;
  font-size: 0.85em;
  color: #e5e5e5;
}

.editor-content :deep(.ProseMirror pre) {
  background: rgba(255,255,255,0.03);
  padding: 16px;
  border-radius: 10px;
  margin: 16px 0;
  overflow-x: auto;
  border: 1px solid rgba(255,255,255,0.06);
  position: relative;
}

.editor-content :deep(.ProseMirror pre code) {
  background: none;
  padding: 0;
  color: #a3a3a3;
}

.editor-content :deep(.copy-code-btn) {
  position: absolute;
  top: 8px;
  right: 8px;
  background: rgba(255,255,255,0.06);
  border: none;
  color: #737373;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
  font-size: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.editor-content :deep(.copy-code-btn:hover) {
  background: rgba(255,255,255,0.12);
  color: #e5e5e5;
}

.editor-content :deep(.ProseMirror strong) {
  color: #ffffff;
  font-weight: 600;
}

.editor-content :deep(.ProseMirror a) {
  color: #a3a3a3;
  text-decoration: underline;
  text-underline-offset: 3px;
}

@media (max-width: 768px) {
  .editor-toolbar {
    padding: 6px;
    gap: 1px;
  }
  
  .editor-toolbar button {
    width: 28px;
    height: 28px;
    font-size: 0.75rem;
  }
  
  .editor-content {
    padding: 16px;
    min-height: 300px;
  }
  
  .table-dialog {
    left: 4px;
    right: 4px;
    min-width: auto;
  }
}
</style>