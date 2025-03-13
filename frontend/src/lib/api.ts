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

export function login(): void {
	window.location.href = `${import.meta.env.VITE_WEB_API_URL}/auth/google/login`;
}

export async function logout(): Promise<void> {
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

export async function createJSONFile(
	fileName: string,
): Promise<FileMetadata | undefined> {
	try {
		const res = await fetch(`${import.meta.env.VITE_WEB_API_URL}/jsonfiles`, {
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
			`${import.meta.env.VITE_WEB_API_URL}/jsonfiles/${fileId}`,
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
			`${import.meta.env.VITE_WEB_API_URL}/jsonfiles/${fileId}/metadata`,
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
		const res = await fetch(`${import.meta.env.VITE_WEB_API_URL}/jsonfiles`, {
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
			`${import.meta.env.VITE_WEB_API_URL}/jsonfiles/${fileId}`,
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
			`${import.meta.env.VITE_WEB_API_URL}/jsonfiles/${fileId}`,
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
			`${import.meta.env.VITE_WEB_API_URL}/jsonfiles/${fileId}`,
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
