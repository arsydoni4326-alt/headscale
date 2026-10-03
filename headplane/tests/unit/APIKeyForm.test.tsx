import { describe, it, expect, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ChakraProvider } from '@chakra-ui/react'
import APIKeyForm from '@/routes/settings/profile/components/APIKeyForm'
import * as api from '@/server/headscale/api'

// Mock the API
vi.mock('@/server/headscale/api', () => ({
  api: {
    updateSettings: vi.fn(),
  },
}))

const renderWithChakra = (ui: React.ReactElement) => {
  return render(<ChakraProvider>{ui}</ChakraProvider>)
}

describe('APIKeyForm', () => {
  it('renders API key input field', () => {
    renderWithChakra(<APIKeyForm />)
    expect(screen.getByLabelText(/headscale api key/i)).toBeInTheDocument()
  })

  it('validates API key format', async () => {
    const user = userEvent.setup()
    renderWithChakra(<APIKeyForm />)

    const input = screen.getByLabelText(/headscale api key/i)
    await user.type(input, 'invalid')
    await user.tab()

    await waitFor(() => {
      expect(screen.getByText(/api key must start with "hs_"/i)).toBeInTheDocument()
    })
  })

  it('shows and hides API key', async () => {
    const user = userEvent.setup()
    renderWithChakra(<APIKeyForm />)

    const input = screen.getByLabelText(/headscale api key/i) as HTMLInputElement
    const toggleButton = screen.getByLabelText(/show api key/i)

    expect(input.type).toBe('password')

    await user.click(toggleButton)
    expect(input.type).toBe('text')

    await user.click(toggleButton)
    expect(input.type).toBe('password')
  })

  it('submits valid API key', async () => {
    const user = userEvent.setup()
    const mockUpdateSettings = vi.spyOn(api.api, 'updateSettings')
    mockUpdateSettings.mockResolvedValue({ success: true })

    renderWithChakra(<APIKeyForm />)

    const input = screen.getByLabelText(/headscale api key/i)
    const saveButton = screen.getByRole('button', { name: /save api key/i })

    await user.type(input, 'hs_valid_key_12345')
    await user.click(saveButton)

    await waitFor(() => {
      expect(mockUpdateSettings).toHaveBeenCalledWith({ apiKey: 'hs_valid_key_12345' })
    })
  })

  it('disables save button when input is empty', () => {
    renderWithChakra(<APIKeyForm />)
    const saveButton = screen.getByRole('button', { name: /save api key/i })
    expect(saveButton).toBeDisabled()
  })
})
