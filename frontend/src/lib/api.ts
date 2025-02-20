export const fetchHelloWebApi = async () => {
    try {
        const res = await fetch(`${import.meta.env.VITE_WEB_API_URL}`, {
            credentials: "include"
        })
        if (!res.ok) {
            throw new Error("Failed to fetch hello from web")
        }
        return res.json()
    } catch (e) {
        console.log(e)
    }
}
export const fetchHelloPublicApi = async () => {
    try {
        const res = await fetch(`${import.meta.env.VITE_PUBLIC_API_URL}`)
        if (!res.ok) {
            throw new Error("Failed to fetch hello from public")
        }
        return res.json()
    } catch (e) {
        console.log(e)
    }
}
