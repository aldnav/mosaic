import { createTheme, style } from '@vanilla-extract/css';

export const [solsTheme, vars] = createTheme({
    color: {
        brand: '#232129'
    },
    font: {
        body: '-apple-system, Roboto, sans-serif, serif'
    }
});

export const solsHeader = style({
    color: vars.color.brand,
    fontFamily: vars.font.body,
    padding: 20
});