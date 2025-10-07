import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { authApi, DEMO_TENANT_ID } from '../shared/api/auth'

interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  roles: string[]
  permissions: string[]
}

interface AuthContextType {
  user: User | null
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  isLoading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const init = async () => {
      const token = localStorage.getItem('access_token')
      const refreshToken = localStorage.getItem('refresh_token')
      
      if (!token || !refreshToken) {
        setIsLoading(false)
        return
      }
      
      try {
        const response = await authApi.me()
        if (response && response.user) {
          setUser({
            id: response.user.id,
            email: response.user.email,
            firstName: response.user.firstName,
            lastName: response.user.lastName,
            roles: response.user.roles || [],
            permissions: response.user.permissions || []
          })
        }
      } catch (e) {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        setUser(null)
      } finally {
        setIsLoading(false)
      }
    }
    void init()
  }, [])

  const login = async (email: string, password: string) => {
    try {
      const response = await authApi.login({ email, password, tenant_id: DEMO_TENANT_ID })
      localStorage.setItem('access_token', response.access_token)
      localStorage.setItem('refresh_token', response.refresh_token)
      
      // Загружаем полные данные пользователя с правами
      const userResponse = await authApi.me()
      if (userResponse && userResponse.user) {
        setUser({
          id: userResponse.user.id,
          email: userResponse.user.email,
          firstName: userResponse.user.firstName,
          lastName: userResponse.user.lastName,
          roles: userResponse.user.roles || [],
          permissions: userResponse.user.permissions || []
        })
      }
    } catch (error) {
      throw error
    }
  }

  const logout = () => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    setUser(null)
  }

  const value = {
    user,
    login,
    logout,
    isLoading,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
