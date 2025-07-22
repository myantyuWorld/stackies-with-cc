import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import UserProfile from '../UserProfile.vue'

describe('UserProfile', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders correctly', () => {
    const wrapper = mount(UserProfile)
    expect(wrapper.exists()).toBe(true)
  })

  it('has user profile functionality', () => {
    const wrapper = mount(UserProfile)
    expect(wrapper.vm).toBeDefined()
  })
})