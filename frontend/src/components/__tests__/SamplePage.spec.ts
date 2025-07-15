import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SamplePage from '../SamplePage.vue'

describe('SamplePage', () => {
  it('renders proper title', () => {
    const wrapper = mount(SamplePage)
    expect(wrapper.text()).toContain('サンプル画面')
  })

  it('renders proper message', () => {
    const wrapper = mount(SamplePage)
    expect(wrapper.text()).toContain('こんにちは、Vue3 + TypeScriptです！')
  })

  it('initializes counter to 0', () => {
    const wrapper = mount(SamplePage)
    expect(wrapper.find('.counter-value').text()).toBe('0')
  })

  it('increments counter when +1 button is clicked', async () => {
    const wrapper = mount(SamplePage)
    const incrementButton = wrapper.find('.btn-primary')
    
    await incrementButton.trigger('click')
    expect(wrapper.find('.counter-value').text()).toBe('1')
    
    await incrementButton.trigger('click')
    expect(wrapper.find('.counter-value').text()).toBe('2')
  })

  it('resets counter when reset button is clicked', async () => {
    const wrapper = mount(SamplePage)
    const incrementButton = wrapper.find('.btn-primary')
    const resetButton = wrapper.find('.btn-secondary')
    
    // Increment counter first
    await incrementButton.trigger('click')
    await incrementButton.trigger('click')
    expect(wrapper.find('.counter-value').text()).toBe('2')
    
    // Reset counter
    await resetButton.trigger('click')
    expect(wrapper.find('.counter-value').text()).toBe('0')
  })

  it('renders counter section with proper heading', () => {
    const wrapper = mount(SamplePage)
    expect(wrapper.find('h2').text()).toBe('カウンター')
  })

  it('has proper button labels', () => {
    const wrapper = mount(SamplePage)
    const buttons = wrapper.findAll('.btn')
    
    expect(buttons[0].text()).toBe('+1')
    expect(buttons[1].text()).toBe('リセット')
  })
})