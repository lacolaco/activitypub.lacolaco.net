/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{html,ts,scss}'],
  theme: {
    extend: {
      container: {
        center: true,
        padding: {
          DEFAULT: '1rem',
        },
      },
      fontFamily: {
        sans: ['"M PLUS 1p"', 'sans-serif'],
      },
      textColor: {
        black: 'rgb(25, 25, 25)',
      },
      backgroundColor: {
        white: 'rgba(255, 255, 255, 0.9)',
        panel: 'rgb(237, 237, 237)',
      },
    },
  },
  plugins: [require('@tailwindcss/typography'), require('@tailwindcss/forms')],
};
