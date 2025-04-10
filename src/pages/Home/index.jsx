import { Box } from "@mui/material"

import { AboutMe } from "../../components/AboutMe.jsx"
import { Projects } from "../../components/Projects.jsx"
import { Skills } from "../../components/Skills.jsx"

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
