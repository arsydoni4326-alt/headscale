import { describe, it, expect, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ChakraProvider } from '@chakra-ui/react'
import ProfileNameForm from '@/routes/settings/profile/components/ProfileNameForm'
import * as api from '@/server/headscale/api'

vi.mock('@/server/headscale/api', () => ({
  api: {
    updateSettings: vi.fn(),
    getSettings: vi.fn(),
  },
}))

const renderWithChakra = (ui: React.ReactElement) => {
  return render(<ChakraProvider>{ui}</ChakraProvider>)
}

describe('ProfileNameForm', () => {
  it('renders profile name input', () => {
    renderWithChakra(<ProfileNameForm />)
    expect(screen.getByLabelText(/display name/i)).toBeInTheDocument()
  })

  it('enables save button when name is changed', async () => {
    const user = userEvent.setup()
    renderWithChakra(<ProfileNameForm />)

    const input = screen.getByLabelText(/display name/i)
    const saveButton = screen.getByRole('button', { name: /save profile name/i })

    expect(saveButton).toBeDisabled()

    await user.type(input, 'John Doe')
    expect(saveButton).toBeEnabled()
  })

  it('submits profile name', async () => {
    const user = userEvent.setup()
    const mockUpdateSettings = vi.spyOn(api.api, 'updateSettings')
    mockUpdateSettings.mockResolvedValue({ success: true })

    renderWithChakra(<ProfileNameForm />)

    await user.type(screen.getByLabelText(/display name/i), 'Jane Smith')
    await user.click(screen.getByRole('button', { name: /save profile name/i }))

    await waitFor(() => {
      expect(mockUpdateSettings).toHaveBeenCalledWith({ profileName: 'Jane Smith' })
    })
  })
})
