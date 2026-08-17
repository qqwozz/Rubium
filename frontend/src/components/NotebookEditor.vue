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
          <i class="fas fa-minus"></i>
        </button>
        <button @click="editor.chain().focus().deleteColumn().run()" title="Удалить колонку">
          <i class="fas fa-minus"></i>
        </button>
        <button @click="editor.chain().focus().deleteTable().run()" title="Удалить таблицу">
          <i class="fas fa-trash"></i>
        </button>
      </template>
      
      <button @click="showFormulaInput = !showFormulaInput" title="Формула">
        <i class="fas fa-square-root-variable"></i>
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
        <button @click="showTableDialog = false">✕</button>
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
      <button @click="insertFormula">Вставить</button>
      <button @click="showFormulaInput = false">✕</button>
    </div>
    
    <editor-content :editor="editor" class="editor-content" @paste="handlePaste" />
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
import Image from '@tiptap/extension-image'
import katex from 'katex'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const showTableDialog = ref(false)
const tableRows = ref(3)
const tableCols = ref(3)
const showFormulaInput = ref(false)
const formula = ref('')
const isInTable = ref(false)

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit,
    Placeholder.configure({ placeholder: 'Пиши конспект...' }),
    Typography,
    Table.configure({ resizable: true }),
    TableRow,
    TableHeader,
    TableCell,
    Image
  ],
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  },
  onSelectionUpdate: ({ editor }) => {
    isInTable.value = editor.isActive('table')
  }
})

function insertTable() {
  editor.value.chain().focus().insertTable({ rows: tableRows.value, cols: tableCols.value, withHeaderRow: true }).run()
  showTableDialog.value = false
}

function insertFormula() {
  if (formula.value) {
    try {
      const cleanFormula = formula.value.replace(/\\\\/g, '\\')
      const html = katex.renderToString(cleanFormula, { displayMode: false, throwOnError: false })
      editor.value.chain().focus().insertContent(html).run()
      formula.value = ''
      showFormulaInput.value = false
    } catch (e) {
      console.error(e)
    }
  }
}

function handlePaste(event) {
  const items = event.clipboardData?.items
  if (items) {
    for (const item of items) {
      if (item.type.startsWith('image/')) {
        const file = item.getAsFile()
        const reader = new FileReader()
        reader.onload = (e) => {
          editor.value.chain().focus().setImage({ src: e.target.result }).run()
        }
        reader.readAsDataURL(file)
        event.preventDefault()
        return
      }
    }
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
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 16px;
  overflow: hidden;
  position: relative;
}

.editor-toolbar {
  display: flex;
  gap: 4px;
  padding: 8px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  flex-wrap: wrap;
}

.editor-toolbar button {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: #94A3B8;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.85rem;
}

.editor-toolbar button:hover {
  background: rgba(255,255,255,0.06);
  color: #F1F5F9;
}

.editor-toolbar button.active {
  background: rgba(167,139,250,0.15);
  color: #A78BFA;
}

.table-dialog {
  position: absolute;
  top: 40px;
  left: 8px;
  background: #1a1a2e;
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 12px;
  padding: 16px;
  z-index: 20;
  min-width: 200px;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-weight: 600;
}

.dialog-header button {
  background: none;
  border: none;
  color: #94A3B8;
  cursor: pointer;
}

.dialog-body {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 8px 12px;
  margin-bottom: 12px;
}

.dialog-body label {
  font-size: 0.8rem;
  color: #94A3B8;
  align-self: center;
}

.dialog-body input {
  padding: 6px 10px;
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 6px;
  color: #F1F5F9;
  font-family: inherit;
  outline: none;
  width: 60px;
}

.btn-insert {
  width: 100%;
  padding: 8px;
  background: #A78BFA;
  color: #0F0F1A;
  border: none;
  border-radius: 8px;
  font-family: inherit;
  font-weight: 600;
  cursor: pointer;
}

.formula-input {
  padding: 12px 16px;
  background: rgba(255,255,255,0.02);
  border-bottom: 1px solid rgba(255,255,255,0.06);
  display: flex;
  gap: 8px;
}

.formula-input input {
  flex: 1;
  padding: 8px 12px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 8px;
  color: #F1F5F9;
  font-family: 'JetBrains Mono', monospace;
  outline: none;
}

.formula-input button {
  padding: 8px 16px;
  background: #A78BFA;
  color: #0F0F1A;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-family: inherit;
  font-weight: 600;
}

.editor-content {
  padding: 16px;
  min-height: 400px;
  color: #F1F5F9;
  line-height: 1.7;
  outline: none;
}

.editor-content :deep(.ProseMirror) {
  outline: none;
  min-height: 400px;
}

.editor-content :deep(.ProseMirror p) {
  margin-bottom: 8px;
}

.editor-content :deep(.ProseMirror h2) {
  margin: 16px 0 8px;
  font-size: 1.3rem;
}

.editor-content :deep(.ProseMirror ul),
.editor-content :deep(.ProseMirror ol) {
  margin-left: 20px;
  margin-bottom: 8px;
}

.editor-content :deep(.ProseMirror table) {
  border-collapse: collapse;
  width: 100%;
  margin: 16px 0;
}

.editor-content :deep(.ProseMirror th),
.editor-content :deep(.ProseMirror td) {
  border: 1px solid rgba(255,255,255,0.15);
  padding: 8px 12px;
  min-width: 50px;
}

.editor-content :deep(.ProseMirror img) {
  max-width: 100%;
  border-radius: 8px;
  margin: 8px 0;
}

.editor-content :deep(.ProseMirror .katex) {
  font-size: 1.1em;
}

.editor-content :deep(.katex-mathml) {
  display: none;
}
</style>