import * as React from "react";
import { galleryWrapper, grid, gridImgThumbnail } from "../styles/gallery.css";

export interface ImagesProps {
  images: [
    {
      full: {
        src: string;
      };
    }
  ];
}

const Gallery = ({ images }: ImagesProps) => {
  const imageElements = images.map((image, id) => {
    return (
      <div key={id} className={grid}>
        <img className={gridImgThumbnail} src={image.full.src} />
      </div>
    );
  });
  return <div className={galleryWrapper}>{imageElements}</div>;
};

export default Gallery;
