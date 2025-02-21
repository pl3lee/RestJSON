import { useQuery } from '@tanstack/react-query'
import { fetchHelloWebApi, fetchHelloPublicApi } from './lib/api'

function App() {

    const { data: helloWebApiData, isLoading: helloWebApiLoading } = useQuery({ queryKey: ['helloweb'], queryFn: fetchHelloWebApi })
    const { data: helloPublicApiData, isLoading: helloPublicApiLoading } = useQuery({ queryKey: ['hellopublic'], queryFn: fetchHelloPublicApi })

    return (
        <div className="text-red-500">
            {helloWebApiLoading ? "Loading web api..." : helloWebApiData}
            {helloPublicApiLoading ? "Loading public api..." : helloPublicApiData}
        </div>
    )
}

export default App
