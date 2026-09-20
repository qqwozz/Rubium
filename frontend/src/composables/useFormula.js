const API_URL = import.meta.env.VITE_FORMULA_API || 'http://localhost:5090'

export function useFormula() {
  async function convert(text) {
    const res = await fetch(`${API_URL}/formula/convert`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ text })
    })
    if (!res.ok) throw new Error('convert failed')
    return await res.json()
  }
  return { convert }
}