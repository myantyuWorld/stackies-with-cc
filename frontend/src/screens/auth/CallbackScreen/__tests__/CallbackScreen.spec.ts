import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import CallbackScreen from '../CallbackScreen.vue'

// Mock vue-router
vi.mock('vue-router', () => ({
  useRouter: vi.fn(() => ({
    push: vi.fn()
  })),
  useRoute: vi.fn(() => ({
    query: {}
  }))
}))

describe('CallbackScreen', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders correctly', () => {
    const wrapper = mount(CallbackScreen)
    expect(wrapper.exists()).toBe(true)
  })

  it('shows processing state', () => {
    const wrapper = mount(CallbackScreen)
    expect(wrapper.find('[data-testid="processing-message"]').exists()).toBe(true)
  })
})