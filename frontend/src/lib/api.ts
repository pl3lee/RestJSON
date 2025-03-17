export type User = {
	id: string;
	email: string;
	name: string;
};

export type FileMetadata = {
	id: string;
	userId: string;
	fileName: string;
	modifiedAt: string;
};

export type ApiKey = {
	apiKey: string;
};

export type ApiKeyMetadata = {
	hash: string;
	name: string;
	createdAt: string;
	lastUsedAt: string;
};

export type Route = {
	method: string;
	url: string;
	description: string;
};

export async function fetchPublic() {
	try {
		const res = await fetch(
			`${import.meta.env.VITE_API_URL}/public/5ea0eb37-32fd-4f43-b8e3-c27eaba15e27/posts/1`,
			{
				headers: {
					Authorization:
						"Bearer 5e234fee8a6a491a97c47d1273e7c44df5e71c04e0223e6fbfed8114c910f835",
				},
			},
		);

		if (!res.ok) {
			throw new Error("failed to fetch from public");
		}
		const data = await res.json();
		return data;
	} catch (e) {
		console.error(e);
		return undefined;
	}
}

export async function fetchMe(): Promise<User | undefined> {
	try {
		const res = await fetch(`${import.meta.env.VITE_API_URL}/me`, {
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

export function login(): void {
	window.location.href = `${import.meta.env.VITE_API_URL}/auth/google/login`;
}

export async function logout(): Promise<void> {
	try {
		const res = await fetch(`${import.meta.env.VITE_API_URL}/logout`, {
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

export async function createJSONFile(
	fileName: string,
): Promise<FileMetadata | undefined> {
	try {
		const res = await fetch(`${import.meta.env.VITE_API_URL}/jsonfiles`, {
			method: "POST",
			credentials: "include",
			body: JSON.stringify({
				fileName,
			}),
		});

		if (!res.ok) {
			throw new Error("Failed to create JSON");
		}
		const data: FileMetadata = await res.json();
		return data;
	} catch (e) {
		console.error(e);
		return undefined;
	}
}

export async function getJSONFile<T>(fileId: string): Promise<T | undefined> {
	try {
		const res = await fetch(
			`${import.meta.env.VITE_API_URL}/jsonfiles/${fileId}`,
			{
				method: "GET",
				credentials: "include",
			},
		);

		if (!res.ok) {
			throw new Error("Failed to get JSON");
		}
		const data: T = await res.json();
		return data;
	} catch (e) {
		console.error(e);
		return undefined;
	}
}

export async function getJSONMetadata(
	fileId: string,
): Promise<FileMetadata | undefined> {
	try {
		const res = await fetch(
			`${import.meta.env.VITE_API_URL}/jsonfiles/${fileId}/metadata`,
			{
				method: "GET",
				credentials: "include",
			},
		);

		if (!res.ok) {
			throw new Error("Failed to get JSON metadata");
		}
		const data: FileMetadata = await res.json();
		return data;
	} catch (e) {
		console.error(e);
		return undefined;
	}
}

export async function getAllJSONMetadata(): Promise<
	FileMetadata[] | undefined
> {
	try {
		const res = await fetch(`${import.meta.env.VITE_API_URL}/jsonfiles`, {
			method: "GET",
			credentials: "include",
		});

		if (!res.ok) {
			throw new Error("Failed to get all JSON files");
		}
		const data: FileMetadata[] = await res.json();
		return data;
	} catch (e) {
		console.error(e);
		return undefined;
	}
}

export async function renameJSONFile({
	name,
	fileId,
}: {
	name: string;
	fileId: string;
}): Promise<FileMetadata | undefined> {
	try {
		const res = await fetch(
			`${import.meta.env.VITE_API_URL}/jsonfiles/${fileId}`,
			{
				method: "PATCH",
				credentials: "include",
				body: JSON.stringify({
					fileName: name,
				}),
			},
		);

		if (!res.ok) {
			throw new Error("Failed to rename JSON");
		}
		const data: FileMetadata = await res.json();
		return data;
	} catch (e) {
		console.error(e);
		return undefined;
	}
}

export async function updateJSONFile<T>({
	fileId,
	contents,
}: {
	fileId: string;
	contents: T;
}): Promise<T | undefined> {
	try {
		const res = await fetch(
			`${import.meta.env.VITE_API_URL}/jsonfiles/${fileId}`,
			{
				method: "PUT",
				credentials: "include",
				body: JSON.stringify(contents),
			},
		);

		if (!res.ok) {
			throw new Error("Failed to get JSON");
		}
		const data: T = await res.json();
		return data;
	} catch (e) {
		console.error(e);
		return undefined;
	}
}

export async function deleteJSONFile(fileId: string): Promise<void> {
	try {
		const res = await fetch(
			`${import.meta.env.VITE_API_URL}/jsonfiles/${fileId}`,
			{
				method: "DELETE",
				credentials: "include",
			},
		);

		if (!res.ok) {
			throw new Error("Failed to delete JSON");
		}
		return;
	} catch (e) {
		console.error(e);
	}
}

export async function createApiKey(name: string): Promise<ApiKey | undefined> {
	try {
		const res = await fetch(`${import.meta.env.VITE_API_URL}/apikeys`, {
			method: "POST",
			credentials: "include",
			body: JSON.stringify({
				name,
			}),
		});

		if (!res.ok) {
			throw new Error("Failed to create api key");
		}
		const data: ApiKey = await res.json();
		return data;
	} catch (e) {
		console.error(e);
		return undefined;
	}
}

export async function getAllApiKeys(): Promise<ApiKeyMetadata[] | undefined> {
	try {
		const res = await fetch(`${import.meta.env.VITE_API_URL}/apikeys`, {
			method: "GET",
			credentials: "include",
		});

		if (!res.ok) {
			throw new Error("Failed to get all api keys");
		}
		const data: ApiKeyMetadata[] = await res.json();
		return data;
	} catch (e) {
		console.error(e);
		return undefined;
	}
}

export async function deleteApiKey(keyHash: string): Promise<void> {
	try {
		const res = await fetch(
			`${import.meta.env.VITE_API_URL}/apikeys/${keyHash}`,
			{
				method: "DELETE",
				credentials: "include",
			},
		);

		if (!res.ok) {
			throw new Error("Failed to delete api key");
		}
	} catch (e) {
		console.error(e);
	}
}

export async function getDynamicRoutes(fileId: string) {
	try {
		const res = await fetch(
			`${import.meta.env.VITE_API_URL}/jsonfiles/${fileId}/routes`,
			{
				method: "GET",
				credentials: "include",
			},
		);

		if (!res.ok) {
			throw new Error("Failed to delete JSON");
		}
		const data: Route[] = await res.json();
		return data;
	} catch (e) {
		console.error(e);
		return undefined;
	}
}
