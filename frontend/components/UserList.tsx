import type { JSX } from "react";
import { useEffect, useState } from "react";

import { Avatar } from "#/components/ui/avater";
import { assetUrl, listUsers } from "#/lib/api";
import type { Identified, User } from "#/lib/api";

export function UserList(): JSX.Element {
	const [users, setUsers] = useState<Array<User & Identified> | undefined>(
		undefined,
	);
	const [error, setError] = useState<string | undefined>(undefined);

	useEffect(() => {
		async function load(): Promise<void> {
			try {
				setUsers(await listUsers());
			} catch (error) {
				setError(
					error instanceof Error ? error.message : "Failed to load users",
				);
			}
		}

		void load();
	}, []);

	if (error) {
		return <p className="mt-4 text-sm text-red-600">{error}</p>;
	}

	if (!users) {
		return <p className="mt-4 text-sm text-zinc-500">Loading users…</p>;
	}

	return (
		<ul className="mt-4 flex flex-col gap-2">
			{users.map((user) => (
				<li
					key={user.id}
					className="flex gap-4 items-center rounded-xl border border-zinc-200 bg-white px-4 py-3"
				>
					<Avatar
						name={user.name}
						src={assetUrl(user.picture_url)}
					/>
					<div>
						<p className="text-sm font-medium text-zinc-900">{user.name}</p>
						<p className="text-xs text-zinc-500">{user.email}</p>
					</div>
				</li>
			))}
		</ul>
	);
}
