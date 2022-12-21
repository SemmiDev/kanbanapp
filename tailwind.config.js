/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ['./views/**/*.{html,js}'],
    theme: {
        theme: {
            container: {
                center: true,
                padding: '2rem',
                width: '100%',
            },
            table: {},
        },
        variants: {
            extend: {
                backgroundColor: ['even'],
            },
        },
        extend: {},
    },
    plugins: [require('daisyui')],
};
