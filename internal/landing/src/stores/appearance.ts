import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type Theme = 'light' | 'dark' | 'system'
export type AccentColor = 'blue' | 'green' | 'purple' | 'red' | 'orange'

export const useAppearanceStore = defineStore(
  'appearance',
  () => {
    // State
    const theme = ref<Theme>('system')
    const accentColor = ref<AccentColor>('blue')

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

    // Generate complete color palette based on accent color
    const colorPalette = computed(() => {
      const baseHue = {
        blue: 220,
        green: 142,
        purple: 262,
        red: 0,
        orange: 25,
      }[accentColor.value]

      const isDark = themeClass.value === 'dark'

      // Special handling for green theme to improve contrast
      const getPrimaryForeground = () => {
        if (accentColor.value === 'green') {
          // Use a very dark green for better contrast on the green primary
          return isDark ? '142 15% 95%' : '142 20% 8%'
        }
        // Default white for other colors
        return isDark ? '0 0% 100%' : '0 0% 100%'
      }

      return {
        // Primary colors
        primary: `${baseHue} 91% ${isDark ? '65%' : '56%'}`,
        // Use custom foreground for green, white for others
        primaryForeground: getPrimaryForeground(),

        // Background colors - tinted with accent
        background: isDark
          ? `${baseHue} 15% 4%` // Very dark with subtle accent tint
          : `${baseHue} 25% 98%`, // Very light with subtle accent tint
        foreground: isDark
          ? `${baseHue} 5% 95%` // Light text with subtle accent
          : `${baseHue} 15% 8%`, // Dark text with subtle accent

        // Card colors
        card: isDark ? `${baseHue} 20% 6%` : `${baseHue} 30% 97%`,
        cardForeground: isDark ? `${baseHue} 5% 90%` : `${baseHue} 15% 12%`,

        // Secondary colors
        secondary: isDark ? `${baseHue} 25% 12%` : `${baseHue} 35% 92%`,
        secondaryForeground: isDark ? `${baseHue} 8% 85%` : `${baseHue} 20% 18%`,

        // Muted colors
        muted: isDark ? `${baseHue} 20% 10%` : `${baseHue} 30% 94%`,
        mutedForeground: isDark ? `${baseHue} 8% 60%` : `${baseHue} 15% 45%`,

        // Accent colors
        accent: isDark ? `${baseHue} 30% 15%` : `${baseHue} 40% 90%`,
        accentForeground: isDark ? `${baseHue} 10% 80%` : `${baseHue} 25% 25%`,

        // Border and input
        border: isDark ? `${baseHue} 15% 18%` : `${baseHue} 25% 85%`,
        input: isDark ? `${baseHue} 15% 18%` : `${baseHue} 25% 85%`,
        ring: `${baseHue} 91% ${isDark ? '65%' : '56%'}`,

        // Sidebar colors
        sidebar: isDark ? `${baseHue} 25% 8%` : `${baseHue} 35% 95%`,
        sidebarForeground: isDark ? `${baseHue} 8% 88%` : `${baseHue} 18% 15%`,
        sidebarPrimary: `${baseHue} 91% ${isDark ? '65%' : '56%'}`,
        // Use the same custom foreground logic for sidebar
        sidebarPrimaryForeground: getPrimaryForeground(),
        sidebarAccent: isDark ? `${baseHue} 30% 20%` : `${baseHue} 40% 88%`,
        sidebarAccentForeground: isDark ? `${baseHue} 10% 85%` : `${baseHue} 25% 20%`,
        sidebarBorder: isDark ? `${baseHue} 20% 15%` : `${baseHue} 30% 82%`,
        sidebarRing: `${baseHue} 91% ${isDark ? '65%' : '56%'}`,
      }
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

        // Reapply colors when theme changes since palette depends on theme
        applyAccentColor()
      }
    }

    // Apply accent color to CSS variables
    const applyAccentColor = () => {
      if (typeof document !== 'undefined') {
        const root = document.documentElement
        const palette = colorPalette.value

        // Apply all color variables
        root.style.setProperty('--primary', `hsl(${palette.primary})`)
        root.style.setProperty('--primary-foreground', `hsl(${palette.primaryForeground})`)
        root.style.setProperty('--background', `hsl(${palette.background})`)
        root.style.setProperty('--foreground', `hsl(${palette.foreground})`)
        root.style.setProperty('--card', `hsl(${palette.card})`)
        root.style.setProperty('--card-foreground', `hsl(${palette.cardForeground})`)
        root.style.setProperty('--secondary', `hsl(${palette.secondary})`)
        root.style.setProperty('--secondary-foreground', `hsl(${palette.secondaryForeground})`)
        root.style.setProperty('--muted', `hsl(${palette.muted})`)
        root.style.setProperty('--muted-foreground', `hsl(${palette.mutedForeground})`)
        root.style.setProperty('--accent', `hsl(${palette.accent})`)
        root.style.setProperty('--accent-foreground', `hsl(${palette.accentForeground})`)
        root.style.setProperty('--border', `hsl(${palette.border})`)
        root.style.setProperty('--input', `hsl(${palette.input})`)
        root.style.setProperty('--ring', `hsl(${palette.ring})`)

        // Apply sidebar colors
        root.style.setProperty('--sidebar', `hsl(${palette.sidebar})`)
        root.style.setProperty('--sidebar-foreground', `hsl(${palette.sidebarForeground})`)
        root.style.setProperty('--sidebar-primary', `hsl(${palette.sidebarPrimary})`)
        root.style.setProperty(
          '--sidebar-primary-foreground',
          `hsl(${palette.sidebarPrimaryForeground})`,
        )
        root.style.setProperty('--sidebar-accent', `hsl(${palette.sidebarAccent})`)
        root.style.setProperty(
          '--sidebar-accent-foreground',
          `hsl(${palette.sidebarAccentForeground})`,
        )
        root.style.setProperty('--sidebar-border', `hsl(${palette.sidebarBorder})`)
        root.style.setProperty('--sidebar-ring', `hsl(${palette.sidebarRing})`)
      }
    }

    // Initialize appearance settings
    const initializeAppearance = () => {
      applyTheme()
      applyAccentColor()

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
    }

    return {
      // State
      theme,
      accentColor,

      // Computed
      themeClass,
      colorPalette,

      // Actions
      setTheme,
      setAccentColor,
      initializeAppearance,
      resetToDefaults,
    }
  },
  {
    persist: {
      key: 'meridian-landing-appearance',
      storage: typeof window !== 'undefined' ? localStorage : undefined,
    },
  },
)
