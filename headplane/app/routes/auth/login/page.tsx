import { useState } from 'react'
import { Navigate } from 'react-router-dom'
import {
  Box,
  Container,
  VStack,
  Heading,
  FormControl,
  FormLabel,
  Input,
  Button,
  Text,
  useToast,
  Card,
  CardBody,
  InputGroup,
  InputRightElement,
  IconButton,
  useColorMode,
} from '@chakra-ui/react'
import { ViewIcon, ViewOffIcon } from '@chakra-ui/icons'
import { useAuth } from '@/contexts/AuthContext'
import { api } from '@/server/headscale/api'

export default function LoginPage() {
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const { isAuthenticated, login } = useAuth()
  const toast = useToast()
  const { colorMode } = useColorMode()

  if (isAuthenticated) {
    return <Navigate to="/settings/profile" replace />
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)

    try {
      const response = await api.login({ password })
      api.setSessionToken(response.token)
      login(response.token)
      toast({
        title: 'Login successful',
        status: 'success',
        duration: 3000,
        isClosable: true,
      })
    } catch (error) {
      toast({
        title: 'Login failed',
        description: error instanceof Error ? error.message : 'Invalid password',
        status: 'error',
        duration: 5000,
        isClosable: true,
      })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Box
      minH="100vh"
      bg={colorMode === 'dark' ? 'gray.900' : 'gray.50'}
      display="flex"
      alignItems="center"
      justifyContent="center"
    >
      <Container maxW="md">
        <Card>
          <CardBody>
            <VStack spacing={6} align="stretch">
              <Box textAlign="center">
                <Heading size="lg" mb={2}>
                  Headplane
                </Heading>
                <Text color="gray.600">Sign in to access settings</Text>
              </Box>

              <form onSubmit={handleSubmit}>
                <VStack spacing={4}>
                  <FormControl isRequired>
                    <FormLabel>Password</FormLabel>
                    <InputGroup>
                      <Input
                        type={showPassword ? 'text' : 'password'}
                        placeholder="Enter your password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        autoFocus
                      />
                      <InputRightElement>
                        <IconButton
                          aria-label={showPassword ? 'Hide password' : 'Show password'}
                          icon={showPassword ? <ViewOffIcon /> : <ViewIcon />}
                          onClick={() => setShowPassword(!showPassword)}
                          variant="ghost"
                          size="sm"
                        />
                      </InputRightElement>
                    </InputGroup>
                  </FormControl>

                  <Button
                    type="submit"
                    colorScheme="brand"
                    width="full"
                    isLoading={isLoading}
                  >
                    Sign In
                  </Button>
                </VStack>
              </form>

              <Box
                p={4}
                bg={colorMode === 'dark' ? 'gray.700' : 'blue.50'}
                borderRadius="md"
              >
                <Text fontSize="sm" color={colorMode === 'dark' ? 'gray.300' : 'blue.800'}>
                  <strong>Demo credentials:</strong>
                  <br />
                  Password: <code>password123</code>
                </Text>
              </Box>
            </VStack>
          </CardBody>
        </Card>
      </Container>
    </Box>
  )
}
