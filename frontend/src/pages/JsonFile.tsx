import { getJSON } from "@/lib/api";
import Editor from "@monaco-editor/react";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";

export function JsonFile() {
	const { fileId } = useParams();
	const { data: jsonFile, isLoading: jsonFileLoading } = useQuery({
		queryKey: [`jsonfile-${fileId}`],
		queryFn: () => getJSON(fileId!),
		enabled: !!fileId,
	});
	console.log(jsonFile);
	if (jsonFileLoading) {
		return <div>Loading...</div>;
	}
	const jsonString = jsonFile ? JSON.stringify(jsonFile, null, 2) : "";

	return (
		<Editor height="90vh" defaultLanguage="json" defaultValue={jsonString} />
	);
}
