import * as path from "node:path"
import preact from "@preact/preset-vite"
import { defineConfig } from "vite"

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [preact()],
	resolve: {
		alias: {
			"@": path.resolve(__dirname, "./"),
		},
	},
})
