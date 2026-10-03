import { test, expect } from '@playwright/test'

test.describe('Settings Page', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to login page
    await page.goto('/')
    
    // Login with mock credentials
    await page.fill('input[type="password"]', 'password123')
    await page.click('button[type="submit"]')
    
    // Wait for redirect to settings page
    await page.waitForURL('**/settings/profile')
  })

  test('displays all settings sections', async ({ page }) => {
    await expect(page.getByRole('heading', { name: 'Settings' })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Account' })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Integration' })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Preferences' })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Profile' })).toBeVisible()
  })

  test('changes password successfully', async ({ page }) => {
    // Fill in password change form
    await page.fill('input[type="password"]', 'password123')
    await page.fill('input[type="password"]', 'newpassword123')
    await page.fill('input[type="password"]', 'newpassword123')
    
    // Submit form
    await page.click('button:has-text("Change Password")')
    
    // Wait for success toast
    await expect(page.getByText('Password changed')).toBeVisible()
  })

  test('saves API key', async ({ page }) => {
    // Find API key input by label
    const apiKeyInput = page.locator('input').filter({ has: page.locator('label:has-text("Headscale API Key")') }).first()
    
    // Fill in API key
    await apiKeyInput.fill('hs_test_key_12345678')
    
    // Click save button
    await page.click('button:has-text("Save API Key")')
    
    // Wait for success toast
    await expect(page.getByText('API key saved')).toBeVisible()
  })

  test('switches theme', async ({ page }) => {
    // Click dark theme button
    await page.click('button:has-text("Dark")')
    
    // Wait for success toast
    await expect(page.getByText('Theme updated')).toBeVisible()
    
    // Verify theme changed (check for dark mode class or attribute)
    // This depends on how Chakra UI implements color mode
  })

  test('updates profile name', async ({ page }) => {
    // Find display name input
    const displayNameInput = page.locator('input[placeholder="Enter your name"]')
    
    // Fill in name
    await displayNameInput.fill('John Doe')
    
    // Click save button
    await page.click('button:has-text("Save Profile Name")')
    
    // Wait for success toast
    await expect(page.getByText('Profile name saved')).toBeVisible()
  })

  test('displays session information', async ({ page }) => {
    await expect(page.getByText('Session Information')).toBeVisible()
    await expect(page.getByText('Status')).toBeVisible()
    await expect(page.getByText('Active')).toBeVisible()
  })

  test('validates password requirements', async ({ page }) => {
    // Try to submit with short password
    const newPasswordInput = page.locator('label:has-text("New Password")').locator('..').locator('input')
    await newPasswordInput.fill('short')
    await newPasswordInput.blur()
    
    // Check for validation error
    await expect(page.getByText(/password must be at least 8 characters/i)).toBeVisible()
  })

  test('validates password confirmation', async ({ page }) => {
    // Fill different passwords
    const inputs = page.locator('input[type="password"]')
    await inputs.nth(1).fill('password123')
    await inputs.nth(2).fill('different123')
    await inputs.nth(2).blur()
    
    // Check for validation error
    await expect(page.getByText(/passwords do not match/i)).toBeVisible()
  })

  test('validates API key format', async ({ page }) => {
    // Find API key input
    const apiKeyInput = page.locator('label:has-text("Headscale API Key")').locator('..').locator('input')
    
    // Enter invalid format
    await apiKeyInput.fill('invalid_key')
    await apiKeyInput.blur()
    
    // Check for validation error
    await expect(page.getByText(/api key must start with "hs_"/i)).toBeVisible()
  })

  test('logout functionality', async ({ page }) => {
    // Click logout button
    await page.click('button:has-text("Logout")')
    
    // Verify redirected to login page
    await page.waitForURL('**/auth/login')
    await expect(page.getByText('Sign in to access settings')).toBeVisible()
  })

  test('theme toggle in header', async ({ page }) => {
    // Find theme toggle button in header
    const themeToggle = page.locator('button[aria-label="Toggle color mode"]')
    
    // Toggle theme
    await themeToggle.click()
    
    // Button should be visible after toggle
    await expect(themeToggle).toBeVisible()
  })
})
