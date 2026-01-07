import { useState, useContext, type ReactNode } from 'react'
import { en } from '../locales/en'
import { th } from '../locales/th'
import { LanguageContext } from './LanguageContextType'


export const LanguageProvider = ({ children }: { children: ReactNode }) => {
  const [language, setLanguage] = useState<'en' | 'th'>('en')

  const t = language === 'en' ? en : th

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  )
}