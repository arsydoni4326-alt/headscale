import {
  Box,
  Container,
  Heading,
  VStack,
  Divider,
  Card,
  CardBody,
  useColorMode,
} from '@chakra-ui/react'
import APIKeyForm from './components/APIKeyForm'
import PasswordChangeForm from './components/PasswordChangeForm'
import ThemeSelector from './components/ThemeSelector'
import ProfileNameForm from './components/ProfileNameForm'
import SessionInfo from './components/SessionInfo'

export default function SettingsProfilePage() {
  const { colorMode } = useColorMode()

  return (
    <Container maxW="container.lg" py={8}>
      <VStack spacing={8} align="stretch">
        <Box>
          <Heading size="xl" mb={2}>
            Settings
          </Heading>
          <Box color="gray.600">
            Manage your Headplane account settings and preferences
          </Box>
        </Box>

        {/* Account Section */}
        <Card>
          <CardBody>
            <VStack spacing={6} align="stretch">
              <Box>
                <Heading size="md" mb={4}>
                  Account
                </Heading>
                <VStack spacing={6} align="stretch">
                  <PasswordChangeForm />
                  <Divider />
                  <SessionInfo />
                </VStack>
              </Box>
            </VStack>
          </CardBody>
        </Card>

        {/* Integration Section */}
        <Card>
          <CardBody>
            <VStack spacing={6} align="stretch">
              <Box>
                <Heading size="md" mb={4}>
                  Integration
                </Heading>
                <APIKeyForm />
              </Box>
            </VStack>
          </CardBody>
        </Card>

        {/* Preferences Section */}
        <Card>
          <CardBody>
            <VStack spacing={6} align="stretch">
              <Box>
                <Heading size="md" mb={4}>
                  Preferences
                </Heading>
                <ThemeSelector />
              </Box>
            </VStack>
          </CardBody>
        </Card>

        {/* Profile Section */}
        <Card>
          <CardBody>
            <VStack spacing={6} align="stretch">
              <Box>
                <Heading size="md" mb={4}>
                  Profile
                </Heading>
                <ProfileNameForm />
              </Box>
            </VStack>
          </CardBody>
        </Card>
      </VStack>
    </Container>
  )
}
