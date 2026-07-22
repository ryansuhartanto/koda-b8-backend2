import { Avatar as BaseAvatar } from "@base-ui/react";
import type { ComponentProps, JSX } from "react";

import { cn, getInitials } from "#/lib/utils";

type AvatarProps = ComponentProps<typeof BaseAvatar.Root> & {
	name?: string;
	src?: string;
};

export function Avatar({
	name,
	src,
	className,
	...props
}: AvatarProps): JSX.Element {
	return (
		<BaseAvatar.Root
			className={cn(
				"inline-flex size-9 shrink-0 items-center justify-center overflow-hidden rounded-full bg-zinc-100 text-xs font-medium text-zinc-600",
				className,
			)}
			{...props}
		>
			<BaseAvatar.Image
				src={src}
				className="size-full object-cover"
			/>
			<BaseAvatar.Fallback>{getInitials(name)}</BaseAvatar.Fallback>
		</BaseAvatar.Root>
	);
}
