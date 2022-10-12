import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { User } from '../models/user.model';

@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {
  private currentUserSubject: BehaviorSubject<User | null>;

  constructor(private http: HttpClient) {
    this.currentUserSubject = new BehaviorSubject<any>(JSON.parse(<string>localStorage.getItem('currentUser')));
  }

  public get currentUserValue(): User | null {
    return this.currentUserSubject.value;
  }

  hasAnyRole(roleNames: string[]): boolean {
    if (this.currentUserValue == null) return false;
    return this.currentUserValue.roles.some(userRole => roleNames.some(requiredRole => requiredRole.toUpperCase() == userRole.name.toUpperCase()));
  }

  login(email: string, password: string): Observable<User> {
    return this.http.post<any>(`${environment.ACCOUNT_API_URL}/login`, { email, password })
      .pipe(map(resp => {
        resp.createdAt = new Date(resp.createdAt);
        localStorage.setItem('currentUser', JSON.stringify(resp));
        this.currentUserSubject.next(resp);
        return resp;
      }));
  }

  logout(): void {
    localStorage.removeItem('currentUser');
    this.currentUserSubject.next(null);
  }

  requestRegister(registerInformation: any): Observable<any> {
    return this.http.post(`${environment.ACCOUNT_API_URL}/register`, registerInformation).pipe(map(response => {
      return response;
    }))
  }
}
