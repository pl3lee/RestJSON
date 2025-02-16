import './App.css'
import { useQuery } from '@tanstack/react-query'
import { fetchHello } from './lib/api'

function App() {

    const { data: helloData, isLoading: helloLoading } = useQuery({ queryKey: ['hello'], queryFn: fetchHello })

    return (
        <div className="text-red-500">
            {helloLoading ? "Loading..." : helloData}
        </div>
    )
}

export default App
