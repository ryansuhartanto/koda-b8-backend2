export type User = {
	id: number;
	name: string;
	email: string;
};

type Credentials = {
	email: string;
	password: string;
};

type Registration = {
	name: string;
} & Credentials;

const URL = "http://localhost:8080/";
const TOKEN = "hello";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
	const res = await fetch(`${URL}users${path}`, {
		...init,
		headers: {
			"Content-Type": "application/json",
			...(init?.headers as Record<string, string>),
		},
	});

	if (!res.ok) {
		const message = (await res.json()) as string;
		throw new Error(message);
	}

	return res.json() as Promise<T>;
}

export async function register(data: Registration): Promise<User> {
	return request("/register", {
		method: "POST",
		body: JSON.stringify(data),
	});
}

export async function login(data: Credentials): Promise<User> {
	return request("/login", {
		method: "POST",
		body: JSON.stringify(data),
	});
}

export async function listUsers(): Promise<User[]> {
	return request("/", {
		method: "GET",
		headers: { Authorization: TOKEN },
	});
}

export async function updateUser(
	id: number,
	data: Registration,
): Promise<User> {
	return request(`/${id}`, {
		method: "PATCH",
		headers: { Authorization: TOKEN },
		body: JSON.stringify(data),
	});
}

export async function deleteUser(id: number): Promise<void> {
	await fetch(`users/${id}`, {
		method: "DELETE",
		headers: { Authorization: TOKEN },
	});
}
