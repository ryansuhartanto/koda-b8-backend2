import oxfmt from "@kekkon-nexus/config/oxfmt";
import oxlint from "@kekkon-nexus/config/oxlint";
import { defineConfig } from "vite-plus";

export default defineConfig({
	staged: {
		"*": "vp check --fix --no-error-on-unmatched-pattern",
		"go.mod,go.sum,*.go": "go vet -fix",
	},
	fmt: {
		...oxfmt,
	},
	lint: {
		extends: [oxlint],
		jsPlugins: [{ name: "vite-plus", specifier: "vite-plus/oxlint-plugin" }],

		options: {
			typeAware: true,
			typeCheck: true,
		},
		rules: { "vite-plus/prefer-vite-plus-imports": "error" },
	},
});
