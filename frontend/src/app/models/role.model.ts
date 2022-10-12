import {UserPermission} from "./user-permission.model";

export interface Role {
  id: number,
  name: string;
  permissions: UserPermission[];
}

export function getRoleBadgeClasses(role: Role): string {
  if (role.name.toUpperCase() == 'SUPER ADMIN') return 'badge bg-secondary text-white';
  if (role.name.toUpperCase() == 'MEMBER') return 'badge bg-light text-dark'
  return 'badge bg-dark text-light';
}
