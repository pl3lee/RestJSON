import { useTheme } from "@/components/theme-provider";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle } from "@/components/ui/card";
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
	getJSONFile,
	getJSONMetadata,
	renameJSONFile,
	updateJSONFile,
} from "@/lib/api";
import Editor from "@monaco-editor/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AlertCircle, Check, Pencil } from "lucide-react";
import { useState } from "react";
import { useParams } from "react-router";
import { toast } from "sonner";
import { useDebouncedCallback } from "use-debounce";

export function JsonFile() {
	const { fileId } = useParams();
	const queryClient = useQueryClient();
	const { theme } = useTheme();
	const [saved, setSaved] = useState(true);
	const { data: jsonFile, isLoading: jsonFileLoading } = useQuery({
		queryKey: [`jsonfile-${fileId}`],
		queryFn: async () => await getJSONFile(fileId!),
		enabled: !!fileId,
	});
	const jsonString = jsonFile ? JSON.stringify(jsonFile, null, 2) : "";
	const { data: jsonMetadata, isLoading: jsonMetadataLoading } = useQuery({
		queryKey: [`jsonmetadata-${fileId}`],
		queryFn: async () => await getJSONMetadata(fileId!),
		enabled: !!fileId,
	});
	const [isDialogOpen, setIsDialogOpen] = useState(false);
	const [newFileName, setNewFileName] = useState("");

	const renameMutation = useMutation({
		mutationFn: renameJSONFile,
		onSuccess: () => {
			queryClient.invalidateQueries({
				queryKey: [`jsonmetadata-${fileId}`],
			});
			setIsDialogOpen(false);
		},
	});

	const updateMutation = useMutation({
		mutationFn: (variables: { fileId: string; contents: unknown }) =>
			updateJSONFile(variables),
		onSuccess: () => {
			queryClient.invalidateQueries({
				queryKey: [`jsonfile-${fileId}`],
			});
			setSaved(true);
			toast.success("Changes saved successfully.");
		},
	});

	const handleRename = () => {
		if (newFileName === "") {
			toast.error("File name cannot be empty!");
			return;
		}
		if (newFileName === jsonMetadata?.fileName) {
			toast.error("New name cannot be the same as previous.");
			return;
		}
		renameMutation.mutate({ fileId: fileId!, name: newFileName });
	};
	const handleEditorChange = (value: string | undefined) => {
		setSaved(false);
		if (value !== undefined) {
			const normalizedValue = JSON.stringify(JSON.parse(value));
			const normalizedCurrentContents = JSON.stringify(jsonFile);
			if (normalizedValue === normalizedCurrentContents) {
				setSaved(true);
				return;
			}
		}
		debouncedHandleEditorChange(value);
	};
	const debouncedHandleEditorChange = useDebouncedCallback(
		(value: string | undefined) => {
			if (value !== undefined) {
				try {
					const contents = JSON.parse(value);
					updateMutation.mutate({ fileId: fileId!, contents });
				} catch (e) {
					console.error(e);
					return;
				}
			}
		},
		1000,
	);
	if (jsonFileLoading || jsonMetadataLoading) {
		return <div>Loading...</div>;
	}

	return (
		<div className="flex flex-col gap-2">
			<Card>
				<CardHeader className="flex flex-row items-center justify-between">
					<div className="flex flex-row justify-start items-center gap-2">
						<Button onClick={() => setIsDialogOpen(true)}>
							<Pencil />
						</Button>
						<CardTitle>
							<h1>{jsonMetadata?.fileName}</h1>
						</CardTitle>
					</div>
					<div className="flex flex-row justify-end items-center gap-2">
						{saved ? (
							<>
								<Check className="h-4 w-4 text-green-500" />
								<span className="text-sm text-green-500 font-medium">
									Saved
								</span>
							</>
						) : (
							<>
								<AlertCircle className="h-4 w-4 text-amber-500" />
								<span className="text-sm text-amber-500 font-medium">
									Unsaved
								</span>
							</>
						)}
					</div>
				</CardHeader>
				<Editor
					height="90vh"
					defaultLanguage="json"
					defaultValue={jsonString}
					theme={theme === "light" ? theme : "vs-dark"}
					onChange={handleEditorChange}
				/>
			</Card>
			<Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>Rename File</DialogTitle>
					</DialogHeader>
					<form
						id="renameForm"
						onSubmit={(e) => {
							e.preventDefault();
							handleRename();
						}}
					>
						<Input
							placeholder="New file name"
							value={newFileName}
							onChange={(e) => setNewFileName(e.target.value)}
							className="w-full"
						/>
					</form>
					<DialogFooter>
						<Button type="submit" form="renameForm">
							Save
						</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>
		</div>
	);
}
