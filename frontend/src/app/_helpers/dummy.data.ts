import {User} from "../models/user.model";


export function NewUser(): User {
  return {
    bannedAt: null,
    permissions: [],
    email: "",
    firstName: "",
    id: 0,
    lastName: "",
    token: "",
    username: "",
    roles: [],
    createdAt: new Date(Date.now())
  };
}
