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
import { getJSONFile, getJSONMetadata, renameJSONFile } from "@/lib/api";
import Editor from "@monaco-editor/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Pencil } from "lucide-react";
import { useState } from "react";
import { useParams } from "react-router";

export function JsonFile() {
	const { fileId } = useParams();
	const queryClient = useQueryClient();
	const { theme } = useTheme();
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
		<div className="flex flex-col gap-2">
			<Card>
				<CardHeader className="flex flex-row items-center gap-2">
					<Button onClick={() => setIsDialogOpen(true)}>
						<Pencil />
					</Button>
					<CardTitle>
						<h1>{jsonMetadata?.fileName}</h1>
					</CardTitle>
				</CardHeader>
				<Editor
					height="90vh"
					defaultLanguage="json"
					defaultValue={jsonString}
					theme={theme === "light" ? theme : "vs-dark"}
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
