import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Route, Routes } from "react-router";
import { ThemeProvider } from "./components/theme-provider.tsx";
import { AppLayout } from "./layouts/AppLayout.tsx";
import { Account } from "./pages/Account.tsx";
import { App } from "./pages/App.tsx";
import { Auth } from "./pages/Auth.tsx";
import { JsonFile } from "./pages/JsonFile.tsx";
import { Landing } from "./pages/Landing.tsx";

const queryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <QueryClientProvider client={queryClient}>
            <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
                <BrowserRouter>
                    <Routes>
                        <Route index element={<Landing />} />
                        <Route path="auth" element={<Auth />} />
                        <Route path="app" element={<AppLayout />}>
                            <Route index element={<App />} />
                            <Route
                                path="jsonfile/:fileId"
                                element={<JsonFile />}
                            />
                            <Route path="account" element={<Account />} />
                        </Route>
                    </Routes>
                </BrowserRouter>
            </ThemeProvider>
        </QueryClientProvider>
    </StrictMode>,
);
