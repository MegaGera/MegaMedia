'use client';

import { useState, useEffect } from 'react';
import { Card, CardBody, CardFooter, Image, Button,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  useDisclosure,
  Input
} from "@nextui-org/react";

import clsx from 'clsx';

import { deleteMegageraImage, fetchMegageraLogos, uploadMegageraImage, updateMegageraImageName } from '@/app/lib/data';
import { Image as ModalImage } from '@/app/modals/image';
import { CONFIG } from '@/app/constants';

export default function MegageraItems() {
  const [data, setData] = useState<ModalImage[] | null>(null);

  const { isOpen, onOpen, onClose } = useDisclosure();

  const [modalImage, setModalImage] = useState({} as ModalImage);
  const [newImage, setNewImage] = useState<FormData | null>(null);
  const [newImageSrc, setNewImageSrc] = useState<string | null>(null);
  const [focusPreviousImage, setFocusPreviousImage] = useState<string | null>(null);
  const [showModalView, setShowModalView] = useState<string | null>("main");
  const [editName, setEditName] = useState<string>("");
  const [isEditingName, setIsEditingName] = useState<boolean>(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetchMegageraLogos();
        const logos = await response.json();
        setData(logos);

      } catch (error) {
        console.error('Error fetching Megagera logos:', error);
      }
    };
  
    fetchData();
  }, []);
  

  const handleOpen = (image: ModalImage) => {
    setModalImage(image);
    setNewImage(null);
    setNewImageSrc(null);
    setFocusPreviousImage(null);
    setShowModalView("main");
    setEditName(image.name);
    setIsEditingName(false);
    onOpen();
  };

  const getSrcImageCacheBuster = (url: string) => {
    const cacheBuster = new Date().getTime();
    return `${CONFIG.megamediaServerApiUrl}/megagera/${url}?cb=${cacheBuster}`;  
  }

  const getDateByPreviousTeamID = (previousTeamID: string) => {
    const parts = previousTeamID.split('_');
    if (parts.length >= 5) {
      const datePart = parts[3];
      const year = datePart.substring(0, 4);
      const month = datePart.substring(4, 6);
      const day = datePart.substring(6, 8);
      return `${day}/${month}/${year}`;
    }
    return null;
  }

  const uploadImage = async (image: FormData) => {
    try {
      await uploadMegageraImage(modalImage.id, image);
      setShowModalView("main");
    } catch (error) {
      console.error(error);
    }
  }

  const deleteImage = async (name: string) => {
    try {
      await deleteMegageraImage(modalImage.id, name);
      setShowModalView("main");
    } catch (error) {
      console.error(error);
    }
  }

  const updateImageName = async () => {
    try {
      await updateMegageraImageName(modalImage.id, editName);
      setModalImage({...modalImage, name: editName});
      setIsEditingName(false);
      // Refresh the data
      const response = await fetchMegageraLogos();
      const logos = await response.json();
      setData(logos);
    } catch (error) {
      console.error(error);
    }
  }

  return (
  <>
    <div className="mt-5 gap-2 grid grid-cols-3 sm:grid-cols-5">

      {data && data.map((image: ModalImage, index: number) => (
        <Card className='w-[90%] m-auto h-full' key={index} isPressable shadow="sm" onPress={() => handleOpen(image)}>
          <CardBody className="overflow-visible p-0 m-auto">
            <Image
              alt={image.name}
              className="w-full object-cover h-full p-2"
              radius="lg"
              shadow="sm"
              src={getSrcImageCacheBuster(image.url.toString())}
              width="100%" />
          </CardBody>
          <CardFooter className="text-small grid">
            <b>{image.name}</b>
            <p className="text-default-500">{image.widthPX}px x {image.heightPX}px</p>
          </CardFooter>
        </Card>
      ))}

      <Modal isOpen={isOpen} size='md' onClose={onClose}>
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader className="flex flex-col gap-1 m-auto">
                {isEditingName ? (
                  <div className="flex gap-2 w-full">
                    <Input
                      value={editName}
                      onChange={(e) => setEditName(e.target.value)}
                      placeholder="Enter new name"
                      className="flex-1"
                    />
                    <Button 
                      color="primary" 
                      size="sm" 
                      onPress={updateImageName}
                      isDisabled={!editName.trim()}
                    >
                      Save
                    </Button>
                    <Button 
                      color="default" 
                      variant="light" 
                      size="sm" 
                      onPress={() => {
                        setIsEditingName(false);
                        setEditName(modalImage.name);
                      }}
                    >
                      Cancel
                    </Button>
                  </div>
                ) : (
                  <div className="flex justify-between items-center w-full">
                    <span>{modalImage.name}</span>
                    <Button 
                      color="primary" 
                      variant="light" 
                      size="sm" 
                      onPress={() => setIsEditingName(true)}
                    >
                      Edit Name
                    </Button>
                  </div>
                )}
              </ModalHeader>
                {showModalView === "main" && (
                  <>
                    <ModalBody className='gap-0 grid grid-cols-1'>
                      <p className="text-warning-500 italic mb-2 text-center">Main</p>
                      <Image
                        alt={modalImage.name}
                        className="w-full object-cover h-full p-2"
                        radius="lg"
                        shadow="sm"
                        src={getSrcImageCacheBuster(modalImage.url.toString())}
                        width="100%" />
                      <div>
                        <Button className="w-full mt-2 col-span-2" color="primary" onPress={() => document.getElementById('fileInput-main')?.click()}>
                          <h3 className="m-auto">Upload a new picture</h3>
                        </Button>
                        <input
                          id="fileInput-main"
                          type="file"
                          accept="image/*"
                          style={{ display: 'none' }}
                          onChange={(e) => {
                            const file = e.target.files?.[0];
                            if (file) {
                              const reader = new FileReader();
                              reader.onload = (event) => {
                                const newImageSrc = event.target?.result as string;
                                const formData = new FormData();
                                formData.append('image', file);
                                setNewImage(formData);
                                setNewImageSrc(newImageSrc);
                                setShowModalView("replace");
                              };
                              reader.readAsDataURL(file);
                            }
                          }}
                        />
                      </div>
                    </ModalBody>
                  </>
                )}
                {showModalView === "replace" && newImage && (
                  <>
                    <ModalBody className='gap-2 grid grid-cols-2'>
                      <h3 className="m-auto">Old image</h3>
                      <h3 className="m-auto">New image</h3>
                      <Image
                        alt={modalImage.name}
                        className="w-full object-cover h-full p-2"
                        radius="lg"
                        shadow="sm"
                        src={getSrcImageCacheBuster(modalImage.url.toString())}
                        width="100%" />
                      <Image
                        alt="New Image"
                        className="w-full object-cover h-full p-2"
                        radius="lg"
                        shadow="sm"
                        src={newImageSrc as string}
                        width="100%" />
                      <Button className="w-full mt-2 col-span-2" color="primary" onPress={() => document.getElementById('fileInput-replace')?.click()}>
                        <h3 className="m-auto">Upload another picture</h3>
                      </Button>
                      <input
                        id="fileInput-replace"
                        type="file"
                        accept="image/*"
                        style={{ display: 'none' }}
                        onChange={(e) => {
                          const file = e.target.files?.[0];
                          if (file) {
                            const reader = new FileReader();
                            reader.onload = (event) => {
                              const newImageSrc = event.target?.result as string;
                              const formData = new FormData();
                              formData.append('image', file);
                              setNewImage(formData);
                              setNewImageSrc(newImageSrc);
                              setShowModalView("replace");
                            };
                            reader.readAsDataURL(file);
                          }
                        }}
                      />
                    </ModalBody>
                  </>
                )}
                {showModalView === "previous" && (
                  <>
                    <ModalBody className={clsx('gap-2 grid grid-cols-3')}>
                      <Image
                        alt={modalImage.name}
                        className="w-full object-cover h-full p-2"
                        radius="lg"
                        shadow="sm"
                        src={getSrcImageCacheBuster(modalImage.url.toString())}
                        width="100%" 
                        onClick={() => {
                          setShowModalView("main");
                        }}
                        />

                      {modalImage.previous && modalImage.previous.length > 0 && modalImage.previous.map((previous_url: string) => (
                        <Image
                          key={previous_url}
                          alt={previous_url}
                          className="w-full object-cover h-full p-2"
                          radius="lg"
                          shadow="sm"
                          src={getSrcImageCacheBuster(previous_url.toString())}
                          width="100%"
                          onClick={() => {
                            setFocusPreviousImage(previous_url.toString());
                            setShowModalView("focusPrevious");
                          }}
                          />
                      ))}
                    </ModalBody>
                  </>
                )}
                {showModalView === "focusPrevious" && focusPreviousImage && (
                  <>
                    <ModalBody className='gap-0 grid grid-cols-1'>
                      <p className="text-warning-500 italic mb-2 text-center">Previous from: {getDateByPreviousTeamID(focusPreviousImage as string)}</p>
                      <Image
                        alt={focusPreviousImage as string}
                        className="w-full object-cover h-full p-2"
                        radius="lg"
                        shadow="sm"
                        src={getSrcImageCacheBuster(focusPreviousImage as string)}
                        width="100%" />
                    </ModalBody>
                  </>
                )}
              <ModalFooter>
                {showModalView === "main" && modalImage.previous && modalImage.previous.length > 0 && (
                  <Button color="primary" onPress={() => {
                    setShowModalView("previous");
                  }}>
                    See previous
                  </Button>
                )}
                {showModalView === "replace" && newImage && (
                  <Button color="primary" onPress={() => {
                    uploadImage(newImage);
                  }}>
                    Replace
                  </Button>
                )}
                {(showModalView === "previous" || showModalView === "replace") && (
                  <Button color="primary" onPress={() => {
                    setShowModalView("main");
                  }}>
                    Back
                  </Button>
                )}
                {showModalView === "focusPrevious" && (
                  <>
                    <Button color="primary" onPress={() => {
                      setFocusPreviousImage(null);
                      setShowModalView("previous");
                    } }>
                      Back
                    </Button>
                    <Button color="danger" onPress={() => {
                      deleteImage(focusPreviousImage as string);
                    } }>
                      Delete
                    </Button>
                  </>
                )}
                <Button color="danger" variant="light" onPress={() => {
                  setNewImage(null);
                  setNewImageSrc(null);
                  onClose();
                }}>
                  Cancel
                </Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>

    </div>

    </>
  );
};

