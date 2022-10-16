import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {environment} from "../../environments/environment";
import {map} from "rxjs/operators";
import {User} from "../models/user.model";

@Injectable({
  providedIn: 'root'
})
export class UsersService {

  constructor(protected http: HttpClient) { }

  getUser(id: number, forceRefresh: boolean = false): Observable<any> {

    return this.http.get(`${environment.ACCOUNT_API_URL}/users/${id}`).pipe(map( resp => {
      return resp
    }));
  }

  getAllUsers(): Observable<User[]> {
    return this.http.get<any>(`${environment.ACCOUNT_API_URL}/users`).pipe(
      map(resp => {
        return resp.users;
      })
    );
  }

  updateUserDetails(id: number, username: string, email: string, firstName: string, lastName: string): Observable<any> {
    return this.http.put(`${environment.ACCOUNT_API_URL}/users/${id}`, {
      username, email, firstName, lastName,
    });
  }


  updateUserPassword(id: number, currentPassword: string, newPassword: string): Observable<any> {
    return this.http.put(`${environment.ACCOUNT_API_URL}/users/${id}/password`, {
      currentPassword, newPassword
    });
  }

  banUser(userId: number): Observable<any> {
    return this.http.post(`${environment.ACCOUNT_API_URL}/users/${userId}/ban`, {});
  }

  unBanUser(userId: number): Observable<any> {
    return this.http.post(`${environment.ACCOUNT_API_URL}/users/${userId}/unban`, {});
  }
}
