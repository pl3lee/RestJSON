import { Button } from "@/components/ui/button";
import { createJSON, getJSONFiles } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "react-router";
import { useAuth } from "../hooks/useAuth";

export function App() {
	const navigate = useNavigate();
	const { user, isLoading: isLoadingUser, logout, isLoggedIn } = useAuth();
	const { data: jsonFiles, isLoading: isLoadingFiles } = useQuery({
		queryKey: ["jsonfiles"],
		queryFn: getJSONFiles,
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
			Welcome {user.name}!
			<Button type="button" onClick={() => createJSON()}>
				Create JSON
			</Button>
			<div className="flex flex-col gap-2">
				{!isLoadingFiles &&
					jsonFiles?.map((file: { ID: string; FileName: string }) => (
						<Button
							key={file.ID}
							type="button"
							onClick={() => navigate(`/app/jsonfile/${file.ID}`)}
						>
							{file.FileName}
						</Button>
					))}
			</div>
			<Button type="button" onClick={() => logout()}>
				Logout
			</Button>
		</div>
	);
}
