import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { createJSONFile, deleteJSONFile, getAllJSONMetadata } from "@/lib/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns";
import { CalendarDays, File, Trash } from "lucide-react";
import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { toast } from "sonner";
import { useAuth } from "../hooks/useAuth";

export function App() {
	const navigate = useNavigate();
	const queryClient = useQueryClient();
	const [newFileName, setNewFileName] = useState("");
	const { user, isLoading: isLoadingUser, isLoggedIn } = useAuth();
	const { data: jsonFiles, isLoading: isLoadingFiles } = useQuery({
		queryKey: ["jsonfiles"],
		queryFn: getAllJSONMetadata,
		enabled: !!user,
	});
	const createMutation = useMutation({
		mutationFn: createJSONFile,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["jsonfiles"] });
			setNewFileName("");
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
	const deleteMutation = useMutation({
		mutationFn: deleteJSONFile,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["jsonfiles"] });
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
	if (isLoadingUser) {
		return <div>Loading user...</div>;
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
			<title>Your JSON files - RestJSON</title>
			<form
				className="flex flex-row w-full gap-2"
				onSubmit={(e) => {
					e.preventDefault();
					if (newFileName === "") {
						toast.error("File name cannot be empty!");
						return;
					}
					createMutation.mutate(newFileName);
				}}
			>
				<Input
					value={newFileName}
					onChange={(e) => setNewFileName(e.target.value)}
					placeholder="Enter new file name"
				/>
				<Button type="submit" disabled={createMutation.isPending}>
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
											deleteMutation.mutate(file.id);
										}}
										disabled={deleteMutation.isPending}
									>
										<Trash />
									</Button>
								</CardContent>
							</Card>
						)))}
			</div>
		</div>
	);
}
