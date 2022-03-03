import { createTheme, style } from '@vanilla-extract/css';

export const [baseTheme, vars] = createTheme({
});

export const galleryWrapper = style({
    display: 'flex',
    maxWidth: 940,
    padding: 30,
    flexFlow: 'row wrap',
    justifyContent: 'center',
});

export const grid = style({
    width: 200,
    height: 200,
    margin: '10px 1em',
    paddingLeft: 10,
    paddingRight: 10,
    flex: '0 1 20%'
});

export const gridImgThumbnail = style({
    maxWidth: '100%',
    width: '100%',
    margin: 0,
    objectFit: 'cover',
    alignSelf: 'center',
    transition: 'opacity .25s linear'
});