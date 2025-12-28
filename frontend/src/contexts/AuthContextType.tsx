import { createContext } from "react"
import type { LoginPayload, Token, User } from "../types/user"

interface AuthContextType {
    user: User | null;
    token: Token | null
    login: (userData: LoginPayload) => void
    logout: () => void
    isPending: boolean
    isLoginModalOpen: boolean
    setIsLoginModalOpen: React.Dispatch<React.SetStateAction<boolean>>
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)
