import { useEffect } from 'react'
import {
  FormControl,
  FormLabel,
  HStack,
  Button,
  useColorMode,
  useToast,
  Text,
} from '@chakra-ui/react'
import { MoonIcon, SunIcon } from '@chakra-ui/icons'
import { api } from '@/server/headscale/api'

export default function ThemeSelector() {
  const { colorMode, setColorMode } = useColorMode()
  const toast = useToast()

  useEffect(() => {
    // Load theme preference from settings
    const loadTheme = async () => {
      try {
        const settings = await api.getSettings()
        if (settings.theme) {
          setColorMode(settings.theme)
        }
      } catch (error) {
        // Ignore errors on initial load
      }
    }
    loadTheme()
  }, [setColorMode])

  const handleThemeChange = async (theme: 'light' | 'dark') => {
    try {
      await api.updateSettings({ theme })
      setColorMode(theme)
      toast({
        title: 'Theme updated',
        description: `Switched to ${theme} mode`,
        status: 'success',
        duration: 2000,
        isClosable: true,
      })
    } catch (error) {
      toast({
        title: 'Error updating theme',
        description: error instanceof Error ? error.message : 'Failed to update theme',
        status: 'error',
        duration: 5000,
        isClosable: true,
      })
    }
  }

  return (
    <FormControl>
      <FormLabel>Theme</FormLabel>
      <HStack spacing={4}>
        <Button
          leftIcon={<SunIcon />}
          onClick={() => handleThemeChange('light')}
          variant={colorMode === 'light' ? 'solid' : 'outline'}
          colorScheme={colorMode === 'light' ? 'brand' : 'gray'}
        >
          Light
        </Button>
        <Button
          leftIcon={<MoonIcon />}
          onClick={() => handleThemeChange('dark')}
          variant={colorMode === 'dark' ? 'solid' : 'outline'}
          colorScheme={colorMode === 'dark' ? 'brand' : 'gray'}
        >
          Dark
        </Button>
      </HStack>
      <Text fontSize="sm" color="gray.500" mt={2}>
        Your theme preference is saved and will persist across sessions
      </Text>
    </FormControl>
  )
}
