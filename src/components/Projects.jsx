import { Box, Typography } from "@mui/material"

import { ProjectCard } from "./ProjectCard"

const projects = [
	{
		name: "flask-simpleldap",
		description: "LDAP authentication extension for the Flask web framework.",
		link: "https://github.com/alexferl/Flask-SimpleLDAP",
	},
	{
		name: "tinysyslog",
		description:
			"A tiny and simple syslog server with log rotation written in Go.",
		link: "https://github.com/alexferl/tinysyslog",
	},
	{
		name: "vyper",
		description: "Python configuration with (more) fangs.",
		link: "https://github.com/alexferl/vyper",
	},
]

export function Projects() {
	return (
		<Box sx={{ mt: 4 }}>
			<Typography variant="h4" gutterBottom>
				Projects
			</Typography>
			<Box
				sx={{
					display: "flex",
					flexWrap: "wrap",
					gap: 2,
					justifyContent: "center",
				}}
			>
				{projects.map((project) => (
					<ProjectCard
						key={project.name}
						name={project.name}
						description={project.description}
						link={project.link}
					/>
				))}
			</Box>
		</Box>
	)
}
