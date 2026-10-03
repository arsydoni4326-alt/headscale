import { describe, it, expect, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ChakraProvider } from '@chakra-ui/react'
import ThemeSelector from '@/routes/settings/profile/components/ThemeSelector'
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

describe('ThemeSelector', () => {
  it('renders light and dark theme buttons', () => {
    renderWithChakra(<ThemeSelector />)
    expect(screen.getByRole('button', { name: /light/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /dark/i })).toBeInTheDocument()
  })

  it('switches theme on button click', async () => {
    const user = userEvent.setup()
    const mockUpdateSettings = vi.spyOn(api.api, 'updateSettings')
    mockUpdateSettings.mockResolvedValue({ success: true })

    renderWithChakra(<ThemeSelector />)

    const darkButton = screen.getByRole('button', { name: /dark/i })
    await user.click(darkButton)

    await waitFor(() => {
      expect(mockUpdateSettings).toHaveBeenCalledWith({ theme: 'dark' })
    })
  })
})
