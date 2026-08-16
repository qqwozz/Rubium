import { defineStore } from 'pinia'
import { createClient } from '@supabase/supabase-js'
import { supabase } from '../api/supabase'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    profile: null,
    loading: true
  }),
  
  getters: {
    isAuthenticated: (state) => !!state.user,
    isAdmin: (state) => state.profile?.status === 'admin',
    userName: (state) => {
      if (!state.profile) return ''
      return state.profile.first_name || state.user?.email?.split('@')[0] || 'ученик'
    }
  },
  
  actions: {
    async init() {
      try {
        const { data: { user } } = await supabase.auth.getUser()
        
        if (user) {
          this.user = user
          
          const { data: profile } = await supabase
            .from('rubium_users')
            .select('*')
            .eq('auth_id', user.id)
            .maybeSingle()
          
          if (profile) {
            this.profile = profile
          }
        }
      } catch (e) {
        console.error('Auth init error:', e)
      } finally {
        this.loading = false
      }
    },
    
    async login(email, password) {
      const { data, error } = await supabase.auth.signInWithPassword({ email, password })
      if (error) throw error
      
      this.user = data.user
      await this.init()
    },
    
    async register(email, password, firstName, lastName) {
      const { data, error } = await supabase.auth.signUp({
        email,
        password,
        options: {
          data: { first_name: firstName, last_name: lastName }
        }
      })
      if (error) throw error
      
      this.user = data.user
    },
    
    async logout() {
      await supabase.auth.signOut()
      this.user = null
      this.profile = null
      window.location.href = '/login'
    }
  }
})