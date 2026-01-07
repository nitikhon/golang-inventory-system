import { createContext } from 'react'
import { en } from '../locales/en'

type Dictionary = typeof en

interface LanguageContextType {
  language: 'en' | 'th'
  setLanguage: React.Dispatch<React.SetStateAction<'en' | 'th'>>
  t: Dictionary
}

export const LanguageContext = createContext<LanguageContextType | undefined>(undefined)
