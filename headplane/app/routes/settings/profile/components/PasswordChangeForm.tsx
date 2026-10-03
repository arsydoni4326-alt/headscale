import { useState } from 'react'
import {
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

export default function PasswordChangeForm() {
  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showPasswords, setShowPasswords] = useState({
    current: false,
    new: false,
    confirm: false,
  })
  const [isLoading, setIsLoading] = useState(false)
  const [errors, setErrors] = useState({
    current: '',
    new: '',
    confirm: '',
  })
  const toast = useToast()

  const validateNewPassword = (password: string): boolean => {
    if (!password) {
      setErrors((prev) => ({ ...prev, new: 'New password is required' }))
      return false
    }
    if (password.length < 8) {
      setErrors((prev) => ({ ...prev, new: 'Password must be at least 8 characters' }))
      return false
    }
    setErrors((prev) => ({ ...prev, new: '' }))
    return true
  }

  const validateConfirmPassword = (password: string): boolean => {
    if (password !== newPassword) {
      setErrors((prev) => ({ ...prev, confirm: 'Passwords do not match' }))
      return false
    }
    setErrors((prev) => ({ ...prev, confirm: '' }))
    return true
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!currentPassword) {
      setErrors((prev) => ({ ...prev, current: 'Current password is required' }))
      return
    }

    if (!validateNewPassword(newPassword) || !validateConfirmPassword(confirmPassword)) {
      return
    }

    setIsLoading(true)
    try {
      await api.changePassword({ currentPassword, newPassword })
      toast({
        title: 'Password changed',
        description: 'Your password has been changed successfully.',
        status: 'success',
        duration: 3000,
        isClosable: true,
      })
      setCurrentPassword('')
      setNewPassword('')
      setConfirmPassword('')
      setErrors({ current: '', new: '', confirm: '' })
    } catch (error) {
      toast({
        title: 'Error changing password',
        description: error instanceof Error ? error.message : 'Failed to change password',
        status: 'error',
        duration: 5000,
        isClosable: true,
      })
      setErrors((prev) => ({ ...prev, current: 'Invalid current password' }))
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <VStack spacing={4} align="stretch">
        <FormControl isInvalid={!!errors.current} isRequired>
          <FormLabel>Current Password</FormLabel>
          <InputGroup>
            <Input
              type={showPasswords.current ? 'text' : 'password'}
              value={currentPassword}
              onChange={(e) => {
                setCurrentPassword(e.target.value)
                setErrors((prev) => ({ ...prev, current: '' }))
              }}
            />
            <InputRightElement>
              <IconButton
                aria-label={showPasswords.current ? 'Hide password' : 'Show password'}
                icon={showPasswords.current ? <ViewOffIcon /> : <ViewIcon />}
                onClick={() =>
                  setShowPasswords((prev) => ({ ...prev, current: !prev.current }))
                }
                variant="ghost"
                size="sm"
              />
            </InputRightElement>
          </InputGroup>
          {errors.current && <FormErrorMessage>{errors.current}</FormErrorMessage>}
        </FormControl>

        <FormControl isInvalid={!!errors.new} isRequired>
          <FormLabel>New Password</FormLabel>
          <InputGroup>
            <Input
              type={showPasswords.new ? 'text' : 'password'}
              value={newPassword}
              onChange={(e) => {
                setNewPassword(e.target.value)
                setErrors((prev) => ({ ...prev, new: '' }))
              }}
              onBlur={() => newPassword && validateNewPassword(newPassword)}
            />
            <InputRightElement>
              <IconButton
                aria-label={showPasswords.new ? 'Hide password' : 'Show password'}
                icon={showPasswords.new ? <ViewOffIcon /> : <ViewIcon />}
                onClick={() =>
                  setShowPasswords((prev) => ({ ...prev, new: !prev.new }))
                }
                variant="ghost"
                size="sm"
              />
            </InputRightElement>
          </InputGroup>
          {errors.new ? (
            <FormErrorMessage>{errors.new}</FormErrorMessage>
          ) : (
            <FormHelperText>Must be at least 8 characters long</FormHelperText>
          )}
        </FormControl>

        <FormControl isInvalid={!!errors.confirm} isRequired>
          <FormLabel>Confirm New Password</FormLabel>
          <InputGroup>
            <Input
              type={showPasswords.confirm ? 'text' : 'password'}
              value={confirmPassword}
              onChange={(e) => {
                setConfirmPassword(e.target.value)
                setErrors((prev) => ({ ...prev, confirm: '' }))
              }}
              onBlur={() => confirmPassword && validateConfirmPassword(confirmPassword)}
            />
            <InputRightElement>
              <IconButton
                aria-label={showPasswords.confirm ? 'Hide password' : 'Show password'}
                icon={showPasswords.confirm ? <ViewOffIcon /> : <ViewIcon />}
                onClick={() =>
                  setShowPasswords((prev) => ({ ...prev, confirm: !prev.confirm }))
                }
                variant="ghost"
                size="sm"
              />
            </InputRightElement>
          </InputGroup>
          {errors.confirm && <FormErrorMessage>{errors.confirm}</FormErrorMessage>}
        </FormControl>

        <Button
          type="submit"
          colorScheme="brand"
          isLoading={isLoading}
          isDisabled={!currentPassword || !newPassword || !confirmPassword}
        >
          Change Password
        </Button>
      </VStack>
    </form>
  )
}

