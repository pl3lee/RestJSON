import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { createJSONFile, getAllJSONMetadata } from "@/lib/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns";
import { CalendarDays, File, Trash } from "lucide-react";
import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { useAuth } from "../hooks/useAuth";

export function App() {
	const navigate = useNavigate();
	const queryClient = useQueryClient();
	const [newFileName, setNewFileName] = useState("");
	const { user, isLoading: isLoadingUser, logout, isLoggedIn } = useAuth();
	const { data: jsonFiles, isLoading: isLoadingFiles } = useQuery({
		queryKey: ["jsonfiles"],
		queryFn: getAllJSONMetadata,
		enabled: !!user,
	});
	const mutation = useMutation({
		mutationFn: createJSONFile,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["jsonfiles"] });
			setNewFileName("");
		},
	});
	if (isLoadingUser) {
		return <div>Loading...</div>;
	}
	if (!isLoggedIn) {
		navigate("/auth");
		return null;
	}
	const formatDate = (dateString: string) => {
		const date = new Date(dateString);
		return format(date, "MMM d, yyyy 'at' h:mm a");
	};

	return (
		<div className="flex flex-col gap-5">
			<form
				className="flex flex-row w-full gap-2"
				onSubmit={async (e) => {
					e.preventDefault();
					mutation.mutate(newFileName);
				}}
			>
				<Input
					value={newFileName}
					onChange={(e) => setNewFileName(e.target.value)}
					placeholder="Enter new file name"
				/>
				<Button type="submit" disabled={mutation.isPending}>
					Create JSON File
				</Button>
			</form>
			<div className="flex flex-col gap-2">
				{!isLoadingFiles &&
					jsonFiles &&
					(jsonFiles.length === 0
						? "No JSON files created yet."
						: jsonFiles?.map((file) => (
								<Card key={file.id}>
									<CardContent className="flex flex-row justify-between items-center">
										<Link
											to={`/app/jsonfile/${file.id}`}
											className="w-full h-full"
										>
											<div className="flex items-center gap-2">
												<File className="h-5 w-5 text-primary" />
												<span className="font-medium">{file.fileName}</span>
											</div>
											<div className="mt-2 text-sm text-foreground">
												<div className="flex items-center gap-1">
													<CalendarDays className="h-3.5 w-3.5" />
													<span>Modified: {formatDate(file.modifiedAt)}</span>
												</div>
											</div>
										</Link>
										<Button
											variant="destructive"
											type="button"
											onClick={(e) => {
												e.stopPropagation();
												console.log("delete");
											}}
										>
											<Trash />
										</Button>
									</CardContent>
								</Card>
							)))}
			</div>
			<Button type="button" onClick={() => logout()}>
				Logout
			</Button>
		</div>
	);
}
