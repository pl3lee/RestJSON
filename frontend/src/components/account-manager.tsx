import { deleteAccount } from "@/lib/api/auth";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { toast } from "sonner";
import { Button } from "./ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "./ui/card";
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "./ui/dialog";
import { Input } from "./ui/input";

export function AccountManager() {
	const queryClient = useQueryClient();
	const [isDialogOpen, setIsDialogOpen] = useState(false);
	const [confirmationText, setConfirmationText] = useState("");

	const deleteAccountMutation = useMutation({
		mutationFn: deleteAccount,
		onSuccess: () => {
			queryClient.resetQueries({
				queryKey: undefined,
				exact: false,
				throwOnError: false,
				cancelRefetch: true,
			});
			toast.success("Account deleted successfully.");
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});

	const handleDelete = () => {
		if (confirmationText === "delete my account") {
			deleteAccountMutation.mutate();
			setIsDialogOpen(false);
			setConfirmationText("");
		} else {
			toast.error("Confirmation text does not match.");
		}
	};

	return (
		<div className="container mx-auto py-10">
			<Card>
				<CardHeader className="flex flex-row items-center justify-between">
					<CardTitle className="text-2xl">Your Account</CardTitle>
				</CardHeader>
				<CardContent>
					<Card className="border-destructive">
						<CardHeader>
							<CardTitle className="text-xl text-destructive">
								Delete Account
							</CardTitle>
							<CardDescription>
								Permanently remove your account and all of its contents. This
								action is not reversible, please continue with caution.
							</CardDescription>
						</CardHeader>
						<CardContent>
							<Button
								variant="destructive"
								onClick={() => setIsDialogOpen(true)}
							>
								Delete Your Account
							</Button>
						</CardContent>
					</Card>
				</CardContent>
			</Card>

			<Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>Confirm Account Deletion</DialogTitle>
					</DialogHeader>
					<p className="text-sm text-muted-foreground">
						This action cannot be undone. To confirm, type{" "}
						<strong>"delete my account"</strong> below.
					</p>
					<Input
						placeholder="Type here..."
						value={confirmationText}
						onChange={(e) => setConfirmationText(e.target.value)}
					/>
					<DialogFooter>
						<Button variant="secondary" onClick={() => setIsDialogOpen(false)}>
							Cancel
						</Button>
						<Button
							variant="destructive"
							onClick={handleDelete}
							disabled={confirmationText !== "delete my account"}
						>
							Delete Account
						</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>
		</div>
	);
}
