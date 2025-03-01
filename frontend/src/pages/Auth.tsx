import { Button } from "@/components/ui/button";
import { login } from "@/lib/api";
import { useGoogleLogin } from "@react-oauth/google";
import { useNavigate } from "react-router";

export function Auth() {
	return (
		<div className="flex flex-col gap-2">
			<Button onClick={login}>Login with google</Button>
		</div>
	);
}
