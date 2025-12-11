'use client';
import { useState } from 'react';
import { Button, Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, Input, useDisclosure } from "@nextui-org/react";
import { createMegageraImage } from '@/app/lib/data';

interface MegageraTopBarProps {
  onImageAdded: () => void;
}

export default function MegageraTopBar({ onImageAdded }: MegageraTopBarProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [imageName, setImageName] = useState('');
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSelectedFile(file);
    }
  };

  const handleUpload = async () => {
    if (!imageName.trim() || !selectedFile) {
      alert('Please provide both image name and select a file');
      return;
    }

    setIsUploading(true);
    try {
      const formData = new FormData();
      formData.append('image', selectedFile);
      formData.append('name', imageName.trim());

      const response = await createMegageraImage(imageName.trim(), formData);

      if (response.ok) {
        const result = await response.json();
        console.log('Image created successfully:', result);
        onImageAdded();
        onClose();
        setImageName('');
        setSelectedFile(null);
        // Reset file input
        const fileInput = document.getElementById('newImageFile') as HTMLInputElement;
        if (fileInput) fileInput.value = '';
      } else {
        const errorText = await response.text();
        console.error('Server error:', errorText);
        alert('Failed to create new image: ' + errorText);
      }
    } catch (error) {
      console.error('Error creating new image:', error);
      alert('Error creating new image: ' + error);
    } finally {
      setIsUploading(false);
    }
  };

  const handleClose = () => {
    onClose();
    setImageName('');
    setSelectedFile(null);
    const fileInput = document.getElementById('newImageFile') as HTMLInputElement;
    if (fileInput) fileInput.value = '';
  };

  return (
    <>
      <div className="mb-5 flex justify-end">
        <Button 
          color="primary" 
          onPress={onOpen}
          className="font-semibold"
        >
          Add New Image
        </Button>
      </div>

      <Modal isOpen={isOpen} onClose={handleClose} size="md">
        <ModalContent>
          {() => (
            <>
              <ModalHeader className="flex flex-col gap-1">
                Add New Image
              </ModalHeader>
              <ModalBody>
                <Input
                  label="Image Name"
                  placeholder="Enter image name"
                  value={imageName}
                  onChange={(e) => setImageName(e.target.value)}
                  isRequired
                />
                <div>
                  <input
                    id="newImageFile"
                    type="file"
                    accept="image/*"
                    onChange={handleFileChange}
                    className="hidden"
                  />
                  <Button
                    color="primary"
                    variant="bordered"
                    onPress={() => document.getElementById('newImageFile')?.click()}
                    className="w-full"
                  >
                    {selectedFile ? selectedFile.name : 'Select Image File'}
                  </Button>
                </div>
              </ModalBody>
              <ModalFooter>
                <Button color="danger" variant="light" onPress={handleClose}>
                  Cancel
                </Button>
                <Button 
                  color="primary" 
                  onPress={handleUpload}
                  isLoading={isUploading}
                  isDisabled={!imageName.trim() || !selectedFile}
                >
                  Create Image
                </Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>
    </>
  );
}
