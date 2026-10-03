import { useState } from 'react'
import {
  Box,
  FormControl,
  FormLabel,
  Input,
  Button,
  VStack,
  useToast,
  InputGroup,
  InputRightElement,
  IconButton,
  FormHelperText,
  FormErrorMessage,
} from '@chakra-ui/react'
import { ViewIcon, ViewOffIcon } from '@chakra-ui/icons'
import { api } from '@/server/headscale/api'

export default function APIKeyForm() {
  const [apiKey, setApiKey] = useState('')
  const [showApiKey, setShowApiKey] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const toast = useToast()

  const validateApiKey = (key: string): boolean => {
    if (!key) {
      setError('API key is required')
      return false
    }
    if (!key.startsWith('hs_')) {
      setError('API key must start with "hs_"')
      return false
    }
    if (key.length < 10) {
      setError('API key is too short')
      return false
    }
    setError('')
    return true
  }

  const handleSave = async () => {
    if (!validateApiKey(apiKey)) {
      return
    }

    setIsLoading(true)
    try {
      await api.updateSettings({ apiKey })
      toast({
        title: 'API key saved',
        description: 'Your Headscale API key has been saved successfully.',
        status: 'success',
        duration: 3000,
        isClosable: true,
      })
      setApiKey('')
      setShowApiKey(false)
    } catch (error) {
      toast({
        title: 'Error saving API key',
        description: error instanceof Error ? error.message : 'Failed to save API key',
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
      <FormControl isInvalid={!!error}>
        <FormLabel>Headscale API Key</FormLabel>
        <InputGroup>
          <Input
            type={showApiKey ? 'text' : 'password'}
            placeholder="hs_..."
            value={apiKey}
            onChange={(e) => {
              setApiKey(e.target.value)
              setError('')
            }}
            onBlur={() => apiKey && validateApiKey(apiKey)}
          />
          <InputRightElement>
            <IconButton
              aria-label={showApiKey ? 'Hide API key' : 'Show API key'}
              icon={showApiKey ? <ViewOffIcon /> : <ViewIcon />}
              onClick={() => setShowApiKey(!showApiKey)}
              variant="ghost"
              size="sm"
            />
          </InputRightElement>
        </InputGroup>
        {error ? (
          <FormErrorMessage>{error}</FormErrorMessage>
        ) : (
          <FormHelperText>
            Your API key will be encrypted and stored securely. Enter a new key to update.
          </FormHelperText>
        )}
      </FormControl>

      <Button
        colorScheme="brand"
        onClick={handleSave}
        isLoading={isLoading}
        isDisabled={!apiKey || !!error}
      >
        Save API Key
      </Button>
    </VStack>
  )
}
