import { DeleteFileButton } from "@/components/delete-file-button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { useAuth } from "@/hooks/useAuth";
import { createJSONFile, getAllJSONMetadata } from "@/lib/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns";
import { AlertCircle, CalendarDays, File } from "lucide-react";
import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { toast } from "sonner";

export function App() {
	const navigate = useNavigate();
	const queryClient = useQueryClient();
	const [newFileName, setNewFileName] = useState("");

	const {
		user,
		isLoading: isLoadingUser,
		isLoggedIn,
		isError: userError,
	} = useAuth();
	const {
		data: jsonFiles,
		isLoading: isLoadingFiles,
		isError: filesError,
	} = useQuery({
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
	if (!isLoggedIn || userError) {
		navigate("/auth");
		return null;
	}
	const formatDate = (dateString: string) => {
		const date = new Date(dateString);
		return format(date, "MMM d, yyyy 'at' h:mm a");
	};

	if (filesError) {
		return (
			<Alert variant="destructive">
				<AlertCircle className="h-4 w-4" />
				<AlertTitle>Error</AlertTitle>
				<AlertDescription>
					Cannot fetch your files, please refresh or try again later.
				</AlertDescription>
			</Alert>
		);
	}

	return (
		<div className="flex flex-col gap-5">
			<title>Your JSON files - RestJSON</title>
			{isLoadingUser ? (
				<Skeleton className="h-8 w-full" />
			) : (
				<>
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
						{isLoadingFiles ? (
							<>
								<Skeleton className="h-24 w-full" />
								<Skeleton className="h-24 w-full" />
								<Skeleton className="h-24 w-full" />
							</>
						) : (
							!isLoadingFiles &&
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
															<span>
																Modified: {formatDate(file.modifiedAt)}
															</span>
														</div>
													</div>
												</Link>
												<DeleteFileButton fileId={file.id} />
											</CardContent>
										</Card>
									)))
						)}
					</div>
				</>
			)}
		</div>
	);
}
