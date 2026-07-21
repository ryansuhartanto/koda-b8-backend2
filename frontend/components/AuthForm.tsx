import { Field } from "@base-ui/react/field";
import type { SubmitEvent, JSX } from "react";
import { useState } from "react";

import { Button } from "#/components/ui/button";
import { login, register } from "#/lib/api";
import type { User } from "#/lib/api";

const cls = {
	field: "flex flex-col gap-1.5",
	label: "text-sm font-medium text-zinc-700",
	control:
		"rounded-lg border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-accent focus:ring-1 focus:ring-accent",
	error: "text-xs text-red-600",
};

function fieldValue(data: FormData, key: string): string {
	const value = data.get(key);
	return typeof value === "string" ? value : "";
}

export function AuthForm({
	onAuth,
}: {
	onAuth: (user: User) => void;
}): JSX.Element {
	const [mode, setMode] = useState<"login" | "register">("login");
	const [error, setError] = useState<string | undefined>(undefined);
	const [pending, setPending] = useState(false);

	const idleLabel = mode === "login" ? "Sign in" : "Create account";
	const submitLabel = pending ? "…" : idleLabel;

	async function handleSubmit(
		event: SubmitEvent<HTMLFormElement>,
	): Promise<void> {
		event.preventDefault();
		setError(undefined);
		setPending(true);

		const data = new FormData(event.currentTarget);
		const email = fieldValue(data, "email");
		const password = fieldValue(data, "password");

		try {
			const user =
				mode === "login"
					? await login({ email, password })
					: await register({ email, password, name: fieldValue(data, "name") });
			onAuth(user);
		} catch (error) {
			setError(error instanceof Error ? error.message : "Something went wrong");
		} finally {
			setPending(false);
		}
	}

	return (
		<div className="mx-auto mt-24 w-full max-w-sm rounded-2xl border border-zinc-200 bg-white p-8 shadow-sm">
			<h1 className="text-lg font-semibold text-zinc-900">
				{mode === "login" ? "Sign in" : "Create an account"}
			</h1>

			<form
				onSubmit={(event) => {
					void handleSubmit(event);
				}}
				className="mt-6 flex flex-col gap-4"
			>
				{mode === "register" && (
					<Field.Root
						name="name"
						className={cls.field}
					>
						<Field.Label className={cls.label}>Name</Field.Label>
						<Field.Control
							required
							className={cls.control}
						/>
						<Field.Error className={cls.error} />
					</Field.Root>
				)}

				<Field.Root
					name="email"
					className={cls.field}
				>
					<Field.Label className={cls.label}>Email</Field.Label>
					<Field.Control
						type="email"
						required
						className={cls.control}
					/>
					<Field.Error className={cls.error} />
				</Field.Root>

				<Field.Root
					name="password"
					className={cls.field}
				>
					<Field.Label className={cls.label}>Password</Field.Label>
					<Field.Control
						type="password"
						required
						className={cls.control}
					/>
					<Field.Error className={cls.error} />
				</Field.Root>

				{error && <p className="text-sm text-red-600">{error}</p>}

				<Button
					type="submit"
					disabled={pending}
					className="mt-2 w-full"
				>
					{submitLabel}
				</Button>
			</form>

			<Button
				variant="ghost"
				className="mt-4 w-full"
				onClick={() => {
					setMode(mode === "login" ? "register" : "login");
					setError(undefined);
				}}
			>
				{mode === "login"
					? "Need an account? Register"
					: "Have an account? Sign in"}
			</Button>
		</div>
	);
}
