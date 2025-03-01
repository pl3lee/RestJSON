import { ModeToggle } from "@/components/mode-toggle";
import { Button } from "@/components/ui/button";
import { Link } from "react-router";

function Home() {
	return (
		<div className="text-red-500 flex flex-col align-middle">
			<ModeToggle />
			<Button asChild>
				<Link to="/auth">Get Started</Link>
			</Button>
		</div>
	);
}

export default Home;
