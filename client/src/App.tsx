import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { fetchHello } from './lib/api'

function App() {
    const queryClient = useQueryClient()


    const { data: helloData, isLoading: helloLoading } = useQuery({ queryKey: ['hello'], queryFn: fetchHello })

    return (
        <div className="text-red-500">
            {helloLoading ? "Loading..." : helloData}
        </div>
    )
}

export default App
