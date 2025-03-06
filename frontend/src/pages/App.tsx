import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router";
import { useAuth } from "../hooks/useAuth";

export function App() {
	const navigate = useNavigate();
	const { user, isLoading, logout, isLoggedIn } = useAuth();
	if (isLoading) {
		return <div>Loading...</div>;
	}
	if (!isLoggedIn) {
		navigate("/auth");
	}
	return (
		<div className="flex flex-col gap-5">
			Welcome {user.name}!
			<Button type="button" onClick={() => logout()}>
				Logout
			</Button>
		</div>
	);
}
