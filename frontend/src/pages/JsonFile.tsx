import { JsonFileEditor } from "@/components/json-file-editor";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { getJSONMetadata, renameJSONFile } from "@/lib/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Pencil } from "lucide-react";
import { useState } from "react";
import { useParams } from "react-router";
import { toast } from "sonner";

export function JsonFile() {
	const { fileId } = useParams();
	const queryClient = useQueryClient();
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
			toast.success("Renamed file successfully");
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
	if (jsonMetadataLoading) {
		return <div>Loading...</div>;
	}

	return (
		<div className="flex flex-col gap-2">
			<div className="flex flex-row justify-start items-center gap-2">
				<Button onClick={() => setIsDialogOpen(true)}>
					<Pencil />
				</Button>
				<h1>{jsonMetadata?.fileName}</h1>
			</div>

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
			<JsonFileEditor fileId={fileId!} />
		</div>
	);
}
