import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import Home from './pages/Home.tsx'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter, Route, Routes } from 'react-router'
import { Auth } from './pages/Auth.tsx'
import { GoogleOAuthProvider } from '@react-oauth/google'
import { App } from './pages/App.tsx'
import { ThemeProvider } from './components/theme-provider.tsx'

const queryClient = new QueryClient()

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <QueryClientProvider client={queryClient}>
            <GoogleOAuthProvider clientId={import.meta.env.VITE_GOOGLE_CLIENT_ID}>
                <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
                    <BrowserRouter>
                        <Routes>
                            <Route path="/" element={<Home />} />
                            <Route path="/auth" element={<Auth />} />
                            <Route path="/app" element={<App />} />
                        </Routes>
                    </BrowserRouter>
                </ThemeProvider>
            </GoogleOAuthProvider>
        </QueryClientProvider>
    </StrictMode>,
)
