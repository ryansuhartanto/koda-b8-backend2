import type { JSX } from "react";
import { useEffect, useState } from "react";

import { Avatar } from "#/components/ui/avater";
import { Button } from "#/components/ui/button";
import { assetUrl, listUsers } from "#/lib/api";
import type { Identified, User } from "#/lib/api";

const limit = 10;

export function UserList(): JSX.Element {
	const [users, setUsers] = useState<Array<User & Identified> | undefined>(
		undefined,
	);
	const [error, setError] = useState<string | undefined>(undefined);
	const [offset, setOffset] = useState(0);
	const [search, setSearch] = useState("");
	const [query, setQuery] = useState("");

	useEffect(() => {
		const timeout = setTimeout(() => {
			setQuery(search);
			setOffset(0);
		}, 300);

		return () => {
			clearTimeout(timeout);
		};
	}, [search]);

	useEffect(() => {
		async function load(): Promise<void> {
			try {
				setUsers(await listUsers({ limit, offset, query }));
			} catch (error) {
				setError(
					error instanceof Error ? error.message : "Failed to load users",
				);
			}
		}

		void load();
	}, [offset, query]);

	return (
		<>
			<input
				type="search"
				placeholder="Search users…"
				value={search}
				onChange={(event) => {
					setSearch(event.currentTarget.value);
				}}
				className="mt-4 w-full rounded-lg border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-accent focus:ring-1 focus:ring-accent"
			/>

			{error && <p className="mt-4 text-sm text-red-600">{error}</p>}

			{!error && !users && (
				<p className="mt-4 text-sm text-zinc-500">Loading users…</p>
			)}

			{!error && users && (
				<>
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
								<div className="grow flex flex-col gap-2">
									<div>
										<p className="text-sm font-medium text-zinc-900">
											{user.name}
										</p>
										<p className="text-xs text-zinc-500">{user.email}</p>
									</div>
									<div className="grid grid-cols-[1fr_auto] text-[0.625rem] text-zinc-400">
										{user.created_at && (
											<>
												<span>Created</span>
												<span className="text-right">
													{new Date(user.created_at).toLocaleString()}
												</span>
											</>
										)}
										{user.updated_at && (
											<>
												<span>Updated</span>
												<span className="text-right">
													{new Date(user.updated_at).toLocaleString()}
												</span>
											</>
										)}
										{user.profile_updated_at && (
											<>
												<span>Profile</span>
												<span className="text-right">
													{new Date(user.profile_updated_at).toLocaleString()}
												</span>
											</>
										)}
									</div>
								</div>
							</li>
						))}
					</ul>
					<div className="mt-4 flex items-center justify-between gap-2">
						<Button
							variant="ghost"
							disabled={offset === 0}
							onClick={() => {
								setOffset((prev) => Math.max(0, prev - limit));
							}}
						>
							Previous
						</Button>
						<Button
							variant="ghost"
							disabled={users.length < limit}
							onClick={() => {
								setOffset((prev) => prev + limit);
							}}
						>
							Next
						</Button>
					</div>
				</>
			)}
		</>
	);
}
