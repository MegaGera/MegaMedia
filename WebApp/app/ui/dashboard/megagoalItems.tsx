'use client';

import {
  PlusCircleIcon
} from '@heroicons/react/24/outline';

import { useState, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';
import { Card, CardBody, CardFooter, Image, Button,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  useDisclosure
} from "@nextui-org/react";

import clsx from 'clsx';

import { useTeams, uploadTeamImage, squaredTeamImage, deleteTeamImage } from '@/app/lib/data';
import { Team } from '@/app/modals/team';
import { CONFIG } from '@/app/constants';

interface MegagoalItemsProps {
  order: number;
}

export default function MegagoalItems({ order }: MegagoalItemsProps) {
  const searchParams = useSearchParams();
  const countryFilter = searchParams?.get('country') || null;
  const leagueFilter = searchParams?.get('league_id') || null;
  const seasonFilter = searchParams?.get('season') || null;
  const [showTeams, setShowTeams] = useState<Team[]>([]);
  const [showN, setShowN] = useState(20);

  const { data: teams } = useTeams(countryFilter, leagueFilter, seasonFilter);
  const { isOpen, onOpen, onClose } = useDisclosure();

  const [modalTeam, setModalTeam] = useState({} as Team);
  const [newImage, setNewImage] = useState<File | null>(null);
  const [newImageSrc, setNewImageSrc] = useState<string | null>(null);
  const [focusPreviousImage, setFocusPreviousImage] = useState<string | null>(null);
  const [showModalView, setShowModalView] = useState<string | null>("main");

  useEffect(() => {
    setShowTeams(teams)
    setShowTeams((prevTeams) => {
      if (order === 1) {
        return [...teams].sort((a: Team, b: Team) => a.team.id - b.team.id);
      } else if (order === 2) {
        return [...teams].sort((a: Team, b: Team) => b.team.id - a.team.id);
      } else if (order === 3) {
        return [...teams].sort((a: Team, b: Team) => a.team.name.localeCompare(b.team.name));
      } else if (order === 4) {
        return [...teams].sort((a: Team, b: Team) => b.team.name.localeCompare(a.team.name));
      }
      return prevTeams;
    });
  }, [teams, order]);

  const handleOpen = (team: Team) => {
    setModalTeam(team);
    setNewImage(null);
    setNewImageSrc(null);
    setFocusPreviousImage(null);
    setShowModalView("main");
    onOpen();
  };

  const getSrcImageByTeamID = (teamID: string) => {
    return `${CONFIG.megamediaServerApiUrl}/megagoal/teams/team_${teamID}.png`;
  }

  const getSrcImageByTeamIDCacheBuster = (teamID: string) => {
    const cacheBuster = new Date().getTime();
    return `${CONFIG.megamediaServerApiUrl}/megagoal/teams/team_${teamID}.png?cb=${cacheBuster}`;  
  }

  const getSrcImageByPreviousTeamID = (previousTeamID: string) => {
    return `${CONFIG.megamediaServerApiUrl}/megagoal/teams/${previousTeamID}`;
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

  const loadMoreTeams = () => {
    setShowN(prevShowTeams => prevShowTeams + 20);
  }

  const uploadImage = async (teamID: string, image: File) => {
    try {
      await uploadTeamImage(teamID, image);
    } catch (error) {
      console.error(error);
    }
  }

  const squaredImage = async (teamID: string) => {
    try {
      const r = await squaredTeamImage(teamID);
      if (r.status == 200) {
        setModalTeam((prevModalTeam) => ({
          ...prevModalTeam,
          team: {
            ...prevModalTeam.team,
            id: Number(teamID),
          },
        }));
        setShowModalView("main");
      }
    } catch (error) {
      console.error(error);
    }
  }

  const deleteImage = async (teamID: string) => {
    try {
      await deleteTeamImage(teamID);
      setShowModalView("main");
    } catch (error) {
      console.error(error);
    }
  }

  return (
  <>
    <div className="mt-5 gap-2 grid grid-cols-3 sm:grid-cols-10">

      {showTeams && showTeams.slice(0, showN).map((team: Team, index: number) => (
        <Card className='w-[90%] m-auto h-full' key={index} isPressable shadow="sm" onPress={() => handleOpen(team)}>
          <CardBody className="overflow-visible p-0 m-auto">
            <Image
              alt={team.team.name}
              className="w-full object-cover h-full p-2"
              radius="lg"
              shadow="sm"
              // src={getSrcImageByTeamID(team.team.id.toString())}
              src={getSrcImageByTeamIDCacheBuster(team.team.id.toString())}
              width="100%" />
          </CardBody>
          <CardFooter className="text-small grid">
            <b>{team.team.name}</b>
            <p className="text-default-500">{team.team.id}</p>
          </CardFooter>
        </Card>
      ))}
      {showTeams && teams.length > showN && (
        <div className="col-span-10 flex justify-center m-auto mt-4">
          <Button size="lg" isIconOnly aria-label="Load more" color="primary" onPress={loadMoreTeams}>
            <PlusCircleIcon className='p-1'/>
          </Button>
        </div>
      )}

      <Modal isOpen={isOpen} size='md' onClose={onClose}>
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader className="flex flex-col gap-1 m-auto">{modalTeam.team.name}</ModalHeader>
                {showModalView === "main" && (
                  <>
                    <ModalBody className='gap-0 grid grid-cols-1'>
                      <p className="text-warning-500 italic mb-2 text-center">Main</p>
                      <Image
                        alt={modalTeam.team.name}
                        className="w-full object-cover h-full p-2"
                        radius="lg"
                        shadow="sm"
                        src={getSrcImageByTeamIDCacheBuster(modalTeam.team.id.toString())}
                        width="100%" />
                      <div>
                        <Button className="w-full mt-2 col-span-2" color="primary" onPress={() => document.getElementById('fileInput')?.click()}>
                          <h3 className="m-auto">Upload a new picture</h3>
                        </Button>
                        <input
                          id="fileInput"
                          type="file"
                          accept="image/*"
                          style={{ display: 'none' }}
                          onChange={(e) => {
                            const file = e.target.files?.[0];
                            if (file) {
                              const reader = new FileReader();
                              reader.onload = (event) => {
                                const newImageSrc = event.target?.result as string;
                                setNewImage(file);
                                setNewImageSrc(newImageSrc);
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
                    <ModalBody className='gap-2 grid-cols-2'>
                      <h3 className="m-auto">Old image</h3>
                      <h3 className="m-auto">New image</h3>
                      <Image
                        alt={modalTeam.team.name}
                        className="w-full object-cover h-full p-2"
                        radius="lg"
                        shadow="sm"
                        src={getSrcImageByTeamIDCacheBuster(modalTeam.team.id.toString())}
                        width="100%" />
                      <Image
                        alt="New Image"
                        className="w-full object-cover h-full p-2"
                        radius="lg"
                        shadow="sm"
                        src={newImageSrc as string}
                        width="100%" />
                      <Button className="w-full mt-2 col-span-2" color="primary" onPress={() => document.getElementById('fileInput')?.click()}>
                        <h3 className="m-auto">Upload another picture</h3>
                      </Button>
                      <input
                        id="fileInput"
                        type="file"
                        accept="image/*"
                        style={{ display: 'none' }}
                        onChange={(e) => {
                          const file = e.target.files?.[0];
                          if (file) {
                            const reader = new FileReader();
                            reader.onload = (event) => {
                              const newImageSrc = event.target?.result as string;
                              setNewImage(file);
                              setNewImageSrc(newImageSrc);
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
                        alt={modalTeam.team.name}
                        className="w-full object-cover h-full p-2"
                        radius="lg"
                        shadow="sm"
                        src={getSrcImageByTeamIDCacheBuster(modalTeam.team.id.toString())}
                        width="100%" 
                        onClick={() => {
                          setShowModalView("main");
                        }}
                        />

                      {modalTeam.previous.length > 0 && modalTeam.previous.slice(0, modalTeam.previous.length).map((team_id: string, index: number) => (
                        <Image
                          key={team_id}
                          alt={team_id}
                          className="w-full object-cover h-full p-2"
                          radius="lg"
                          shadow="sm"
                          src={getSrcImageByPreviousTeamID(team_id.toString())}
                          width="100%"
                          onClick={() => {
                            setFocusPreviousImage(team_id.toString());
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
                        src={getSrcImageByPreviousTeamID(focusPreviousImage as string)}
                        width="100%" />
                    </ModalBody>
                  </>
                )}
              <ModalFooter>
                {showModalView === "main" && modalTeam.previous && modalTeam.previous.length > 0 && (
                  <Button color="primary" onPress={() => {
                    setShowModalView("previous");
                  }}>
                    See previous
                  </Button>
                )}
                {showModalView === "main" && (
                  <Button color="primary" onPress={() => {
                    squaredImage(modalTeam.team.id.toString());
                  }}>
                    Squared
                  </Button>
                )}
                {showModalView === "replace" && newImage && (
                  <Button color="primary" onPress={() => {
                    uploadImage(modalTeam.team.id.toString(), newImage as File);
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

