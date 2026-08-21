<template>
  <router-view />
</template>

<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from './stores/auth'

const handleEmailConfirmation = async () => {
  const params = new URLSearchParams(window.location.search)
  const tokenHash = params.get('token_hash')
  const type = params.get('type')

  if (tokenHash && type === 'signup') {
    const { error } = await supabase.auth.verifyOtp({
      token_hash: tokenHash,
      type: 'signup'
    })
    if (!error) {
      window.location.href = '/' // Redirect to Home
    }
  }
}

const auth = useAuthStore()

onMounted(() => {
  auth.init()
})
</script>