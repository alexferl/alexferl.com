import EmailIcon from "@mui/icons-material/Email"
import GitHubIcon from "@mui/icons-material/GitHub"
import LinkedInIcon from "@mui/icons-material/LinkedIn"
import { AppBar, Box, IconButton, Toolbar, Typography } from "@mui/material"

export function Header() {
	return (
		<Box sx={{ flexGrow: 1 }}>
			<AppBar
				position="static"
				sx={{
					background: "inherit",
					borderBottom: "1px solid #333",
				}}
			>
				<Toolbar>
					<Box
						sx={{
							width: "100%",
							maxWidth: 1200,
							mx: "auto",
							display: "flex",
							justifyContent: "space-between",
							alignItems: "center",
						}}
					>
						<Typography variant="h6" component="div">
							Alexandre Ferland
						</Typography>

						<Box>
							<IconButton
								size="large"
								color="inherit"
								href="mailto:me@alexferl.com"
								target="_blank"
							>
								<EmailIcon />
							</IconButton>
							<IconButton
								size="large"
								color="inherit"
								href="https://github.com/alexferl"
								target="_blank"
							>
								<GitHubIcon />
							</IconButton>
							<IconButton
								size="large"
								color="inherit"
								edge="end"
								href="https://www.linkedin.com/in/alexferl/"
								target="_blank"
							>
								<LinkedInIcon />
							</IconButton>
						</Box>
					</Box>
				</Toolbar>
			</AppBar>
		</Box>
	)
}
