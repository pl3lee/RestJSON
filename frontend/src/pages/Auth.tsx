import { Button } from "@/components/ui/button";
import { useGoogleLogin } from "@react-oauth/google";
import { useNavigate } from "react-router";
import { login } from "@/lib/api";

export function Auth() {
    return <div className="flex flex-col gap-2">
        <Button onClick={login}>Login with google</Button>
    </div>
}
