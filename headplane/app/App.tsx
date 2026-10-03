import { Routes, Route, Navigate } from 'react-router-dom'
import { Box } from '@chakra-ui/react'
import Layout from './components/Layout'
import LoginPage from './routes/auth/login/page'
import SettingsProfilePage from './routes/settings/profile/page'
import { AuthProvider } from './contexts/AuthContext'

function App() {
  return (
    <AuthProvider>
      <Box minH="100vh">
        <Routes>
          <Route path="/auth/login" element={<LoginPage />} />
          <Route path="/" element={<Layout />}>
            <Route index element={<Navigate to="/settings/profile" replace />} />
            <Route path="settings/profile" element={<SettingsProfilePage />} />
          </Route>
        </Routes>
      </Box>
    </AuthProvider>
  )
}

export default App
