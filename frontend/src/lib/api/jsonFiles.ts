import type { FileMetadata, Route } from "@/lib/types";

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
