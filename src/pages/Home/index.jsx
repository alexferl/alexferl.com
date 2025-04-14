import { Box } from "@mui/material"

import { AboutMe } from "@/src/components/AboutMe.jsx"
import { Projects } from "@/src/components/Projects.jsx"
import { Skills } from "@/src/components/Skills.jsx"

export function Home() {
	return (
		<Box sx={{ display: "flex", justifyContent: "center" }}>
			<Box sx={{ width: "100%" }}>
				<AboutMe />
				<Skills />
				<Projects />
			</Box>
		</Box>
	)
}
