import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { getAllJSONMetadata } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../hooks/useAuth";

export function App() {
	const navigate = useNavigate();
	const [newFileName, setNewFileName] = useState("");
	const { user, isLoading: isLoadingUser, logout, isLoggedIn } = useAuth();
	const { data: jsonFiles, isLoading: isLoadingFiles } = useQuery({
		queryKey: ["jsonfiles"],
		queryFn: getAllJSONMetadata,
		enabled: !!user,
	});
	if (isLoadingUser) {
		return <div>Loading...</div>;
	}
	if (!isLoggedIn) {
		navigate("/auth");
		return null;
	}

	return (
		<div className="flex flex-col gap-5">
			<form>
				<Input
					value={newFileName}
					onChange={(e) => setNewFileName(e.target.value)}
				/>
				<Button type="button" onClick={() => console.log("hello")}>
					Create JSON
				</Button>
			</form>
			<div className="flex flex-col gap-2">
				{!isLoadingFiles &&
					jsonFiles?.map((file) => (
						<Button
							key={file.id}
							type="button"
							onClick={() => navigate(`/app/jsonfile/${file.id}`)}
						>
							{file.fileName}
						</Button>
					))}
			</div>
			<Button type="button" onClick={() => logout()}>
				Logout
			</Button>
		</div>
	);
}
