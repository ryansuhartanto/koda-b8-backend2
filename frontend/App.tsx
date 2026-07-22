import type { JSX } from "react";
import { useState } from "react";

import { AuthForm } from "#/components/AuthForm";
import { EditUserForm } from "#/components/EditUserForm";
import { Avatar } from "#/components/ui/avater";
import { Button } from "#/components/ui/button";
import { UserList } from "#/components/UserList";
import { assetUrl } from "#/lib/api";
import type { Identified, User } from "#/lib/api";

// oxlint-disable-next-line import/no-unassigned-import
import "#/index.css";

const STORAGE_KEY = "user";

function readStoredUser(): (User & Identified) | undefined {
	const raw = localStorage.getItem(STORAGE_KEY);
	return raw ? (JSON.parse(raw) as User & Identified) : undefined;
}

export function App(): JSX.Element {
	const [user, setUser] = useState<(User & Identified) | undefined>(
		readStoredUser,
	);
	const [editing, setEditing] = useState(false);

	function handleAuth(nextUser: User & Identified) {
		localStorage.setItem(STORAGE_KEY, JSON.stringify(nextUser));
		setUser(nextUser);
	}

	function handleLogout() {
		localStorage.removeItem(STORAGE_KEY);
		setUser(undefined);
	}

	function handleSave(nextUser: User & Identified) {
		localStorage.setItem(STORAGE_KEY, JSON.stringify(nextUser));
		setUser(nextUser);
		setEditing(false);
	}

	if (!user) {
		return <AuthForm onAuth={handleAuth} />;
	}

	return (
		<div className="mx-auto mt-16 w-full max-w-sm">
			<div className="flex items-center justify-between">
				<div className="flex gap-4 items-center">
					<Avatar
						// name={user.name}
						src={assetUrl(user.picture_url)}
					/>
					<div>
						<p className="text-sm font-medium text-zinc-900">{user.name}</p>
						<p className="text-xs text-zinc-500">{user.email}</p>
					</div>
				</div>
				<div className="flex gap-2">
					{!editing && (
						<Button
							variant="ghost"
							onClick={() => setEditing(true)}
						>
							Edit
						</Button>
					)}
					<Button
						variant="ghost"
						onClick={handleLogout}
					>
						Log out
					</Button>
				</div>
			</div>

			{editing && (
				<EditUserForm
					user={user}
					onSave={handleSave}
					onDelete={handleLogout}
					onCancel={() => setEditing(false)}
				/>
			)}

			<UserList />
		</div>
	);
}
