/** @type {import('tailwindcss').Config} */
export default {
    content: ["./src/**/*.{tsx,ts}"],
    theme: {
        extend: {
            colors: {
                theme: {
                    inverted: {
                        DEFAULT: "var(--color-inverted)"
                    },
                    base: {
                        DEFAULT: "var(--color-base)",
                        50: "var(--color-base-50)",
                        200: "var(--color-base-200)"
                    },
                    primary: {
                        DEFAULT: "var(--color-primary)",
                        dark: "var(--color-primary-dark)",
                        light: "var(--color-primary-light)"
                    },
                    disabled: {
                        DEFAULT: "var(--color-disabled)"
                    }
                }
            }
        }
    },
    plugins: []
};
