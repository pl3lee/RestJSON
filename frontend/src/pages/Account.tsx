import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { createApiKey, getAllApiKeys } from "@/lib/api";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useState } from "react";

export function Account() {
	const [apiKey, setApiKey] = useState<string | null>(null);
	const { data: apiKeysMetadata, isLoading: apiKeysMetadataLoading } = useQuery(
		{
			queryKey: ["apikeysmetadata"],
			queryFn: getAllApiKeys,
		},
	);
	const createApiKeyMutation = useMutation({
		mutationFn: createApiKey,
		onSuccess: (data) => setApiKey(data!.apiKey),
	});

	return (
		<div>
			<Card>
				<CardHeader>
					<CardTitle>API Key</CardTitle>
				</CardHeader>
				<CardContent>
					{apiKeysMetadata.length == 0
						? "No API Keys yet"
						: apiKeysMetadata.map((apiKey) => {
							<div className="flex flex-row gap-2"></div>;
						})}
					{!apiKey ? (
						<Button
							onClick={() => createApiKeyMutation.mutate()}
							disabled={createApiKeyMutation.isPending}
						>
							{createApiKeyMutation.isPending ? "Loading" : "Generate API Key"}
						</Button>
					) : (
						<div>
							<h2>Your API Key:</h2>
							<Input readOnly value={apiKey} />
						</div>
					)}
				</CardContent>
			</Card>
		</div>
	);
}
