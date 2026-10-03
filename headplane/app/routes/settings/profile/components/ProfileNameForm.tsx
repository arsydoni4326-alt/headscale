import { useState, useEffect } from 'react'
import {
  FormControl,
  FormLabel,
  Input,
  Button,
  VStack,
  useToast,
  FormHelperText,
} from '@chakra-ui/react'
import { api } from '@/server/headscale/api'

export default function ProfileNameForm() {
  const [profileName, setProfileName] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [hasChanges, setHasChanges] = useState(false)
  const toast = useToast()

  useEffect(() => {
    const loadProfileName = async () => {
      try {
        const settings = await api.getSettings()
        setProfileName(settings.profileName || '')
      } catch (error) {
        // Ignore errors on initial load
      }
    }
    loadProfileName()
  }, [])

  const handleSave = async () => {
    setIsLoading(true)
    try {
      await api.updateSettings({ profileName })
      toast({
        title: 'Profile name saved',
        description: 'Your profile name has been updated successfully.',
        status: 'success',
        duration: 3000,
        isClosable: true,
      })
      setHasChanges(false)
    } catch (error) {
      toast({
        title: 'Error saving profile name',
        description: error instanceof Error ? error.message : 'Failed to save profile name',
        status: 'error',
        duration: 5000,
        isClosable: true,
      })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <VStack spacing={4} align="stretch">
      <FormControl>
        <FormLabel>Display Name</FormLabel>
        <Input
          placeholder="Enter your name"
          value={profileName}
          onChange={(e) => {
            setProfileName(e.target.value)
            setHasChanges(true)
          }}
        />
        <FormHelperText>
          This name is used for display purposes only and is optional.
        </FormHelperText>
      </FormControl>

      <Button
        colorScheme="brand"
        onClick={handleSave}
        isLoading={isLoading}
        isDisabled={!hasChanges}
      >
        Save Profile Name
      </Button>
    </VStack>
  )
}
