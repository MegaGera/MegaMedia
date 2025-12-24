'use client';

import {
  PlusCircleIcon
} from '@heroicons/react/24/outline';

import { useState, useEffect } from 'react';
import { Card, CardBody, CardFooter, Image, Button,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  useDisclosure
} from "@nextui-org/react";

import { useTopLeagues, uploadLeagueImage } from '@/app/lib/data';
import { League } from '@/app/modals/league';
import { CONFIG } from '@/app/constants';

interface MegagoalLeaguesItemsProps {
  order: number;
}

export default function MegagoalLeaguesItems({ order }: MegagoalLeaguesItemsProps) {
  const [showLeagues, setShowLeagues] = useState<League[]>([]);
  const [showN, setShowN] = useState(20);

  const { data: leagues } = useTopLeagues();
  const { isOpen, onOpen, onClose } = useDisclosure();

  const [modalLeague, setModalLeague] = useState({} as League);
  const [showModalView, setShowModalView] = useState<string | null>("main");
  const [imageType, setImageType] = useState<'sm' | 'back'>('sm');
  const [newImage, setNewImage] = useState<FormData | null>(null);
  const [newImageSrc, setNewImageSrc] = useState<string | null>(null);

  useEffect(() => {
    if (leagues) {
      if (order === 1) {
        setShowLeagues([...leagues].sort((a: League, b: League) => a.league_id - b.league_id));
      } else if (order === 2) {
        setShowLeagues([...leagues].sort((a: League, b: League) => b.league_id - a.league_id));
      } else if (order === 3) {
        setShowLeagues([...leagues].sort((a: League, b: League) => a.league_name.localeCompare(b.league_name)));
      } else if (order === 4) {
        setShowLeagues([...leagues].sort((a: League, b: League) => b.league_name.localeCompare(a.league_name)));
      } else {
        setShowLeagues(leagues);
      }
    }
  }, [leagues, order]);

  const handleOpen = (league: League) => {
    setModalLeague(league);
    setShowModalView("main");
    setImageType('sm');
    setNewImage(null);
    setNewImageSrc(null);
    onOpen();
  };

  const getSrcImageByLeagueIDCacheBuster = (leagueID: string, type: 'sm' | 'back' = 'sm') => {
    const cacheBuster = new Date().getTime();
    return `${CONFIG.megamediaServerApiUrl}/megagoal/leagues/${type}/${leagueID}.png?cb=${cacheBuster}`;  
  }


  const loadMoreLeagues = () => {
    setShowN(prevShowLeagues => prevShowLeagues + 20);
  }

  const uploadImage = async (leagueID: string, imageType: 'sm' | 'back', image: FormData) => {
    try {
      await uploadLeagueImage(leagueID, imageType, image);
      setShowModalView("main");
      setNewImage(null);
      setNewImageSrc(null);
    } catch (error) {
      console.error(error);
    }
  }

  return (
  <>
    <div className="mt-5 gap-2 grid grid-cols-3 sm:grid-cols-10">

      {showLeagues && showLeagues.slice(0, showN).map((league: League, index: number) => (
        <Card className='w-[90%] m-auto h-full' key={index} isPressable shadow="sm" onPress={() => handleOpen(league)}>
          <CardBody className="overflow-visible p-0 m-auto">
            <Image
              alt={league.league_name}
              className="w-full object-cover h-full p-2"
              radius="lg"
              shadow="sm"
              src={getSrcImageByLeagueIDCacheBuster(league.league_id.toString())}
              width="100%" />
          </CardBody>
          <CardFooter className="text-small grid">
            <b>{league.league_name}</b>
            <p className="text-default-500">{league.league_id}</p>
          </CardFooter>
        </Card>
      ))}
      {showLeagues && leagues && leagues.length > showN && (
        <div className="col-span-10 flex justify-center m-auto mt-4">
          <Button size="lg" isIconOnly aria-label="Load more" color="primary" onPress={loadMoreLeagues}>
            <PlusCircleIcon className='p-1'/>
          </Button>
        </div>
      )}

      <Modal isOpen={isOpen} size='md' onClose={onClose}>
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader className="flex flex-col gap-1 m-auto">{modalLeague.league_name || 'League'}</ModalHeader>
                {showModalView === "main" && modalLeague.league_name && (
                  <>
                    <ModalBody className='gap-0 grid grid-cols-1'>
                      <p className="text-warning-500 italic mb-2 text-center">Main ({imageType === 'sm' ? 'Small' : 'Back'})</p>
                      <Image
                        alt={modalLeague.league_name || 'League'}
                        className="w-full object-cover h-full p-2"
                        radius="lg"
                        shadow="sm"
                        src={getSrcImageByLeagueIDCacheBuster(modalLeague.league_id.toString(), imageType)}
                        width="100%" />
                      <div>
                        <Button className="w-full mt-2 col-span-2" color="primary" onPress={() => document.getElementById(`fileInput-main-${imageType}`)?.click()}>
                          <h3 className="m-auto">Upload a new {imageType === 'sm' ? 'small' : 'back'} picture</h3>
                        </Button>
                        <input
                          id={`fileInput-main-${imageType}`}
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
                        alt={modalLeague.league_name || 'League'}
                        className="w-full object-cover h-full p-2"
                        radius="lg"
                        shadow="sm"
                        src={getSrcImageByLeagueIDCacheBuster(modalLeague.league_id.toString(), imageType)}
                        width="100%" />
                      <Image
                        alt="New Image"
                        className="w-full object-cover h-full p-2"
                        radius="lg"
                        shadow="sm"
                        src={newImageSrc as string}
                        width="100%" />
                      <Button className="w-full mt-2 col-span-2" color="primary" onPress={() => document.getElementById(`fileInput-replace-${imageType}`)?.click()}>
                        <h3 className="m-auto">Upload another picture</h3>
                      </Button>
                      <input
                        id={`fileInput-replace-${imageType}`}
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
              <ModalFooter>
                {showModalView === "main" && (
                  <Button 
                    color="primary" 
                    variant={imageType === 'sm' ? 'solid' : 'bordered'}
                    onPress={() => {
                      setImageType(imageType === 'sm' ? 'back' : 'sm');
                    }}
                  >
                    {imageType === 'sm' ? 'Show Back' : 'Show Small'}
                  </Button>
                )}
                {showModalView === "replace" && newImage && (
                  <Button color="primary" onPress={() => {
                    uploadImage(modalLeague.league_id.toString(), imageType, newImage);
                  }}>
                    Replace
                  </Button>
                )}
                {(showModalView === "main" || showModalView === "replace") && (
                  <Button color="danger" variant="light" onPress={() => {
                    if (showModalView === "replace") {
                      setNewImage(null);
                      setNewImageSrc(null);
                      setShowModalView("main");
                    } else {
                      onClose();
                    }
                  }}>
                    {showModalView === "replace" ? 'Back' : 'Cancel'}
                  </Button>
                )}
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>

    </div>

    </>
  );
};

