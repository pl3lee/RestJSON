import { JsonFileEditor } from "@/components/json-file-editor";
import JsonFileTopbar from "@/components/json-file-topbar";
import { getJSONMetadata } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { useParams } from "react-router";

export function JsonFile() {
	const { fileId } = useParams();
	const [saved, setSaved] = useState(true);
	const { data: jsonMetadata, isLoading: jsonMetadataLoading } = useQuery({
		queryKey: [`jsonmetadata-${fileId}`],
		queryFn: async () => await getJSONMetadata(fileId!),
		enabled: !!fileId,
	});

	return (
		<div className="flex flex-col gap-2">
			{jsonMetadataLoading ? (
				<title>Loading JSON File - RestJSON</title>
			) : (
				<title>{`${jsonMetadata!.fileName} - RestJSON`}</title>
			)}
			<JsonFileTopbar fileId={fileId!} saved={saved} />

			<JsonFileEditor fileId={fileId!} setSaved={setSaved} />
		</div>
	);
}
