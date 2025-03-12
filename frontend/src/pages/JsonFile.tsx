import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { getJSONFile, getJSONMetadata, renameJSONFile } from "@/lib/api";
import Editor from "@monaco-editor/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { useParams } from "react-router";

export function JsonFile() {
	const { fileId } = useParams();
	const queryClient = useQueryClient();
	const { data: jsonFile, isLoading: jsonFileLoading } = useQuery({
		queryKey: [`jsonfile-${fileId}`],
		queryFn: async () => await getJSONFile(fileId!),
		enabled: !!fileId,
	});
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

	const handleRename = () => {
		if (newFileName !== jsonMetadata?.fileName && newFileName !== "") {
			renameMutation.mutate({ fileId: fileId!, name: newFileName });
		}
	};
	console.log(jsonFile);
	if (jsonFileLoading || jsonMetadataLoading) {
		return <div>Loading...</div>;
	}
	const jsonString = jsonFile ? JSON.stringify(jsonFile, null, 2) : "";

	return (
		<div className="container mx-auto p-4">
			<Card>
				<CardHeader>
					<CardTitle>
						<h1>{jsonMetadata?.fileName}</h1>
					</CardTitle>
					<Button onClick={() => setIsDialogOpen(true)}>Rename</Button>
				</CardHeader>
				<CardContent>
					<Editor
						height="90vh"
						defaultLanguage="json"
						defaultValue={jsonString}
					/>
				</CardContent>
			</Card>
			<Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>Rename File</DialogTitle>
					</DialogHeader>
					<Input
						placeholder="New file name"
						value={newFileName}
						onChange={(e) => setNewFileName(e.target.value)}
						className="w-full"
					/>
					<DialogFooter>
						<Button onClick={handleRename}>Save</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>
		</div>
	);
}
