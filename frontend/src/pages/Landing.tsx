import { ModeToggle } from "@/components/mode-toggle";
import { Typewriter } from "@/components/typewriter";
import { Button } from "@/components/ui/button";
import { ArrowRight, Code, Edit, FileJson, Globe, Key } from "lucide-react";
import { Link } from "react-router";

export function Landing() {
	const handleScroll = (id: string) => {
		const element = document.getElementById(id);
		if (element) {
			element.scrollIntoView({ behavior: "smooth" });
		}
	};
	return (
		<div className="flex min-h-screen flex-col">
			<header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
				<div className="container mx-auto px-4 flex h-16 items-center justify-between">
					<div className="flex items-center gap-2">
						<FileJson className="h-6 w-6" />
						<span className="text-xl font-bold">RestJSON</span>
					</div>

					<nav className="hidden md:flex gap-6">
						<Link
							to="#features"
							onClick={(e) => {
								e.preventDefault();
								handleScroll("features");
							}}
							className="text-sm font-medium transition-colors hover:text-primary"
						>
							Features
						</Link>
						<Link
							to="#how-it-works"
							onClick={(e) => {
								e.preventDefault();
								handleScroll("how-it-works");
							}}
							className="text-sm font-medium transition-colors hover:text-primary"
						>
							How It Works
						</Link>
						<Link
							to="#use-cases"
							onClick={(e) => {
								e.preventDefault();
								handleScroll("use-cases");
							}}
							className="text-sm font-medium transition-colors hover:text-primary"
						>
							Use Cases
						</Link>
					</nav>
					<div className="flex items-center gap-4">
						<ModeToggle />
						<Button variant="outline" asChild>
							<Link to="/auth">Sign In</Link>
						</Button>
						<Button asChild className="hidden md:block">
							<Link to="/auth">Get Started</Link>
						</Button>
					</div>
				</div>
			</header>
			<main className="flex-1">
				{/* Hero Section */}
				<section className="w-full py-12 md:py-24 lg:py-32 xl:py-48">
					<div className="container mx-auto px-4 md:px-6">
						<div className="grid gap-8 md:gap-12 lg:grid-cols-[1fr_1fr]">
							<div className="flex flex-col justify-center space-y-4">
								<div className="space-y-2">
									<h1 className="text-3xl font-bold tracking-tighter sm:text-5xl xl:text-6xl/none">
										<div className="flex flex-col">
											<span>Create JSON APIs for</span>
											<div className="h-[1.2em] overflow-hidden">
												<Typewriter
													words={[
														"frontend prototyping",
														"mobile apps",
														"hackathons",
														"MVPs",
														"testing",
														"learning",
													]}
												/>
											</div>
										</div>
									</h1>
									<p className="max-w-[600px] text-muted-foreground md:text-xl">
										RestJSON lets you create, edit, and deploy JSON-based APIs
										in seconds. No installation required, just edit and go.
									</p>
								</div>
								<div className="flex flex-col gap-2 min-[400px]:flex-row">
									<Button size="lg" className="gap-1" asChild>
										<Link to="/auth">
											Get Started <ArrowRight className="h-4 w-4" />
										</Link>
									</Button>
									<Button
										size="lg"
										variant="outline"
										asChild
										onClick={(e) => {
											e.preventDefault();
											handleScroll("how-it-works");
										}}
									>
										<Link to="#how-it-works">How It Works</Link>
									</Button>
								</div>
							</div>
							<div className="rounded-xl border bg-muted/50 p-4 md:p-6 lg:p-8 w-full overflow-hidden">
								<div className="flex flex-col space-y-2">
									<div className="flex items-center gap-2 text-sm text-muted-foreground">
										<span className="font-medium">users.json</span>
									</div>
									<div className="w-full overflow-auto">
										<pre className="rounded-md bg-muted p-4 text-sm w-full overflow-x-auto">
											{`{
  "users": [
    { "id": 1, "name": "John Doe", "email": "john@example.com" },
    { "id": 2, "name": "Jane Smith", "email": "jane@example.com" }
  ],
  "posts": [
    { "id": 1, "userId": 1, "title": "Hello World", "body": "..." }
  ]
}`}
										</pre>
									</div>

									<div className="mt-4 text-sm text-muted-foreground">
										<p>Your API is ready at:</p>
										<div className="overflow-x-auto">
											<code className="rounded bg-muted px-1 py-0.5 font-mono text-sm">
												{`${import.meta.env.VITE_API_URL}/public/abc123/users`}
											</code>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>
				</section>

				{/* Features Section */}
				<section
					id="features"
					className="w-full py-12 md:py-24 lg:py-32 bg-muted"
				>
					<div className="container px-4 md:px-6 mx-auto">
						<div className="flex flex-col items-center justify-center space-y-4 text-center">
							<div className="space-y-2">
								<h2 className="text-3xl font-bold tracking-tighter md:text-4xl/tight">
									Why Choose RestJSON?
								</h2>
								<p className="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
									All the power of JSON Server with none of the setup. Create,
									edit, and deploy JSON APIs in seconds.
								</p>
							</div>
						</div>
						<div className="mx-auto grid max-w-5xl items-center gap-6 py-12 lg:grid-cols-3">
							<div className="flex flex-col items-center space-y-4 rounded-lg border p-6">
								<div className="rounded-full bg-primary/10 p-3">
									<Edit className="h-6 w-6 text-primary" />
								</div>
								<h3 className="text-xl font-bold">Easy Editing</h3>
								<p className="text-center text-muted-foreground">
									Edit your JSON files directly in the browser with our
									intuitive editor.
								</p>
							</div>
							<div className="flex flex-col items-center space-y-4 rounded-lg border p-6">
								<div className="rounded-full bg-primary/10 p-3">
									<Globe className="h-6 w-6 text-primary" />
								</div>
								<h3 className="text-xl font-bold">Instant API</h3>
								<p className="text-center text-muted-foreground">
									Your JSON structure automatically becomes RESTful API
									endpoints.
								</p>
							</div>
							<div className="flex flex-col items-center space-y-4 rounded-lg border p-6">
								<div className="rounded-full bg-primary/10 p-3">
									<Key className="h-6 w-6 text-primary" />
								</div>

								<h3 className="text-xl font-bold">Secure Access</h3>
								<p className="text-center text-muted-foreground">
									Control access to your API with API keys and authentication.
								</p>
							</div>
						</div>
					</div>
				</section>

				{/* How It Works Section */}

				<section id="how-it-works" className="w-full py-12 md:py-24 lg:py-32">
					<div className="container px-4 md:px-6 mx-auto">
						<div className="flex flex-col items-center justify-center space-y-4 text-center">
							<div className="space-y-2">
								<h2 className="text-3xl font-bold tracking-tighter md:text-4xl/tight">
									How It Works
								</h2>

								<p className="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
									Get your API up and running in minutes, not hours.
								</p>
							</div>
						</div>
						<div className="mx-auto grid max-w-5xl gap-8 py-12">
							<div className="grid gap-8 md:grid-cols-[1fr_2fr]">
								<div className="flex items-center justify-center rounded-lg border bg-muted p-8">
									<div className="flex h-16 w-16 items-center justify-center rounded-full bg-primary text-3xl font-bold text-primary-foreground">
										1
									</div>
								</div>
								<div className="flex flex-col justify-center space-y-4">
									<h3 className="text-xl font-bold">Sign In</h3>
									<p className="text-muted-foreground">
										Create an account or sign in to get started with RestJSON.
									</p>
								</div>
							</div>
							<div className="grid gap-8 md:grid-cols-[1fr_2fr]">
								<div className="flex items-center justify-center rounded-lg border bg-muted p-8">
									<div className="flex h-16 w-16 items-center justify-center rounded-full bg-primary text-3xl font-bold text-primary-foreground">
										2
									</div>
								</div>
								<div className="flex flex-col justify-center space-y-4">
									<h3 className="text-xl font-bold">Create JSON File</h3>
									<p className="text-muted-foreground">
										Create a new JSON file or upload an existing one to get
										started.
									</p>
								</div>
							</div>
							<div className="grid gap-8 md:grid-cols-[1fr_2fr]">
								<div className="flex items-center justify-center rounded-lg border bg-muted p-8">
									<div className="flex h-16 w-16 items-center justify-center rounded-full bg-primary text-3xl font-bold text-primary-foreground">
										3
									</div>
								</div>
								<div className="flex flex-col justify-center space-y-4">
									<h3 className="text-xl font-bold">Edit JSON File</h3>
									<p className="text-muted-foreground">
										Use our intuitive editor to structure your data however you
										want.
									</p>
								</div>
							</div>
							<div className="grid gap-8 md:grid-cols-[1fr_2fr]">
								<div className="flex items-center justify-center rounded-lg border bg-muted p-8">
									<div className="flex h-16 w-16 items-center justify-center rounded-full bg-primary text-3xl font-bold text-primary-foreground">
										4
									</div>
								</div>
								<div className="flex flex-col justify-center space-y-4">
									<h3 className="text-xl font-bold">Create API Key</h3>
									<p className="text-muted-foreground">
										Generate an API key to secure access to your endpoints.
									</p>
								</div>
							</div>
							<div className="grid gap-8 md:grid-cols-[1fr_2fr]">
								<div className="flex items-center justify-center rounded-lg border bg-muted p-8">
									<div className="flex h-16 w-16 items-center justify-center rounded-full bg-primary text-3xl font-bold text-primary-foreground">
										5
									</div>
								</div>
								<div className="flex flex-col justify-center space-y-4">
									<h3 className="text-xl font-bold">Call API</h3>
									<p className="text-muted-foreground">
										Use your API endpoints from any application or service.
									</p>
								</div>
							</div>
						</div>
					</div>
				</section>

				{/* Code Example Section */}
				<section className="w-full py-12 md:py-24 lg:py-32 bg-muted">
					<div className="container px-4 md:px-6 mx-auto">
						<div className="flex flex-col items-center justify-center space-y-4 text-center">
							<div className="space-y-2">
								<h2 className="text-3xl font-bold tracking-tighter md:text-4xl/tight">
									Simple to Use
								</h2>
								<p className="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
									Just a few lines of code to get your data.
								</p>
							</div>
						</div>
						<div className="mx-auto max-w-3xl mt-8">
							<div className="rounded-xl border bg-card p-6">
								<div className="flex items-center gap-2 text-sm text-muted-foreground mb-2">
									<Code className="h-4 w-4" />

									<span className="font-medium">Example: Fetching users</span>
								</div>
								<pre className="rounded-md bg-muted p-4 overflow-auto text-sm">
									{`// Fetch all users
fetch('${import.meta.env.VITE_API_URL}/public/FILE_ID/users', {
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY'
  }
})
.then(response => response.json())
.then(data => console.log(data));

// Get a specific user with id 1
fetch('${import.meta.env.VITE_API_URL}/public/FILE_ID/users/1', {
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY'
  }
})
.then(response => response.json())
.then(data => console.log(data));`}
								</pre>
							</div>
						</div>
					</div>
				</section>

				{/* Use Cases Section */}
				<section id="use-cases" className="w-full py-12 md:py-24 lg:py-32">
					<div className="container px-4 md:px-6 mx-auto">
						<div className="flex flex-col items-center justify-center space-y-4 text-center">
							<div className="space-y-2">
								<h2 className="text-3xl font-bold tracking-tighter md:text-4xl/tight">
									Use Cases
								</h2>
								<p className="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
									RestJSON is perfect for a variety of scenarios.
								</p>
							</div>
						</div>
						<div className="mx-auto grid max-w-5xl gap-8 py-12 md:grid-cols-2 lg:grid-cols-3">
							<div className="flex flex-col space-y-3 rounded-lg border p-6">
								<h3 className="text-xl font-bold">Frontend Prototyping</h3>
								<p className="text-muted-foreground">
									Quickly create mock APIs for your frontend applications
									without setting up a backend.
								</p>
							</div>
							<div className="flex flex-col space-y-3 rounded-lg border p-6">
								<h3 className="text-xl font-bold">Mobile App Development</h3>
								<p className="text-muted-foreground">
									Test your mobile app against real API endpoints with
									customizable data.
								</p>
							</div>

							<div className="flex flex-col space-y-3 rounded-lg border p-6">
								<h3 className="text-xl font-bold">Hackathons & MVPs</h3>
								<p className="text-muted-foreground">
									Build functional prototypes in hours instead of days by using
									RestJSON as your backend.
								</p>
							</div>
							<div className="flex flex-col space-y-3 rounded-lg border p-6">
								<h3 className="text-xl font-bold">Education & Learning</h3>
								<p className="text-muted-foreground">
									Teach API concepts without the complexity of setting up
									servers and databases.
								</p>
							</div>
							<div className="flex flex-col space-y-3 rounded-lg border p-6">
								<h3 className="text-xl font-bold">Small Projects</h3>
								<p className="text-muted-foreground">
									Perfect for small projects that need a simple data store
									without the overhead of a full backend.
								</p>
							</div>
							<div className="flex flex-col space-y-3 rounded-lg border p-6">
								<h3 className="text-xl font-bold">API Mocking</h3>
								<p className="text-muted-foreground">
									Create mock APIs that match your production API structure for
									testing and development.
								</p>
							</div>
						</div>
					</div>
				</section>

				{/* CTA Section */}
				<section className="w-full py-12 md:py-24 lg:py-32 bg-muted">
					<div className="container grid items-center justify-center gap-4 px-4 text-center md:px-6 mx-auto">
						<div className="space-y-3">
							<h2 className="text-3xl font-bold tracking-tighter md:text-4xl/tight">
								Ready to simplify your API development?
							</h2>
							<p className="mx-auto max-w-[600px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
								Click on the button below to get started.
							</p>
						</div>
						<div className="mx-auto flex flex-col gap-2 min-[400px]:flex-row justify-center">
							<Button size="lg" className="gap-1" asChild>
								<Link to="/auth">
									Get Started <ArrowRight className="h-4 w-4" />
								</Link>
							</Button>
						</div>
					</div>
				</section>
			</main>
			<footer className="w-full border-t py-6">
				<div className="container mx-auto px-4 flex flex-col items-center justify-between gap-4 md:flex-row md:gap-0">
					<div className="flex items-center gap-2">
						<FileJson className="h-5 w-5" />
						<p className="text-sm font-medium">
							RestJSON © {new Date().getFullYear()}
						</p>
					</div>
					<nav className="flex gap-4 sm:gap-6">
						<Link
							to="#"
							className="text-sm font-medium hover:underline underline-offset-4"
						>
							Terms
						</Link>
						<Link
							to="#"
							className="text-sm font-medium hover:underline underline-offset-4"
						>
							Privacy
						</Link>

						<Link
							to="#"
							className="text-sm font-medium hover:underline underline-offset-4"
						>
							Contact
						</Link>
					</nav>
				</div>
			</footer>
		</div>
	);
}
