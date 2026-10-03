import { useEffect, useState } from 'react'
import {
  Box,
  Text,
  VStack,
  HStack,
  Badge,
  Divider,
  Skeleton,
} from '@chakra-ui/react'
import { useAuth } from '@/contexts/AuthContext'

export default function SessionInfo() {
  const { sessionToken } = useAuth()
  const [sessionAge, setSessionAge] = useState<string>('')

  useEffect(() => {
    if (!sessionToken) return

    // Extract timestamp from mock token (format: mock-session-token-{timestamp})
    const parts = sessionToken.split('-')
    const timestamp = parseInt(parts[parts.length - 1], 10)
    
    if (!isNaN(timestamp)) {
      const sessionStart = new Date(timestamp)
      
      const updateAge = () => {
        const now = new Date()
        const diff = now.getTime() - sessionStart.getTime()
        const minutes = Math.floor(diff / 60000)
        const hours = Math.floor(minutes / 60)
        
        if (hours > 0) {
          setSessionAge(`${hours} hour${hours !== 1 ? 's' : ''} ago`)
        } else if (minutes > 0) {
          setSessionAge(`${minutes} minute${minutes !== 1 ? 's' : ''} ago`)
        } else {
          setSessionAge('Just now')
        }
      }
      
      updateAge()
      const interval = setInterval(updateAge, 60000) // Update every minute
      
      return () => clearInterval(interval)
    }
  }, [sessionToken])

  return (
    <Box>
      <Text fontWeight="semibold" mb={3}>
        Session Information
      </Text>
      <VStack spacing={3} align="stretch">
        <HStack justify="space-between">
          <Text fontSize="sm" color="gray.600">
            Status
          </Text>
          <Badge colorScheme="green">Active</Badge>
        </HStack>
        
        <Divider />
        
        <HStack justify="space-between">
          <Text fontSize="sm" color="gray.600">
            Logged in
          </Text>
          {sessionAge ? (
            <Text fontSize="sm">{sessionAge}</Text>
          ) : (
            <Skeleton height="20px" width="100px" />
          )}
        </HStack>
        
        <Divider />
        
        <HStack justify="space-between">
          <Text fontSize="sm" color="gray.600">
            Session expires
          </Text>
          <Text fontSize="sm">In 24 hours</Text>
        </HStack>
      </VStack>
    </Box>
  )
}
