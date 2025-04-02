import { AccountManager } from "@/components/account-manager";
import { ApiKeysManager } from "@/components/api-keys-manager";
import { SubscriptionsManager } from "@/components/subscriptions-manager";

export function Account() {
	return (
		<div>
			<title>Account - RestJSON</title>
			<ApiKeysManager />
			<SubscriptionsManager />
			<AccountManager />
		</div>
	);
}
