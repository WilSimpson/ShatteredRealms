import { Role } from "./role.model";
import {UserPermission} from "./user-permission.model";

export interface User {
    id: number;
    firstName: string;
    lastName: string;
    email: string;
    username: string;
    roles: Role[];
    permissions: UserPermission[];
    token: string;
    createdAt: Date;
    bannedAt: Date | null;
}
