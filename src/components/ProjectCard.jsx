import { Card, CardActionArea, CardContent, Typography } from "@mui/material"

export function ProjectCard({ name, description, link }) {
	return (
		<Card
			sx={{
				width: 280,
				flexShrink: 0,
			}}
		>
			<CardActionArea
				component="a"
				href={link}
				target="_blank"
				rel="noopener noreferrer"
			>
				<CardContent
					sx={{
						textAlign: "center",
					}}
				>
					<Typography variant="h6">{name}</Typography>
					<Typography variant="body2" color="text.secondary">
						{description}
					</Typography>
				</CardContent>
			</CardActionArea>
		</Card>
	)
}
