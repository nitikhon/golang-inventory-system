import { createContext, useState, useContext, type ReactNode } from 'react'
import { en } from '../locales/en'
import { th } from '../locales/th'

// 1. Define Type ของ Dictionary (เอามาจาก en)
type Dictionary = typeof en

// 2. Define Context Type
interface LanguageContextType {
  language: 'en' | 'th'
  setLanguage: React.Dispatch<React.SetStateAction<"en" | "th">>
  t: Dictionary // นี่คือตัวแปรที่เราจะเรียกใช้ใน Component!
}

// 3. Create Context
export const LanguageContext = createContext<LanguageContextType | undefined>(undefined)

// 4. Provider Logic
export const LanguageProvider = ({ children }: { children: ReactNode }) => {
  const [language, setLanguage] = useState<'en' | 'th'>('en')

  // Logic เลือก dictionary ตาม language state
  // const t = ... ? ... : ...
  const t = language === 'en' ? en : th

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  )
}

// 5. Custom Hook เพื่อให้เรียกใช้ง่ายๆ
export const useTranslation = () => {
    const context = useContext(LanguageContext)
    if (!context) throw new Error('useTranslation must be used within LanguageProvider')
    return context
}