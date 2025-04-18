export type User = {
    id: string;
    email: string;
    name: string;
};

export type FileMetadata = {
    id: string;
    userId: string;
    fileName: string;
    modifiedAt: string;
};

export type ApiKey = {
    apiKey: string;
};

export type ApiKeyMetadata = {
    hash: string;
    name: string;
    createdAt: string;
    lastUsedAt: string;
};

export type Route = {
    method: string;
    url: string;
    description: string;
};
