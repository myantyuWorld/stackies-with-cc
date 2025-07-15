import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import LoginScreen from '../LoginScreen.vue'

describe('LoginScreen', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders correctly', () => {
    const wrapper = mount(LoginScreen)
    expect(wrapper.exists()).toBe(true)
  })

  it('contains login form', () => {
    const wrapper = mount(LoginScreen)
    expect(wrapper.find('[data-testid="login-form"]').exists()).toBe(true)
  })
})