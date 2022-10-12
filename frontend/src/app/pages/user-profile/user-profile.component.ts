import {Component, OnInit, TemplateRef} from '@angular/core';
import {User} from "../../models/user.model";
import {ActivatedRoute, Router} from "@angular/router";
import {AuthenticationService} from "../../_services/authentication.service";
import {UsersService} from "../../_services/users.service";
import {MdbNotificationService} from "mdb-angular-ui-kit/notification";
import {AlertComponent} from "../../_components/alert/alert.component";
import {getRoleBadgeClasses, Role} from "../../models/role.model";

@Component({
  selector: 'app-user-profile',
  templateUrl: './user-profile.component.html',
  styleUrls: ['./user-profile.component.scss']
})
export class UserProfileComponent implements OnInit {
  user: User;
  loading: boolean = true;

  constructor(protected route: ActivatedRoute,
              protected router: Router,
              protected notificationService: MdbNotificationService,
              protected authService: AuthenticationService,
              protected usersService: UsersService) {
  }

  ngOnInit(): void {
    const id = Number(this.route.snapshot.paramMap.get('user'));
    this.usersService.getUser(id).subscribe((user) => {
      this.user = user;
    }, (error) => {
      this.router.navigate(['/']).then(() => {
        if (error.status == 404) {
          this.notificationService.open(AlertComponent, {
            data: {
              message: 'User not found',
              color: 'danger',
            },
            stacking: true
          })
        } else if (error.status == 401) {
          this.notificationService.open(AlertComponent, {
            data: {
              message: 'Unauthorized',
              color: 'danger',
            },
            stacking: true
          })
        } else {
          this.notificationService.open(AlertComponent, {
            data: {
              message: 'Unknown server error. Please try again later.',
              color: 'danger',
            },
            stacking: true
          })
        }
      })
    })
  }

  getRoleHTMLClass(role: Role): string {
    return getRoleBadgeClasses(role);
  }
}
