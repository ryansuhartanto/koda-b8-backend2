export type Identified = {
	id: number;
};

export type User = {
	name: string;
} & Credentials;

type Credentials = {
	email: string;
	password: string;
};

const URL = "http://localhost:8080";
const TOKEN = "hello";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
	const res = await fetch(`${URL}${path}`, init);

	if (!res.ok) {
		const message = (await res.json()) as string;
		throw new Error(message);
	}

	return res.json() as Promise<T>;
}

const encoder = new TextEncoder();

export async function register(data: User): Promise<User> {
	data.password = encoder.encode(data.password).toBase64();
	return request("/auth/register", {
		method: "POST",
		body: new URLSearchParams(data),
	});
}

export async function login(data: Credentials): Promise<User> {
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

export async function editUser(id: number, data: User): Promise<User> {
	data.password = encoder.encode(data.password).toBase64();
	return request(`/users/${id}`, {
		method: "PATCH",
		headers: {
			Authorization: TOKEN,
		},
		body: new URLSearchParams(data),
	});
}

export async function deleteUser(id: number): Promise<void> {
	await fetch(`/users/${id}`, {
		method: "DELETE",
		headers: { Authorization: TOKEN },
	});
}
