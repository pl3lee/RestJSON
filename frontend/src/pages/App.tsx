import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { createJSON, getJSON } from "@/lib/api";
import { useState } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../hooks/useAuth";

export function App() {
	const navigate = useNavigate();
	const [fileId, setFileId] = useState<string>("");
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
			<Button type="button" onClick={() => createJSON()}>
				Create JSON
			</Button>
			<form
				onSubmit={async (e) => {
					e.preventDefault();
					const jsonData = await getJSON(fileId);
					console.log(jsonData);
				}}
			>
				<Input
					value={fileId}
					onChange={(e) => {
						setFileId(e.target.value);
					}}
				/>
				<Button type="submit">Fetch json</Button>
			</form>
			<Button type="button" onClick={() => logout()}>
				Logout
			</Button>
		</div>
	);
}
