import { useGoogleLogin } from "@react-oauth/google";

export function Auth() {
    const login = useGoogleLogin({
        flow: "auth-code",
        redirect_uri: `${import.meta.env.VITE_BASE_URL}/auth`,
        onSuccess: async (codeResponse) => {
            try {
                const response = await fetch(
                    `${import.meta.env.VITE_WEB_API_URL}/auth/google/callback`,
                    {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                        },
                        body: JSON.stringify({ code: codeResponse.code }),
                        credentials: "include",
                    }
                );
                if (!response.ok) {
                    throw new Error("Login failed");
                }
                const data = await response.json();
                console.log("auth data from api", data);
            } catch (error) {
                console.error("Login error:", error);
            }
        },
        onError: (error) => {
            console.error("Google login error:", error);
        },
    });
    return <div className="flex flex-col gap-2">
        <button onClick={login}>Login with google</button>
    </div>
}
