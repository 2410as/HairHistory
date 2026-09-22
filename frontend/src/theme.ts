import { defineConfig, defineTokens } from '@chakra-ui/react'

const SANS =
  "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif"

const colors = defineTokens.colors({
  white: { value: '#ffffff' },
  black: { value: '#000000' },
  gray: {
    0: { value: '#ffffff' },
    50: { value: '#fafafa' },
    100: { value: '#f5f5f5' },
    200: { value: '#ececec' },
    300: { value: '#e5e5e5' },
    400: { value: '#dcdcdc' },
    500: { value: '#999999' },
    600: { value: '#666666' },
    700: { value: '#333333' },
    800: { value: '#1a1a1a' },
    900: { value: '#000000' },
  },
})

export const theme = defineConfig({
  theme: {
    tokens: {
      colors,
      fonts: {
        heading: { value: SANS },
        body: { value: SANS },
        mono: { value: "ui-monospace, SFMono-Regular, Menlo, Consolas, monospace" },
      },
      lineHeights: {
        tight: { value: '1.3' },
        normal: { value: '1.7' },
      },
      shadows: {
        'card-hover': { value: '0 8px 24px rgba(0, 0, 0, 0.08)' },
        'card-rest': { value: 'none' },
      },
    },
    semanticTokens: {
      colors: {
        'bg.base': { value: '{colors.white}' },
        'bg.surface': { value: '{colors.white}' },
        'bg.elevated': { value: '{colors.gray.100}' },
        'text.primary': { value: '{colors.black}' },
        'text.secondary': { value: '{colors.gray.600}' },
        'text.tertiary': { value: '{colors.gray.600}' },
        'text.placeholder': { value: '{colors.gray.500}' },
        'border.subtle': { value: '{colors.gray.300}' },
        'border.strong': { value: '{colors.gray.400}' },
        'button.bg': { value: '{colors.black}' },
        'button.text': { value: '{colors.white}' },
        'button.hover': { value: '{colors.gray.700}' },
        'button.active': { value: '{colors.gray.800}' },
        'focus.outline': { value: '{colors.black}' },
      },
      shadows: {
        'card-hover': { value: '{shadows.card-hover}' },
        'card-rest': { value: '{shadows.card-rest}' },
      },
    },
  },
})

export default theme
