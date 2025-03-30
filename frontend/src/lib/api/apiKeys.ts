import type { ApiKey, ApiKeyMetadata } from "@/lib/types";

export async function createApiKey(name: string): Promise<ApiKey> {
	const res = await fetch(`${import.meta.env.VITE_API_URL}/apikeys`, {
		method: "POST",
		credentials: "include",
		body: JSON.stringify({
			name,
		}),
	});
	if (!res.ok) {
		const errorData = await res.json();
		const error = errorData.error;
		throw new Error(error);
	}

	const data: ApiKey = await res.json();
	return data;
}

export async function getAllApiKeys(): Promise<ApiKeyMetadata[]> {
	const res = await fetch(`${import.meta.env.VITE_API_URL}/apikeys`, {
		method: "GET",
		credentials: "include",
	});

	if (!res.ok) {
		const errorData = await res.json();
		const error = errorData.error;
		throw new Error(error);
	}
	const data: ApiKeyMetadata[] = await res.json();
	return data;
}

export async function deleteApiKey(keyHash: string): Promise<void> {
	const res = await fetch(
		`${import.meta.env.VITE_API_URL}/apikeys/${keyHash}`,
		{
			method: "DELETE",
			credentials: "include",
		},
	);
	if (!res.ok) {
		const errorData = await res.json();
		const error = errorData.error;
		throw new Error(error);
	}
}
