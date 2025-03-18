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

export async function createJSONFile(fileName: string): Promise<FileMetadata> {
	const res = await fetch(`${import.meta.env.VITE_API_URL}/jsonfiles`, {
		method: "POST",
		credentials: "include",
		body: JSON.stringify({
			fileName,
		}),
	});

	if (!res.ok) {
		const errorData = await res.json();
		const error = errorData.error;
		throw new Error(error);
	}
	const data: FileMetadata = await res.json();
	return data;
}

export async function getJSONFile<T>(fileId: string): Promise<T> {
	const res = await fetch(
		`${import.meta.env.VITE_API_URL}/jsonfiles/${fileId}`,
		{
			method: "GET",
			credentials: "include",
		},
	);
	if (!res.ok) {
		const errorData = await res.json();
		const error = errorData.error;
		throw new Error(error);
	}

	const data: T = await res.json();
	return data;
}

export async function getJSONMetadata(fileId: string): Promise<FileMetadata> {
	const res = await fetch(
		`${import.meta.env.VITE_API_URL}/jsonfiles/${fileId}/metadata`,
		{
			method: "GET",
			credentials: "include",
		},
	);

	if (!res.ok) {
		const errorData = await res.json();
		const error = errorData.error;
		throw new Error(error);
	}

	const data: FileMetadata = await res.json();
	return data;
}

export async function getAllJSONMetadata(): Promise<FileMetadata[]> {
	const res = await fetch(`${import.meta.env.VITE_API_URL}/jsonfiles`, {
		method: "GET",
		credentials: "include",
	});
	if (!res.ok) {
		const errorData = await res.json();
		const error = errorData.error;
		throw new Error(error);
	}

	const data: FileMetadata[] = await res.json();
	return data;
}

export async function renameJSONFile({
	name,
	fileId,
}: {
	name: string;
	fileId: string;
}): Promise<FileMetadata> {
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
		const errorData = await res.json();
		const error = errorData.error;
		throw new Error(error);
	}

	const data: FileMetadata = await res.json();
	return data;
}

export async function updateJSONFile<T>({
	fileId,
	contents,
}: {
	fileId: string;
	contents: T;
}): Promise<T> {
	const res = await fetch(
		`${import.meta.env.VITE_API_URL}/jsonfiles/${fileId}`,
		{
			method: "PUT",
			credentials: "include",
			body: JSON.stringify(contents),
		},
	);
	if (!res.ok) {
		const errorData = await res.json();
		const error = errorData.error;
		throw new Error(error);
	}

	const data: T = await res.json();
	return data;
}

export async function deleteJSONFile(fileId: string): Promise<void> {
	const res = await fetch(
		`${import.meta.env.VITE_API_URL}/jsonfiles/${fileId}`,
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

	return;
}

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

export async function getDynamicRoutes(fileId: string): Promise<Route[]> {
	const res = await fetch(
		`${import.meta.env.VITE_API_URL}/jsonfiles/${fileId}/routes`,
		{
			method: "GET",
			credentials: "include",
		},
	);

	if (!res.ok) {
		const errorData = await res.json();
		const error = errorData.error;
		throw new Error(error);
	}
	const data: Route[] = await res.json();
	return data;
}
