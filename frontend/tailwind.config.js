/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        accent: {
          DEFAULT: 'var(--color-accent)',
          hover: 'var(--color-accent-hover)',
          light: 'var(--color-accent-light)',
        },
        success: { DEFAULT: 'var(--color-success)', light: 'var(--color-success-light)' },
        warning: { DEFAULT: 'var(--color-warning)', light: 'var(--color-warning-light)' },
        danger: { DEFAULT: 'var(--color-danger)', light: 'var(--color-danger-light)' },
      },
      fontFamily: {
        sans: 'var(--font-sans)',
        mono: 'var(--font-mono)',
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
