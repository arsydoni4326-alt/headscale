import { Outlet, Navigate } from 'react-router-dom'
import {
  Box,
  Container,
  Flex,
  Heading,
  Button,
  useColorMode,
  IconButton,
  HStack,
  Spacer,
} from '@chakra-ui/react'
import { MoonIcon, SunIcon } from '@chakra-ui/icons'
import { useAuth } from '../contexts/AuthContext'

export default function Layout() {
  const { isAuthenticated, logout } = useAuth()
  const { colorMode, toggleColorMode } = useColorMode()

  if (!isAuthenticated) {
    return <Navigate to="/auth/login" replace />
  }

  return (
    <Box minH="100vh">
      {/* Header */}
      <Box bg={colorMode === 'dark' ? 'gray.800' : 'white'} boxShadow="sm" py={4}>
        <Container maxW="container.xl">
          <Flex align="center">
            <Heading size="md" color="brand.600">
              Headplane Settings
            </Heading>
            <Spacer />
            <HStack spacing={4}>
              <IconButton
                aria-label="Toggle color mode"
                icon={colorMode === 'light' ? <MoonIcon /> : <SunIcon />}
                onClick={toggleColorMode}
                variant="ghost"
              />
              <Button onClick={logout} variant="outline" size="sm">
                Logout
              </Button>
            </HStack>
          </Flex>
        </Container>
      </Box>

      {/* Main content */}
      <Container maxW="container.xl" py={8}>
        <Outlet />
      </Container>
    </Box>
  )
}
