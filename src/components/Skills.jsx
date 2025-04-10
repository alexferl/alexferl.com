import { Box, Chip, Typography } from "@mui/material"

const skills = {
	"Programming Languages": ["Go", "Python", "JavaScript", "Shell Scripting"],
	"Backend Development": ["REST APIs", "Microservices", "gRPC"],
	"DevOps Tools": [
		"Docker",
		"Kubernetes",
		"Pulumi",
		"CI/CD Pipelines",
		"Prometheus/Grafana",
	],
	"Cloud Platforms": ["GCP", "Linode", "OCI"],
	"Databases & Search Technologies": ["MongoDB", "Redis", "Elasticsearch"],
	"Security Fundamentals": [
		"Access Management",
		"Encryption",
		"Least Privilege",
		"Secure Software Development",
	],
	"Other Skills": ["Linux", "Networking Fundamentals", "Version Control (Git)"],
}

export function Skills() {
	return (
		<Box sx={{ mt: 4 }}>
			<Typography variant="h4" gutterBottom>
				Skills
			</Typography>
			<Box sx={{ display: "flex", flexDirection: "column", gap: 3 }}>
				{Object.entries(skills).map(([category, items]) => (
					<Box key={category}>
						<Typography variant="h5" gutterBottom>
							{category}
						</Typography>
						<Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
							{items.map((skill) => (
								<Chip key={skill} label={skill} variant="outlined" />
							))}
						</Box>
					</Box>
				))}
			</Box>
		</Box>
	)
}
