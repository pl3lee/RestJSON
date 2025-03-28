import CodeBlock from "@/components/code-block";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import {
	type Route,
	getDynamicRoutes,
	getJSONMetadata,
	renameJSONFile,
} from "@/lib/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AlertCircle, Check, Code, Edit2, FileJson, Key } from "lucide-react";
import { useState } from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";
import { DeleteFileButton } from "./delete-file-button";

interface JsonFileTopbarProps {
	fileId: string;
	saved: boolean;
}

export default function JsonFileTopbar({ fileId, saved }: JsonFileTopbarProps) {
	const navigate = useNavigate();
	const queryClient = useQueryClient();
	const { data: jsonMetadata, isLoading: jsonMetadataLoading } = useQuery({
		queryKey: [`jsonmetadata-${fileId}`],
		queryFn: async () => await getJSONMetadata(fileId!),
		enabled: !!fileId,
	});

	const renameMutation = useMutation({
		mutationFn: renameJSONFile,
		onSuccess: () => {
			queryClient.invalidateQueries({
				queryKey: [`jsonmetadata-${fileId}`],
			});
			toast.success("Renamed file successfully");
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
	const [isRenaming, setIsRenaming] = useState(false);

	const [nameInput, setNameInput] = useState(
		jsonMetadata ? jsonMetadata.fileName : "",
	);

	const handleRenameKeyDown = (e: React.KeyboardEvent | React.FocusEvent) => {
		if (e.type === "keydown" && (e as React.KeyboardEvent).key === "Escape") {
			setNameInput(jsonMetadata!.fileName);
			setIsRenaming(false);
		} else if (
			e.type === "blur" ||
			(e.type === "keydown" && (e as React.KeyboardEvent).key === "Enter")
		) {
			if (nameInput === "") {
				toast.error("File name cannot be empty!");
				return;
			}
			renameMutation.mutate({
				name: nameInput,
				fileId,
			});
			setIsRenaming(false);
		}
	};

	const endpoint = `${import.meta.env.VITE_API_URL}/public/${fileId}`;
	const { data: routesData, isLoading: routesLoading } = useQuery({
		queryFn: async () => await getDynamicRoutes(fileId),
		queryKey: [`dynamic-${fileId}`],
	});
	const buildRoutesString = (routes: Route[]) => {
		return routes
			.map(
				(route, index) =>
					`${index + 1}. ${route.method}    ${route.url}    ${route.description}`,
			)
			.join("\n");
	};

	return (
		<div className="flex justify-between w-full gap-2 px-4 py-2 border-b bg-background items-center flex-row">
			<div className="flex items-center gap-2">
				<FileJson className="h-5 w-5 text-muted-foreground" />

				{isRenaming ? (
					<Input
						value={nameInput}
						onChange={(e) => setNameInput(e.target.value)}
						onBlur={handleRenameKeyDown}
						onKeyDown={handleRenameKeyDown}
						className="h-8"
						autoFocus
					/>
				) : jsonMetadataLoading ? (
					<Skeleton className="h-8 w-32" />
				) : (
					<div className="flex items-center gap-2 flex-grow max-w-[30dvw]">
						<span className="text-sm font-medium md:text-base text-nowrap overflow-hidden text-ellipsis">
							{jsonMetadata?.fileName}
						</span>
						<Button
							variant="ghost"
							size="icon"
							onClick={() => setIsRenaming(true)}
							className="h-4 w-4"
						>
							<Edit2 className="h-4 w-4" />
							<span className="sr-only">Rename</span>
						</Button>
					</div>
				)}
				<div className="flex items-center gap-2 z-50">
					{saved ? (
						<div className="flex items-center gap-1">
							<Check className="h-4 w-4 text-green-500" />
							<span className="text-sm text-green-500 font-medium">Saved</span>
						</div>
					) : (
						<Badge variant="outline" className="">
							<AlertCircle className="h-4 w-4 text-amber-500" />
							<span className="text-sm text-amber-500 font-medium">
								Unsaved
							</span>
						</Badge>
					)}
				</div>
			</div>

			<div className="flex items-center gap-2">
				<Dialog>
					<DialogTrigger asChild>
						<Button variant="outline" size="sm" className="gap-1">
							<Code className="h-4 w-4" />
							API
						</Button>
					</DialogTrigger>
					<DialogContent className="overflow-y-auto max-h-screen">
						<DialogHeader>
							<DialogTitle>API Endpoint</DialogTitle>
						</DialogHeader>
						<div className="space-y-4 py-4">
							<p className="text-sm text-muted-foreground">
								Use this endpoint to access your JSON data programmatically,
								including your API key in the Authorization header:
							</p>
							<CodeBlock code={endpoint} />

							<div className="space-y-2">
								<h4 className="text-sm font-medium">Example usage:</h4>
								<CodeBlock
									code={`const res = await fetch("${endpoint}", {
	headers: {
		Authorization: "Bearer YOUR_API_KEY"
	}
})
const data = await res.json()
console.log(data)
`}
								/>
							</div>
							{routesLoading ? (
								<Skeleton className="h-24 w-full" />
							) : (
								routesData && (
									<div className="space-y-2">
										<h4 className="text-sm font-medium">
											Dynamic routes generated:
										</h4>
										<CodeBlock code={buildRoutesString(routesData)} />
									</div>
								)
							)}
							<Alert>
								<Key className="h-4 w-4" />
								<AlertTitle>API Key</AlertTitle>
								<AlertDescription>
									You can get your API key in your account page
								</AlertDescription>
							</Alert>
						</div>
					</DialogContent>
				</Dialog>
				{jsonMetadata && (
					<DeleteFileButton
						fileId={jsonMetadata!.id}
						onDeleteSuccess={() => navigate("/app")}
					/>
				)}
			</div>
		</div>
	);
}
