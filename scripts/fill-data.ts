// Seeds the backend with fake users + profile pictures from randomuser.me.
// Usage: bun scripts/fill-data.ts [count=10]

// oxlint-disable typescript/no-unsafe-return no-console
const API = "http://localhost:8080";
const count = Number(Bun.argv[2] ?? 10);

type RandomUser = {
	name: { first: string; last: string };
	email: string;
	picture: { large: string };
};

const encoder = new TextEncoder();

async function register(
	name: string,
	email: string,
): Promise<{ id: number; jwt: string }> {
	const password = encoder.encode("Password123!").toBase64();
	const res = await fetch(`${API}/auth/register`, {
		method: "POST",
		body: new URLSearchParams({ name, email, password }),
	});

	if (!res.ok) {
		throw new Error(`register ${email}: ${await res.text()}`);
	}

	return res.json();
}

async function uploadPicture(
	id: number,
	jwt: string,
	url: string,
): Promise<void> {
	const image = await fetch(url);
	const res = await fetch(`${API}/users/${id}/picture`, {
		method: "PUT",
		headers: {
			"Authorization": `Bearer ${jwt}`,
			"Content-Type": image.headers.get("Content-Type") ?? "image/jpeg",
		},
		body: await image.arrayBuffer(),
	});

	if (!res.ok) {
		throw new Error(`picture for user ${id}: ${await res.text()}`);
	}
}

async function main(): Promise<void> {
	const res = await fetch(
		`https://randomuser.me/api/?results=${count}&inc=name,email,picture`,
	);
	const { results } = (await res.json()) as { results: RandomUser[] };

	for (const person of results) {
		const name = `${person.name.first} ${person.name.last}`;

		try {
			const { id, jwt } = await register(name, person.email);
			await uploadPicture(id, jwt, person.picture.large);
			console.log(`created ${name} <${person.email}>`);
		} catch (error) {
			console.error(
				error instanceof Error ? error.message : `failed for ${name}`,
			);
		}
	}
}

await main();
