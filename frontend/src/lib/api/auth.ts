import type { User } from "@/lib/types";

export async function fetchMe(): Promise<User> {
	const res = await fetch(`${import.meta.env.VITE_API_URL}/me`, {
		credentials: "include",
	});

	if (!res.ok) {
		const errorData = await res.json();
		const error = errorData.error;
		throw new Error(error);
	}

	const user: User = await res.json();
	return user;
}

export function login(): void {
	window.location.href = `${import.meta.env.VITE_API_URL}/auth/google/login`;
}

export async function logout(): Promise<void> {
	const res = await fetch(`${import.meta.env.VITE_API_URL}/logout`, {
		method: "PUT",
		credentials: "include",
	});

	if (!res.ok) {
		const errorData = await res.json();
		const error = errorData.error;
		throw new Error(error);
	}

	return;
}
