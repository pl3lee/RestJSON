import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { deleteJSONFile } from "@/lib/api/jsonFiles";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Trash } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

export function DeleteFileButton({
	fileId,
	onDeleteSuccess = () => { },
}: { fileId: string; onDeleteSuccess?: () => void }) {
	const queryClient = useQueryClient();
	const [isDialogOpen, setIsDialogOpen] = useState(false);

	const deleteMutation = useMutation({
		mutationFn: deleteJSONFile,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["jsonfiles"] });
			toast.success("Deleted file successfully");
			onDeleteSuccess();
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});

	return (
		<>
			<Button
				variant="destructive"
				type="button"
				onClick={(e) => {
					e.stopPropagation();
					setIsDialogOpen(true);
				}}
				disabled={deleteMutation.isPending}
			>
				<Trash />
			</Button>
			<Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>Confirm Deletion</DialogTitle>
					</DialogHeader>
					<p>
						Are you sure you want to delete this file? This action cannot be
						undone.
					</p>
					<DialogFooter>
						<Button variant="secondary" onClick={() => setIsDialogOpen(false)}>
							Cancel
						</Button>
						<Button
							variant="destructive"
							onClick={() => deleteMutation.mutate(fileId)}
							disabled={deleteMutation.isPending}
						>
							Delete
						</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>
		</>
	);
}
