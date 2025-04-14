type checkoutResponse = {
    checkoutUrl: string
}
export async function checkout(priceId: string): Promise<void> {
    const res = await fetch(
        `${import.meta.env.VITE_API_URL}/subscriptions/checkout`,
        {
            method: "POST",
            credentials: "include",
            body: JSON.stringify({
                priceId,
            }),
        },
    );

    if (!res.ok) {
        const errorData = await res.json();
        const error = errorData.error;
        throw new Error(error);
    }
    const data: checkoutResponse = await res.json()
    const checkoutUrl = data.checkoutUrl;

    window.location.href = checkoutUrl;
}

export async function paymentSuccessful(): Promise<void> {
    const res = await fetch(
        `${import.meta.env.VITE_API_URL}/subscriptions/checkout`,
        {
            method: "POST",
            credentials: "include",
        },
    );

    if (!res.ok) {
        const errorData = await res.json();
        const error = errorData.error;
        throw new Error(error);
    }
    return;
}
