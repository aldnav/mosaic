import * as React from "react";
import { solsHeader, solsTheme } from "../styles/index.css";
import { useState } from "react";
import { useForm, SubmitHandler } from "react-hook-form";
import Gallery, { ImagesProps } from "../components/gallery";

const MAX_NUMBER_OF_PHOTOS = 9;
const MAX_FILE_SIZE = 2 * 1000000;

const checkValidFileSize = (
  fileList: FileList,
  maxFileSize = MAX_FILE_SIZE
) => {
  let invalidFiles: number = 0;
  for (let aFile of fileList) {
    if (aFile?.size > maxFileSize) {
      ++invalidFiles;
    }
  }
  return invalidFiles === 0 || "Maximum of 2MB only.";
};

const checkValidFileType = (
  fileList: FileList,
  validFileTypes: string[] = ["image/jpeg", "image/png"]
) => {
  let invalidFiles: number = 0;
  for (let aFile of fileList) {
    if (!validFileTypes.includes(aFile?.type)) {
      ++invalidFiles;
    }
  }
  return invalidFiles === 0 || "Only PNG and JPEG";
};

// markup
const IndexPage = () => {
  // form defaults
  const [formData, setFormData] = useState({
    mainPhotoData: "",
  });
  const [tilePhotoData, setTilePhotoData] = useState("");
  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm();
  const [imageForPreview, setImageForPreview] = useState<ImagesProps[]>([]);
  const onSubmit: SubmitHandler<any> = (data) => console.log(data);

  React.useEffect(() => {
    const subscription = watch((data) => {
      renderGallery(data?.tilePhotoData ?? []);
    });
    return () => {
      subscription.unsubscribe();
    };
  }, [watch]);

  const renderGallery = (files = []) => {
    let images = [];
    if (files.length === 0) {
      setImageForPreview([]);
      return;
    }

    function readFile(file) {
      return new Promise(function (resolve, reject) {
        let fr = new FileReader();
        fr.onload = function () {
          resolve(fr.result);
        };
        fr.onerror = function () {
          reject(fr);
        };
        fr.readAsDataURL(file);
      });
    }

    let readers = [];
    for (let aFile of files) {
      readers.push(readFile(aFile));
    }

    Promise.all(readers).then((values) => {
      setImageForPreview(
        values.map((x) => ({
          full: {
            src: x,
          },
        }))
      );
    });
  };

  return (
    <main className={solsTheme}>
      <title>Mosaic</title>
      <h1 className={solsHeader}>Mosaic</h1>
      <form onSubmit={handleSubmit(onSubmit)}>
        {/* <input
          type="file"
          required
          {...register("mainPhotoData", {
            validate: {
              acceptedFormats: files => ['image/jpeg', 'image/png'].includes(files[0]?.type) || 'Only PNG and JPEG',
              lessThan2MB: files => files[0]?.size <= MAX_FILE_SIZE || 'Maximum of 2MB only',
            }
          })}
        />
        {errors?.mainPhotoData && <p>{errors.mainPhotoData.message}</p>} */}
        <input
          type="file"
          multiple
          required
          {...register("tilePhotoData", {
            validate: {
              allowedNumberOfPhotos: (files) =>
                files.length <= MAX_NUMBER_OF_PHOTOS ||
                `Only ${MAX_NUMBER_OF_PHOTOS} photos are allowed`,
              acceptedFormats: (files) => checkValidFileType(files),
              lessThan2MB: (files) => checkValidFileSize(files),
            },
          })}
        />
        {errors?.tilePhotoData && <p>{errors.tilePhotoData.message}</p>}
        <input type="submit" />
      </form>
      <Gallery images={imageForPreview} />
    </main>
  );
};

export default IndexPage;
