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
			<button type="button" onClick={() => logout()}>
				Logout
			</button>
		</div>
	);
}
