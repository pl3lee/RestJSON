import { createApiKey, deleteApiKey, getAllApiKeys } from "@/lib/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AlertTriangle, Check, Key, Trash2 } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import CodeBlock from "./code-block";
import { Alert, AlertDescription } from "./ui/alert";
import { Button } from "./ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "./ui/card";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "./ui/dialog";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "./ui/table";

export default function ApiKeysManager() {
	const queryClient = useQueryClient();
	const { data: apiKeysMetadata, isLoading: apiKeysMetadataLoading } = useQuery(
		{
			queryKey: ["apikeysmetadata"],
			queryFn: getAllApiKeys,
		},
	);
	const createApiKeyMutation = useMutation({
		mutationFn: createApiKey,
		onSuccess: (data) => {
			queryClient.invalidateQueries({
				queryKey: ["apikeysmetadata"],
			});

			// Set the new API key to display in the dialog
			setNewApiKey(data!.apiKey);

			// Close the create dialog and open the new key dialog
			setIsCreateDialogOpen(false);
			setIsNewKeyDialogOpen(true);

			// Reset the form
			setNewKeyName("");
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});

	const deleteApiKeyMutation = useMutation({
		mutationFn: deleteApiKey,
		onSuccess: () => {
			queryClient.invalidateQueries({
				queryKey: ["apikeysmetadata"],
			});
			toast.success("API Key successfully deleted");
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
	const [newKeyName, setNewKeyName] = useState("");
	const [newApiKey, setNewApiKey] = useState<string | null>(null);

	const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
	const [isNewKeyDialogOpen, setIsNewKeyDialogOpen] = useState(false);

	// Format date to a more readable format
	const formatDate = (dateString: string | null) => {
		if (!dateString) return "Never used";
		const date = new Date(dateString);
		return new Intl.DateTimeFormat("en-US", {
			year: "numeric",
			month: "short",
			day: "numeric",
			hour: "2-digit",
			minute: "2-digit",
		}).format(date);
	};

	const handleCreateApiKey = () => {
		if (!newKeyName.trim()) {
			toast.error("Please provide a name for your API key");
			return;
		}
		createApiKeyMutation.mutate(newKeyName);
	};

	if (apiKeysMetadataLoading) {
		return <>Loading...</>;
	}

	return (
		<div className="container mx-auto py-10">
			<Card>
				<CardHeader className="flex flex-row items-center justify-between">
					<div>
						<CardTitle className="text-2xl">API Keys</CardTitle>
						<CardDescription>
							Manage your API keys for accessing the API.
						</CardDescription>
					</div>

					<Dialog
						open={isCreateDialogOpen}
						onOpenChange={setIsCreateDialogOpen}
					>
						<DialogTrigger asChild>
							<Button>
								<Key className="h-4 w-4" />
								Create API Key
							</Button>
						</DialogTrigger>
						<DialogContent>
							<DialogHeader>
								<DialogTitle>Create a new API key</DialogTitle>
								<DialogDescription>
									Give your API key a name to help you identify it later.
								</DialogDescription>
							</DialogHeader>
							<div className="grid gap-4 py-4">
								<div className="grid gap-2">
									<Label htmlFor="name">API Key Name</Label>
									<Input
										id="name"
										placeholder="e.g., Production API Key"
										value={newKeyName}
										onChange={(e) => setNewKeyName(e.target.value)}
									/>
								</div>
							</div>
							<DialogFooter>
								<Button
									variant="outline"
									onClick={() => setIsCreateDialogOpen(false)}
								>
									Cancel
								</Button>
								<Button onClick={handleCreateApiKey}>Create API Key</Button>
							</DialogFooter>
						</DialogContent>
					</Dialog>

					<Dialog
						open={isNewKeyDialogOpen}
						onOpenChange={setIsNewKeyDialogOpen}
					>
						<DialogContent>
							<DialogHeader>
								<DialogTitle>Your new API key</DialogTitle>
								<DialogDescription>
									Please copy your API key now. You won't be able to see it
									again.
								</DialogDescription>
							</DialogHeader>

							<Alert variant="destructive" className="my-4">
								<AlertTriangle className="h-4 w-4" />
								<AlertDescription>
									This API key will only be displayed once and cannot be
									retrieved later.
								</AlertDescription>
							</Alert>
							<div className="relative mt-2">
								<CodeBlock code={newApiKey!} />
							</div>
							<DialogFooter>
								<Button
									onClick={() => {
										setIsNewKeyDialogOpen(false);
										setNewApiKey(null);
									}}
									className="mt-4"
								>
									<Check className="mr-2 h-4 w-4" />
									I've copied my API key
								</Button>
							</DialogFooter>
						</DialogContent>
					</Dialog>
				</CardHeader>
				<CardContent>
					{apiKeysMetadata!.length === 0 ? (
						<div className="text-center py-6 text-muted-foreground">
							No API keys found. Create your first API key to get started.
						</div>
					) : (
						<Table>
							<TableHeader>
								<TableRow>
									<TableHead>Name</TableHead>
									<TableHead>Created</TableHead>
									<TableHead>Last used</TableHead>
									<TableHead className="text-right">Actions</TableHead>
								</TableRow>
							</TableHeader>
							<TableBody>
								{apiKeysMetadata!.map((key) => (
									<TableRow key={key.hash}>
										<TableCell className="font-medium">{key.name}</TableCell>
										<TableCell>{formatDate(key.createdAt)}</TableCell>
										<TableCell>{formatDate(key.lastUsedAt)}</TableCell>
										<TableCell className="text-right">
											<Button
												variant="destructive"
												size="sm"
												onClick={() => deleteApiKeyMutation.mutate(key.hash)}
											>
												<Trash2 className="h-4 w-4" />
												Delete
											</Button>
										</TableCell>
									</TableRow>
								))}
							</TableBody>
						</Table>
					)}
				</CardContent>
				<CardFooter className="border-t px-6 py-4">
					<p className="text-sm text-muted-foreground">
						API keys provide full access to JSON files in your account. Keep
						them secure and never share them in public repositories or
						client-side code.
					</p>
				</CardFooter>
			</Card>
		</div>
	);
}
