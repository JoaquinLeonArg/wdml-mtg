const { nextui } = require("@nextui-org/react");

module.exports = {
  darkMode: 'class',
  content: [
    './src/**/*.{html,js,ts,tsx}',
    "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  plugins: [nextui()],
  theme: {
    extend: {
      colors: {
        primary: {
          "50": "#eff6ff",
          "100": "#dbeafe",
          "200": "#bfdbfe",
          "300": "#93c5fd",
          "400": "#60a5fa",
          "500": "#3b82f6",
          "600": "#2563eb",
          "700": "#1d4ed8",
          "800": "#1e40af",
          "900": "#1e3a8a",
          "950": "#172554"
        },
        secondary: {
          50: '#fdf8e8',
          100: '#fbf1d0',
          200: '#f7e3a1',
          300: '#f3d572',
          400: '#efc743',
          500: '#ebb914',
          600: '#bc9410',
          700: '#8d6f0c',
          800: '#5e4a08',
          900: '#2f2504',
          950: '#171202',
        },
        rarity: {
          "common": "#666666",
          "uncommon": "#9BCDFF",
          "rare": "#C6AC6E",
          "mythic": "#D66525"
        },
        mana: {
          "white": "#F8E7B9",
          "blue": "#B3CEEA",
          "black": "#A69F9D",
          "red": "#EBA082",
          "green": "#C4D3CA",
        }
      }
    },
    fontFamily: {
      'body': [
        'Inter',
        'ui-sans-serif',
        'system-ui',
        '-apple-system',
        'system-ui',
        'Segoe UI',
        'Roboto',
        'Helvetica Neue',
        'Arial',
        'Noto Sans',
        'sans-serif',
        'Apple Color Emoji',
        'Segoe UI Emoji',
        'Segoe UI Symbol',
        'Noto Color Emoji'
      ],
      'sans': [
        'Inter',
        'ui-sans-serif',
        'system-ui',
        '-apple-system',
        'system-ui',
        'Segoe UI',
        'Roboto',
        'Helvetica Neue',
        'Arial',
        'Noto Sans',
        'sans-serif',
        'Apple Color Emoji',
        'Segoe UI Emoji',
        'Segoe UI Symbol',
        'Noto Color Emoji'
      ]
    }
  }
}