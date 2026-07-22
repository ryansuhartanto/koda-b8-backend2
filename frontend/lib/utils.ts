import type { ClassValue } from "clsx";
import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]): string {
	return twMerge(clsx(inputs));
}

export function getInitials(name = ""): string {
	const names = name.split(" ");
	return names
		.map((n) => n[0])
		.join("")
		.toUpperCase()
		.slice(0, 2);
}
