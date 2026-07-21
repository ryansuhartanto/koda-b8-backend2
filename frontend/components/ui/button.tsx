import { Button as BaseButton } from "@base-ui/react/button";
import type { ComponentProps, JSX } from "react";

import { cn } from "#/lib/cn";

type ButtonProps = Omit<ComponentProps<typeof BaseButton>, "className"> & {
	variant?: "primary" | "ghost";
	className?: string;
};

const variantClasses = {
	primary: "bg-accent text-accent-foreground hover:opacity-90",
	ghost: "text-accent hover:bg-accent/10",
};

export function Button({
	variant = "primary",
	className,
	...props
}: ButtonProps): JSX.Element {
	return (
		<BaseButton
			className={cn(
				"rounded-full px-4 py-2 text-sm font-medium transition-colors disabled:opacity-50",
				variantClasses[variant],
				className,
			)}
			{...props}
		/>
	);
}
