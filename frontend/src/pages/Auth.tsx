import { Button } from "@/components/ui/button";
import { login } from "@/lib/api";

export function Auth() {
	return (
		<div className="flex flex-col gap-2">
			<title>Login - RestJSON</title>
			<Button onClick={login}>Login with google</Button>
		</div>
	);
}
