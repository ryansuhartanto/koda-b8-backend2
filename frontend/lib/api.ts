export type Identified = {
	id: number;
};

export type User = {
	name: string;
	picture_url?: string;
} & Credentials;

type Credentials = {
	email: string;
	password: string;
};

const URL = "http://localhost:8080";
const TOKEN = "hello";

export function assetUrl(url?: string): string {
	return `${URL}/${url}`;
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
	const res = await fetch(`${URL}${path}`, init);

	if (!res.ok) {
		const message = (await res.json()) as string;
		throw new Error(message);
	}

	return res.json() as Promise<T>;
}

const encoder = new TextEncoder();

export async function register(data: User): Promise<User & Identified> {
	data.password = encoder.encode(data.password).toBase64();
	return request("/auth/register", {
		method: "POST",
		body: new URLSearchParams(data),
	});
}

export async function login(data: Credentials): Promise<User & Identified> {
	data.password = encoder.encode(data.password).toBase64();
	return request("/auth/login", {
		method: "POST",
		body: new URLSearchParams(data),
	});
}

export async function listUsers(): Promise<Array<User & Identified>> {
	return request("/users/", {
		method: "GET",
		headers: {
			Authorization: TOKEN,
		},
	});
}

export async function editUser(
	id: number,
	data: User,
): Promise<User & Identified> {
	data.password = encoder.encode(data.password).toBase64();
	return request(`/users/${id}`, {
		method: "PATCH",
		headers: {
			Authorization: TOKEN,
		},
		body: new URLSearchParams(data),
	});
}

export async function updateUserPicture(
	id: number,
	file?: Blob,
): Promise<void> {
	const res = await fetch(`${URL}/users/${id}/picture`, {
		method: "PUT",
		headers: {
			Authorization: TOKEN,
			...(file && { "Content-Type": file.type }),
		},
		body: file,
	});

	if (!res.ok) {
		const message = (await res.json()) as string;
		throw new Error(message);
	}
}

export async function deleteUser(id: number): Promise<void> {
	await fetch(`/users/${id}`, {
		method: "DELETE",
		headers: { Authorization: TOKEN },
	});
}
