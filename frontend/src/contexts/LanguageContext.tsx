import { createContext, useState, useContext, type ReactNode } from 'react'
import { en } from '../locales/en'
import { th } from '../locales/th'

type Dictionary = typeof en

interface LanguageContextType {
  language: 'en' | 'th'
  setLanguage: React.Dispatch<React.SetStateAction<"en" | "th">>
  t: Dictionary
}

export const LanguageContext = createContext<LanguageContextType | undefined>(undefined)

export const LanguageProvider = ({ children }: { children: ReactNode }) => {
  const [language, setLanguage] = useState<'en' | 'th'>('en')

  const t = language === 'en' ? en : th

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  )
}

export const useTranslation = () => {
    const context = useContext(LanguageContext)
    if (!context) throw new Error('useTranslation must be used within LanguageProvider')
    return context
}