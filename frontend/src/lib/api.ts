export const fetchHelloWebApi = async () => {
	try {
		const res = await fetch(`${import.meta.env.VITE_WEB_API_URL}`, {
			credentials: "include",
		});
		if (!res.ok) {
			throw new Error("Failed to fetch hello from web");
		}
		return res.json();
	} catch (e) {
		console.error(e);
	}
};
export type User = {
	email: string;
	name: string;
	id: string;
};

export async function fetchMe(): Promise<User | undefined> {
	try {
		const res = await fetch(`${import.meta.env.VITE_WEB_API_URL}/me`, {
			credentials: "include",
		});

		if (!res.ok) {
			throw new Error("Failed to fetch user information");
		}
		const user: User = await res.json();
		return user;
	} catch (e) {
		console.error(e);
		return undefined;
	}
}

export async function login() {
	window.location.href = `${import.meta.env.VITE_WEB_API_URL}/auth/google/login`;
}

export async function logout() {
	try {
		const res = await fetch(`${import.meta.env.VITE_WEB_API_URL}/logout`, {
			method: "PUT",
			credentials: "include",
		});

		if (!res.ok) {
			throw new Error("Failed to logout");
		}
		return;
	} catch (e) {
		console.error(e);
	}
}

export async function createJSON() {
	try {
		const res = await fetch(`${import.meta.env.VITE_WEB_API_URL}/jsonfile`, {
			method: "POST",
			credentials: "include",
			body: JSON.stringify({
				hello: "world",
				hello2: "world2",
			}),
		});

		if (!res.ok) {
			throw new Error("Failed to create JSON");
		}
		return res.json();
	} catch (e) {
		console.error(e);
	}
}

export async function getJSON(fileId: string) {
	try {
		const res = await fetch(
			`${import.meta.env.VITE_WEB_API_URL}/jsonfile/${fileId}`,
			{
				method: "GET",
				credentials: "include",
			},
		);

		if (!res.ok) {
			throw new Error("Failed to get JSON");
		}
		return res.json();
	} catch (e) {
		console.error(e);
	}
}

export async function getJSONFiles() {
	try {
		const res = await fetch(`${import.meta.env.VITE_WEB_API_URL}/jsonfiles`, {
			method: "GET",
			credentials: "include",
		});

		if (!res.ok) {
			throw new Error("Failed to get JSON files");
		}
		return res.json();
	} catch (e) {
		console.error(e);
	}
}
