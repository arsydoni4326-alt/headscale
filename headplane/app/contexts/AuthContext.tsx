import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { useNavigate } from 'react-router-dom'

interface AuthContextType {
  isAuthenticated: boolean
  sessionToken: string | null
  login: (token: string) => void
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [sessionToken, setSessionToken] = useState<string | null>(null)
  const navigate = useNavigate()

  useEffect(() => {
    // Check for existing session token in localStorage
    const token = localStorage.getItem('sessionToken')
    if (token) {
      setSessionToken(token)
    }
  }, [])

  const login = (token: string) => {
    localStorage.setItem('sessionToken', token)
    setSessionToken(token)
    navigate('/settings/profile')
  }

  const logout = () => {
    localStorage.removeItem('sessionToken')
    setSessionToken(null)
    navigate('/auth/login')
  }

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated: !!sessionToken,
        sessionToken,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
