import {Injectable, OnInit} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Role} from "../models/role.model";
import {Observable, of} from "rxjs";
import {environment} from "../../environments/environment";
import {map} from "rxjs/operators";
import {UserPermission} from "../models/user-permission.model";

const ROLES_KEY = 'roles';
const PERMISSIONS_KEY = 'permission';

@Injectable({
  providedIn: 'root'
})
export class AuthorizationService {

  constructor(private http: HttpClient) {
  }

  public getAllRoles(forceRefresh: boolean = true): Observable<Role[]> {
    if (!forceRefresh) {
      let rolesString = localStorage.getItem(ROLES_KEY);
      if (rolesString) {
        return of(JSON.parse(rolesString));
      }
    }

    return this.http.get<any>(`${environment.ACCOUNT_API_URL}/authorization/roles`)
      .pipe(map(resp => {
        localStorage.setItem(ROLES_KEY, JSON.stringify(resp.roles));
        return resp.roles;
      }));
  }

  public getAllPermissions(forceRefresh: boolean = true): Observable<UserPermission[]> {
    if (!forceRefresh) {
      let permissionsString = localStorage.getItem(PERMISSIONS_KEY);
      if (permissionsString) {
        return of(JSON.parse(permissionsString));
      }
    }

    return this.http.get<any>(`${environment.ACCOUNT_API_URL}/authorization/permissions`)
      .pipe(map(resp => {
        localStorage.setItem(PERMISSIONS_KEY, JSON.stringify(resp.permissions));
        return resp.permissions;
      }));
  }

  public addRoles(roles: Role[], userId: number) {
    return this.http.post(
      `${environment.ACCOUNT_API_URL}/users/${userId}/authorization/add`,
      {
        roles: roles.map(role => ({id: role.id})),
      }
    );
  }

  public remRoles(roles: Role[], userId: number) {
    return this.http.post(
      `${environment.ACCOUNT_API_URL}/users/${userId}/authorization/remove`,
      {
        roles: roles.map(role => ({id: role.id})),
      }
    );
  }

  public getRole(id: number): Observable<Role> {
    return this.http.get<Role>(`${environment.ACCOUNT_API_URL}/authorization/roles/${id}`).pipe(
      map(resp => {
        return resp
      })
    );
  }

  public updateRoleName(id: number, name: string) {
    return this.http.put(
      `${environment.ACCOUNT_API_URL}/authorization/roles/${id}`,
      {id, name},
    );
  }

  public updateRolePermissions(id: number, permissions: UserPermission[]) {
    return this.http.put(
      `${environment.ACCOUNT_API_URL}/authorization/roles/${id}`,
      {id, permissions},
    );
  }

  public createRole(name: string, permissions: UserPermission[]) {
    return this.http.post(
      `${environment.ACCOUNT_API_URL}/authorization/roles`,
      {name, permissions},
    );
  }

  deleteRole(id: number) {
    return this.http.delete(`${environment.ACCOUNT_API_URL}/authorization/roles/${id}`);
  }
}
