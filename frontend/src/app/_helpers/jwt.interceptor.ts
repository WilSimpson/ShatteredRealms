import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor, HttpResponse
} from '@angular/common/http';
import {catchError, EMPTY, Observable} from 'rxjs';
import { AuthenticationService } from '../_services/authentication.service';
import {Router} from "@angular/router";
import {MdbNotificationService} from "mdb-angular-ui-kit/notification";
import {AlertComponent} from "../_components/alert/alert.component";

@Injectable()
export class JwtInterceptor implements HttpInterceptor {
  constructor(
    private authenticationService: AuthenticationService,
    private router: Router,
    private notificationService: MdbNotificationService,
  ) { }

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const currentUser = this.authenticationService.currentUserValue;
    if (currentUser && currentUser.token) {
      request = request.clone({
        setHeaders: {
          Authorization: `Bearer ${currentUser.token}`
        }
      });
    }
    return next.handle(request).pipe(catchError(err => {
      console.log('error', err);
      if ((err.status == 401 && err.error.message == 'invalid authentication token')) {
        this.authenticationService.logout();
        this.router.navigate(['/']).then(() => {
          this.notificationService.open(AlertComponent, {
            data: {
              message: 'You session has timed out. Please sign in again.',
              color: 'warning',
              fade: true,
            },
            stacking: true,
            position: "top-center",
          })
        });
      }

      return EMPTY;
    }));
  }
}
