import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { fetchMe, logout } from "../lib/api";

export function useAuth() {
	const queryClient = useQueryClient();
	const {
		data: user,
		isLoading,
		error,
	} = useQuery({
		queryKey: ["auth"],
		queryFn: fetchMe,
		retry: false,
		staleTime: 1000 * 60 * 5, // 5 mins
	});

	const logoutMutation = useMutation({
		mutationFn: logout,
		onSuccess: () => {
			queryClient.resetQueries({
				queryKey: undefined,
				excat: false,
				throwOnError: false,
				cancelRefetch: true,
			});
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});

	return {
		user,
		isLoading,
		error,
		logout: logoutMutation.mutate,
		isLoggedIn: !!user,
	};
}
