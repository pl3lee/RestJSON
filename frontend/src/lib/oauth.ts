export function generateGoogleUrl() {
    const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID;
    const redirectUri = encodeURIComponent(`${import.meta.env.VITE_BASE_URL as string}/auth`);
    const responseType = "token";
    const scope = encodeURIComponent("profile email");
    const params = new URLSearchParams({
        client_id: clientId,
        redirect_uri: redirectUri,
        response_type: responseType,
        scope: scope,
    })
    const url = `https://accounts.google.com/o/oauth2/v2/auth?${params.toString()}`;
    console.log(url);
    return url;
}

