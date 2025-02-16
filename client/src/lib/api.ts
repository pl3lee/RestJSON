export const fetchHello = async () => {
    try {
        const res = await fetch(import.meta.env.VITE_BACKEND_URL, {
            credentials: "include"
        })
        if (!res.ok) {
            throw new Error("Failed to fetch hello")
        }
        return res.json()
    } catch (e) {
        console.log(e)
    }
}
