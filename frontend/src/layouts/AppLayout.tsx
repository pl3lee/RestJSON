import { ModeToggle } from "@/components/mode-toggle";
import { Button } from "@/components/ui/button";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Skeleton } from "@/components/ui/skeleton";
import { useAuth } from "@/hooks/useAuth";
import { LogOut, User } from "lucide-react";
import { Link, Outlet, useNavigate } from "react-router";

export function AppLayout() {
	const navigate = useNavigate();
	const { user, error, isLoading, isLoggedIn, logout } = useAuth();
	if (!isLoading && !isLoggedIn) {
		navigate("/auth");
		return null;
	}
	return (
		<div className="flex flex-col h-screen">
			<header className="flex items-center justify-between border-b bg-background px-4 py-3 shadow-sm">
				<div className="flex items-center">
					<h1 className="text-xl font-bold">
						<Link to="/app">Web JSON</Link>
					</h1>
				</div>
				<div className="flex items-center gap-4">
					<ModeToggle />
					<DropdownMenu>
						<DropdownMenuTrigger asChild>
							<Button
								className="flex items-center gap-2 outline-none"
								variant="ghost"
							>
								<div className="flex flex-col items-end">
									{isLoading ? (
										<>
											<Skeleton className="w-24 h-6" />
											<Skeleton className="w-36 h-6" />
										</>
									) : error ? (
										<span>Error loading user data</span>
									) : (
										<>
											<span className="font-medium">{user?.name}</span>
											<span className="text-xs text-muted-foreground">
												{user?.email}
											</span>
										</>
									)}
								</div>
								{isLoading ? (
									<Skeleton className="h-8 w-8 rounded-full" />
								) : (
									<div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground">
										<User className="h-4 w-4" />
									</div>
								)}
							</Button>
						</DropdownMenuTrigger>
						<DropdownMenuContent align="end" className="w-56">
							<DropdownMenuLabel>My Account</DropdownMenuLabel>
							<DropdownMenuSeparator />
							<DropdownMenuItem onClick={() => logout}>
								<LogOut className="mr-2 h-4 w-4" />
								<span>Log out</span>
							</DropdownMenuItem>
						</DropdownMenuContent>
					</DropdownMenu>
				</div>
			</header>
			<main className="flex-grow p-4">
				<Outlet />
			</main>
		</div>
	);
}
