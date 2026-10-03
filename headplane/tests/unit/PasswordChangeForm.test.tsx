import { describe, it, expect, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ChakraProvider } from '@chakra-ui/react'
import PasswordChangeForm from '@/routes/settings/profile/components/PasswordChangeForm'
import * as api from '@/server/headscale/api'

vi.mock('@/server/headscale/api', () => ({
  api: {
    changePassword: vi.fn(),
  },
}))

const renderWithChakra = (ui: React.ReactElement) => {
  return render(<ChakraProvider>{ui}</ChakraProvider>)
}

describe('PasswordChangeForm', () => {
  it('renders all password fields', () => {
    renderWithChakra(<PasswordChangeForm />)
    expect(screen.getByLabelText(/current password/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/^new password$/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/confirm new password/i)).toBeInTheDocument()
  })

  it('validates password length', async () => {
    const user = userEvent.setup()
    renderWithChakra(<PasswordChangeForm />)

    const newPasswordInput = screen.getByLabelText(/^new password$/i)
    await user.type(newPasswordInput, 'short')
    await user.tab()

    await waitFor(() => {
      expect(screen.getByText(/password must be at least 8 characters/i)).toBeInTheDocument()
    })
  })

  it('validates password confirmation match', async () => {
    const user = userEvent.setup()
    renderWithChakra(<PasswordChangeForm />)

    const newPasswordInput = screen.getByLabelText(/^new password$/i)
    const confirmPasswordInput = screen.getByLabelText(/confirm new password/i)

    await user.type(newPasswordInput, 'password123')
    await user.type(confirmPasswordInput, 'password456')
    await user.tab()

    await waitFor(() => {
      expect(screen.getByText(/passwords do not match/i)).toBeInTheDocument()
    })
  })

  it('submits form with valid passwords', async () => {
    const user = userEvent.setup()
    const mockChangePassword = vi.spyOn(api.api, 'changePassword')
    mockChangePassword.mockResolvedValue({ success: true })

    renderWithChakra(<PasswordChangeForm />)

    await user.type(screen.getByLabelText(/current password/i), 'oldpass123')
    await user.type(screen.getByLabelText(/^new password$/i), 'newpass123')
    await user.type(screen.getByLabelText(/confirm new password/i), 'newpass123')
    await user.click(screen.getByRole('button', { name: /change password/i }))

    await waitFor(() => {
      expect(mockChangePassword).toHaveBeenCalledWith({
        currentPassword: 'oldpass123',
        newPassword: 'newpass123',
      })
    })
  })

  it('disables submit button when fields are empty', () => {
    renderWithChakra(<PasswordChangeForm />)
    const submitButton = screen.getByRole('button', { name: /change password/i })
    expect(submitButton).toBeDisabled()
  })
})
