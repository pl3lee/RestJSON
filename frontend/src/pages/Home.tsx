import { ModeToggle } from "@/components/mode-toggle";
import { Button } from "@/components/ui/button";
import { fetchPublic } from "@/lib/api";
import { Link } from "react-router";

function Home() {
	return (
		<div className="text-red-500 flex flex-col align-middle">
			<title>RestJSON</title>
			<ModeToggle />
			<Button asChild>
				<Link to="/auth">Get Started</Link>
			</Button>
			<Button onClick={async () => await fetchPublic()}>
				Fetch from public
			</Button>
		</div>
	);
}

export default Home;
