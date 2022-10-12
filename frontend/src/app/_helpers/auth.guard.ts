import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot} from '@angular/router';
import {AuthenticationService} from '../_services/authentication.service';
import {MdbNotificationService} from "mdb-angular-ui-kit/notification";
import {AlertComponent} from "../_components/alert/alert.component";

@Injectable({
  providedIn: 'root'
})
export class AuthGuard implements CanActivate {
  constructor(
    private router: Router,
    private authenticationService: AuthenticationService,
    private notificationService: MdbNotificationService,
  ) {
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean {
    const currentUser = this.authenticationService.currentUserValue;
    if (currentUser) {
      if (!this.authenticationService.hasAnyRole(route.data['roles'])) {
        this.router.navigate(['/']).then(() =>
          this.notificationService.open(AlertComponent, {
            data: {
              message: 'Access denied',
              color: 'warning',
            },
            stacking: true,
            position: "top-center",
          })
        );

        return false;
      }

      return true;
    }

    this.router.navigate(['/login'], {queryParams: {returnUrl: state.url}});
    this.notificationService.open(AlertComponent, {
      data: {
        message: 'Please login to view this page',
        color: 'info',
      },
      stacking: true,
      position: "top-center",
    })
    return false;
  }
}
