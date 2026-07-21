import type { JSX } from "react";
import { useState } from "react";

import { AuthForm } from "#/components/AuthForm";
import { Button } from "#/components/ui/button";
import { UserList } from "#/components/UserList";
import type { User } from "#/lib/api";

// oxlint-disable-next-line import/no-unassigned-import
import "#/index.css";

const STORAGE_KEY = "user";

function readStoredUser(): User | undefined {
	const raw = localStorage.getItem(STORAGE_KEY);
	return raw ? (JSON.parse(raw) as User) : undefined;
}

export function App(): JSX.Element {
	const [user, setUser] = useState<User | undefined>(readStoredUser);

	function handleAuth(nextUser: User) {
		localStorage.setItem(STORAGE_KEY, JSON.stringify(nextUser));
		setUser(nextUser);
	}

	function handleLogout() {
		localStorage.removeItem(STORAGE_KEY);
		setUser(undefined);
	}

	if (!user) {
		return <AuthForm onAuth={handleAuth} />;
	}

	return (
		<div className="mx-auto mt-16 w-full max-w-sm">
			<div className="flex items-center justify-between">
				<div>
					<p className="text-sm font-medium text-zinc-900">{user.name}</p>
					<p className="text-xs text-zinc-500">{user.email}</p>
				</div>
				<Button
					variant="ghost"
					onClick={handleLogout}
				>
					Log out
				</Button>
			</div>
			<UserList />
		</div>
	);
}
