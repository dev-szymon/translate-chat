/** @type {import('tailwindcss').Config} */
export default {
    content: ["./src/**/*.{tsx,ts}"],
    theme: {
        extend: {
            colors: {
                theme: {
                    primary: {
                        DEFAULT: "var(--color-primary-base)",
                        dark: "var(--color-primary-dark)",
                        light: "var(--color-primary-light)"
                    },
                    secondary: {
                        DEFAULT: "var(--color-secondary-base)"
                    }
                }
            },
            textColor: {
                theme: {
                    base: "var(--color-text-base)",
                    inverted: "var(--color-text-inverted)"
                }
            }
        }
    },
    plugins: []
};
