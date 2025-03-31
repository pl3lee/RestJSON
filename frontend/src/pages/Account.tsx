import { AccountManager } from "@/components/account-manager";
import ApiKeysManager from "@/components/api-keys-manager";

export function Account() {
	return (
		<div>
			<title>Account - RestJSON</title>
			<ApiKeysManager />
			<AccountManager />
		</div>
	);
}
