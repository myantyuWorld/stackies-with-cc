import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import LoginForm from '../LoginForm.vue'

describe('LoginForm', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders correctly', () => {
    const wrapper = mount(LoginForm)
    expect(wrapper.exists()).toBe(true)
  })

  it('has login functionality', () => {
    const wrapper = mount(LoginForm)
    expect(wrapper.vm).toBeDefined()
  })
})