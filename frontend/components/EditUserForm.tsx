import { Field } from "@base-ui/react/field";
import type { ChangeEvent, SubmitEvent, JSX } from "react";
import { useState } from "react";

import { Avatar } from "#/components/ui/avater";
import { Button } from "#/components/ui/button";
import { deleteUser, editUser, updateUserPicture } from "#/lib/api";
import type { Identified, User } from "#/lib/api";

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

export function EditUserForm({
	user,
	onSave,
	onDelete,
	onCancel,
}: {
	user: User & Identified;
	onSave: (user: User & Identified) => void;
	onDelete: () => void;
	onCancel: () => void;
}): JSX.Element {
	const [error, setError] = useState<string | undefined>(undefined);
	const [pending, setPending] = useState(false);
	const [previewUrl, setPreviewUrl] = useState<string | undefined>(undefined);

	function handlePictureChange(event: ChangeEvent<HTMLInputElement>): void {
		const file = event.currentTarget.files?.[0];
		if (!file) {
			return;
		}

		setPreviewUrl((prev) => {
			if (prev) {
				URL.revokeObjectURL(prev);
			}
			return URL.createObjectURL(file);
		});
	}

	async function handleSubmit(
		event: SubmitEvent<HTMLFormElement>,
	): Promise<void> {
		event.preventDefault();
		setError(undefined);
		setPending(true);

		const data = new FormData(event.currentTarget);
		const picture = data.get("picture");

		try {
			if (picture instanceof File && picture.size > 0) {
				await updateUserPicture(user.id, picture);
			}

			const updated = await editUser(user.id, {
				name: fieldValue(data, "name"),
				email: fieldValue(data, "email"),
				password: fieldValue(data, "password"),
			});

			onSave(updated);
		} catch (error) {
			setError(error instanceof Error ? error.message : "Something went wrong");
		} finally {
			setPending(false);
		}
	}

	async function handleDelete(): Promise<void> {
		if (!globalThis.confirm("Delete your account? This cannot be undone.")) {
			return;
		}

		setError(undefined);
		setPending(true);

		try {
			await deleteUser(user.id);
			onDelete();
		} catch (error) {
			setError(error instanceof Error ? error.message : "Something went wrong");
		} finally {
			setPending(false);
		}
	}

	return (
		<form
			onSubmit={(event) => {
				void handleSubmit(event);
			}}
			className="mt-4 flex flex-col gap-4 rounded-2xl border border-zinc-200 bg-white p-6"
		>
			<Field.Root
				name="picture"
				className={cls.field}
			>
				<Field.Label className="group flex w-fit cursor-pointer items-center gap-4">
					<Avatar
						name={user.name}
						src={previewUrl}
						className="size-12 text-base ring-0 ring-accent transition-all group-hover:ring-1"
					/>
					<span className="text-sm font-medium text-accent">
						{previewUrl ? "Change photo" : "Upload photo"}
					</span>
				</Field.Label>
				<Field.Control
					type="file"
					accept="image/*"
					className="sr-only"
					onChange={handlePictureChange}
				/>
				<Field.Error className={cls.error} />
			</Field.Root>

			<Field.Root
				name="name"
				className={cls.field}
			>
				<Field.Label className={cls.label}>Name</Field.Label>
				<Field.Control
					required
					defaultValue={user.name}
					className={cls.control}
				/>
				<Field.Error className={cls.error} />
			</Field.Root>

			<Field.Root
				name="email"
				className={cls.field}
			>
				<Field.Label className={cls.label}>Email</Field.Label>
				<Field.Control
					type="email"
					required
					defaultValue={user.email}
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

			<div className="flex items-center justify-between gap-2">
				<Button
					type="button"
					variant="ghost"
					disabled={pending}
					onClick={onCancel}
				>
					Cancel
				</Button>
				<div className="flex gap-2">
					<Button
						type="button"
						variant="ghost"
						className="text-red-600 hover:bg-red-600/10"
						disabled={pending}
						onClick={() => {
							void handleDelete();
						}}
					>
						Delete account
					</Button>
					<Button
						type="submit"
						disabled={pending}
					>
						{pending ? "…" : "Save"}
					</Button>
				</div>
			</div>
		</form>
	);
}
