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
					{apiKeysMetadataLoading ? (
						<p>Loading...</p>
					) : apiKeysMetadata?.length === 0 ? (
						<p>No API Keys yet</p>
					) : (
						apiKeysMetadata?.map((key) => (
							<div key={key.hash} className="flex flex-col gap-2 p-2 border-b">
								<p>
									<strong>Created At:</strong>{" "}
									{new Date(key.createdAt).toLocaleString()}
								</p>
								<p>
									<strong>Last Used At:</strong>{" "}
									{new Date(key.lastUsedAt).toLocaleString()}
								</p>
							</div>
						))
					)}
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
							<p className="text-red-500">
								Please save this API key securely. You won't be able to see it
								again.
							</p>
						</div>
					)}
				</CardContent>
			</Card>
		</div>
	);
}
