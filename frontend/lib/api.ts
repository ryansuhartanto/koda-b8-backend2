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

const URL = "";
const TOKEN = "hello";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
	const res = await fetch(`/users${path}`, {
		headers: { "Content-Type": "application/json" },
		...init,
	});

	if (!res.ok) {
		const message = (await res.json()) as string;
		throw new Error(message);
	}

	return res.json() as Promise<T>;
}

export async function register(data: Registration): Promise<User> {
	return request(`${URL}/register`, {
		method: "POST",
		body: JSON.stringify(data),
	});
}

export async function login(data: Credentials): Promise<User> {
	return request(`${URL}/login`, {
		method: "POST",
		body: JSON.stringify(data),
	});
}

export async function listUsers(): Promise<User[]> {
	return request(`${URL}/`, {
		method: "GET",
		headers: { Authorization: TOKEN },
	});
}
