import { Container, CssBaseline } from "@mui/material"
import { ThemeProvider, createTheme } from "@mui/material/styles"
import { render } from "preact"
import { LocationProvider, Route, Router } from "preact-iso"

import { Header } from "./components/Header.jsx"
import { Home } from "./pages/Home/index.jsx"
import { NotFound } from "./pages/_404.jsx"

const theme = createTheme({
	palette: {
		mode: "dark",
	},
})

export function App() {
	return (
		<ThemeProvider theme={theme}>
			<LocationProvider>
				<CssBaseline enableColorScheme />
				<Header />
				<Container sx={{ maxWidth: 1024, p: 3 }}>
					<main>
						<Router>
							<Route path="/" component={Home} />
							<Route default component={NotFound} />
						</Router>
					</main>
				</Container>
			</LocationProvider>
		</ThemeProvider>
	)
}

render(<App />, document.getElementById("app"))
