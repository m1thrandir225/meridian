import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type Theme = 'light' | 'dark' | 'system'
export type AccentColor = 'blue' | 'green' | 'purple' | 'red' | 'orange'
export type FontSize = 'small' | 'medium' | 'large'
export type MessageDisplayMode = 'cozy' | 'compact'

export const useAppearanceStore = defineStore(
  'appearance',
  () => {
    // State
    const theme = ref<Theme>('system')
    const accentColor = ref<AccentColor>('blue')
    const fontSize = ref<FontSize>('medium')
    const messageDisplayMode = ref<MessageDisplayMode>('cozy')

    // Computed values for CSS classes
    const themeClass = computed(() => {
      if (theme.value === 'system') {
        // Use system preference
        if (typeof window !== 'undefined') {
          return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
        }
        return 'light'
      }
      return theme.value
    })

    const accentColorClass = computed(() => {
      const colorMap = {
        blue: '220 91% 56%', // HSL values for shadcn compatibility
        green: '142 71% 45%',
        purple: '262 83% 58%',
        red: '0 84% 60%',
        orange: '25 95% 53%',
      }
      return colorMap[accentColor.value]
    })

    const fontSizeClass = computed(() => {
      const sizeMap = {
        small: '14px',
        medium: '16px',
        large: '18px',
      }
      return sizeMap[fontSize.value]
    })

    const messageDisplayClass = computed(() => {
      return messageDisplayMode.value === 'compact' ? 'compact' : 'cozy'
    })

    // Actions
    const setTheme = (newTheme: Theme) => {
      theme.value = newTheme
      applyTheme()
    }

    const setAccentColor = (newColor: AccentColor) => {
      accentColor.value = newColor
      applyAccentColor()
    }

    const setFontSize = (newSize: FontSize) => {
      fontSize.value = newSize
      applyFontSize()
    }

    const setMessageDisplayMode = (newMode: MessageDisplayMode) => {
      messageDisplayMode.value = newMode
    }

    // Apply theme to document
    const applyTheme = () => {
      if (typeof document !== 'undefined') {
        const root = document.documentElement

        // Remove existing theme classes
        root.classList.remove('light', 'dark')

        // Add new theme class
        root.classList.add(themeClass.value)

        // Update color-scheme for better browser integration
        root.style.colorScheme = themeClass.value
      }
    }

    // Apply accent color to CSS variables
    const applyAccentColor = () => {
      if (typeof document !== 'undefined') {
        const root = document.documentElement
        root.style.setProperty('--primary', `hsl(${accentColorClass.value})`)
      }
    }

    // Apply font size to CSS variables
    const applyFontSize = () => {
      if (typeof document !== 'undefined') {
        const root = document.documentElement
        root.style.setProperty('--font-size-base', fontSizeClass.value)
      }
    }

    // Initialize appearance settings
    const initializeAppearance = () => {
      applyTheme()
      applyAccentColor()
      applyFontSize()

      // Listen for system theme changes
      if (typeof window !== 'undefined' && theme.value === 'system') {
        const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
        mediaQuery.addEventListener('change', applyTheme)
      }
    }

    // Reset to defaults
    const resetToDefaults = () => {
      setTheme('system')
      setAccentColor('blue')
      setFontSize('medium')
      setMessageDisplayMode('cozy')
    }

    return {
      // State
      theme,
      accentColor,
      fontSize,
      messageDisplayMode,

      // Computed
      themeClass,
      accentColorClass,
      fontSizeClass,
      messageDisplayClass,

      // Actions
      setTheme,
      setAccentColor,
      setFontSize,
      setMessageDisplayMode,
      initializeAppearance,
      resetToDefaults,
    }
  },
  {
    persist: {
      key: 'meridian-appearance',
      storage: typeof window !== 'undefined' ? localStorage : undefined,
    },
  },
)
