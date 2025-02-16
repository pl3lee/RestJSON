export const fetchHello = async () => {
    try {
        const res = await fetch("http://localhost:3000", {
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
